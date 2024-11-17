package graceful

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
)

func ShutdownGraceful(log *slog.Logger, server *echo.Echo) chan bool {
	var pending int32 = 0
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan bool, 1)
	go func() {
		sig := <-sigs
		log.Warn("Signal received: " + sig.String())
		signal.Reset(syscall.SIGINT, syscall.SIGTERM)
		gracefulBye(server, &pending, done, log)
	}()
	return done
}
func gracefulBye(httpServer *echo.Echo, pending *int32, done chan bool, log *slog.Logger) {
	log.Warn("quit: shutting down ...")
	defer log.Warn("quit: shutdown completed")
	defer func() {
		done <- true
	}()
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := httpServer.Shutdown(ctx)
	if err != nil {
		panic("Error shutdown: " + err.Error())
	}
	for {
		pendingRequests := atomic.LoadInt32(pending)
		if pendingRequests == 0 {
			break
		}

		log.Warn("waiting for", pendingRequests, "pending requests to complete ...")
		time.Sleep(1 * time.Second)
	}
}
