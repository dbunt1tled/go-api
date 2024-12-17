package rmq

import (
	"fmt"
	"go_echo/internal/config/env"
	"go_echo/internal/config/logger"
	"go_echo/internal/util/hasher"
	"go_echo/internal/util/helper"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

type RabbitClient struct {
	sendCon *amqp091.Connection
	recCon  *amqp091.Connection
	sendCh  *amqp091.Channel
	recCh   *amqp091.Channel
}

func (rcl *RabbitClient) Publish(queueName string, body string) {
	log := logger.GetLoggerInstance()
	r := false
	for {
		for {
			_, err := rcl.channel(false, r)
			if err == nil {
				break
			}
		}
		q, err := rcl.sendCh.QueueDeclare(
			queueName,
			true,
			false,
			false,
			false,
			amqp091.Table{"x-queue-mode": "lazy"},
		)
		if err != nil {
			log.Error(fmt.Sprintf("Failed to declare queue %s: %s. Trying resend...", queueName, err.Error()))
			r = true
			continue
		}
		err = rcl.sendCh.Publish(
			"",
			q.Name,
			false,
			false,
			amqp091.Publishing{
				MessageId:    helper.Must(hasher.UUIDVv7()).String(),
				DeliveryMode: amqp091.Persistent,
				ContentType:  "text/plain",
				Body:         []byte(body),
			})
		if err != nil {
			log.Error(fmt.Sprintf("Failed to publish in queue %s: %s. Trying republish...", queueName, err.Error()))
			r = true
			continue
		}
		break
	}
}

func (rcl *RabbitClient) Consume(queueName string, f func(string) error) {
	log := logger.GetLoggerInstance()
	for {
		for {
			_, err := rcl.channel(true, true)
			if err == nil {
				break
			}
		}
		log.Info(fmt.Sprintf("Consumer %s connected to RabbitMQ", queueName))
		q, err := rcl.recCh.QueueDeclare(
			queueName,
			true,
			false,
			false,
			false,
			amqp091.Table{"x-queue-mode": "lazy"},
		)
		if err != nil {
			log.Error(fmt.Sprintf("Failed to declare queue %s: %s. Trying reconnect...", queueName, err.Error()))
			time.Sleep(time.Second * 1)
			continue
		}
		conClose := rcl.recCon.NotifyClose(make(chan *amqp091.Error))
		conBlocked := rcl.recCon.NotifyBlocked(make(chan amqp091.Blocking))
		chClose := rcl.recCh.NotifyClose(make(chan *amqp091.Error))
		msgs, err := rcl.recCh.Consume(
			q.Name,
			helper.Must(hasher.UUIDVv7()).String(),
			false,
			false,
			true,
			false,
			nil,
		)
		if err != nil {
			log.Error(fmt.Sprintf("Failed to consume queue %s: %s. Trying reconnect...", queueName, err.Error()))
			time.Sleep(time.Second * 1)
			continue
		}
		shouldStop := false
		for {
			if shouldStop {
				break
			}
			select {
			case <-conClose:
				log.Error(fmt.Sprintf("Consumer %s (queue %s) connection closed", queueName, q.Name))
				shouldStop = true
				break
			case <-conBlocked:
				log.Error(fmt.Sprintf("Consumer %s (queue %s) connection blocked", queueName, q.Name))
				shouldStop = true
				break
			case <-chClose:
				log.Error(fmt.Sprintf("Consumer %s (queue %s) channel closed", queueName, q.Name))
				shouldStop = true
				break
			case d := <-msgs:
				log.Info(fmt.Sprintf("Consumer %s (queue %s) received message: %s", queueName, q.Name, d.Body))
				err = f(string(d.Body))
				if err != nil {
					log.Error(fmt.Sprintf(
						"Consumer %s (queue %s) failed to process message(%s): %s",
						queueName,
						q.Name,
						d.Body,
						err.Error(),
					))
					_ = d.Ack(false)
					break
				}
				_ = d.Ack(true)
			}
		}
	}
}
func (rcl *RabbitClient) connect(isRec bool, reconect bool) (*amqp091.Connection, error) {
	if reconect {
		if isRec {
			rcl.recCon = nil
		} else {
			rcl.sendCon = nil
		}
	}

	if isRec && rcl.recCon != nil {
		return rcl.recCon, nil
	} else if !isRec && rcl.sendCon != nil {
		return rcl.sendCon, nil
	}

	cfg := env.GetConfigInstance()
	log := logger.GetLoggerInstance()
	amqpURL := fmt.Sprintf(
		"amqp://%s:%s@%s:%d/%s",
		cfg.RabbitMQ.User,
		cfg.RabbitMQ.Password,
		cfg.RabbitMQ.Host,
		cfg.RabbitMQ.Port,
		cfg.RabbitMQ.Vhost,
	)
	conn, err := amqp091.Dial(amqpURL)
	if err != nil {
		log.Error("Failed to connect to RabbitMQ: " + err.Error())
		time.Sleep(time.Second * 1)
		return nil, err
	}
	if isRec {
		rcl.recCon = conn
		return rcl.recCon, nil
	}
	rcl.sendCon = conn
	return rcl.sendCon, nil
}

func (rcl *RabbitClient) channel(isRec bool, recreate bool) (*amqp091.Channel, error) { //nolint:unparam
	log := logger.GetLoggerInstance()
	if recreate {
		if isRec {
			rcl.recCh = nil
		} else {
			rcl.sendCh = nil
		}
	}
	if isRec && rcl.recCon == nil {
		rcl.recCh = nil
	}
	if !isRec && rcl.sendCon == nil {
		rcl.recCh = nil
	}
	if isRec && rcl.recCh != nil {
		return rcl.recCh, nil
	} else if !isRec && rcl.sendCh != nil {
		return rcl.sendCh, nil
	}
	for {
		_, err := rcl.connect(isRec, recreate)
		if err == nil {
			break
		}
	}
	var err error
	if isRec {
		rcl.recCh, err = rcl.recCon.Channel()
	} else {
		rcl.sendCh, err = rcl.sendCon.Channel()
	}
	if err != nil {
		log.Error("Failed to create channel: " + err.Error())
		time.Sleep(1 * time.Second)
		return nil, err
	}
	if isRec {
		return rcl.recCh, err
	}
	return rcl.sendCh, err
}

func (rcl *RabbitClient) Close() {
	if rcl.sendCh != nil {
		_ = rcl.sendCh.Close()
	}
	if rcl.recCh != nil {
		_ = rcl.recCh.Close()
	}
	if rcl.sendCon != nil {
		_ = rcl.sendCon.Close()
	}
	if rcl.recCon != nil {
		_ = rcl.recCon.Close()
	}
}
