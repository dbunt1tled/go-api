package authhandle

import (
	"database/sql"
	"go_echo/app/auth/model/token"
	"go_echo/app/auth/service/auth"
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
	FirstName       string `json:"FirstName" validate:"required"`
	SecondName      string `json:"secondName" validate:"required"`
	Email           string `json:"email" validate:"required, email"`
	PhoneNumber     string `json:"phoneNumber" validate:"required, phone"`
	Password        string `json:"password" validate:"required, passwd, eqfield=PasswordConfirm"`
	PasswordConfirm string `json:"passwordConfirm" validate:"required"`
}

func Register(c echo.Context) error {
	var (
		err    error
		u      *user.User
		req    RegisterRequest
		tokens *token.Tokens
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

	u, err = service.UserRepository{}.Create(
		req.FirstName,
		req.SecondName,
		req.Email,
		req.PhoneNumber,
		req.Password,
		user_status.Pending,
		"",
		[]string{},
		sql.NullTime{},
	)
	if err != nil {
		return jsonerror.ErrorUnprocessableEntity(c, err, app_error.Err422SignupUserNotFoundError)
	}

	tokens, err = auth.GetAuthTokens(*u)
	if err != nil {
		return jsonerror.ErrorUnprocessableEntity(c, err, app_error.Err422SignupAuthTokensError)
	}

	return helper.JSONAPIModel(c.Response(), tokens, http.StatusOK)
}
