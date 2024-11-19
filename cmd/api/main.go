package main

import (
	"fmt"
	"go_echo/app/user/service"
	"go_echo/internal/config"
	"go_echo/internal/config/env"
	"go_echo/internal/lib/graceful"
	"go_echo/internal/lib/handler"
	"go_echo/internal/lib/profiler"
	"go_echo/internal/router"
	"go_echo/internal/storage"
	"go_echo/internal/util/builder"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/pkg/errors"
)

func main() {
	var locale *i18n.Localizer
	cfg := env.GetConfigInstance()
	config.InitLogger(cfg.Env, cfg.Debug)
	log := config.GetLoggerInstance()
	bundle := config.SetupLocale()
	profiler.SetProfiler()
	storage.GetInstance()
	defer storage.Close()
	httpServer := echo.New()
	httpServer.HideBanner = true
	httpServer.Debug = cfg.Debug
	httpServer.HTTPErrorHandler = handler.APIErrorHandler
	router.SetupRoutes(httpServer, locale, bundle)

	done := graceful.ShutdownGraceful(log, httpServer)
	u, err := service.UserRepository{}.List([]builder.FilterCondition{
		{Field: "id", Type: builder.In, Value: []interface{}{1, 2}},
	}, []builder.SortOrder{
		{Field: "id", Order: builder.Desc},
	})
	if err != nil {
		log.Error(err.Error())
	} else {
		log.Debug(fmt.Sprintf("%+v\n", u))
	}

	go func() {
		log.Debug("Start listening on address: " + cfg.HTTPServer.Address)
		if err := httpServer.Start(cfg.HTTPServer.Address); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("shutting down the server")
		}
	}()
	<-done
}
