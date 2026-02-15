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
	db := postgres.New(
		config.Get().DB.Main.DSN,
		config.Get().Debug,
	)
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

	if err := migrator.Lock(ctx); err != nil {
		panic(err)
	}
	defer func(migrator *migrate.Migrator, ctx context.Context) {
		err := migrator.Unlock(ctx)
		if err != nil {
			panic(err)
		}
	}(migrator, ctx)

	group, err = migrator.Migrate(ctx)
	if err != nil {
		panic(err)
	}

	if group == nil {
		log.Logger().Error("No migration apply.", nil)
	}

	log.Logger().Infof("Success migration %v", group)
}
