package main

import (
	"context"
	"go_echo/internal/config/env"
	"go_echo/internal/config/locale"
	"go_echo/internal/config/logger"
	"go_echo/internal/config/mailer"
	"go_echo/internal/config/validate"
	"go_echo/internal/lib/handler"
	"go_echo/internal/lib/profiler"
	"go_echo/internal/rmq"
	"go_echo/internal/router"
	"go_echo/internal/storage"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

func main() {
	cfg := env.GetConfigInstance()
	locale.GetLocaleBundleInstance()
	logger.InitLogger(cfg.Env, cfg.Debug)
	log := logger.GetLoggerInstance()
	validate.GetValidateInstance()
	profiler.SetProfiler()
	storage.GetInstance()
	defer storage.Close()
	mailer.GetMailInstance()
	defer mailer.Close()
	// (import|study_anal|sub_system)
	var rc rmq.RabbitClient
	rc.Publish("bb", "aaaa", "{\"Page\":1,\"Fruits\":[\"apple\",\"peach\",\"pear\"]}")
	httpServer := echo.New()
	httpServer.HideBanner = true
	httpServer.Debug = cfg.Debug
	httpServer.HTTPErrorHandler = handler.APIErrorHandler
	router.SetupRoutes(httpServer)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	go func() {
		log.Debug("Start listening on address: " + cfg.HTTPServer.Address)
		if err := httpServer.Start(cfg.HTTPServer.Address); err != nil && !errors.Is(err, http.ErrServerClosed) { //nolint:lll,govet
			log.Error("shutting down the server")
		}
	}()
	<-ctx.Done()
	log.Warn("quit: shutting down ...")
	defer log.Warn("quit: shutdown completed")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		httpServer.Logger.Fatal(err)
	}
}
