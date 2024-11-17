package mware

import (
	_ "go_echo/app/user/model/user"
	"go_echo/internal/lib/jsonerror"
	"go_echo/internal/util/jwt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

const BearerSchema = "Bearer"

func authBearer(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var authToken string
		var isEmpty bool
		authToken, isEmpty = fromAuthHeader(c)
		if isEmpty {
			authToken, isEmpty = fromQueryParam(c)
			if isEmpty {
				return &jsonerror.ErrException{
					Inner:  errors.New("Unauthorized."),
					Code:   40100001,
					Status: http.StatusUnauthorized,
				}
			}
		}

		_, err := jwt.JWToken{}.Decode(authToken, true)
		if err != nil {
			return &jsonerror.ErrException{
				Inner:  errors.New("Unauthorized."),
				Code:   40100002,
				Status: http.StatusUnauthorized,
			}
		}
		//u, err := storage.GetUser(mysql.UserFilter{ID: int64(token["iss"].(float64))})
		//if err != nil {
		//	return &jsonerror.ErrException{Inner: errors.New("Unauthorized."), Code: 40100003, Status: http.StatusUnauthorized}
		//}
		//
		//if u.Status != user.StatusActive {
		//	return &jsonerror.ErrException{Inner: errors.New("Unauthorized."), Code: 40100004, Status: http.StatusUnauthorized}
		//}

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
