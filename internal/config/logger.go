package config

import (
	"go_echo/internal/lib/logger/handlers/pretty"
	"log/slog"
	"os"
	"sync"
)

const (
	EnvProd = "prod"
	EnvDev  = "dev"
)

var (
	instance *slog.Logger //nolint:gochecknoglobals // singleton
	once     sync.Once    //nolint:gochecknoglobals // singleton
)

func InitLogger(env string, debug bool) *slog.Logger {
	once.Do(func() {
		instance = setupLogger(env, debug)
	})
	return GetLoggerInstance()
}
func GetLoggerInstance() *slog.Logger {
	if instance == nil {
		panic("Singleton is not initialized. Call InitSingleton first.")
	}
	return instance
}

func setupLogger(env string, debug bool) *slog.Logger {
	var log *slog.Logger

	switch env {
	case EnvProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	case EnvDev:
		log = PrettyLogHandler(env, debug)
	default:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	}

	return log
}

func PrettyLogHandler(env string, debug bool) *slog.Logger {
	opts := pretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	logger := slog.New(opts.NewPrettyHandler(os.Stdout))
	if debug == true {
		// logger = logger.With(slog.String("env", env))
	}
	return logger
}
