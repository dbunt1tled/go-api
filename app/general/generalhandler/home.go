package generalhandler

import (
	"bytes"
	"go_echo/internal/config/app_error"
	"go_echo/internal/lib/jsonerror"
	"go_echo/internal/util/helper"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

func Home(c echo.Context) error {
	var doc bytes.Buffer
	err := helper.
		GetTemplate("general/home.gohtml").
		Execute(&doc, helper.MakeTemplateData(map[string]any{}))
	if err != nil {
		return jsonerror.ErrorUnprocessableEntity(
			c,
			errors.Wrap(err, "Error render template"),
			app_error.Err422HomeGeneralError,
		)
	}

	return c.HTML(http.StatusOK, doc.String())
}
