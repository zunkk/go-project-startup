package errcode

import "github.com/zunkk/go-sidecar/errcode"

var (
	ErrRequestParameter = errcode.NewCustomError(10002, "error request parameter")
	ErrAuthCode         = errcode.NewCustomError(10003, "error auth token")
)
