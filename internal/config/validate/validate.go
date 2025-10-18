package validate

import (
	"github.com/dbunt1tled/go-api/internal/lib/rule"
	"sync"

	"github.com/go-playground/validator/v10"
)

var (
	validInstance *validator.Validate //nolint:gochecknoglobals // singleton
	m             sync.Once           //nolint:gochecknoglobals // singleton
)

func GetValidateInstance() *validator.Validate {
	m.Do(func() {
		validInstance = InitValidateInstance()
	})
	return validInstance
}

func InitValidateInstance() *validator.Validate {
	v := validator.New()
	_ = v.RegisterValidation("regex", rule.Regex)        //nolint:errcheck // ignore error
	_ = v.RegisterValidation("passwd", rule.Password)    //nolint:errcheck // ignore error
	_ = v.RegisterValidation("unique_db", rule.UniqueDB) //nolint:errcheck // ignore error
	return v
}
