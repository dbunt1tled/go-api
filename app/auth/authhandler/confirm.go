package authhandler

import (
	"net/http"
	"time"

	"github.com/dbunt1tled/go-api/app/auth/service/auth"
	"github.com/dbunt1tled/go-api/app/user/model/user"
	"github.com/dbunt1tled/go-api/app/user/service"
	"github.com/dbunt1tled/go-api/internal/config/app_error"
	"github.com/dbunt1tled/go-api/internal/config/validate"
	"github.com/dbunt1tled/go-api/internal/lib/jsonerror"
	"github.com/dbunt1tled/go-api/internal/util/helper"
	"github.com/dbunt1tled/go-api/internal/util/jwt"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

type ConfirmRequest struct {
	Token string `query:"token" json:"token" validate:"required,jwt"`
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
				helper.ValidationErrorString(c, e),
				app_error.Err422ConfirmValidateRuleError,
			)
		}
		return jsonerror.ErrorUnprocessableEntity(c, err, app_error.Err422ConfirmValidateError)
	}

	token, err = jwt.JWToken{}.Decode(req.Token, true)
	if err != nil || token["sub"] != auth.ConfirmTokenSubject {
		return jsonerror.ErrorUnprocessableEntity(c, err, app_error.Err422ConfirmTokenError)
	}
	u, err = service.UserRepository{}.ByID(c.Request().Context(), int64(token["iss"].(float64))) //nolint:nolintlint,errcheck
	if err != nil {
		return jsonerror.ErrorUnprocessableEntity(c, err, app_error.Err422ConfirmUserError)
	}

	if u.Status != user.Pending {
		return jsonerror.ErrorUnprocessableEntityString(
			c,
			"Users already confirmed your account",
			app_error.Err422ConfirmUserStatusError,
		)
	}

	status := user.Active
	date := time.Now()
	u, err = service.UserRepository{}.Update(c.Request().Context(), u.ID, service.UpdateUserParams{
		Status:      &status,
		ConfirmedAt: &date,
	})
	if err != nil {
		return jsonerror.ErrorUnprocessableEntity(c, err, app_error.Err422ConfirmUserNotFoundError)
	}

	return helper.JSONAPIModel(c.Response(), u, http.StatusOK)
}
