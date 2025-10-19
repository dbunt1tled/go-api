package service

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/dbunt1tled/go-api/app/usernotification/model/usernotification"
	"github.com/dbunt1tled/go-api/internal/util/builder"
	"github.com/dbunt1tled/go-api/internal/util/builder/page"
	"github.com/dbunt1tled/go-api/internal/util/helper"
	"github.com/dbunt1tled/go-api/internal/util/type/json"

	"github.com/pkg/errors"
)

const (
	UserNotificationTableName = "user_notifications"
)

type UserNotificationRepository struct {
}

type UserNotificationParams struct {
	Data   *json.JsonField
	UserID *int64
	Status *usernotification.Status
}

func (r UserNotificationRepository) ByID(ctx context.Context, id int64) (*usernotification.UserNotification, error) {
	return builder.ByID[usernotification.UserNotification](
		ctx,
		UserNotificationTableName,
		id,
		castUserNotificationRow,
	)
}

func (r UserNotificationRepository) Paginator(
	ctx context.Context,
	filter *[]page.FilterCondition,
	sorts *[]page.SortOrder,
	paginator *page.Pagination,
) (page.Paginate[usernotification.UserNotification], error) {
	return builder.Paginator[usernotification.UserNotification](
		ctx,
		UserNotificationTableName,
		filter,
		sorts,
		paginator,
		castUserNotificationRows,
	)
}

func (r UserNotificationRepository) One(
	ctx context.Context,
	filter *[]page.FilterCondition,
	sorts *[]page.SortOrder,
) (*usernotification.UserNotification, error) {
	var _validFields = map[string]bool{
		"id":      true,
		"user_id": true,
		"status":  true,
	}

	if filter != nil && len(*filter) > 0 {
		if err := builder.ValidateFilter(*filter, _validFields); err != nil {
			return nil, err
		}
	}

	return builder.One(ctx, UserNotificationTableName, filter, sorts, castUserNotificationRow)
}

func (r UserNotificationRepository) List(
	ctx context.Context,
	filter *[]page.FilterCondition,
	sorts *[]page.SortOrder,
) ([]*usernotification.UserNotification, error) {
	var _validFields = map[string]bool{
		"id":      true,
		"user_id": true,
		"status":  true,
	}
	var err error
	if filter != nil && len(*filter) > 0 {
		if err = builder.ValidateFilter(*filter, _validFields); err != nil {
			return nil, err
		}
	}
	return builder.List(ctx, UserNotificationTableName, filter, sorts, castUserNotificationRows, nil)
}

func (r UserNotificationRepository) Count(
	ctx context.Context,
	filter *[]page.FilterCondition,
) (int, error) {
	var _validFields = map[string]bool{
		"id":      true,
		"user_id": true,
		"status":  true,
	}
	var (
		err error
	)
	if filter != nil && len(*filter) > 0 {
		if err = builder.ValidateFilter(*filter, _validFields); err != nil {
			return 0, err
		}
	}
	return builder.Count(ctx, UserNotificationTableName, filter)
}

func (r UserNotificationRepository) Create(ctx context.Context, params UserNotificationParams) (*usernotification.UserNotification, error) {
	var (
		columns []string
		values  []string
		args    []interface{}
	)

	if params.UserID != nil {
		columns = append(columns, "user_id")
		values = append(values, "?")
		args = append(args, *params.UserID)
	}

	if params.Data != nil {
		columns = append(columns, "data")
		values = append(values, "?")
		args = append(args, *params.Data)
	}

	columns = append(columns, "status")
	values = append(values, "?")
	if params.Status != nil {
		args = append(args, *params.Status)
	} else {
		args = append(args, usernotification.New)
	}

	if len(columns) == 0 {
		return nil, errors.New("no fields to insert")
	}

	columns = append(columns, "created_at")
	values = append(values, "?")
	args = append(args, time.Now())

	columns = append(columns, "updated_at")
	values = append(values, "?")
	args = append(args, time.Now())

	smt, err := builder.GetDB().Prepare(
		fmt.Sprintf(
			"INSERT INTO user_notifications (%s) VALUES (%s)",
			strings.Join(columns, ", "),
			strings.Join(values, ", ")))
	if err != nil {
		return nil, errors.Wrap(err, "create user prepare error")
	}
	defer smt.Close()
	res, err := smt.ExecContext(ctx, args...)
	if err != nil {
		return nil, errors.Wrap(err, "create user error")
	}
	return helper.Must(r.ByID(ctx, helper.Must(res.LastInsertId()))), nil
}

func (r UserNotificationRepository) Update(
	ctx context.Context,
	id int64, params UserNotificationParams,
) (*usernotification.UserNotification, error) {

	var (
		setClauses []string
		args       []interface{}
	)

	if params.UserID != nil {
		setClauses = append(setClauses, "user_id = ?")
		args = append(args, *params.UserID)
	}

	if params.Data != nil {
		setClauses = append(setClauses, "data = ?")
		args = append(args, *params.Data)
	}

	if params.Status != nil {
		setClauses = append(setClauses, "status = ?")
		args = append(args, *params.Status)
	}

	if len(setClauses) == 0 {
		return helper.Must(r.ByID(ctx, id)), nil
	}

	setClauses = append(setClauses, "updated_at = ?")
	args = append(args, time.Now())

	args = append(args, id)
	query := fmt.Sprintf("UPDATE user_notifications SET %s WHERE id = ?", strings.Join(setClauses, ", "))

	smt, err := builder.GetDB().Prepare(query)
	if err != nil {
		return nil, errors.Wrap(err, "update user prepare error")
	}
	defer smt.Close()
	_, err = smt.ExecContext(ctx, args...)
	if err != nil {
		return nil, errors.Wrap(err, "update user error")
	}

	return helper.Must(r.ByID(ctx, id)), nil
}

func castUserNotificationRow(row *sql.Row) (*usernotification.UserNotification, error) {
	un := usernotification.UserNotification{}
	return builder.ScanStructRow(un, row)
}

func castUserNotificationRows(rows *sql.Rows) (*usernotification.UserNotification, error) {
	un := usernotification.UserNotification{}
	return builder.ScanStructRows(un, rows)
}
