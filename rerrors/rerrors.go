package rerrors

import (
	"fmt"

	// "github.com/pkg/errors"
	"errors"
)

type ErrorType string

const (
	NetworkError   ErrorType = "network-error"
	NetworkTimeout ErrorType = "network-timeout"
	URLError       ErrorType = "error-url"
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
	// err := errors.Errorf(format, args...)
	err := fmt.Errorf(format, args...)
	return &Error{
		ErrType: errType,
		Err:     err,
	}
}

func Wrapf(errType ErrorType, err error, format string, args ...interface{}) *Error {
	// err = errors.Wrapf(err, format, args...)
	s := fmt.Sprintf(format, args...)
	err = fmt.Errorf("%s|%s", err.Error(), s)
	return &Error{
		ErrType: errType,
		Err:     err,
	}
}

func Wrap(errType ErrorType, err error, msg string) *Error {
	// err = errors.Wrap(err, msg)
	err = fmt.Errorf("%s|%s", err.Error(), msg)
	return &Error{
		ErrType: errType,
		Err:     err,
	}
}