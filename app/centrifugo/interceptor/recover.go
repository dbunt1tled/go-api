package interceptor

import (
	"context"
	"fmt"
	"go_echo/internal/config/logger"
	"runtime"

	"google.golang.org/grpc"
)

func RecoverInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ any, err error) {
		defer func() {
			if r := recover(); r != nil {
				log := logger.GetLoggerInstance()
				stack := make([]byte, 64<<10)
				stack = stack[:runtime.Stack(stack, false)]
				err = &PanicError{Method: info.FullMethod, Panic: r, Stack: stack}
				log.ErrorContext(ctx, "RecoverInterceptor Panic",err)
			}
		}()

		return handler(ctx, req)
	}
}

type PanicError struct {
	Method string
	Panic  any
	Stack  []byte
}

func (e *PanicError) Error() string {
	return fmt.Sprintf("[PANIC RECOVER] Method %s: %v\n\n%s", e.Method, e.Panic, e.Stack)
}
