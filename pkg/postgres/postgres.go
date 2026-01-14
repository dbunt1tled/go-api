package postgres

import (
	"database/sql"
	"log"
	"runtime"
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
	maxOpenConns := 4 * runtime.GOMAXPROCS(0) //nolint:mnd // prod
	sqldb.SetMaxOpenConns(maxOpenConns)
	sqldb.SetMaxIdleConns(maxOpenConns)
	sqldb.SetConnMaxLifetime(timeout)         // Connection lifetime
	sqldb.SetConnMaxIdleTime(connMaxLifetime) // Idle connection timeout

	if err := sqldb.Ping(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	return &Postgres{db: bun.NewDB(sqldb, pgdialect.New())}
}

func (p *Postgres) DB() *bun.DB {
	return p.db
}

func (p *Postgres) Close() error {
	return p.db.Close()
}
