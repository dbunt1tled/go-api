package jsonerror

import (
	"errors"
	"fmt"
	"go_echo/internal/config/env"
	error2 "go_echo/internal/config/error"
	"go_echo/internal/util/helper"
	"net/http"
	"runtime"

	"github.com/labstack/echo/v4"
)

type APIError struct {
	ID     int         `json:"id,omitempty" jsonapi:"primary,error"`
	Status int         `json:"status,omitempty" jsonapi:"attr,status"`
	Errors []ErrorData `json:"errors,omitempty" jsonapi:"attr,errors"`
	Source Source      `json:"source,omitempty" jsonapi:"attr,source"`
}

type Source struct {
	Pointer string `json:"pointer,omitempty"  jsonapi:"attr,pointer"`
}

type ErrorData struct {
	Title  string `json:"title,omitempty" jsonapi:"attr,title"`
	Detail string `json:"detail,omitempty" jsonapi:"attr,detail"`
}

func NewError(err error, code int, status int) *APIError {
	var stack *string
	var e error2.StackTracer
	stack = nil
	if errors.As(err, &e) {
		st := ""
		for _, f := range e.StackTrace() {
			st += fmt.Sprintf("%+s:%d\n", f, f)
		}
		if st != "" {
			stack = &st
		}
	}
	return NewErrorString(err.Error(), code, status, stack)
}

func NewErrorString(err string, code int, status int, stack *string) *APIError {
	return NewErrorMap(map[string]interface{}{"api": err}, code, status, stack)
}

func NewErrorMap(err map[string]interface{}, code int, status int, stack *string) *APIError {
	cfg := env.GetConfigInstance()
	pc, _, line, _ := runtime.Caller(2)
	details := runtime.FuncForPC(pc)
	var source Source
	if (cfg.Debug) && details != nil {
		if stack != nil {
			source = Source{Pointer: *stack}
		} else {
			source = Source{Pointer: fmt.Sprintf("%s#%d", details.Name(), line)}
		}
	}
	data := make([]ErrorData, 0, len(err))
	for key, value := range err {
		data = append(data, ErrorData{Title: key, Detail: fmt.Sprintf("%v", value)})
	}
	return &APIError{
		ID:     code,
		Status: status,
		Errors: data,
		Source: source,
	}
}

func ErrorNotFound(c echo.Context, err error, code int) error {
	return errorError(c, err, code, http.StatusNotFound)
}
func ErrorInternal(c echo.Context, err error, code int) error {
	return errorError(c, err, code, http.StatusInternalServerError)
}
func ErrorUnauthorized(c echo.Context, err error, code int) error {
	return errorError(c, err, code, http.StatusUnauthorized)
}
func ErrorUnprocessableEntity(c echo.Context, err error, code int) error {
	return errorError(c, err, code, http.StatusUnprocessableEntity)
}

func ErrorNotFoundMap(c echo.Context, err map[string]interface{}, code int) error {
	return errorMap(c, err, code, http.StatusNotFound)
}
func ErrorInternalMap(c echo.Context, err map[string]interface{}, code int) error {
	return errorMap(c, err, code, http.StatusInternalServerError)
}
func ErrorUnauthorizedMap(c echo.Context, err map[string]interface{}, code int) error {
	return errorMap(c, err, code, http.StatusUnauthorized)
}
func ErrorUnprocessableEntityMap(c echo.Context, err map[string]interface{}, code int) error {
	return errorMap(c, err, code, http.StatusUnprocessableEntity)
}
func ErrorNotFoundString(c echo.Context, err string, code int) error {
	return errorString(c, err, code, http.StatusNotFound)
}
func ErrorInternalString(c echo.Context, err string, code int) error {
	return errorString(c, err, code, http.StatusInternalServerError)
}
func ErrorUnauthorizedString(c echo.Context, err string, code int) error {
	return errorString(c, err, code, http.StatusUnauthorized)
}
func ErrorUnprocessableEntityString(c echo.Context, err string, code int) error {
	return errorString(c, err, code, http.StatusUnprocessableEntity)
}

func errorString(c echo.Context, err string, code int, status int) error {
	return helper.JSONAPIModel(c.Response(), NewErrorString(err, code, status, nil), status)
}

func errorMap(c echo.Context, err map[string]interface{}, code int, status int) error {
	return helper.JSONAPIModel(c.Response(), NewErrorMap(err, code, status, nil), status)
}

func errorError(c echo.Context, err error, code int, status int) error {
	return helper.JSONAPIModel(c.Response(), NewError(err, code, status), status)
}
