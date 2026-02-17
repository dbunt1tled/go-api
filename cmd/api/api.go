package main

import (
	"context"

	"github.com/dbunt1tled/go-api/internal/app"
	"github.com/dbunt1tled/go-api/internal/config"
	"github.com/dbunt1tled/go-api/internal/lib/cent"
	"github.com/dbunt1tled/go-api/pkg/f"
	"github.com/dbunt1tled/go-api/pkg/log"
	"github.com/dbunt1tled/go-api/pkg/mailer"
	"github.com/dbunt1tled/go-api/pkg/postgres"
	"github.com/dbunt1tled/go-api/pkg/rmq"
)

func main() {
	config.Load()
	log.Load(config.Get().Name, config.Get().Env, config.Get().Log.Level, config.Get().Log.File)
	db := postgres.New(
		config.Get().DB.Main.DSN,
		config.Get().Debug,
	)

	mail := mailer.NewMailer(
		config.Get().Mailer.Host,
		config.Get().Mailer.Port,
		config.Get().Mailer.Username,
		config.Get().Mailer.Password,
		config.Get().Mailer.Address,
	)
	rmp := f.Must(rmq.NewProducer(config.Get().AMQP.URL, 0))
	cs := cent.NewCentrifugoService(config.Get().Centrifugo.APIKey, config.Get().Centrifugo.APIUrl)
	cfg := config.NewServiceConfig(db, mail, rmp, cs)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	application := app.NewApp(cfg)
	app.Router(application)
	application.Run(ctx)
}
