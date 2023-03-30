package errcode

import (
	"fmt"

	"github.com/pkg/errors"
)

const (
	// UnknownErrorCode means unknown error
	UnknownErrorCode = 10001
)

var (
	ErrRequestParameter = NewCustomError(10002, "error request parameter")
	ErrAuthCode         = NewCustomError(10003, "error auth token")
)

type CustomError struct {
	msg  string
	code uint32
}

func NewCustomError(code uint32, msg string) *CustomError {
	return &CustomError{
		msg:  msg,
		code: code,
	}
}

func (e *CustomError) Error() string {
	return e.msg
}

func (e *CustomError) Wrap(errMsg string) *CustomError {
	return &CustomError{
		msg:  fmt.Sprintf("%s: %s", e.msg, errMsg),
		code: e.code,
	}
}

func DecodeError(customErr error) uint32 {
	rootErr := errors.Cause(customErr)
	if ce, ok := rootErr.(*CustomError); ok {
		return ce.code
	}

	return UnknownErrorCode
}
