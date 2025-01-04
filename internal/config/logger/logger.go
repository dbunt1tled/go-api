package logger

import (
	"fmt"
	"go_echo/internal/lib/logger/handlers/pretty"
	"log/slog"
	"os"
	"sync"
)

const (
	EnvProd = "prod"
	EnvDev  = "dev"
)

type AppLogger struct {
	*slog.Logger
	AdditionalLogs bool
}

var (
	logInstance *AppLogger //nolint:gochecknoglobals // singleton
	m           sync.Once  //nolint:gochecknoglobals // singleton
)

func InitLogger(env string, debug bool) *AppLogger {
	m.Do(func() {
		logInstance = setupLogger(env, debug)
	})
	return GetLoggerInstance()
}
func GetLoggerInstance() *AppLogger {
	if logInstance == nil {
		panic("Singleton is not initialized. Call InitSingleton first.")
	}
	return logInstance
}

func setupLogger(env string, debug bool) *AppLogger {
	log := &AppLogger{
		AdditionalLogs: debug,
	}
	switch env {
	case EnvProd:
		log.Logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	case EnvDev:
		log.Logger = PrettyLogHandler(env, debug)
	default:
		log.Logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
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

func (l *AppLogger) Fatalf(msg string, args ...interface{}) {
	if !l.AdditionalLogs {
		return
	}
	l.Logger.Error(fmt.Sprintf(msg, args...))
}
func (l *AppLogger) Errorf(msg string, args ...interface{}) {
	if !l.AdditionalLogs {
		return
	}
	l.Logger.Error(fmt.Sprintf(msg, args...))
}
func (l *AppLogger) Warnf(msg string, args ...interface{}) {
	if !l.AdditionalLogs {
		return
	}
	l.Logger.Warn(fmt.Sprintf(msg, args...))
}
func (l *AppLogger) Infof(msg string, args ...interface{}) {
	if !l.AdditionalLogs {
		return
	}
	l.Logger.Info(fmt.Sprintf(msg, args...))
}
func (l *AppLogger) Debugf(msg string, args ...interface{}) {
	if !l.AdditionalLogs {
		return
	}
	l.Logger.Debug(msg, args...)
}
