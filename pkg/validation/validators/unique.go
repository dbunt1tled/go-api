package validators

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/uptrace/bun"
)

type UniqueFieldValidator struct {
	db *bun.DB
}

func NewUniqueFieldValidator(db *bun.DB) *UniqueFieldValidator {
	return &UniqueFieldValidator{db: db}
}

func (uv *UniqueFieldValidator) Validate(fl validator.FieldLevel) bool {
	// Format: unique_db=table.column[.exclude_id]
	params := strings.Split(fl.Param(), ".")
	if len(params) < 2 {
		return false
	}

	tableName := params[0]
	columnName := params[1]

	var query strings.Builder
	query.WriteString(fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s = ?", tableName, columnName))
	args := []interface{}{fl.Field().Interface()}
	if len(params) > 2 && params[2] == "exclude_id" {
		currentStruct := fl.Parent()
		if currentStruct.Kind() == reflect.Ptr {
			currentStruct = currentStruct.Elem()
		}

		idField := currentStruct.FieldByName("ID")
		if !idField.IsValid() || idField.IsZero() {
			return true // If no ID field, can't exclude anything
		}

		query.WriteString(" AND id != $2")
		args = append(args, idField.Interface())
	}
	query.WriteString(" LIMIT 1")

	var count int
	err := uv.db.QueryRowContext(context.Background(), query.String(), args...).Scan(&count)
	if err != nil {
		return false
	}

	return count == 0
}
