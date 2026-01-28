package auth

import (
	"github.com/bytedance/sonic"
	"github.com/dbunt1tled/go-api/internal/jobs/rmqmail"
	"github.com/dbunt1tled/go-api/internal/jobs/rmqmail/handlers"
	"github.com/dbunt1tled/go-api/internal/modules/user"
	"github.com/dbunt1tled/go-api/pkg/e"
	"github.com/dbunt1tled/go-api/pkg/http"
	"github.com/dbunt1tled/go-api/pkg/log"
	"github.com/dbunt1tled/go-api/pkg/rmq"
	"github.com/dbunt1tled/go-api/pkg/storage"
	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
)

type Controller struct {
	http.BaseController

	authService *Service
	rMProducer  *rmq.Producer
	userService *user.Service
}

func NewAuthController(
	authService *Service,
	rMProducer *rmq.Producer,
	userService *user.Service,
) *Controller {
	return &Controller{
		BaseController: http.NewBaseController(),
		authService:    authService,
		rMProducer:     rMProducer,
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
		data                   []byte
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
	data, err = sonic.ConfigFastest.Marshal(handlers.MailUserConfirmationJobMessage{
		UserID: u.ID,
		Token:  confirmToken,
	})
	if err != nil {
		return e.NewUnprocessableEntityErrorWrap(
			"Generate marshal token error.",
			e.Err422CreateConfirmTokenMarshalError,
			err,
		)
	}
	err = ac.rMProducer.Publish(
		c.Request().Context(),
		rmqmail.MailExchange,
		rmqmail.MailRoutingKey,
		handlers.UserConfirmationEmailJobAction,
		data,
	)
	if err != nil {
		return e.NewUnprocessableEntityErrorWrap(
			"Generate confirm message error.",
			e.Err422RegisterUserGenerateConfirmMsgError,
			err,
		)
	}
	log.Logger().Infof("User registered successfully. Confirm token: %s", confirmToken)

	return ac.JSON200(c, user.NewUserResource(u))
}

func (ac *Controller) Confirm(c *echo.Context) error {
	var (
		err   error
		id    uuid.UUID
		token map[string]interface{}
		u     *user.User
	)
	req := new(Confirm)

	if err = ac.BindAndValidate(c, req, e.Err422UserConfirmValidateError); err != nil {
		return err
	}
	token, err = ac.authService.DecodeConfirmToken(req.Token)
	if err != nil {
		return e.NewUnprocessableEntityErrorWrap(
			"Token error.",
			e.Err422UserConfirmTokenDecodeError,
			err,
		)
	}
	userId, ok := token["iss"]
	if !ok {
		return e.NewUnprocessableEntityError(
			"Token error.",
			e.Err422UserConfirmTokenError,
		)
	}
	id, err = uuid.Parse(userId.(string))
	if err != nil {
		return e.NewNotFoundErrorWrap(
			"User id error.",
			e.Err422UserConfirmTokenUserIdError,
			err,
		)
	}
	u, err = ac.userService.ByID(c.Request().Context(), id)
	if err != nil {
		return e.NewNotFoundErrorWrap(
			"User not found.",
			e.Err404UserConfirmTokenUserNotFoundError,
			err,
		)
	}
	if u.Status != user.Pending {
		return e.NewUnprocessableEntityError(
			"User is not in pending state.",
			e.Err422UserConfirmTokenUserNotPendingError,
		)
	}
	u.Status = user.Active
	u, err = ac.userService.Update(c.Request().Context(), u)
	if err != nil {
		return e.NewUnprocessableEntityErrorWrap(
			"User update error.",
			e.Err422UserConfirmTokenUserUpdateError,
			err,
		)
	}
	return ac.JSON200(c, user.NewUserResource(u))
}
