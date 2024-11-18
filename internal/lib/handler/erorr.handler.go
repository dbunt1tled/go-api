package handler

import (
	"errors"
	"go_echo/internal/config"
	"go_echo/internal/lib/jsonerror"
	"go_echo/internal/util/helper"
	"net/http"

	"github.com/labstack/echo/v4"
)

func APIErrorHandler(err error, c echo.Context) {

	log := config.GetLoggerInstance()
	status := http.StatusInternalServerError
	message := err.Error()
	code := 0

	var he *echo.HTTPError
	if errors.As(err, &he) {
		status = he.Code
		message = he.Error()
	}

	var exception *jsonerror.ErrException
	if errors.As(err, &exception) {
		status = exception.Status
		message = exception.Error()
		code = exception.Code
	}

	log.Warn(err.Error())
	helper.JSONAPIModel(c.Response(), jsonerror.NewErrorString(message, code, status, nil), status)
}
