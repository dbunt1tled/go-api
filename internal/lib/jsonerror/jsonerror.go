package jsonerror

import (
	"fmt"
	"net/http"
	"runtime"

	"github.com/google/jsonapi"
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
	return NewErrorString(err.Error(), code, status)
}

func NewErrorString(err string, code int, status int) *APIError {
	return NewErrorMap(map[string]interface{}{"api": err}, code, status)
}

func NewErrorMap(err map[string]interface{}, code int, status int) *APIError {
	pc, _, line, _ := runtime.Caller(2)
	details := runtime.FuncForPC(pc)
	source := Source{Pointer: fmt.Sprintf("%s#%d", details.Name(), line)}
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

func ErrorNotFound(c echo.Context, err error, code int) {
	e := jsonapi.MarshalPayload(c.Response(), NewError(err, code, http.StatusNotFound))
	if e != nil {
		c.JSON(http.StatusNotFound, e.Error())
		return
	}
	c.JSON(http.StatusNotFound, c.Response())
}
func ErrorInternal(c echo.Context, err error, code int) {
	e := jsonapi.MarshalPayload(c.Response(), NewError(err, code, http.StatusInternalServerError))
	if e != nil {
		http.Error(c.Response(), e.Error(), http.StatusInternalServerError)
	}
}
func ErrorUnauthorized(c echo.Context, err error, code int) {
	e := jsonapi.MarshalPayload(c.Response(), NewError(err, code, http.StatusUnauthorized))
	if e != nil {
		c.JSON(http.StatusUnauthorized, e.Error())
		return
	}
	c.JSON(http.StatusUnauthorized, c.Response())
}
func ErrorUnprocessableEntity(c echo.Context, err error, code int) {
	e := jsonapi.MarshalPayload(c.Response(), NewError(err, code, http.StatusUnprocessableEntity))
	if e != nil {
		c.JSON(http.StatusUnprocessableEntity, e.Error())
		return
	}
	c.JSON(http.StatusUnprocessableEntity, c.Response())
}

func ErrorNotFoundMap(c echo.Context, err map[string]interface{}, code int) {
	e := jsonapi.MarshalPayload(c.Response(), NewErrorMap(err, code, http.StatusNotFound))
	if e != nil {
		c.JSON(http.StatusNotFound, e.Error())
		return
	}
	c.JSON(http.StatusNotFound, c.Response())
}
func ErrorInternalMap(c echo.Context, err map[string]interface{}, code int) {
	e := jsonapi.MarshalPayload(c.Response(), NewErrorMap(err, code, http.StatusInternalServerError))
	if e != nil {
		c.JSON(http.StatusInternalServerError, e.Error())
		return
	}
	c.JSON(http.StatusInternalServerError, c.Response())
}
func ErrorUnauthorizedMap(c echo.Context, err map[string]interface{}, code int) {
	e := jsonapi.MarshalPayload(c.Response(), NewErrorMap(err, code, http.StatusUnauthorized))
	if e != nil {
		c.JSON(http.StatusUnauthorized, e.Error())
		return
	}
	c.JSON(http.StatusUnauthorized, c.Response())
}
func ErrorUnprocessableEntityMap(c echo.Context, err map[string]interface{}, code int) {
	e := jsonapi.MarshalPayload(c.Response(), NewErrorMap(err, code, http.StatusUnprocessableEntity))
	if e != nil {
		c.JSON(http.StatusUnprocessableEntity, e.Error())
		return
	}
	c.JSON(http.StatusUnprocessableEntity, c.Response())
}

func ErrorNotFoundString(c echo.Context, err string, code int) {
	e := jsonapi.MarshalPayload(c.Response(), NewErrorString(err, code, http.StatusNotFound))
	if e != nil {
		c.JSON(http.StatusUnprocessableEntity, e.Error())
		return
	}
	c.JSON(http.StatusUnprocessableEntity, c.Response())
}
func ErrorInternalString(c echo.Context, err string, code int) {
	e := jsonapi.MarshalPayload(c.Response(), NewErrorString(err, code, http.StatusInternalServerError))
	if e != nil {
		c.JSON(http.StatusInternalServerError, e.Error())
		return
	}
	c.JSON(http.StatusInternalServerError, c.Response())
}
func ErrorUnauthorizedString(c echo.Context, err string, code int) {
	e := jsonapi.MarshalPayload(c.Response(), NewErrorString(err, code, http.StatusUnauthorized))
	if e != nil {
		c.JSON(http.StatusUnauthorized, e.Error())
		return
	}
	c.JSON(http.StatusUnauthorized, c.Response())
}
func ErrorUnprocessableEntityString(c echo.Context, err string, code int) {
	e := jsonapi.MarshalPayload(c.Response(), NewErrorString(err, code, http.StatusUnprocessableEntity))
	if e != nil {
		c.JSON(http.StatusUnprocessableEntity, e.Error())
		return
	}
	c.JSON(http.StatusUnprocessableEntity, c.Response())
}
