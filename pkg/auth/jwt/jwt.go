package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims interface {
	jwt.Claims
	Init(registeredClaims jwt.RegisteredClaims)
}

type BaseClaims struct {
	jwt.RegisteredClaims
}

func (c *BaseClaims) Init(registeredClaims jwt.RegisteredClaims) {
	c.RegisteredClaims = registeredClaims
}

func GenerateWithHMACKey(hmacKey string, validDuration time.Duration, id string, customClaims CustomClaims) (token string, expiredDate int64, err error) {
	notBefore := time.Now()
	expiresAt := notBefore.Add(validDuration)
	customClaims.Init(jwt.RegisteredClaims{
		Subject:   id,
		ExpiresAt: jwt.NewNumericDate(expiresAt),
		NotBefore: jwt.NewNumericDate(notBefore),
	})
	tk := jwt.NewWithClaims(jwt.SigningMethodHS256, customClaims)

	tokenString, err := tk.SignedString([]byte(hmacKey))
	if err != nil {
		return "", 0, err
	}
	return tokenString, expiresAt.Unix(), nil
}

func ParseWithHMACKey(hmacKey string, token string, res CustomClaims) (id string, err error) {
	_, err = jwt.ParseWithClaims(token, res, func(token *jwt.Token) (interface{}, error) {
		return []byte(hmacKey), nil
	})
	if err != nil {
		return "", err
	}
	id, _ = res.GetSubject()
	return id, nil
}
