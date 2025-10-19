package rmq

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/dbunt1tled/go-api/internal/config/env"
	"github.com/dbunt1tled/go-api/internal/config/logger"
	"github.com/dbunt1tled/go-api/internal/util/hasher"
	"github.com/dbunt1tled/go-api/internal/util/helper"

	"github.com/wagslane/go-rabbitmq"
)

const (
	MailExchange = "mail-go-exchange"
	MailQueue    = "mail-go-queue"
	MaxTry       = 10
	Delay        = 100 * time.Millisecond
)

var (
	rabbitMQInstance *rabbitmq.Conn //nolint:gochecknoglobals // singleton
	m                sync.Once      //nolint:gochecknoglobals // singleton
)

func GetRMQInstance() *rabbitmq.Conn {
	m.Do(func() {
		var err error
		cfg := env.GetConfigInstance()
		log := logger.GetLoggerInstance()
		amqpURL := fmt.Sprintf(
			"amqp://%s:%s@%s:%d/%s",
			cfg.RabbitMQ.User,
			cfg.RabbitMQ.Password,
			cfg.RabbitMQ.Host,
			cfg.RabbitMQ.Port,
			strings.TrimSuffix(cfg.RabbitMQ.Vhost, "/"),
		)
		rabbitMQInstance, err = rabbitmq.NewConn(
			amqpURL,
			rabbitmq.WithConnectionOptionsLogger(log),
		)
		if err != nil {
			panic("failed to create rabbitmq client: " + err.Error())
		}
	})

	return rabbitMQInstance
}

func Publish(exchangeName string, queueName string, action string, body string) {
	log := logger.GetLoggerInstance()
	publisher, err := rabbitmq.NewPublisher(
		GetRMQInstance(),
		rabbitmq.WithPublisherOptionsExchangeName(exchangeName),
		rabbitmq.WithPublisherOptionsExchangeDeclare,
		rabbitmq.WithPublisherOptionsLogger(log),
	)
	if err != nil {
		panic(fmt.Sprintf("Failed to create publisher: %s.", err.Error()))
	}
	defer publisher.Close()

	// mandatory - we don't care if there is no queue // false
	// immediate - we don't care if there is no consumer on the queue // false
	// persistent - save on disk and restore on restart
	err = publisher.Publish(
		[]byte(body),
		[]string{queueName},
		rabbitmq.WithPublishOptionsContentType("application/json"),
		rabbitmq.WithPublishOptionsExchange(exchangeName),
		rabbitmq.WithPublishOptionsType(action),
		rabbitmq.WithPublishOptionsPersistentDelivery,
		rabbitmq.WithPublishOptionsMessageID(helper.Must(hasher.UUIDVv7()).String()),
	)
	if err != nil {
		log.Error(fmt.Sprintf("Failed to publish in queue %s:", queueName), err)
	}
}

func Consume(
	exchangeName string,
	queueName string,
	f func(d rabbitmq.Delivery) error,
	sleep *time.Duration,
) {
	log := logger.GetLoggerInstance()
	delay := Delay
	if sleep != nil {
		delay = *sleep
	}

	// durable -queue stored after server restart // true
	// auto delete - delete when unused // false
	// exclusive - only one consumer and deleted when queue is empty // false
	// no-wait - don't wait for queue to be declared if true - fast, but you don't know if it has error // false
	consumer, err := rabbitmq.NewConsumer(
		GetRMQInstance(),
		queueName,
		rabbitmq.WithConsumerOptionsQueueDurable,
		rabbitmq.WithConsumerOptionsExchangeDeclare,
		rabbitmq.WithConsumerOptionsRoutingKey(queueName),
		rabbitmq.WithConsumerOptionsExchangeName(exchangeName),
		rabbitmq.WithConsumerOptionsExchangeArgs(rabbitmq.Table{"x-queue-mode": "lazy"}),
	)
	if err != nil {
		panic(fmt.Sprintf("Failed to create consumer: %s.", err.Error()))
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		log.Info("Consumer received signal: " + sig.String())
		consumer.Close()
	}()
	err = consumer.Run(func(d rabbitmq.Delivery) rabbitmq.Action {
		time.Sleep(delay)
		return worker(d, f, exchangeName, queueName, log)
	})
	if err != nil {
		panic(fmt.Sprintf("Failed to consume: %s.", err.Error()))
	}
}

func worker(
	msg rabbitmq.Delivery,
	f func(d rabbitmq.Delivery) error,
	exchangeName string,
	queueName string,
	log *logger.AppLogger,
) rabbitmq.Action {
	var err error
	startTime := time.Now()
	log.Info(fmt.Sprintf(
		"Consumer %s (queue %s) START processing message (%s): %s",
		exchangeName,
		queueName,
		msg.MessageId,
		msg.Body,
	))
	if err = f(msg); err != nil {
		log.Info(fmt.Sprintf("Consumer %s (queue %s) Error Handler: %s", exchangeName, queueName, err.Error()))
		return rabbitmq.NackRequeue
	}
	log.Info(fmt.Sprintf(
		"Consumer %s (queue %s) FINISH processing message (%s): %s --- %s",
		exchangeName,
		queueName,
		msg.MessageId,
		msg.Body,
		helper.RuntimeStatistics(startTime, false),
	))
	return rabbitmq.Ack
}

func Close() error {
	return GetRMQInstance().Close()
}
