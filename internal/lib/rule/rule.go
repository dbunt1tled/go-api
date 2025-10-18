package rule

import (
	sql2 "database/sql"
	"fmt"
	"github.com/dbunt1tled/go-api/internal/util/builder"
	"github.com/dbunt1tled/go-api/internal/util/helper"
	"regexp"
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
)

const passwordMinLen = 7

func Regex(fl validator.FieldLevel) bool {
	// "^[a-zA-Z0-9]$"
	return regexp.MustCompile(fl.Param()).MatchString(fl.Field().String())
}

func Password(fl validator.FieldLevel) bool {
	s := fl.Field().String()
	var (
		hasMinLen  = false
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)
	if len(s) >= passwordMinLen {
		hasMinLen = true
	}
	for _, char := range s {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}
	return hasMinLen && hasUpper && hasLower && hasNumber && hasSpecial
}

func UniqueDB(fl validator.FieldLevel) bool {
	var (
		table        string
		field        string
		excludeField string
		excludeValue string
		sql          string
		res          *sql2.Rows
		smt          *sql2.Stmt
		err          error
	)

	value := fl.Field().String()
	helper.AssignFromSlice(strings.Split(fl.Param(), "#"), &table, &field, &excludeField, &excludeValue)
	if table == "" || field == "" {
		return false
	}

	if excludeField != "" {
		sql = fmt.Sprintf("SELECT `%s` FROM %s WHERE %s = ? AND %s != ? LIMIT 1", field, table, field, excludeField)
	} else {
		sql = fmt.Sprintf("SELECT `%s` FROM %s WHERE %s = ? LIMIT 1", field, table, field)
	}

	smt, err = builder.GetDB().Prepare(sql)
	if err != nil {
		return false
	}
	defer smt.Close()
	if excludeField != "" {
		res, err = smt.Query(value, value)
	} else {
		res, err = smt.Query(value)
	}

	if err != nil {
		return false
	}
	defer res.Close()

	return !res.Next()
}
