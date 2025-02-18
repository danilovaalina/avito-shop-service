package service

import (
	"time"

	"github.com/cockroachdb/errors"
	"github.com/golang-jwt/jwt/v5"
)

const defaultTokenLifetime = 24 * time.Hour

type JWT struct {
	lifetime time.Duration
}

type JWTOption func(*JWT)

func WithLifetime(lifetime time.Duration) JWTOption {
	return func(jwt *JWT) {
		if lifetime > 0 {
			jwt.lifetime = lifetime
		}
	}
}

func NewJWT(opts ...JWTOption) *JWT {
	j := &JWT{
		lifetime: defaultTokenLifetime,
	}

	for _, opt := range opts {
		opt(j)
	}

	return j
}

func (j *JWT) CreateToken(username string) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(j.lifetime).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secret := []byte("secret")

	s, err := token.SignedString(secret)
	if err != nil {
		return "", errors.WithStack(err)
	}

	return s, nil
}
