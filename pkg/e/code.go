package e

const (
	// Err0System System error.
	Err0System = iota
)

const (
	// 401 Unauthorized errors.
	_ = 40100000 + iota
	Err401AuthEmptyTokenError
	Err401TokenNotFoundError
	Err401TokenEmptyError
	Err401TokenError
	Err401TokenSubjectError
	Err401TokenUserIdError
	Err401UserNotFoundError
	Err401UserNotActiveError
	Err401RefreshEmptyTokenError
	Err401TokenRefreshSubjectError
	Err401TokenRefreshUserIdError
	Err401TokenRefreshUserError
	Err401RefreshUserNotFoundError
	Err401RefreshUserNotActiveError
	Err401SystemEmptyTokenError
	Err401SystemTokenError
)

const (
	// 404 Not Found errors.
	_ = 40400000 + iota
	Err404NotFoundDefault
	Err404UserNotFound
	Err404URLExpired
)

const (
	// 422 Unprocessable Entity errors.
	_ = 42200000 + iota
	Err422HomeTemplateError
	Err422HomeTemplateExecuteError
	Err422LoginValidateError
	Err422LoginUserNotFoundError
	Err422LoginAccessTokenError
	Err422LoginRefreshTokenError
	Err422RegisterValidateError
	Err422RegisterUserCreationError
	Err422RegisterUserPasswordError
	Err422LoginUserPasswordError
	Err422LoginUserPasswordWrongError
	Err422CreateConfirmTokenError
	Err422SendConfirmEmailError
	Err422ConfirmValidateError
	Err422TokenError
	Err422TokenEmptyError
	Err422ConfirmTokenError
	Err422TokenConfirmUserIdError
	Err422TokenConfirmUserError
	Err422TokenConfirmUserNotFoundError
	Err422TokenConfirmUserNotPendingError
	Err422TokenConfirmUserUpdateError
	Err422TokenSubjectError
	Err422UserListValidateError
	Err422UserListError
)
