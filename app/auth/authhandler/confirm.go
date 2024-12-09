package authhandler

import (
	"go_echo/app/auth/service/auth"
	"go_echo/app/user/model/user"
	"go_echo/app/user/service"
	"go_echo/internal/config/app_error"
	"go_echo/internal/config/validate"
	"go_echo/internal/lib/jsonerror"
	"go_echo/internal/util/helper"
	"go_echo/internal/util/jwt"
	"go_echo/internal/util/type/user_status"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

type ConfirmRequest struct {
	Token string `query:"token" json:"token" validate:"required"`
}

func Confirm(c echo.Context) error {
	var (
		err   error
		token map[string]interface{}
		u     *user.User
		req   ConfirmRequest
	)

	if err = c.Bind(&req); err != nil {
		return jsonerror.ErrorUnprocessableEntity(
			c,
			errors.Wrap(err, "Error validate request"),
			app_error.Err422ConfirmRequestError,
		)
	}

	if err = validate.GetValidateInstance().Struct(req); err != nil {
		var e validator.ValidationErrors
		if errors.As(err, &e) {
			return jsonerror.ErrorUnprocessableEntityMap(
				c,
				helper.ValidationErrorString(e),
				app_error.Err422ConfirmValidateRuleError,
			)
		}
		return jsonerror.ErrorUnprocessableEntity(c, err, app_error.Err422ConfirmValidateError)
	}

	token, err = jwt.JWToken{}.Decode(req.Token, true)
	if err != nil || token["sub"] != auth.ConfirmTokenSubject {
		return jsonerror.ErrorUnprocessableEntity(c, err, app_error.Err422ConfirmTokenError)
	}
	u, err = service.UserRepository{}.ByID(int64(token["iss"].(float64))) //nolint:errcheck
	if err != nil {
		return jsonerror.ErrorUnprocessableEntity(c, err, app_error.Err422ConfirmUserError)
	}

	if u.Status != user_status.Pending.Val() {
		return jsonerror.ErrorUnprocessableEntityString(c, "Users already confirmed your account", app_error.Err422ConfirmUserStatusError)
	}

	status := user_status.Active
	date := time.Now()
	u, err = service.UserRepository{}.Update(u.ID, service.UpdateUserParams{
		Status:      &status,
		ConfirmedAt: &date,
	})
	if err != nil {
		return jsonerror.ErrorUnprocessableEntity(c, err, app_error.Err422ConfirmUserNotFoundError)
	}

	return helper.JSONAPIModel(c.Response(), u, http.StatusOK)
}
