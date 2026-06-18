package apperror

import (
	"errors"
	"fmt"
)

type Code string

const (
	CodeInvalidArgument Code = "invalid_argument"
	CodeNotFound        Code = "not_found"
	CodeUnavailable     Code = "unavailable"
	CodeInternal        Code = "internal"
	CodeDublicate       Code = "dublicate"
)

type Error struct {
	Code Code
	Msg  string
	Err  error
}

func (err *Error) Error() string {
	if err.Err == nil {
		return err.Msg
	}

	return fmt.Sprintf("%s: %v", err.Msg, err.Err)
}

func New(code Code, msg string) error {
	return &Error{
		Code: code,
		Msg:  msg,
	}
}

func CodeOf(err error) Code {
	var apperr *Error

	if errors.As(err, &apperr) {
		return apperr.Code
	}

	return CodeInternal
}

func WrapCode(code Code, msg string, err error) error {
	return &Error{
		Code: code,
		Msg:  msg,
		Err:  err,
	}
}

func Wrap(msg string, err error) error {
	return WrapCode(CodeOf(err), msg, err)
}

func (err *Error) Unwrap() error {
	return err.Err
}
