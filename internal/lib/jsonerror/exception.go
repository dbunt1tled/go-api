package jsonerror

type ErrException struct {
	Inner  error
	Code   int
	Status int
}

func (e *ErrException) Error() string { return e.Inner.Error() }
