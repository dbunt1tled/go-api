package main

import (
	"context"

	"github.com/dbunt1tled/go-api/internal/config"
	"github.com/dbunt1tled/go-api/pkg/log"
	"github.com/dbunt1tled/go-api/pkg/postgres"
	"github.com/dbunt1tled/go-api/pkg/storage"
	"github.com/uptrace/bun/migrate"
)

func main() {
	var (
		err   error
		group *migrate.MigrationGroup
	)
	config.Load()
	log.Load(config.Get().Name, config.Get().Env, config.Get().Log.Level, config.Get().Log.File)
	db := postgres.New(config.Get().DB.Main.DSN)
	defer func(db *postgres.Postgres) {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}(db)
	ctx := context.Background()
	migrator := storage.NewMigrator(db.DB())
	err = migrator.Init(ctx)
	if err != nil {
		panic(err)
	}
	group, err = migrator.Migrate(ctx)
	if err != nil {
		panic(err)
	}
	log.Logger().Infof("Success migration %v", group)
}
