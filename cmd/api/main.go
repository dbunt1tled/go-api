package main

import (
	"go_echo/internal/config"
	"go_echo/internal/config/env"
	"go_echo/internal/lib/graceful"
	"go_echo/internal/lib/handler"
	"go_echo/internal/router"
	"go_echo/internal/storage"
	"net/http"
	"runtime"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/pkg/errors"
)

var (
	fastEventDuration = 1 * time.Millisecond
	slowEventDuration = 10 * fastEventDuration
)

func main() {
	var locale *i18n.Localizer
	cfg := env.GetConfigInstance()
	config.InitLogger(cfg.Env, cfg.Debug)
	log := config.GetLoggerInstance()
	bundle := config.SetupLocale()
	if cfg.Profiling {
		runtime.SetBlockProfileRate(int(slowEventDuration.Nanoseconds()))
		go func() {
			err := http.ListenAndServe("localhost:6060", nil)
			if err != nil {
				log.Error(err.Error())
			}
		}()
	}

	db, err := storage.Open(cfg.DatabaseDSN)
	if err != nil {
		log.Error(err.Error())
		return
	}
	defer storage.Close(db)
	httpServer := echo.New()
	httpServer.HideBanner = true
	httpServer.Debug = cfg.Debug
	httpServer.HTTPErrorHandler = handler.CustomHTTPErrorHandler

	done := graceful.ShutDown(log, httpServer)
	router.SetupRoutes(httpServer, locale, bundle)
	go func() {
		log.Debug("Start listening on address: " + cfg.HTTPServer.Address)
		if err := httpServer.Start(cfg.HTTPServer.Address); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("shutting down the server")
		}
	}()
	<-done
}
