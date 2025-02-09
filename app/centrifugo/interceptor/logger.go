package interceptor

import (
	"context"
	"go_echo/internal/config/env"
	"go_echo/internal/config/logger"
	"log/slog"
	"time"

	"google.golang.org/grpc"
)

func LoggingInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ any, err error) {
		startTime := time.Now()
		log := logger.GetLoggerInstance()
		cfg := env.GetConfigInstance()
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
		if cfg.Debug.DebugBody {
			attrs = append(attrs, slog.Any("response", resp))
		}
		attrs = append(attrs, slog.String("time", time.Since(startTime).Round(time.Millisecond).String()))
		log.Log(ctx, level, "Incoming Request", attrs...)
		return resp, err
	}
}
