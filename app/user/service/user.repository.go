package service

import (
	"database/sql"
	"go_echo/app/user/model/user"
	"go_echo/internal/util/builder"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

const (
	UserTableName = "users"
)

type UserRepository struct {
}

func (r UserRepository) ByID(id int) (*user.User, error) {
	smt, err := builder.GetDB().Prepare("SELECT * FROM users WHERE id = ? LIMIT 1")
	if err != nil {
		return nil, errors.Wrap(err, "byId user prepare error")
	}
	res, err := smt.Query(id)
	if err != nil {
		return nil, errors.Wrap(err, "byId user error")
	}
	defer res.Close()
	if res.Next() {
		return castUser(res)
	}
	return nil, errors.New("user not found")
}

func (r UserRepository) One(filter []builder.FilterCondition, sorts []builder.SortOrder) (*user.User, error) {
	var _validFields = map[string]bool{
		"id":     true,
		"phone":  true,
		"email":  true,
		"status": true,
	}

	if err := builder.ValidateFilter(filter, _validFields); err != nil {
		return nil, err
	}

	query, args := builder.BuildSQLQuery(UserTableName, filter, sorts, true)

	smt, err := builder.GetDB().Prepare(query)
	if err != nil {
		return nil, errors.Wrap(err, "get user prepare error")
	}
	res, err := smt.Query(args...)
	if err != nil {
		return nil, errors.Wrap(err, "get user error")
	}
	defer res.Close()
	if res.Next() {
		return castUser(res)
	}
	return nil, errors.New("user not found")
}
func (r UserRepository) List(filter []builder.FilterCondition, sorts []builder.SortOrder) (*[]user.User, error) {
	var _validFields = map[string]bool{
		"id":     true,
		"phone":  true,
		"email":  true,
		"status": true,
	}
	var u *user.User
	var res *sql.Rows
	var err error
	if err := builder.ValidateFilter(filter, _validFields); err != nil {
		return nil, err
	}

	query, args := builder.BuildSQLQuery(UserTableName, filter, sorts, false)
	users := make([]user.User, 0)
	smt, err := builder.GetDB().Prepare(query)
	if err != nil {
		return nil, errors.Wrap(err, "list user prepare error")
	}
	res, err = smt.Query(args...)
	if err != nil {
		return nil, errors.Wrap(err, "list user error")
	}
	defer res.Close()
	for res.Next() {
		u, err = castUser(res)
		if err != nil {
			return nil, errors.Wrap(err, "list user cast error")
		}
		users = append(users, *u)
	}
	return &users, nil
}

func castUser(res *sql.Rows) (*user.User, error) {
	u := user.User{}
	err := res.Scan(
		&u.ID,
		&u.FirstName,
		&u.SecondName,
		&u.Email,
		&u.PhoneNumber,
		&u.Password,
		&u.Status,
		&u.ConfirmedAt,
		&u.UpdatedAt,
		&u.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, errors.Wrap(err, "user get by id error")
	}

	return &u, nil
}
