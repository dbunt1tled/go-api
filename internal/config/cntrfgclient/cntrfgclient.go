package cntrfgclient

import (
	"go_echo/internal/config/env"
	"sync"

	"github.com/centrifugal/gocent/v3"
)

var (
	instance *gocent.Client //nolint:gochecknoglobals // singleton
	m        sync.Once      //nolint:gochecknoglobals // singleton
)

func GetInstance() *gocent.Client {
	m.Do(func() {
		cfg := env.GetConfigInstance()
		instance = gocent.New(gocent.Config{
			Addr: cfg.Centrifugo.APIURL + "/api",
			Key:  cfg.Centrifugo.APIKey,
		})
	})
	return instance
}
