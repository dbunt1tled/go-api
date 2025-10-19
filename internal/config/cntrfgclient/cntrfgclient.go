package cntrfgclient

import (
	"sync"

	"github.com/dbunt1tled/go-api/internal/config/env"

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
