package helper

import (
	"go_echo/internal/config/locale"
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/google/jsonapi"
	"github.com/labstack/echo/v4"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func isVarType(value interface{}, targetType reflect.Type) bool {
	return reflect.TypeOf(value) == targetType
}

func IsSliceVarOfType(slice interface{}, elemType reflect.Type) bool {
	t := reflect.TypeOf(slice)
	if t.Kind() != reflect.Slice {
		return false
	}
	return t.Elem() == elemType
}

func ValidationErrorString(validationErrors validator.ValidationErrors) map[string]interface{} {
	var fName string
	localizer := locale.GetLocalizerInstance()
	errors := make(map[string]interface{})
	defaultMessage := localizer.MustLocalize(&i18n.LocalizeConfig{
		MessageID: "default_validation_message",
		TemplateData: map[string]interface{}{
			"Field": "{{.Field}}",
			"Tag":   "{{.Tag}}",
		},
	})
	for _, err := range validationErrors {
		fName = localizer.MustLocalize(&i18n.LocalizeConfig{
			MessageID: err.Field(),
		})
		errors[err.Field()] = localizer.MustLocalize(&i18n.LocalizeConfig{
			TemplateData: map[string]interface{}{
				"Field": fName,
				"Tag":   err.ActualTag(),
			},
			DefaultMessage: &i18n.Message{
				ID:    err.ActualTag(),
				One:   defaultMessage,
				Other: defaultMessage,
			},
		})
	}
	return errors
}

func JSONAPIModel(r *echo.Response, models interface{}, status int) error {
	r.Header().Set(echo.HeaderContentType, jsonapi.MediaType)
	r.WriteHeader(status)
	e := jsonapi.MarshalPayload(r, models)
	if e != nil {
		panic(e.Error()) // TODO logging
	}
	return nil
}
