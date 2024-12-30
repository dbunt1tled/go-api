package rmq

import (
	"fmt"
	"go_echo/internal/config/env"
	"go_echo/internal/config/logger"
	"go_echo/internal/util/hasher"
	"go_echo/internal/util/helper"
	"sync"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

const (
	MailExchange = "mail-go-exchange"
	MailQueue    = "mail-go-queue"
)

var (
	rabbitMQInstance map[string]*RabbitClient //nolint:gochecknoglobals // singleton
	m                sync.Once                //nolint:gochecknoglobals // singleton
)

type RabbitClient struct {
	sendCon *amqp091.Connection
	recCon  *amqp091.Connection
	sendCh  *amqp091.Channel
	recCh   *amqp091.Channel
}

func GetRMQInstance(exchange string) *RabbitClient {
	var (
		val *RabbitClient
		ok  bool
	)
	m.Do(func() {
		rabbitMQInstance = make(map[string]*RabbitClient)
	})
	if val, ok = rabbitMQInstance[exchange]; !ok {
		val = &RabbitClient{}
		rabbitMQInstance[exchange] = val
	}

	return val
}

func (rcl *RabbitClient) Publish(exchangeName string, queueName string, action string, body string) {
	log := logger.GetLoggerInstance()
	r := false
	for {
		for {
			_, err := rcl.channel(false, r)
			if err == nil {
				break
			}
		}

		err := rcl.sendCh.Publish(
			exchangeName,
			queueName,
			false, // mandatory - we don't care if there is no queue
			false, // immediate - we don't care if there is no consumer on the queue
			amqp091.Publishing{
				MessageId:    helper.Must(hasher.UUIDVv7()).String(),
				DeliveryMode: amqp091.Persistent, // save on disk and restore on restart
				ContentType:  "application/json",
				Type:         action,
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

func (rcl *RabbitClient) Consume(
	exchangeName string,
	queueName string,
	f func(d amqp091.Delivery) error,
	concurrency int,
) {
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
			true,  // durable -queue stored after server restart
			false, // delete when unused
			false, // exclusive - only one consumer and deleted when queue is empty
			false, // no-wait - don't wait for queue to be declared if true - fast but you don't know if it has error
			amqp091.Table{"x-queue-mode": "lazy"},
		)
		if err != nil {
			log.Error(fmt.Sprintf("Failed to declare queue %s: %s. Trying reconnect...", queueName, err.Error()))
			time.Sleep(time.Second * 1)
			continue
		}

		err = rcl.recCh.ExchangeDeclare(
			exchangeName,
			"direct", // тип обменника
			true,     // durable
			false,    // autoDelete
			false,    // internal
			false,    // noWait
			nil,      // args
		)
		if err != nil {
			log.Error(fmt.Sprintf("Failed to declare exchange %s: %s. Trying reconnect...", queueName, err.Error()))
			time.Sleep(time.Second * 1)
			continue
		}

		err = rcl.recCh.QueueBind(
			queueName,
			queueName, // routing key - key in method Publish
			exchangeName,
			false,
			nil,
		)
		if err != nil {
			log.Error(fmt.Sprintf("Failed to bind queue to exchange %s: %s. Trying reconnect...", queueName, err.Error()))
			time.Sleep(time.Second * 1)
			continue
		}
		err = rcl.recCh.Qos(concurrency, 0, false)
		if err != nil {
			log.Error(fmt.Sprintf("Failed to set QoS for queue %s: %s. Trying reconnect...", queueName, err.Error()))
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
				log.Error(fmt.Sprintf("Consumer %s (queue %s) connection closed", exchangeName, q.Name))
				shouldStop = true
				break
			case <-conBlocked:
				log.Error(fmt.Sprintf("Consumer %s (queue %s) connection blocked", exchangeName, q.Name))
				shouldStop = true
				break
			case <-chClose:
				log.Error(fmt.Sprintf("Consumer %s (queue %s) channel closed", exchangeName, q.Name))
				shouldStop = true
				break
			case d := <-msgs:
				time.Sleep(time.Millisecond * 100)
				go func() {
					worker(d, f, exchangeName, q.Name)
				}()
			}
		}
	}
}

func worker(msg amqp091.Delivery, f func(d amqp091.Delivery) error, exchangeName string, queueName string) {
	var err error
	startTime := time.Now()
	log := logger.GetLoggerInstance()
	log.Info(fmt.Sprintf(
		"Consumer %s (queue %s) START processing message (%s): %s",
		exchangeName,
		queueName,
		msg.MessageId,
		msg.Body,
	))
	if err = f(msg); err == nil {
		err = msg.Ack(false)
		if err != nil {
			log.Info(fmt.Sprintf("Consumer %s (queue %s) Error Ack: %s", exchangeName, queueName, err.Error()))
		}
	} else {
		m := err.Error()
		err = msg.Nack(false, true)
		if err != nil {
			log.Info(fmt.Sprintf(
				"Consumer %s (queue %s) Error Nack %s after handler error: %s",
				exchangeName,
				queueName,
				err.Error(),
				m,
			))
		}
		log.Info(fmt.Sprintf("Consumer %s (queue %s) Error Handler: %s", exchangeName, queueName, m))
	}
	log.Info(fmt.Sprintf(
		"Consumer %s (queue %s) FINISH processing message (%s): %s --- %s",
		exchangeName,
		queueName,
		msg.MessageId,
		msg.Body,
		helper.RuntimeStatistics(startTime, false),
	))
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
