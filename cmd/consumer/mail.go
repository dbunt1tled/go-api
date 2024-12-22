package main

import (
	"fmt"
	"go_echo/app/jobs/rmqmail"
	"go_echo/internal/config/env"
	"go_echo/internal/config/locale"
	"go_echo/internal/config/logger"
	"go_echo/internal/config/mailer"
	"go_echo/internal/config/validate"
	"go_echo/internal/lib/profiler"
	"go_echo/internal/rmq"
	"go_echo/internal/storage"

	"github.com/pkg/errors"
	"github.com/rabbitmq/amqp091-go"
)

const (
	NumWorkers = 6
)

func main() {
	cfg := env.GetConfigInstance()
	locale.GetLocaleBundleInstance()
	logger.InitLogger(cfg.Env, cfg.Debug)
	validate.GetValidateInstance()
	profiler.SetProfiler()
	storage.GetInstance()
	defer storage.Close()
	mailer.GetMailInstance()
	defer mailer.Close()

	jobResolver := rmqmail.NewRMQJobMailResolver()

	rc := rmq.GetRMQInstance(rmq.MailExchange)
	f := func(d amqp091.Delivery) error {
		handler, err := jobResolver.Resolver.Resolve(d.Type)
		if err != nil {
			return fmt.Errorf("mail consumer Error: %s", err.Error())
		}
		if handler == nil {
			return errors.New("mail consumer Error: handler is empty")
		}
		if err = (*handler).Handle(d.Body); err != nil {
			return fmt.Errorf("mail consumer Error processing message: %s", err.Error())
		}
		return nil
	}
	rc.Consume(rmq.MailExchange, rmq.MailQueue, f, NumWorkers)
	rc.Close()
}
