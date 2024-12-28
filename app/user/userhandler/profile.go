package userhandler

import (
	"go_echo/app/user/model/user"
	"go_echo/internal/util/helper"
	"net/http"

	"github.com/labstack/echo/v4"
)

func Profile(c echo.Context) error {
	u := c.Get("user").(*user.User) //nolint:errcheck //auth middleware
	return helper.JSONAPIModel(c.Response(), u, http.StatusOK)
}
