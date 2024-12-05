package middlewares

import (
	"go_echo/internal/config/app_error"
	"go_echo/internal/config/env"
	"go_echo/internal/lib/jsonerror"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

func SystemAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		authToken, isEmpty := fromAuthSystemHeader(c)
		if isEmpty {
			return &jsonerror.ExceptionErr{
				Inner:  errors.New("unauthorized"),
				Code:   app_error.Err401SystemEmptyTokenError,
				Status: http.StatusUnauthorized,
			}
		}
		cfg := env.GetConfigInstance()
		if authToken != cfg.JWT.SystemAPIKey {
			return &jsonerror.ExceptionErr{
				Inner:  errors.New("unauthorized"),
				Code:   app_error.Err401SystemTokenError,
				Status: http.StatusUnauthorized,
			}
		}

		return next(c)
	}
}

func fromAuthSystemHeader(c echo.Context) (string, bool) {
	authHeader := c.Request().Header.Get("X-Api-Key")
	if authHeader == "" {
		return "", true
	}
	authToken := strings.TrimSpace(authHeader)
	if authToken == "" {
		return "", true
	}
	return authToken, false
}
