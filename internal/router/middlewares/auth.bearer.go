package middlewares

import (
	"go_echo/app/user/model/user"
	"go_echo/app/user/service"
	"go_echo/internal/config/app_error"
	"go_echo/internal/lib/jsonerror"
	"go_echo/internal/util/jwt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

const BearerSchema = "Bearer"

func AuthBearer(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var authToken string
		var isEmpty bool
		authToken, isEmpty = fromAuthHeader(c)
		if isEmpty {
			authToken, isEmpty = fromQueryParam(c)
			if isEmpty {
				return &jsonerror.ExceptionErr{
					Inner:  errors.New("unauthorized"),
					Code:   app_error.Err401AuthEmptyTokenError,
					Status: http.StatusUnauthorized,
				}
			}
		}

		token, err := jwt.JWToken{}.Decode(authToken, true)
		if err != nil {
			return &jsonerror.ExceptionErr{
				Inner:  errors.New("unauthorized"),
				Code:   app_error.Err401TokenError,
				Status: http.StatusUnauthorized,
			}
		}
		u, err := service.UserRepository{}.ByID(int64(token["iss"].(float64))) //nolint:nolintlint,errcheck
		if err != nil {
			return &jsonerror.ExceptionErr{
				Inner:  errors.New("unauthorized"),
				Code:   app_error.Err401UserNotFoundError,
				Status: http.StatusUnauthorized,
			}
		}

		if u.Status != user.Active {
			return &jsonerror.ExceptionErr{
				Inner:  errors.New("unauthorized"),
				Code:   app_error.Err401UserNotActiveError,
				Status: http.StatusUnauthorized,
			}
		}
		c.Set("user", u)

		return next(c)
	}
}

func fromAuthHeader(c echo.Context) (string, bool) {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return "", true
	}
	authToken := strings.TrimSpace(strings.Split(authHeader, BearerSchema)[1])
	if authToken == "" {
		return "", true
	}
	return authToken, false
}

func fromQueryParam(c echo.Context) (string, bool) {
	authToken := c.QueryParam("team")
	if authToken == "" {
		return "", true
	}
	return authToken, false
}
