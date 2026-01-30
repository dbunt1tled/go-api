package user_notification

import (
	"github.com/dbunt1tled/go-api/pkg/e"
	"github.com/dbunt1tled/go-api/pkg/http"
	"github.com/dbunt1tled/go-api/pkg/storage"
	"github.com/labstack/echo/v5"
)

type Controller struct {
	http.BaseController

	userNotificationService *Service
}

func NewUserNotificationController(
	userNotificationService *Service,
) *Controller {
	return &Controller{
		BaseController:          http.NewBaseController(),
		userNotificationService: userNotificationService,
	}
}

func (uc *Controller) List(c *echo.Context) error {
	var (
		err   error
		users *storage.Paginator[*UserNotification]
	)
	req := new(ListRequest)

	if err = uc.BindAndValidate(c, req, e.Err422UserListValidateError); err != nil {
		return err
	}

	users, err = uc.userNotificationService.Paginate(
		c.Request().Context(),
		req.Page.Page,
		req.Page.Limit,
		storage.WithFilter(
			storage.NewRule("status", storage.OpIn, req.Status),
		),
		storage.WithSort(req.Sort.Field, req.Sort.Order),
	)

	if err != nil {
		return e.NewUnprocessableEntityError(
			err.Error(),
			e.Err422UserListError,
		)
	}

	return uc.JSON200(c, NewUserNotificationListResponse(users))
}
