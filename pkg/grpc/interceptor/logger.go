package interceptor

import (
	"context"
	"log/slog"
	"time"

	"github.com/dbunt1tled/go-api/internal/config"
	"github.com/dbunt1tled/go-api/pkg/log"

	"google.golang.org/grpc"
)

func LoggingInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		startTime := time.Now()
		resp, err := handler(ctx, req)
		level := slog.LevelInfo
		if err != nil {
			level = slog.LevelError
		}
		var attrs []any
		attrs = append(attrs, slog.String("method", info.FullMethod))
		if err != nil {
			attrs = append(attrs, slog.Any("error", err))
		}
		attrs = append(attrs, slog.Any("request", req))
		if config.Get().Debug {
			attrs = append(attrs, slog.Any("response", resp))
		}
		attrs = append(attrs, slog.String("time", time.Since(startTime).Round(time.Millisecond).String()))
		log.Logger().Log(ctx, level, "Incoming Request", attrs...)
		return resp, err
	}
}
