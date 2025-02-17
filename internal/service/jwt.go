package service

import (
	"time"

	"github.com/cockroachdb/errors"
	"github.com/golang-jwt/jwt/v5"
)

type JWT struct{}

func NewJWT() *JWT {
	return &JWT{}
}

func (j *JWT) CreateToken(username string) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(defaultTokenLifetime).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secret := []byte("secret")

	s, err := token.SignedString(secret)
	if err != nil {
		return "", errors.WithStack(err)
	}

	return s, nil
}
