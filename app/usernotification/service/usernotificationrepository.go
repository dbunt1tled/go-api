package service

import (
	"database/sql"
	"fmt"
	"go_echo/app/usernotification/model/usernotification"
	"go_echo/internal/util/builder"
	"go_echo/internal/util/helper"
	"go_echo/internal/util/type/json"
	"strings"
	"time"

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

func (r UserNotificationRepository) ByID(id int64) (*usernotification.UserNotification, error) {
	smt, err := builder.GetDB().Prepare("SELECT * FROM user_notifications WHERE id = ? LIMIT 1")
	if err != nil {
		return nil, errors.Wrap(err, "byId user prepare error")
	}
	defer smt.Close()
	res, err := smt.Query(id)
	if err != nil {
		return nil, errors.Wrap(err, "byId user error")
	}
	defer res.Close()
	if res.Next() {
		return castUserNotification(res)
	}
	return nil, errors.New("user not found")
}

func (r UserNotificationRepository) One(
	filter *[]builder.FilterCondition,
	sorts *[]builder.SortOrder,
) (*usernotification.UserNotification, error) {
	var _validFields = map[string]bool{
		"id":     true,
		"userId": true,
		"status": true,
	}

	if filter != nil && len(*filter) > 0 {
		if err := builder.ValidateFilter(*filter, _validFields); err != nil {
			return nil, err
		}
	}

	query, args := builder.BuildSQLQuery(UserNotificationTableName, filter, sorts, true)

	smt, err := builder.GetDB().Prepare(query)
	if err != nil {
		return nil, errors.Wrap(err, "get user prepare error")
	}
	defer smt.Close()
	res, err := smt.Query(args...)
	if err != nil {
		return nil, errors.Wrap(err, "get user error")
	}
	defer res.Close()
	if res.Next() {
		return castUserNotification(res)
	}
	return nil, errors.New("user not found")
}
func (r UserNotificationRepository) List(
	filter *[]builder.FilterCondition,
	sorts *[]builder.SortOrder,
) ([]*usernotification.UserNotification, error) {
	var _validFields = map[string]bool{
		"id":      true,
		"user_id": true,
		"status":  true,
	}
	var u *usernotification.UserNotification
	var res *sql.Rows
	var err error
	if filter != nil && len(*filter) > 0 {
		if err = builder.ValidateFilter(*filter, _validFields); err != nil {
			return nil, err
		}
	}
	query, args := builder.BuildSQLQuery(UserNotificationTableName, filter, sorts, false)
	userNotifications := make([]*usernotification.UserNotification, 0)
	smt, err := builder.GetDB().Prepare(query)
	if err != nil {
		return nil, errors.Wrap(err, "list user prepare error")
	}
	defer smt.Close()
	res, err = smt.Query(args...)
	if err != nil {
		return nil, errors.Wrap(err, "list user error")
	}
	defer res.Close()
	for res.Next() {
		u, err = castUserNotification(res)
		if err != nil {
			return nil, errors.Wrap(err, "list user cast error")
		}
		userNotifications = append(userNotifications, u)
	}
	return userNotifications, nil
}

func (r UserNotificationRepository) Create(params UserNotificationParams) (*usernotification.UserNotification, error) {
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
	res, err := smt.Exec(args...)
	if err != nil {
		return nil, errors.Wrap(err, "create user error")
	}
	return helper.Must(r.ByID(helper.Must(res.LastInsertId()))), nil
}

func (r UserNotificationRepository) Update(
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
		return helper.Must(r.ByID(id)), nil
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
	_, err = smt.Exec(args...)
	if err != nil {
		return nil, errors.Wrap(err, "update user error")
	}

	return helper.Must(r.ByID(id)), nil
}

func castUserNotification(res *sql.Rows) (*usernotification.UserNotification, error) {
	u := usernotification.UserNotification{}
	err := res.Scan(
		&u.ID,
		&u.UserID,
		&u.Data,
		&u.Status,
		&u.UpdatedAt,
		&u.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user notification not found")
		}
		return nil, errors.Wrap(err, "user notification get by id error")
	}

	return &u, nil
}
