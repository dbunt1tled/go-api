package postgres

import (
	"database/sql"
	"runtime"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

type Postgres struct {
	db *bun.DB
}

const (
	timeout         = 5 * time.Second
	connMaxLifetime = time.Hour
)

func New(dsn string, debug bool) *Postgres {
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
		panic(err)
	}

	db := bun.NewDB(sqldb, pgdialect.New())
	if debug {
		db = db.WithQueryHook(
			bundebug.NewQueryHook(bundebug.WithVerbose(true)),
		)
	}

	return &Postgres{db: db}
}

func (p *Postgres) DB() *bun.DB {
	return p.db
}

func (p *Postgres) Close() error {
	return p.db.Close()
}
