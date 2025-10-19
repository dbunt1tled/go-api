package logger

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/dbunt1tled/go-api/internal/config/app_error"
)

func Error(err error) []slog.Attr {
	var er app_error.StackTracer
	if errors.As(err, &er) {
		stack := ""
		for _, f := range er.StackTrace() {
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
