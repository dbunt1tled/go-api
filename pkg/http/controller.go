package http

import (
	"net/http"

	"github.com/dbunt1tled/go-api/pkg/e"
	"github.com/dbunt1tled/go-api/pkg/http/dto"
	"github.com/labstack/echo/v5"
)

type BaseController struct {
}

func NewBaseController() BaseController {
	return BaseController{}
}

func (b *BaseController) Bind(c *echo.Context, dst any, code int) error {
	if err := c.Bind(dst); err != nil {
		return e.NewUnprocessableEntityError("invalid body", code)
	}

	return nil
}

func (b *BaseController) BindAndValidate(c *echo.Context, dst any, code int) error {
	var err error
	err = b.Bind(c, dst, code)
	if err != nil {
		return err
	}

	if dst, ok := dst.(dto.SetDefaults); ok {
		dst.SetDefaults()
	}

	if err = c.Validate(dst); err != nil {
		return err
	}

	return nil
}

func (b *BaseController) JSON(c *echo.Context, status int, data any) error {
	return c.JSON(status, data)
}

func (b *BaseController) JSON200(c *echo.Context, data any) error {
	return b.JSON(c, http.StatusOK, data)
}

func (b *BaseController) JSON201(c *echo.Context, data any) error {
	return b.JSON(c, http.StatusCreated, data)
}

func (b *BaseController) Error(msg string, code int, status int) error {
	return e.NewErrNo(msg, code, status)
}
