package entity

import (
	"github.com/zunkk/go-sidecar/auth/jwt"
)

type CustomClaims struct {
	jwt.BaseClaims
}
