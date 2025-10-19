package userhandler

import (
	"net/http"

	"github.com/dbunt1tled/go-api/app/user/model/user"
	"github.com/dbunt1tled/go-api/internal/util/helper"

	"github.com/labstack/echo/v4"
)

func Profile(c echo.Context) error {
	u := c.Get("user").(*user.User) //nolint:errcheck //auth middleware
	return helper.JSONAPIModel(c.Response(), u, http.StatusOK)
}
