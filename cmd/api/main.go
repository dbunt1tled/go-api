package main

import (
	"go_echo/internal/config/env"
	"go_echo/internal/config/locale"
	"go_echo/internal/config/logger"
	"go_echo/internal/lib/graceful"
	"go_echo/internal/lib/handler"
	"go_echo/internal/lib/profiler"
	"go_echo/internal/router"
	"go_echo/internal/storage"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

func main() {
	cfg := env.GetConfigInstance()
	locale.GetLocaleBundleInstance()
	logger.InitLogger(cfg.Env, cfg.Debug)
	log := logger.GetLoggerInstance()
	profiler.SetProfiler()
	storage.GetInstance()
	defer storage.Close()
	httpServer := echo.New()
	httpServer.HideBanner = true
	httpServer.Debug = cfg.Debug
	httpServer.HTTPErrorHandler = handler.APIErrorHandler
	router.SetupRoutes(httpServer)

	done := graceful.ShutdownGraceful(log, httpServer)
	//u, err := service.UserRepository{}.List([]builder.FilterCondition{
	//	{Field: "id", Type: builder.In, Value: []interface{}{1, 2}},
	//}, []builder.SortOrder{
	//	{Field: "id", Order: builder.Desc},
	//})
	//if err != nil {
	//	log.Error(err.Error())
	//} else {
	//	log.Debug(fmt.Sprintf("%+v\n", u))
	//}

	go func() {
		log.Debug("Start listening on address: " + cfg.HTTPServer.Address)
		if err := httpServer.Start(cfg.HTTPServer.Address); err != nil && !errors.Is(err, http.ErrServerClosed) { //nolint:lll,govet
			log.Error("shutting down the server")
		}
	}()
	<-done
}
