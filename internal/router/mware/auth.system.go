package mware

import (
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
				Inner:  errors.New("Unauthorized."),
				Code:   40100005,
				Status: http.StatusUnauthorized,
			}
		}
		cfg := env.GetConfigInstance()
		if authToken != cfg.JWT.SystemAPIKey {
			return &jsonerror.ExceptionErr{
				Inner:  errors.New("Unauthorized."),
				Code:   40100006,
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
