package service

import "github.com/zunkk/go-sidecar/frame"

func init() {
	frame.RegisterComponents(NewUser)
}
