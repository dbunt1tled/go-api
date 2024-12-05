package middlewares

import (
	"github.com/labstack/echo/v4"
)

func Base(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Response().Header().Get("Content-Type") == "" {
			c.Response().Header().Set("Content-Type", "application/json")
		}
		return next(c)
	}
}
