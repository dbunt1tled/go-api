package main

import (
	"context"
	"fmt"
	"time"

	"github.com/dbunt1tled/go-api/app/jobs/rmqmail"
	"github.com/dbunt1tled/go-api/internal/config/env"
	"github.com/dbunt1tled/go-api/internal/config/locale"
	"github.com/dbunt1tled/go-api/internal/config/logger"
	"github.com/dbunt1tled/go-api/internal/config/mailer"
	"github.com/dbunt1tled/go-api/internal/config/validate"
	"github.com/dbunt1tled/go-api/internal/lib/profiler"
	"github.com/dbunt1tled/go-api/internal/rmq"
	"github.com/dbunt1tled/go-api/internal/storage"

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
