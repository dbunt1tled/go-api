package log

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/dbunt1tled/go-api/pkg/e"
	"github.com/dbunt1tled/go-api/pkg/log/logger/handlers/slogpretty"
	"github.com/natefinch/lumberjack/v3"
)

const (
	EnvProd = "prod"
	EnvDev  = "dev"
)

const (
	LevelDebug = slog.LevelDebug
	LevelInfo  = slog.LevelInfo
	LevelWarn  = slog.LevelWarn
	LevelError = slog.LevelError
)

type AppLogger struct {
	*slog.Logger

	level  slog.Level
	cancel func() error
}

var (
	l  atomic.Value //nolint:gochecknoglobals // singleton
	lm sync.Once    //nolint:gochecknoglobals // singleton
)

func Load(name string, env string, debugLevel slog.Level, filePath string) {
	lm.Do(func() {
		log, cfunc := initLogger(name, env, debugLevel, filePath)
		logInstance := &AppLogger{
			Logger: log,
			level:  LevelDebug,
			cancel: cfunc,
		}
		l.Store(logInstance)
	})

}

func Logger() *AppLogger {
	return l.Load().(*AppLogger)
}

func (l *AppLogger) Base() *slog.Logger {
	return l.Logger
}

func Close() error {
	return l.Load().(*AppLogger).cancel()
}

func initLogger(name string, env string, level slog.Level, filePath string) (*slog.Logger, func() error) {
	var (
		log    *slog.Logger
		writer io.Writer
	)

	var cleanup = func() error { return nil }

	if filePath != "" {
		roller, err := lumberjack.NewRoller(filePath,
			10*1024*1024, //nolint:mnd // Temporary
			&lumberjack.Options{
				MaxAge:     24 * time.Hour, //nolint:mnd // Temporary
				MaxBackups: 2,              //nolint:mnd // Temporary
				Compress:   true,
			},
		)
		if err != nil {
			panic(err)
		}

		writer = io.MultiWriter(os.Stdout, roller)

		cleanup = func() error {
			return roller.Close()
		}
	} else {
		writer = os.Stdout
	}

	switch env {
	case EnvDev:
		log = prettyLogHandler(writer, level)
	default:
		log = slog.New(
			slog.NewJSONHandler(writer, &slog.HandlerOptions{Level: level}),
		)
	}

	return log.With(
		slog.String("app", name),
		slog.String("env", env),
	), cleanup
}

func prettyLogHandler(w io.Writer, level slog.Level) *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: level,
		},
	}
	handler := opts.NewPrettyHandler(w)

	return slog.New(handler)
}

// WithContext returns a logger with context values.
func (l *AppLogger) WithContext(ctx context.Context) *slog.Logger {
	// Extract values from context and add them to the logger
	// This is a placeholder implementation that can be expanded
	// to extract specific values from the context
	return l.Logger.With("context", "true")
}

// With returns a logger with the given attributes.
func (l *AppLogger) With(args ...any) *AppLogger {
	newLogger := &AppLogger{
		Logger: l.Logger.With(args...),
		level:  l.level,
	}
	return newLogger
}

// WithGroup returns a logger with the given group.
func (l *AppLogger) WithGroup(name string) *AppLogger {
	newLogger := &AppLogger{
		Logger: l.Logger.WithGroup(name),
		level:  l.level,
	}
	return newLogger
}

// ErrorWithStack logs an error with its stack trace if available.
func (l *AppLogger) ErrorWithStack(msg string, err error) {
	attrs := errorStack(err)
	args := make([]any, 0, len(attrs)*2) //nolint:mnd // dual volume
	for _, attr := range attrs {
		args = append(args, attr.Key, attr.Value.Any())
	}
	l.Logger.Error(msg, args...)
}

// Error logs an error message with structured data.
func (l *AppLogger) Error(msg string, err error, keyvals ...any) {
	args := make([]any, 0, len(keyvals)+2) //nolint:mnd // dual volume
	args = append(args, "error", err.Error())
	args = append(args, keyvals...)
	l.Logger.Error(msg, args...)
}

// ErrorContext logs an error message with context and structured data.
func (l *AppLogger) ErrorContext(ctx context.Context, msg string, err error, keyvals ...any) {
	args := make([]any, 0, len(keyvals)+2) //nolint:mnd // dual volume
	args = append(args, "error", err.Error())
	args = append(args, keyvals...)
	l.Logger.ErrorContext(ctx, msg, args...)
}

// WarnContext logs a warning message with context and structured data.
func (l *AppLogger) WarnContext(ctx context.Context, msg string, keyvals ...any) {
	l.Logger.WarnContext(ctx, msg, keyvals...)
}

// InfoContext logs an info message with context and structured data.
func (l *AppLogger) InfoContext(ctx context.Context, msg string, keyvals ...any) {
	l.Logger.InfoContext(ctx, msg, keyvals...)
}

// DebugContext logs a debug message with context and structured data.
func (l *AppLogger) DebugContext(ctx context.Context, msg string, keyvals ...any) {
	l.Logger.DebugContext(ctx, msg, keyvals...)
}

// Fatalf logs a fatal error message and exits the program.
func (l *AppLogger) Fatalf(msg string, args ...interface{}) {
	l.Logger.Error(fmt.Sprintf(msg, args...))
	os.Exit(1)
}

// Errorf logs an error message.
func (l *AppLogger) Errorf(msg string, args ...interface{}) {
	l.Logger.Error(fmt.Sprintf(msg, args...))
}

// Warnf logs a warning message.
func (l *AppLogger) Warnf(msg string, args ...interface{}) {
	l.Logger.Warn(fmt.Sprintf(msg, args...))
}

// Infof logs an info message.
func (l *AppLogger) Infof(msg string, args ...interface{}) {
	l.Logger.Info(fmt.Sprintf(msg, args...))
}

// Debugf logs a debug message.
func (l *AppLogger) Debugf(msg string, args ...interface{}) {
	l.Logger.Debug(fmt.Sprintf(msg, args...))
}

func errorStack(err error) []slog.Attr {
	stack := e.GetErrTrace(err)
	if stack != nil {
		return []slog.Attr{
			{Key: "stack", Value: slog.StringValue(*stack)},
			{Key: "message", Value: slog.StringValue(err.Error())},
		}
	}
	return []slog.Attr{
		{Key: "message", Value: slog.StringValue(err.Error())},
	}
}
