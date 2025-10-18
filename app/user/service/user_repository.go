package service

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/dbunt1tled/go-api/app/user/model/user"
	"github.com/dbunt1tled/go-api/internal/util/builder"
	"github.com/dbunt1tled/go-api/internal/util/builder/page"
	"github.com/dbunt1tled/go-api/internal/util/hasher"
	"github.com/dbunt1tled/go-api/internal/util/helper"
	"github.com/dbunt1tled/go-api/internal/util/type/roles"
	"strings"
	"time"

	"github.com/pkg/errors"
)

const (
	UserTableName = "users"
)

type UserRepository struct {
}

type UpdateUserParams struct {
	FirstName   *string
	SecondName  *string
	Email       *string
	PhoneNumber *string
	Password    *string
	Status      *user.Status
	Hash        *string
	Roles       *roles.Roles
	ConfirmedAt *time.Time
}

type CreateUserParams struct {
	FirstName   *string
	SecondName  *string
	Email       *string
	PhoneNumber *string
	Password    *string
	Status      *user.Status
	Hash        *string
	Roles       *roles.Roles
	ConfirmedAt *time.Time
}

func (r UserRepository) ByID(ctx context.Context, id int64) (*user.User, error) {
	return builder.ByID[user.User](
		ctx,
		UserTableName,
		id,
		castUserRow,
	)
}
func (r UserRepository) Paginator(
	ctx context.Context,
	filter *[]page.FilterCondition,
	sorts *[]page.SortOrder,
	paginator *page.Pagination,
) (page.Paginate[user.User], error) {
	return builder.Paginator[user.User](
		ctx,
		UserTableName,
		filter,
		sorts,
		paginator,
		castUserRows,
	)
}

func (r UserRepository) One(ctx context.Context, filter *[]page.FilterCondition, sorts *[]page.SortOrder) (*user.User, error) {
	var _validFields = map[string]bool{
		"id":     true,
		"phone":  true,
		"email":  true,
		"status": true,
	}

	if filter != nil && len(*filter) > 0 {
		if err := builder.ValidateFilter(*filter, _validFields); err != nil {
			return nil, err
		}
	}

	return builder.One(ctx, UserTableName, filter, sorts, castUserRow)
}
func (r UserRepository) List(
	ctx context.Context,
	filter *[]page.FilterCondition,
	sorts *[]page.SortOrder,
) ([]*user.User, error) {
	var _validFields = map[string]bool{
		"id":     true,
		"phone":  true,
		"email":  true,
		"status": true,
	}
	var err error
	if filter != nil && len(*filter) > 0 {
		if err = builder.ValidateFilter(*filter, _validFields); err != nil {
			return nil, err
		}
	}

	return builder.List(ctx, UserTableName, filter, sorts, castUserRows, nil)
}

func (r UserRepository) ByIdentity(ctx context.Context, login string) (*user.User, error) {
	res := builder.GetDB().QueryRowContext(
		ctx,
		"SELECT * FROM users WHERE (phone_number = ? OR email = ?) AND status = 1 LIMIT 1;",
		login,
		login,
	)
	return castUserRow(res)
}

//nolint:funlen
func (r UserRepository) Create(ctx context.Context, params CreateUserParams) (*user.User, error) {
	var (
		columns []string
		values  []string
		args    []interface{}
	)

	if params.FirstName != nil {
		columns = append(columns, "first_name")
		values = append(values, "?")
		args = append(args, *params.FirstName)
	}

	if params.SecondName != nil {
		columns = append(columns, "second_name")
		values = append(values, "?")
		args = append(args, *params.SecondName)
	}

	if params.Email != nil {
		columns = append(columns, "email")
		values = append(values, "?")
		args = append(args, *params.Email)
	}

	if params.PhoneNumber != nil {
		columns = append(columns, "phone_number")
		values = append(values, "?")
		args = append(args, *params.PhoneNumber)
	}

	if params.Password != nil {
		columns = append(columns, "password")
		values = append(values, "?")
		args = append(args, helper.Must(hasher.HashArgon(*params.Password)))
	}

	if params.Status != nil {
		columns = append(columns, "status")
		values = append(values, "?")
		args = append(args, *params.Status)
	}

	if params.Hash != nil {
		columns = append(columns, "hash")
		values = append(values, "?")
		args = append(args, *params.Hash)
	}

	if params.Roles != nil {
		columns = append(columns, "roles")
		values = append(values, "?")
		args = append(args, *params.Roles)
	}

	if params.ConfirmedAt != nil {
		columns = append(columns, "confirmed_at")
		values = append(values, "?")
		args = append(args, *params.ConfirmedAt)
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
			"INSERT INTO users (%s) VALUES (%s)",
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

func (r UserRepository) Update(ctx context.Context, id int64, params UpdateUserParams) (*user.User, error) {

	var (
		setClauses []string
		args       []interface{}
	)

	if params.FirstName != nil {
		setClauses = append(setClauses, "first_name = ?")
		args = append(args, *params.FirstName)
	}

	if params.SecondName != nil {
		setClauses = append(setClauses, "second_name = ?")
		args = append(args, *params.SecondName)
	}

	if params.Email != nil {
		setClauses = append(setClauses, "email = ?")
		args = append(args, *params.Email)
	}

	if params.PhoneNumber != nil {
		setClauses = append(setClauses, "phone_number = ?")
		args = append(args, *params.PhoneNumber)
	}

	if params.Password != nil {
		setClauses = append(setClauses, "password = ?")
		args = append(args, helper.Must(hasher.HashArgon(*params.Password)))
	}

	if params.Status != nil {
		setClauses = append(setClauses, "status = ?")
		args = append(args, *params.Status)
	}

	if params.Hash != nil {
		setClauses = append(setClauses, "hash = ?")
		args = append(args, *params.Hash)
	}

	if params.Roles != nil {
		setClauses = append(setClauses, "roles = ?")
		args = append(args, *params.Roles)
	}

	if params.ConfirmedAt != nil {
		setClauses = append(setClauses, "confirmed_at = ?")
		args = append(args, *params.ConfirmedAt)
	}

	if len(setClauses) == 0 {
		return helper.Must(r.ByID(ctx, id)), nil
	}

	setClauses = append(setClauses, "updated_at = ?")
	args = append(args, time.Now())

	args = append(args, id)
	query := fmt.Sprintf("UPDATE users SET %s WHERE id = ?", strings.Join(setClauses, ", "))

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

func castUserRow(row *sql.Row) (*user.User, error) {
	u := user.User{}
	return builder.ScanStructRow(u, row)
}

func castUserRows(rows *sql.Rows) (*user.User, error) {
	u := user.User{}
	return builder.ScanStructRows(u, rows)
}
