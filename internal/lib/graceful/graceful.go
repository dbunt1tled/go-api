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

func ShutDown(log *slog.Logger, httpServer *echo.Echo) chan bool {
	var pending int32 = 0
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan bool, 1)
	go func() {
		sig := <-sigs
		log.Warn("Signal received: " + sig.String())
		signal.Reset(syscall.SIGINT, syscall.SIGTERM)
		gracefulBye(httpServer, &pending, done, log)
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

	httpServer.Shutdown(ctx)
	for {
		pending_ := atomic.LoadInt32(pending)
		if pending_ == 0 {
			break
		}

		log.Warn("waiting for", pending_, "pending requests to complete ...")
		time.Sleep(1 * time.Second)
	}
}
