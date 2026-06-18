package apperror

import (
	"errors"
	"fmt"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func ToGRPCCode(err error) codes.Code {
	switch CodeOf(err) {
	case CodeInternal:
		return codes.Internal
	case CodeInvalidArgument:
		return codes.InvalidArgument
	case CodeDublicate:
		return codes.AlreadyExists
	case CodeNotFound:
		return codes.NotFound
	case CodeUnavailable:
		return codes.Unavailable
	}

	panic("implement me")
}

func FromGRPC(err error, msg string) error {
	switch status.Code(err) {
	case codes.NotFound:
		return WrapCode(CodeNotFound, msg, err)
	case codes.InvalidArgument:
		return WrapCode(CodeInvalidArgument, msg, err)
	case codes.Unavailable:
		return WrapCode(CodeUnavailable, msg, err)
	case codes.AlreadyExists:
		return WrapCode(CodeDublicate, msg, err)
	case codes.Internal:
		return WrapCode(CodeInternal, msg, err)
	default:
		log.Default().Print("unknown error code", "error", err.Error())
		return status.Error(codes.Unknown, err.Error())
	}
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
