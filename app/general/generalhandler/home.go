package generalhandler

import (
	"bytes"
	"net/http"

	"github.com/dbunt1tled/go-api/internal/config/app_error"
	"github.com/dbunt1tled/go-api/internal/lib/jsonerror"
	"github.com/dbunt1tled/go-api/internal/util/helper"

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
