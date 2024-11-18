package helper

import (
	"fmt"
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/google/jsonapi"
	"github.com/labstack/echo/v4"
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

func ValidationErrorString(validationErrors validator.ValidationErrors) []string {
	errors := make([]string, 0)
	for _, err := range validationErrors {
		errors = append(errors, fmt.Sprintf("Field %s: %s - %s", err.Field(), err.ActualTag(), err.Tag()))
		// errors = append(errors, ?err.Translate(ut))
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
