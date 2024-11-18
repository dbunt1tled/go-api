package server

import (
	"go_echo/internal/util/helper"
	"net/http"
	_ "net/http"

	"github.com/labstack/echo/v4"
	_ "github.com/pkg/errors"
)

type HelmResponse struct {
	ID int64  `json:"id,omitempty" jsonapi:"primary,helm"`
	DB string `json:"db" jsonapi:"attr,db"`
}

func Helm(c echo.Context) error {
	return helper.JSONAPIModel(c.Response(), &HelmResponse{ID: 1, DB: "test"}, http.StatusOK)
}
