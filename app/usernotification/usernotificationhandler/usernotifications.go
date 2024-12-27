package usernotificationhandler

import (
	"go_echo/app/user/model/user"
	"go_echo/app/usernotification/model/usernotification"
	"go_echo/app/usernotification/service"
	"go_echo/internal/config/app_error"
	"go_echo/internal/config/validate"
	"go_echo/internal/dto"
	"go_echo/internal/lib/jsonerror"
	"go_echo/internal/util/builder"
	"go_echo/internal/util/helper"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

type UserNotificationRequest struct {
	Status     *usernotification.Status `query:"status" validate:"omitempty,number"`
	Pagination dto.PaginationQuery      `json:"pagination" validate:"required"`
}

func UserNotificationList(c echo.Context) error {
	var (
		err error
		un  []*usernotification.UserNotification
		req UserNotificationRequest
	)
	
	if err = c.Bind(&req); err != nil {
		return jsonerror.ErrorUnprocessableEntity(
			c,
			errors.Wrap(err, "Error validate request"),
			app_error.Err422UserNotificationValidateMapError,
		)
	}

	if err = validate.GetValidateInstance().Struct(req); err != nil {
		var e validator.ValidationErrors
		if errors.As(err, &e) {
			return jsonerror.ErrorUnprocessableEntityMap(
				c,
				helper.ValidationErrorString(e),
				app_error.Err422UserNotificationValidateRuleError,
			)
		}
		return jsonerror.ErrorUnprocessableEntity(c, err, app_error.Err422UserNotificationValidateError)
	}

	un, err = service.UserNotificationRepository{}.List(&[]builder.FilterCondition{
		builder.Eq("user_id", c.Get("user").(*user.User).ID),
		builder.Eq("status", helper.GetVarValue(req.Status)),
	}, builder.GetSortOrder(req.Pagination.Sort))

	return helper.JSONAPIModel(c.Response(), un, http.StatusOK)
}
