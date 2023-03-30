package entity

import (
	"github.com/zunkk/go-project-startup/pkg/auth/jwt"
)

type CustomClaims struct {
	jwt.BaseClaims
}
