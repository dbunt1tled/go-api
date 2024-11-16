package helper

import (
	"fmt"
	"reflect"

	"github.com/go-playground/validator/v10"
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
