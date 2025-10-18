package usernotificationhandler

import (
	"github.com/dbunt1tled/go-api/app/user/model/user"
	"github.com/dbunt1tled/go-api/app/usernotification/model/usernotification"
	"github.com/dbunt1tled/go-api/app/usernotification/service"
	"github.com/dbunt1tled/go-api/internal/config/app_error"
	"github.com/dbunt1tled/go-api/internal/config/validate"
	"github.com/dbunt1tled/go-api/internal/dto"
	"github.com/dbunt1tled/go-api/internal/lib/jsonerror"
	"github.com/dbunt1tled/go-api/internal/util/builder"
	"github.com/dbunt1tled/go-api/internal/util/builder/page"
	"github.com/dbunt1tled/go-api/internal/util/helper"
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
		u   *user.User
		p   page.Paginate[usernotification.UserNotification]
		req UserNotificationRequest
	)

	u = c.Get("user").(*user.User) //nolint:errcheck //auth middleware
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
				helper.ValidationErrorString(c, e),
				app_error.Err422UserNotificationValidateRuleError,
			)
		}
		return jsonerror.ErrorUnprocessableEntity(c, err, app_error.Err422UserNotificationValidateError)
	}

	p, err = service.UserNotificationRepository{}.Paginator(
		c.Request().Context(),
		&[]page.FilterCondition{
			builder.Eq("user_id", u.ID),
			builder.Eq("status", helper.GetVarValue(req.Status)),
		},
		builder.GetSortOrder(req.Pagination.Sort),
		builder.GetPagination(req.Pagination),
	)
	if err != nil {
		return jsonerror.ErrorUnprocessableEntity(c, err, app_error.Err422UserNotificationQueryError)
	}

	return helper.JSONAPIModel(c.Response(), p, http.StatusOK)
}
