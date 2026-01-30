package authctx

import (
	"reflect"

	"github.com/dbunt1tled/go-api/pkg/e"
	"github.com/labstack/echo/v5"
)

type userKeyType struct{}

var UserKey = userKeyType{}

func AuthUser[T any](c *echo.Context) (T, error) {
	var zero T
	u, ok := c.Request().Context().Value(UserKey).(T)
	if !ok {
		return zero, e.NewUnauthorizedError("Unauthorized", e.Err401UserNotAuth1Error)
	}
	v := reflect.ValueOf(u)
	if v.Kind() == reflect.Pointer && v.IsNil() {
		return zero, e.NewUnauthorizedError("Unauthorized", e.Err401UserNotAuth2Error)
	}

	return u, nil
}
