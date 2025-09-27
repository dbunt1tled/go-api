package main

import (
	"context"
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
	"time"

	"github.com/pkg/errors"
	"github.com/wagslane/go-rabbitmq"
)

const (
	NumWorkers = 6
	Duration   = 1 * time.Second
)

func main() {
	cfg := env.GetConfigInstance()
	locale.GetLocaleBundleInstance()
	logger.InitLogger(cfg.Env, cfg.Debug.Debug, cfg.Logger)
	validate.GetValidateInstance()
	profiler.SetProfiler()
	storage.GetInstance()
	defer storage.Close()
	mailer.GetMailInstance()
	defer mailer.Close()
	jobResolver := rmqmail.NewRMQJobMailResolver()
	rmq.GetRMQInstance()
	defer rmq.Close()
	f := func(d rabbitmq.Delivery) error {
		handler, err := jobResolver.Resolver.Resolve(d.Type)
		if err != nil {
			return fmt.Errorf("mail consumer Error: %s", err.Error())
		}
		if handler == nil {
			return errors.New("mail consumer Error: handler is empty")
		}
		if err = (*handler).Handle(context.Background(), d.Body); err != nil {
			return fmt.Errorf("mail consumer Error processing message: %s", err.Error())
		}
		return nil
	}
	sleep := Duration
	rmq.Consume(rmq.MailExchange, rmq.MailQueue, f, &sleep)
}
