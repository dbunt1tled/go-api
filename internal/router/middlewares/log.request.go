package middlewares

import (
	"github.com/dbunt1tled/go-api/internal/config/env"
	"github.com/dbunt1tled/go-api/internal/util/helper"
	"slices"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

//nolint:gocognit
func LogRequest(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cfg := env.GetConfigInstance()
		if !cfg.Debug.Debug {
			return next(c)
		}
		if cfg.Debug.DebugRequest { //nolint:nestif
			loggerMiddleware := middleware.Logger()
			path := c.Path()
			urls := []string{"/", "/system/helm"}
			if cfg.Debug.DebugBody && !slices.Contains(urls, path) && !strings.HasPrefix(path, "/"+cfg.Static.URL) {
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
