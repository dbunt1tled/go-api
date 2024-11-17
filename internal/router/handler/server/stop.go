package server

import (
	"context"
	_ "net/http"
	"time"

	"github.com/labstack/echo/v4"
	_ "github.com/pkg/errors"
)

func ShutDown(c echo.Context) error {
	//log := config.GetLoggerInstance()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()
	return c.Echo().Shutdown(ctx)

	//err := c.Echo().Shutdown(context.Background())
	//if err != nil {
	//	if !errors.Is(err, http.ErrServerClosed) {
	//		log.Warn("Shutting down server via api ")
	//	} else {
	//		log.Info("Error Shutting down server via api: " + err.Error())
	//	}
	//}
	//return nil
}
