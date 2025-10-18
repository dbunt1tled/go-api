package main

import (
	"context"
	"github.com/dbunt1tled/go-api/internal/cache"
	"github.com/dbunt1tled/go-api/internal/config/env"
	"github.com/dbunt1tled/go-api/internal/config/locale"
	"github.com/dbunt1tled/go-api/internal/config/logger"
	"github.com/dbunt1tled/go-api/internal/config/validate"
	"github.com/dbunt1tled/go-api/internal/lib/handler"
	"github.com/dbunt1tled/go-api/internal/lib/profiler"
	"github.com/dbunt1tled/go-api/internal/rmq"
	"github.com/dbunt1tled/go-api/internal/router"
	"github.com/dbunt1tled/go-api/internal/storage"
	"github.com/dbunt1tled/go-api/internal/util/helper"
	"github.com/dbunt1tled/go-api/internal/util/sanitizer"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	_ "github.com/redis/go-redis/v9"
)

func main() {
	cfg := env.GetConfigInstance()
	locale.GetLocaleBundleInstance()
	logger.InitLogger(cfg.Env, cfg.Debug.Debug, cfg.Logger)
	log := logger.GetLoggerInstance()
	validate.GetValidateInstance()
	sanitizer.GetSanitizerInstance()
	profiler.SetProfiler()
	storage.GetInstance()
	cache.GetRedisCache()
	rmq.GetRMQInstance()
	httpServer := echo.New()
	httpServer.HideBanner = true
	httpServer.Debug = cfg.Debug.Debug
	httpServer.HTTPErrorHandler = handler.APIErrorHandler
	router.SetupRoutes(httpServer)
	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		os.Interrupt,
	)
	defer stop()
	go func() {
		if cfg.HTTPServer.TLS.IsSet() {
			log.Debug("(っ◕‿◕)っ Start Server TLS listening on address: " + cfg.HTTPServer.Address)
			err := httpServer.StartTLS(
				cfg.HTTPServer.Address,
				cfg.HTTPServer.TLS.GetCertData(),
				cfg.HTTPServer.TLS.GetKeyData(),
			)
			if !errors.Is(err, http.ErrServerClosed) {
				log.Error("¯\\_(͡° ͜ʖ ͡°)_/¯Shutting down the server", err)
			}
		} else {
			log.Debug("(/◔◡◔)/ Start Server listening on address: " + cfg.HTTPServer.Address)
			if err := httpServer.Start(cfg.HTTPServer.Address); err != nil && !errors.Is(err, http.ErrServerClosed) {
				log.Error("¯\\_(͡° ͜ʖ ͡°)_/¯Shutting down the server", err)
			}
		}
	}()
	<-ctx.Done()
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second) //nolint:mnd // 10 seconds timeout
	defer cancel()
	log.Warn("Quit: shutting down ...")
	defer log.Warn("｡◕‿‿◕｡ Quit: shutdown completed")
	helper.GracefulShutdown(
		log,
		func() error {
			log.Info("㋡ Quit: closing database connection")
			return storage.Close()
		},
		func() error {
			log.Info("㋡ Quit: closing cache connection")
			return cache.GetRedisCache().Close()
		},
		func() error {
			log.Info("㋡ Quit: closing rmq connection")
			return rmq.Close()
		},
		func() error {
			log.Info("㋡ Quit: shutting down the server")
			return httpServer.Shutdown(ctx)
		},
		func() error {
			log.Warn("｡◕‿‿◕｡ Quit: shutdown completed")
			os.Exit(0)
			return nil
		},
	)
}
