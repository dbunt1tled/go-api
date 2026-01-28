package config

import (
	"github.com/dbunt1tled/go-api/pkg/mailer"
	"github.com/dbunt1tled/go-api/pkg/postgres"
	"github.com/dbunt1tled/go-api/pkg/rmq"
)

type ServiceConfig struct {
	DB         *postgres.Postgres
	Mailer     *mailer.Mailer
	RMProducer *rmq.Producer
}

func NewServiceConfig(
	db *postgres.Postgres,
	mail *mailer.Mailer,
	rmp *rmq.Producer,
) *ServiceConfig {
	return &ServiceConfig{
		DB:         db,
		Mailer:     mail,
		RMProducer: rmp,
	}
}
