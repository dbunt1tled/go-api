package storage

import (
	"github.com/dbunt1tled/go-api/pkg/storage/migrations"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/migrate"
)

func NewMigrator(db *bun.DB) *migrate.Migrator {
	m := migrate.NewMigrations()
	migrations.CreateUserTable(m)
	migrations.CreateUserNotificationTable(m)

	return migrate.NewMigrator(db, m)
}
