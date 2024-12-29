package middlewares

import (
	"go_echo/internal/config/env"
	"go_echo/internal/config/logger"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func LogRequest(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		log := logger.GetLoggerInstance()
		cfg := env.GetConfigInstance()
		if !cfg.Debug.Debug {
			return next(c)
		}
		if cfg.Debug.DebugRequest {
			loggerMiddleware := middleware.Logger()
			if cfg.Debug.DebugBody {
				bodyDumpMiddleware := middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
					id := c.Request().Header.Get(echo.HeaderXRequestID)
					if id == "" {
						id = c.Response().Header().Get(echo.HeaderXRequestID)
					}
					log.Info("Request(" + id + "):" + string(reqBody))
					log.Info("Response(" + id + "):" + string(resBody))
				})

				return bodyDumpMiddleware(loggerMiddleware(next))(c)
			}
			return loggerMiddleware(next)(c)
		}
		return next(c)
	}
}
