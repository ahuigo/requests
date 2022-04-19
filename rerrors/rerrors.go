package rerrors

import (
	"fmt"

	"github.com/pkg/errors"
)

type ErrorType string

const (
	NetworkError ErrorType = "error-network"
	URLError     ErrorType = "error-url"
)

type Error struct {
	ErrType ErrorType
	Err     error
	Data    interface{}
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s:%+v", e.ErrType, e.Err)
}

func New(errType ErrorType, msg string) *Error {
	err := errors.New(msg)
	return &Error{
		ErrType: errType,
		Err:     err,
	}
}

func Errorf(errType ErrorType, format string, args ...interface{}) *Error {
	err := errors.Errorf(format, args...)
	return &Error{
		ErrType: errType,
		Err:     err,
	}
}

func Wrapf(errType ErrorType, err error, format string, args ...interface{}) *Error {
	err = errors.Wrapf(err, format, args...)
	return &Error{
		ErrType: errType,
		Err:     err,
	}
}

func Wrap(errType ErrorType, err error, msg string) *Error {
	err = errors.Wrap(err, msg)
	return &Error{
		ErrType: errType,
		Err:     err,
	}
}
