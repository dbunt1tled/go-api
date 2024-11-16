package storage

import (
	"database/sql"
	"go_echo/internal/config"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

type Mysql struct {
	db *sql.DB
}

func Open(storagePath string) (*Mysql, error) {
	db, err := sql.Open("mysql", storagePath)
	if err != nil {
		return nil, errors.Wrap(err, "db open error")
	}
	err = db.Ping()
	if err != nil {
		return nil, errors.Wrap(err, "db ping error")
	}

	return &Mysql{db: db}, nil
}

func Close(db *Mysql) {
	log := config.GetLoggerInstance()
	err := db.db.Close()
	if err != nil {
		log.Error(errors.Wrap(err, "db close error").Error())
	}
}
