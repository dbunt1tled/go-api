package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"sync"

	"github.com/dbunt1tled/go-api/internal/config/env"
	"github.com/dbunt1tled/go-api/internal/lib/logger"
	"github.com/dbunt1tled/go-api/internal/lib/logger/handlers/pretty"
)

// Environment constants
const (
	EnvProd = "prod"
	EnvDev  = "dev"
)

// Log levels
const (
	LevelDebug = slog.LevelDebug
	LevelInfo  = slog.LevelInfo
	LevelWarn  = slog.LevelWarn
	LevelError = slog.LevelError
)

// AppLogger wraps slog.Logger with additional functionality
type AppLogger struct {
	*slog.Logger
	level      slog.Level
	fileWriter *os.File
}

var (
	logInstance *AppLogger //nolint:gochecknoglobals // singleton
	m           sync.Once  //nolint:gochecknoglobals // singleton
)

// InitLogger initializes the logger singleton with the specified environment and log level
// If level is not specified (0), it will default to Info for prod and Debug for dev
func InitLogger(env string, levelOrDebug interface{}, loggerOptions env.Logger) *AppLogger {
	m.Do(func() {
		var level slog.Level

		switch v := levelOrDebug.(type) {
		case slog.Level:
			level = v
		case bool:
			// For backward compatibility with the old debug parameter
			if v {
				level = LevelDebug
			} else {
				if env == EnvProd {
					level = LevelInfo
				} else {
					level = LevelDebug
				}
			}
		default:
			// Default to environment-based level
			level = 0
		}

		logInstance = setupLogger(env, level)

		// Enable file logging if FilePath is not empty
		if loggerOptions.FilePath != "" {
			err := logInstance.EnableFileLogging(loggerOptions.FilePath)
			if err != nil {
				logInstance.Errorf("Failed to enable file logging: %v", err)
			}
		}
	})
	return GetLoggerInstance()
}

// GetLoggerInstance returns the logger singleton instance
// Panics if the logger has not been initialized
func GetLoggerInstance() *AppLogger {
	if logInstance == nil {
		panic("Logger is not initialized. Call InitLogger first.")
	}
	return logInstance
}

// setupLogger creates a new logger with the specified environment and log level
func setupLogger(env string, level slog.Level) *AppLogger {
	// Set default log level based on environment if not specified
	if level == 0 {
		switch env {
		case EnvProd:
			level = LevelInfo
		default:
			level = LevelDebug
		}
	}

	log := &AppLogger{
		level: level,
	}

	switch env {
	case EnvProd:
		log.Logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level}))
	case EnvDev:
		log.Logger = PrettyLogHandler(env, level)
	default:
		log.Logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level}))
	}

	return log
}

// EnableFileLogging enables logging to a file in addition to the console
func (l *AppLogger) EnableFileLogging(filePath string) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	// Close existing file writer if any
	if l.fileWriter != nil {
		if err := l.fileWriter.Close(); err != nil {
			return fmt.Errorf("failed to close existing log file: %w", err)
		}
		l.fileWriter = nil
	}

	// Open log file
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}

	// Create a multi-writer that writes to both stdout and the file
	multiWriter := io.MultiWriter(os.Stdout, file)

	// Create a new handler based on the current logger's format
	var handler slog.Handler
	handler = slog.NewTextHandler(multiWriter, &slog.HandlerOptions{Level: l.level})

	// Update the logger
	l.Logger = slog.New(handler)
	l.fileWriter = file

	l.Infof("File logging enabled: %s", filePath)
	return nil
}

// Close closes the logger's file writer if any
func (l *AppLogger) Close() error {
	if l.fileWriter != nil {
		err := l.fileWriter.Close()
		l.fileWriter = nil
		return err
	}
	return nil
}

// PrettyLogHandler creates a pretty logger for development environments
func PrettyLogHandler(env string, level slog.Level) *slog.Logger {
	opts := pretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: level,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)
	return slog.New(handler.WithAttrs([]slog.Attr{
		slog.String("env", env),
	}))
}

// WithContext returns a logger with context values
func (l *AppLogger) WithContext(ctx context.Context) *slog.Logger {
	// Extract values from context and add them to the logger
	// This is a placeholder implementation that can be expanded
	// to extract specific values from the context
	return l.Logger.With("context", "true")
}

// With returns a logger with the given attributes
func (l *AppLogger) With(args ...any) *AppLogger {
	newLogger := &AppLogger{
		Logger: l.Logger.With(args...),
		level:  l.level,
	}
	return newLogger
}

// WithGroup returns a logger with the given group
func (l *AppLogger) WithGroup(name string) *AppLogger {
	newLogger := &AppLogger{
		Logger: l.Logger.WithGroup(name),
		level:  l.level,
	}
	return newLogger
}

// ErrorWithStack logs an error with its stack trace if available
func (l *AppLogger) ErrorWithStack(msg string, err error) {
	attrs := logger.Error(err)
	args := make([]any, 0, len(attrs)*2)
	for _, attr := range attrs {
		args = append(args, attr.Key, attr.Value.Any())
	}
	l.Logger.Error(msg, args...)
}

// Error logs an error message with structured data
func (l *AppLogger) Error(msg string, err error, keyvals ...any) {
	args := make([]any, 0, len(keyvals)+2)
	args = append(args, "error", err.Error())
	args = append(args, keyvals...)
	l.Logger.Error(msg, args...)
}

// ErrorContext logs an error message with context and structured data
func (l *AppLogger) ErrorContext(ctx context.Context, msg string, err error, keyvals ...any) {
	args := make([]any, 0, len(keyvals)+2)
	args = append(args, "error", err.Error())
	args = append(args, keyvals...)
	l.Logger.ErrorContext(ctx, msg, args...)
}

// WarnContext logs a warning message with context and structured data
func (l *AppLogger) WarnContext(ctx context.Context, msg string, keyvals ...any) {
	l.Logger.WarnContext(ctx, msg, keyvals...)
}

// InfoContext logs an info message with context and structured data
func (l *AppLogger) InfoContext(ctx context.Context, msg string, keyvals ...any) {
	l.Logger.InfoContext(ctx, msg, keyvals...)
}

// DebugContext logs a debug message with context and structured data
func (l *AppLogger) DebugContext(ctx context.Context, msg string, keyvals ...any) {
	l.Logger.DebugContext(ctx, msg, keyvals...)
}

// Fatalf logs a fatal error message and exits the program
func (l *AppLogger) Fatalf(msg string, args ...interface{}) {
	l.Logger.Error(fmt.Sprintf(msg, args...))
	os.Exit(1)
}

// Errorf logs an error message
func (l *AppLogger) Errorf(msg string, args ...interface{}) {
	l.Logger.Error(fmt.Sprintf(msg, args...))
}

// Warnf logs a warning message
func (l *AppLogger) Warnf(msg string, args ...interface{}) {
	l.Logger.Warn(fmt.Sprintf(msg, args...))
}

// Infof logs an info message
func (l *AppLogger) Infof(msg string, args ...interface{}) {
	l.Logger.Info(fmt.Sprintf(msg, args...))
}

// Debugf logs a debug message
func (l *AppLogger) Debugf(msg string, args ...interface{}) {
	l.Logger.Debug(fmt.Sprintf(msg, args...))
}
