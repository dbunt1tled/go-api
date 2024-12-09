package authhandler

import (
	"go_echo/app/auth/model/token"
	"go_echo/app/auth/service/auth"
	"go_echo/app/user/model/user"
	"go_echo/app/user/service"
	"go_echo/internal/config/app_error"
	"go_echo/internal/config/validate"
	"go_echo/internal/lib/jsonerror"
	"go_echo/internal/util/helper"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

type LoginRequest struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func Login(c echo.Context) error {
	var (
		err    error
		u      *user.User
		req    LoginRequest
		tokens *token.Tokens
	)

	if err = c.Bind(&req); err != nil {
		return jsonerror.ErrorUnprocessableEntity(
			c,
			errors.Wrap(err, "Error validate request"),
			app_error.Err422LoginValidateMapError,
		)
	}

	if err = validate.GetValidateInstance().Struct(req); err != nil {
		var e validator.ValidationErrors
		if errors.As(err, &e) {
			return jsonerror.ErrorUnprocessableEntityMap(
				c,
				helper.ValidationErrorString(e),
				app_error.Err422LoginValidateRuleError,
			)
		}
		return jsonerror.ErrorUnprocessableEntity(c, err, app_error.Err422LoginValidateError)
	}

	u, err = service.UserRepository{}.ByIdentity(req.Login, req.Password)
	if err != nil {
		return jsonerror.ErrorUnprocessableEntity(c, err, app_error.Err422LoginUserNotFoundError)
	}

	tokens, err = auth.GetAuthTokens(*u)
	if err != nil {
		return jsonerror.ErrorUnprocessableEntity(c, err, app_error.Err422LoginAuthTokensError)
	}

	return helper.JSONAPIModel(c.Response(), tokens, http.StatusOK)
}
