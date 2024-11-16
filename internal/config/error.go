package config

import (
	"github.com/pkg/errors"
)

type StackTracer interface {
	StackTrace() errors.StackTrace
	Error() string
}
