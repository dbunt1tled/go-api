package er

import (
	"errors"
	"net/http"

	"github.com/dbunt1tled/go-api/internal/config"
	"github.com/dbunt1tled/go-api/pkg/e"
	"github.com/dbunt1tled/go-api/pkg/http/dto"
	"github.com/dbunt1tled/go-api/pkg/log"
	"github.com/dbunt1tled/go-api/pkg/validation"
	"github.com/labstack/echo/v5"
)

func APIErrorHandler(c *echo.Context, err error) {
	status := http.StatusInternalServerError
	message := err.Error()
	code := 0

	var (
		er    error
		errNo *e.ErrNo
		he    *echo.HTTPError
	)
	if errors.As(err, &he) {
		status = he.Code
	}

	if errors.As(err, &errNo) {
		status = errNo.Status
		code = errNo.Code
	}

	vErr := validation.ErrorValidation(err)
	if vErr != nil {
		status = http.StatusUnprocessableEntity
		er = c.JSON(status, dto.Document{
			Errors: vErr,
		})
		log.Logger().ErrorWithStack(er.Error(), er)

		return
	}
	log.Logger().ErrorWithStack(message, err)
	if config.Get().Debug {
		stack := e.GetErrTrace(err)
		er = c.JSON(status, dto.Document{
			Errors: []e.ErrNo{{Status: status, Msg: message, Code: code, Stack: stack}},
		})
		log.Logger().ErrorWithStack(er.Error(), er)

		return
	}

	er = c.JSON(status, dto.Document{
		Errors: []e.ErrNo{{Status: status, Msg: message, Code: code}},
	})
	log.Logger().ErrorWithStack(er.Error(), er)
}
