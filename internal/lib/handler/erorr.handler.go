package handler

import (
	"errors"
	"go_echo/internal/config"
	"go_echo/internal/lib/jsonerror"
	"net/http"

	"github.com/google/jsonapi"
	"github.com/labstack/echo/v4"
)

func CustomHTTPErrorHandler(err error, c echo.Context) {
	log := config.GetLoggerInstance()
	code := http.StatusInternalServerError
	var he *echo.HTTPError
	if errors.As(err, &he) {
		code = he.Code
	}
	log.Error(err.Error())
	e := jsonapi.MarshalPayload(c.Response(), jsonerror.NewError(err, code, code))
	if e != nil {
		log.Error(e.Error())
		c.JSON(code, e.Error())
		return
	}
	c.JSON(code, c.Response())
}
