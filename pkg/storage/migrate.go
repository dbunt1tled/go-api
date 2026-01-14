package storage

import (
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/migrate"
)

var migrations = migrate.NewMigrations()

func NewMigrator(db *bun.DB) *migrate.Migrator {
	return migrate.NewMigrator(db, migrations)
}
