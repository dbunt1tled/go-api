package validation

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/dbunt1tled/go-api/pkg/e"
	"github.com/dbunt1tled/go-api/pkg/validation/validators"
	"github.com/go-playground/validator/v10"
	"github.com/uptrace/bun"
)

var customMessages = map[string]string{
	"required":  "Field %s is required",
	"email":     "Invalid email address for field %s",
	"min":       "Field %s must have a minimum length of %s characters",
	"max":       "Field %s must have a maximum length of %s characters",
	"len":       "Field %s must be exactly %s characters long",
	"number":    "Field %s must be a number",
	"positive":  "Field %s must be a positive number",
	"alphanum":  "Field %s must contain only alphanumeric characters",
	"oneof":     "Invalid value for field %s",
	"passwd":    "Field %s must contain at least 1 letter and 1 number",
	"unique_db": "%s is already taken",
	"eqfield":   "Field %s must be equal to %s",
}

type ReqValidator struct {
	validate *validator.Validate
}

func (r *ReqValidator) Validate(i any) error {
	return r.validate.Struct(i)
}

func ErrorValidation(err error) []e.ErrNo {
	var (
		validationErrors   validator.ValidationErrors
		ok                 bool
		msg, customMessage string
	)
	if errors.As(err, &validationErrors) {
		errorsMap := make([]e.ErrNo, 0, len(validationErrors))
		for _, err := range validationErrors {
			fieldName := err.StructNamespace()
			tag := err.Tag()
			customMessage, ok = customMessages[tag]
			if !ok {
				msg = defaultErrorMessage(err)
			} else {
				msg = formatErrorMessage(customMessage, err, tag)
			}

			errorsMap = append(errorsMap, e.ErrNo{
				Status: http.StatusUnprocessableEntity,
				Msg:    fmt.Sprintf("%s#%s", fieldName, msg),
				Code:   0,
			})
		}

		return errorsMap
	}
	return nil
}

func formatErrorMessage(customMessage string, err validator.FieldError, tag string) string {
	if tag == "min" || tag == "max" || tag == "len" {
		return fmt.Sprintf(customMessage, err.Field(), err.Param())
	}
	return fmt.Sprintf(customMessage, err.Field())
}

func defaultErrorMessage(err validator.FieldError) string {
	return fmt.Sprintf("Field '%s' failed on the '%s' tag", err.Field(), err.Tag())
}

func Validator(db *bun.DB) (*ReqValidator, error) {
	validate := validator.New()
	if validate == nil {
		return nil, errors.New("validator is not initialized")
	}

	if err := validate.RegisterValidation("passwd", validators.Password); err != nil {
		return nil, fmt.Errorf("failed to register password validator: %w", err)
	}

	if err := validate.RegisterValidation("regex", validators.Password); err != nil {
		return nil, fmt.Errorf("failed to register regex validator: %w", err)
	}

	if db != nil {
		uniqueDBValidator := validators.NewUniqueFieldValidator(db)
		if err := validate.RegisterValidation("unique_db", uniqueDBValidator.Validate); err != nil {
			return nil, fmt.Errorf("failed to register unique db validator: %w", err)
		}
	}

	return &ReqValidator{validate: validate}, nil
}
