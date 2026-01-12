package config

import (
	"github.com/dbunt1tled/go-api/pkg/mailer"
	"github.com/dbunt1tled/go-api/pkg/postgres"
)

type ServiceConfig struct {
	DB     *postgres.Postgres
	Mailer *mailer.Mailer
}

func NewServiceConfig(
	db *postgres.Postgres,
	mail *mailer.Mailer,
) *ServiceConfig {
	return &ServiceConfig{
		DB:     db,
		Mailer: mail,
	}
}
