package authhandler

import (
	"go_echo/app/auth/service/auth"
	"go_echo/app/jobs/rmqmail/handlers"
	"go_echo/app/user/model/user"
	"go_echo/app/user/service"
	"go_echo/internal/config/app_error"
	"go_echo/internal/config/validate"
	"go_echo/internal/lib/jsonerror"
	"go_echo/internal/util/helper"
	"go_echo/internal/util/type/user_status"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

type RegisterRequest struct {
	FirstName       string `json:"firstName" validate:"required"`
	SecondName      string `json:"secondName" validate:"required"`
	Email           string `json:"email" validate:"required,email,unique_db=users#email"`
	PhoneNumber     string `json:"phoneNumber" validate:"required,unique_db=users#phone_number"`
	Password        string `json:"password" validate:"required,passwd,eqfield=PasswordConfirm"`
	PasswordConfirm string `json:"passwordConfirm" validate:"required"`
}

func Register(c echo.Context) error {
	var (
		err  error
		u    *user.User
		req  RegisterRequest
		code string
	)
	if err = c.Bind(&req); err != nil {
		return jsonerror.ErrorUnprocessableEntity(
			c,
			errors.Wrap(err, "Error validate request"),
			app_error.Err422SignupValidateMapError,
		)
	}

	if err = validate.GetValidateInstance().Struct(req); err != nil {
		var e validator.ValidationErrors
		if errors.As(err, &e) {
			return jsonerror.ErrorUnprocessableEntityMap(
				c,
				helper.ValidationErrorString(e),
				app_error.Err422SignupValidateRuleError,
			)
		}
		return jsonerror.ErrorUnprocessableEntity(c, err, app_error.Err422SignupValidateError)
	}
	status := user_status.Pending
	u, err = service.UserRepository{}.Create(service.CreateUserParams{
		FirstName:   &req.FirstName,
		SecondName:  &req.SecondName,
		Email:       &req.Email,
		PhoneNumber: &req.PhoneNumber,
		Password:    &req.Password,
		Status:      &status,
		Hash:        nil,
		Roles:       nil,
		ConfirmedAt: nil,
	})
	if err != nil {
		return jsonerror.ErrorUnprocessableEntity(c, err, app_error.Err422SignupUserNotFoundError)
	}
	code, err = auth.GenerateConfirmToken(*u)
	if err != nil {
		return jsonerror.ErrorUnprocessableEntity(c, err, app_error.Err422SignupAuthTokensError)
	}
	handlers.UserConfirmationEmail{}.Send(u.ID, code)
	return helper.JSONAPIModel(c.Response(), u, http.StatusOK)
}
