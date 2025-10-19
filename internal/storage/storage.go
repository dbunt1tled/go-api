package storage

import (
	"database/sql"
	"sync"
	"time"

	"github.com/dbunt1tled/go-api/internal/config/env"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

const (
	maxOpenConns    = 60
	connMaxLifetime = 2
	maxIdleConns    = 30
	connMaxIdleTime = 20
)

var (
	instance *Mysql    //nolint:gochecknoglobals // singleton
	m        sync.Once //nolint:gochecknoglobals // singleton
)

func GetInstance() *Mysql {
	m.Do(func() {
		var err error
		instance, err = Open()
		if err != nil {
			panic(err)
		}
	})
	return instance
}

type Mysql struct {
	db *sql.DB
}

func (m *Mysql) GetDB() *sql.DB {
	return m.db
}

func Open() (*Mysql, error) {
	cfg := env.GetConfigInstance()
	db, err := sql.Open("mysql", cfg.DatabaseDSN)
	if err != nil {
		return nil, errors.Wrap(err, "db open error")
	}
	db.SetConnMaxLifetime(time.Minute * connMaxLifetime)
	db.SetConnMaxIdleTime(time.Minute * connMaxIdleTime)
	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	err = db.Ping()
	if err != nil {
		return nil, errors.Wrap(err, "db ping error")
	}
	return &Mysql{db: db}, nil
}

func Close() error {
	return GetInstance().db.Close()
}
