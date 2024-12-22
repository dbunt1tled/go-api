package cntrfgclient

import (
	"go_echo/internal/config/env"

	"github.com/centrifugal/gocent/v3"
)

var instance *gocent.Client //nolint:gochecknoglobals // singleton

func GetInstance() *gocent.Client {
	if instance == nil {
		cfg := env.GetConfigInstance()
		instance = gocent.New(gocent.Config{
			Addr: cfg.Centrifugo.ApiUrl + "/api",
			Key:  cfg.Centrifugo.APIKey,
		})
	}
	return instance
}
