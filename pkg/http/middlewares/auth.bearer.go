package middlewares

import (
	"context"

	"github.com/dbunt1tled/go-api/internal/modules/auth"
	"github.com/dbunt1tled/go-api/internal/modules/user"
	"github.com/dbunt1tled/go-api/pkg/e"
	"github.com/dbunt1tled/go-api/pkg/hasher"
	"github.com/dbunt1tled/go-api/pkg/http/authctx"
	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
)

type AuthMiddleware struct {
	authService *auth.Service
	userService *user.Service
}

func NewAuthMiddleware(
	authService *auth.Service,
	userService *user.Service,
) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
		userService: userService,
	}
}

func (a *AuthMiddleware) AuthBearer(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		var (
			authToken string
			isEmpty   bool
			token     map[string]any
			u         *user.User
			id        uuid.UUID
			err       error
		)

		authToken, isEmpty = a.authService.TokenFromAuthHeader(c)
		if isEmpty {
			authToken, isEmpty = a.authService.TokenFromQueryParam("token", c)
			if isEmpty {
				return e.NewUnauthorizedError("Unauthorized", e.Err401TokenEmptyError)
			}
		}

		token, err = a.authService.DecodeToken(
			authToken,
			hasher.WithSubject(hasher.AccessTokenSubject),
			hasher.WithExpire(true),
		)
		if err != nil {
			return err
		}
		id, err = uuid.Parse(token["iss"].(string))
		if err != nil {
			return e.NewUnauthorizedError("Unauthorized", e.Err401TokenUserIdError)
		}
		u, err = a.userService.ByID(c.Request().Context(), id)
		if err != nil {
			return e.NewUnauthorizedError("Unauthorized", e.Err401UserNotFoundError)
		}

		if !u.IsActive() {
			return e.NewUnauthorizedError("Unauthorized", e.Err401UserNotActiveError)
		}
		ctx := context.WithValue(c.Request().Context(), authctx.UserKey, u)
		c.SetRequest(c.Request().WithContext(ctx))

		return next(c)
	}
}
