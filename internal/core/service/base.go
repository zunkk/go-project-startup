package service

import "github.com/zunkk/go-project-startup/pkg/frame"

func init() {
	frame.RegisterComponents(NewUser)
}
