package authhandle

import (
	"go_echo/internal/lib/jsonerror"
	"go_echo/internal/util/helper"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

type Request struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func Login(c echo.Context) error {
	var err error
	//var u *user.User
	var req Request

	if err = c.Bind(&req); err != nil {
		return jsonerror.ErrorUnprocessableEntity(c, errors.Wrap(err, "Error validate request"), 42200001)
	}

	if err = validator.New().Struct(req); err != nil {
		var e validator.ValidationErrors
		if errors.As(err, &e) {

			return jsonerror.ErrorUnprocessableEntityMap(c, helper.ValidationErrorString(e), 42200002)
		}
		return jsonerror.ErrorUnprocessableEntity(c, err, 42200003)
	}

	return nil
}
