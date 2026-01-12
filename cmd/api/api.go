package main

import (
	"context"

	"github.com/dbunt1tled/go-api/internal/app"
	"github.com/dbunt1tled/go-api/internal/config"
	"github.com/dbunt1tled/go-api/pkg/log"
	"github.com/dbunt1tled/go-api/pkg/mailer"
	"github.com/dbunt1tled/go-api/pkg/postgres"
)

func main() {
	config.Load()
	log.Load(config.Get().Name, config.Get().Env, config.Get().Log.Level, config.Get().Log.File)
	db := postgres.New(config.Get().DB.Main.DSN)

	mail := mailer.NewMailer(
		config.Get().Mailer.Host,
		config.Get().Mailer.Port,
		config.Get().Mailer.Username,
		config.Get().Mailer.Password,
		config.Get().Mailer.Address,
	)

	cfg := config.NewServiceConfig(db, mail)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	application := app.NewApp(cfg)
	application.Run(ctx)
}
