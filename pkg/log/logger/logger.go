package logger

import (
	"log/slog"

	"github.com/dbunt1tled/go-api/pkg/e"
)

func Error(err error) []slog.Attr {
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
