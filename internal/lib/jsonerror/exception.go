package jsonerror

type ExceptionErr struct {
	Inner  error
	Code   int
	Status int
}

func (e *ExceptionErr) Error() string { return e.Inner.Error() }
