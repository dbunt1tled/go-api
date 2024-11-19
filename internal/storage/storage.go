package storage

import (
	"database/sql"
	"go_echo/internal/config"
	"go_echo/internal/config/env"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

var instance *Mysql //nolint:gochecknoglobals // singleton

func GetInstance() *Mysql {
	if instance == nil {
		var err error
		instance, err = Open()
		if err != nil {
			panic(err)
		}
	}
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
	db.SetConnMaxLifetime(time.Minute * 3) //nolint:mnd // TODO need to choose the optimal value
	db.SetMaxOpenConns(100)                //nolint:mnd // TODO need to choose the optimal value
	db.SetMaxIdleConns(25)                 //nolint:mnd // TODO need to choose the optimal value
	err = db.Ping()
	if err != nil {
		return nil, errors.Wrap(err, "db ping error")
	}
	return &Mysql{db: db}, nil
}

func Close() {
	db := GetInstance()
	log := config.GetLoggerInstance()
	err := db.db.Close()
	if err != nil {
		log.Error(errors.Wrap(err, "db close error").Error())
	}
}
