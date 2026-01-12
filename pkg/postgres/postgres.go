package postgres

import (
	"database/sql"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

type Postgres struct {
	db *bun.DB
}

const (
	timeout         = 5 * time.Second
	maxOpenConns    = 10
	maxIdleConns    = 10
	connMaxLifetime = time.Hour
)

func New(dsn string) *Postgres {
	sqldb := sql.OpenDB(pgdriver.NewConnector(
		pgdriver.WithDSN(dsn),
		pgdriver.WithTimeout(timeout),
		pgdriver.WithDialTimeout(timeout),
		pgdriver.WithReadTimeout(timeout),
		pgdriver.WithWriteTimeout(timeout),
	))
	sqldb.SetMaxOpenConns(maxOpenConns)
	sqldb.SetMaxIdleConns(maxIdleConns)
	sqldb.SetConnMaxLifetime(connMaxLifetime)
	return &Postgres{db: bun.NewDB(sqldb, pgdialect.New())}
}

func (p *Postgres) DB() *bun.DB {
	return p.db
}

func (p *Postgres) Close() error {
	return p.db.Close()
}
