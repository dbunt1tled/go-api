package validate

import (
	"go_echo/internal/lib/rule"

	"github.com/go-playground/validator/v10"
)

var validInstance *validator.Validate //nolint:gochecknoglobals // singleton

func GetValidateInstance() *validator.Validate {
	if validInstance == nil {
		validInstance = InitValidateInstance()
	}
	return validInstance
}

func InitValidateInstance() *validator.Validate {
	v := validator.New()
	v.RegisterValidation("regex", rule.Regex)     //nolint:errcheck // ignore error
	v.RegisterValidation("passwd", rule.Password) //nolint:errcheck // ignore error
	return v
}
