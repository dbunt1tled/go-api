package e

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

type ErrNo struct {
	Msg    string  `json:"message"`
	Code   int     `json:"code"`
	Status int     `json:"status"`
	Stack  *string `json:"stack,omitempty"`
	stack  errors.StackTrace
}

type APIError interface {
	Error() string
}

type StackTracer interface {
	StackTrace() errors.StackTrace
	Error() string
}

func NewErrNo(msg string, code int, status int) error {
	base := &ErrNo{
		Msg:    msg,
		Code:   code,
		Status: status,
	}

	return errors.WithStack(base)
}

func NewErrNoWithStack(msg string, code int, status int, stack *string) *ErrNo {
	return &ErrNo{
		Msg:    msg,
		Code:   code,
		Status: status,
		Stack:  stack,
	}
}

func NewNotFoundError(msg string, code int) APIError {
	return NewErrNo(msg, code, http.StatusNotFound)
}

func NewBadRequestError(msg string, code int) APIError {
	return NewErrNo(msg, code, http.StatusBadRequest)
}
func NewBadRequestErrorWrap(msg string, code int, e error) APIError {
	var message string
	if msg == "" {
		message = e.Error()
	} else {
		message = fmt.Sprintf("%s: %s", msg, e.Error())
	}
	return NewErrNo(message, code, http.StatusUnprocessableEntity)
}

func NewInternalError(msg string, code int) APIError {
	return NewErrNo(msg, code, http.StatusInternalServerError)
}

func NewForbiddenError(msg string, code int) APIError {
	return NewErrNo(msg, code, http.StatusForbidden)
}

func NewUnauthorizedError(msg string, code int) APIError {
	return NewErrNo(msg, code, http.StatusUnauthorized)
}

func NewValidationError(msg string, code int) APIError {
	return NewErrNo(msg, code, http.StatusUnprocessableEntity)
}

func NewUnprocessableEntityError(msg string, code int) APIError {
	return NewErrNo(msg, code, http.StatusUnprocessableEntity)
}

func NewUnprocessableEntityErrorWrap(msg string, code int, e error) APIError {
	var message string
	if msg == "" {
		message = e.Error()
	} else {
		message = fmt.Sprintf("%s: %s", msg, e.Error())
	}
	return NewErrNo(message, code, http.StatusUnprocessableEntity)
}

func (e ErrNo) Error() string {
	return e.Msg
}

func GetErrTrace(err error) *string {
	var er StackTracer
	if errors.As(err, &er) {
		stack := ""
		for _, f := range er.StackTrace() {
			stack += fmt.Sprintf("%+s:%d\n", f, f)
		}
		return &stack
	}
	return nil
}