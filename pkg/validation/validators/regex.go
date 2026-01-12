package validators

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

func Regex(field validator.FieldLevel) bool {
	// "^[a-zA-Z0-9]$"
	return regexp.MustCompile(field.Param()).MatchString(field.Field().String())
}
