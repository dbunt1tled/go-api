package app_error

import (
	"github.com/pkg/errors"
)

const (
	Err422LoginValidateMapError  = 42200001
	Err422LoginValidateRuleError = 42200002
	Err422LoginValidateError     = 42200003
	Err422LoginUserNotFoundError = 42200004
	Err422LoginAuthTokensError   = 42200005

	Err422SignupValidateMapError  = 42200006
	Err422SignupValidateRuleError = 42200007
	Err422SignupValidateError     = 42200008
	Err422SignupUserNotFoundError = 42200009
	Err422SignupAuthTokensError   = 42200010

	Err401AuthEmptyTokenError   = 40100001
	Err401TokenError            = 40100002
	Err401UserNotFoundError     = 40100003
	Err401UserNotActiveError    = 40100004
	Err401SystemEmptyTokenError = 40100005
	Err401SystemTokenError      = 40100005
)

type StackTracer interface {
	StackTrace() errors.StackTrace
	Error() string
}
