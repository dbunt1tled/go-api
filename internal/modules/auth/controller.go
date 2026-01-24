package auth

import (
	"github.com/dbunt1tled/go-api/internal/modules/user"
	"github.com/dbunt1tled/go-api/pkg/e"
	"github.com/dbunt1tled/go-api/pkg/http"
	"github.com/dbunt1tled/go-api/pkg/log"
	"github.com/dbunt1tled/go-api/pkg/storage"
	"github.com/labstack/echo/v5"
)

type Controller struct {
	http.BaseController
	authService *Service
	userService *user.Service
}

func NewAuthController(
	authService *Service,
	userService *user.Service,
) *Controller {
	return &Controller{
		BaseController: http.NewBaseController(),
		authService:    authService,
		userService:    userService,
	}
}

func (ac *Controller) Login(c *echo.Context) error {
	var (
		err             error
		u               *user.User
		ch              bool
		access, refresh string
	)

	req := new(Login)

	if err = ac.BindAndValidate(c, req, e.Err422LoginValidateError); err != nil {
		return err
	}

	u, err = ac.userService.One(
		c.Request().Context(),
		storage.WithFilter(storage.NewRule("status", storage.OpEqual, user.Active)),
	)
	if err != nil {
		return e.NewUnprocessableEntityError(
			"Authorization error, password or login is incorrect.",
			e.Err422LoginUserNotFoundError,
		)
	}

	ch, err = ac.authService.ValidatePassword(req.Password, u.Password)
	if err != nil {
		return e.NewUnprocessableEntityError(
			err.Error(),
			e.Err422LoginUserPasswordError,
		)
	}

	if !ch {
		return e.NewUnprocessableEntityError(
			"Authorization error, password or login is incorrect.",
			e.Err422LoginUserPasswordWrongError,
		)
	}

	access, refresh, err = ac.authService.GenerateAuthTokens(u)
	if err != nil {
		return e.NewUnprocessableEntityError(
			err.Error(),
			e.Err422LoginAccessTokenError,
		)
	}

	return ac.JSON200(c, NewLoginResponse(map[string]interface{}{
		"accessToken":  access,
		"refreshToken": refresh,
	}))

}
func (ac *Controller) Register(c *echo.Context) error {
	var (
		err                    error
		password, confirmToken string
		u                      *user.User
	)
	req := new(Register)

	if err = ac.BindAndValidate(c, req, e.Err422UserListValidateError); err != nil {
		return err
	}

	password, err = ac.authService.GeneratePasswordHash(req.Password)
	if err != nil {
		return e.NewUnprocessableEntityErrorWrap(
			"Password hash error.",
			e.Err422RegisterUserPasswordError,
			err,
		)
	}

	u, err = ac.userService.Create(c.Request().Context(), req.ToUser().WithPassword(password))

	if err != nil {
		return e.NewUnprocessableEntityErrorWrap(
			"User creation error.",
			e.Err422RegisterUserCreationError,
			err,
		)
	}

	confirmToken, err = ac.authService.GenerateConfirmToken(u)
	if err != nil {
		return e.NewUnprocessableEntityErrorWrap(
			"Generate token error.",
			e.Err422CreateConfirmTokenError,
			err,
		)
	}
	log.Logger().Infof("User registered successfully. Confirm token: %s", confirmToken)

	return ac.JSON200(c, user.NewUserResource(u))
}
