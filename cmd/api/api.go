package main

import (
	"context"
	"go_echo/internal/cache"
	"go_echo/internal/config/env"
	"go_echo/internal/config/locale"
	"go_echo/internal/config/logger"
	"go_echo/internal/config/validate"
	"go_echo/internal/lib/handler"
	"go_echo/internal/lib/profiler"
	"go_echo/internal/router"
	"go_echo/internal/storage"
	"go_echo/internal/util/sanitizer"
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
	logger.InitLogger(cfg.Env, cfg.Debug.Debug)
	log := logger.GetLoggerInstance()
	validate.GetValidateInstance()
	sanitizer.GetSanitizerInstance()
	profiler.SetProfiler()
	storage.GetInstance()
	defer storage.Close()
	cache.GetRedisCache()
	defer cache.GetRedisCache().Close()
	httpServer := echo.New()
	httpServer.HideBanner = true
	httpServer.Debug = cfg.Debug.Debug
	httpServer.HTTPErrorHandler = handler.APIErrorHandler
	router.SetupRoutes(httpServer)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	go func() {
		log.Debug("Start listening on address: " + cfg.HTTPServer.Address)
		if err := httpServer.Start(cfg.HTTPServer.Address); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("shutting down the server" + err.Error())
		}
	}()
	<-ctx.Done()
	log.Warn("quit: shutting down ...")
	defer log.Warn("quit: shutdown completed")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) //nolint:mnd // 10 seconds timeout
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		log.ErrorContext(ctx, "error shutting down the server", err)
	}
}
