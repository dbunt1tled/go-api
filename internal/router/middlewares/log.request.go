package middlewares

import (
	"go_echo/internal/config/env"
	"go_echo/internal/util/helper"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func LogRequest(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cfg := env.GetConfigInstance()
		if !cfg.Debug.Debug {
			return next(c)
		}
		if cfg.Debug.DebugRequest {
			loggerMiddleware := middleware.Logger()
			if cfg.Debug.DebugBody {
				bodyDumpMiddleware := middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
					rqB := string(reqBody)
					rsB := string(resBody)
					if rqB == "" {
						rqB = "null"
					}
					if rsB == "" {
						rsB = "null"
					}
					id := helper.RequestID(c)
					m := `{"info":"RequestDump","id":"` +
						id + `","request":` +
						rqB + `,"response":` +
						rsB + `}`
					_, err := c.Logger().Output().Write([]byte(m))
					if err != nil {
						return
					}
				})

				return bodyDumpMiddleware(loggerMiddleware(next))(c)
			}
			return loggerMiddleware(next)(c)
		}
		return next(c)
	}
}
