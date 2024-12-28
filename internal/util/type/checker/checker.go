package checker

import (
	"go_echo/internal/util/builder/page"
	"reflect"
)

func IsVarType(value interface{}, targetType reflect.Type) bool {
	return reflect.TypeOf(value) == targetType
}

func IsSliceVarOfType(slice interface{}, elemType reflect.Type) bool {
	t := reflect.TypeOf(slice)
	if t.Kind() != reflect.Slice {
		return false
	}
	return t.Elem() == elemType
}

func VarToPaginate(v any) *page.PaginateInterface {
	p, ok := v.(page.PaginateInterface)
	if !ok {
		return nil
	}
	return &p
}
