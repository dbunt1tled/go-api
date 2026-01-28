package main

import (
	"context"

	"github.com/dbunt1tled/go-api/internal/config"
	"github.com/dbunt1tled/go-api/internal/jobs/rmqmail"
	"github.com/dbunt1tled/go-api/internal/jobs/rmqmail/handlers"
	"github.com/dbunt1tled/go-api/internal/modules/user"
	"github.com/dbunt1tled/go-api/pkg/log"
	"github.com/dbunt1tled/go-api/pkg/mailer"
	"github.com/dbunt1tled/go-api/pkg/postgres"
	"github.com/dbunt1tled/go-api/pkg/rmq"
	"github.com/uptrace/bun"
)

func main() {
	config.Load()
	log.Load(config.Get().Name, config.Get().Env, config.Get().Log.Level, config.Get().Log.File)
	db := postgres.New(config.Get().DB.Main.DSN)
	defer func(db *bun.DB) {
		err := db.Close()
		if err != nil {
			log.Logger().Errorf("Failed to close db: %v", err)
		}
	}(db.DB())
	userService := user.NewUserService(db.DB())

	mailService := mailer.NewMailer(
		config.Get().Mailer.Host,
		config.Get().Mailer.Port,
		config.Get().Mailer.Username,
		config.Get().Mailer.Password,
		config.Get().Mailer.Address,
	)

	defer func(mailService *mailer.Mailer) {
		err := mailService.Close()
		if err != nil {
			log.Logger().Errorf("Failed to close mailer: %v", err)
		}
	}(mailService)

	consumer, err := rmq.NewConsumer(config.Get().AMQP.URL, rmq.ConsumerConfig{
		QueueName:  rmqmail.MailQueue,
		MaxRetries: rmqmail.MaxRetries,
		RetryDelay: rmqmail.RetryDelay,
		QueueType:  "classic",
		Durable:    true,

		AutoDeclare:  true,
		ExchangeName: rmqmail.MailExchange,
		ExchangeType: "direct",
		RoutingKeys:  []string{rmqmail.MailRoutingKey},
	})
	if err != nil {
		log.Logger().Errorf("Failed to create consumer: %v", err)
		return
	}
	defer func(consumer *rmq.Consumer) {
		err = consumer.Close()
		if err != nil {
			log.Logger().Errorf("Failed to close consumer: %v", err)
		}
	}(consumer)
	jobResolver := rmqmail.NewRMQJobMailResolver(
		handlers.NewUserConfirmationEmailJob(userService, mailService),
	)
	err = consumer.Start(context.Background(), jobResolver, 1, nil)
	if err != nil {
		log.Logger().Errorf("Failed to start consumer: %v", err)
	}
}
