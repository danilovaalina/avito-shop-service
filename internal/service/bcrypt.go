package service

import (
	"github.com/cockroachdb/errors"
	"golang.org/x/crypto/bcrypt"
)

type Bcrypt struct {
}

func NewBcrypt() *Bcrypt {
	return &Bcrypt{}
}

func (b *Bcrypt) Encrypt(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.WithStack(err)
	}

	return string(hash), nil
}

func (b *Bcrypt) Verify(hash string, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
