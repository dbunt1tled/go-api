package middlewares

import (
	"github.com/dbunt1tled/go-api/internal/config/env"
	"strings"

	"github.com/labstack/echo/v4"
)

func Base(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := next(c)
		c.Response().After(func() {
			cfg := env.GetConfigInstance()
			if !(cfg.Static.Enable && strings.HasPrefix(c.Request().URL.Path, "/"+cfg.Static.URL)) &&
				c.Response().Header().Get("Content-Type") == "" {
				c.Response().Header().Set("Content-Type", "application/json")
			}
		})
		return err
	}
}
