package logger

import (
	"fmt"
	"go_echo/internal/config"
	"log/slog"
)

func Error(err error) []slog.Attr {
	if err, ok := err.(config.StackTracer); ok {
		stack := ""
		for _, f := range err.StackTrace() {
			stack += fmt.Sprintf("%+s:%d\n", f, f)
		}
		return []slog.Attr{
			{Key: "stack", Value: slog.StringValue(stack)},
			{Key: "message", Value: slog.StringValue(err.Error())},
		}
	}
	return []slog.Attr{
		{Key: "message", Value: slog.StringValue(err.Error())},
	}
}
