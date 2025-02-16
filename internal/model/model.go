package model

import "github.com/cockroachdb/errors"

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrUserUnauthorized = errors.New("user is not authorized")
	ErrItemNotFound     = errors.New("item not found")
	ErrNegativeBalance  = errors.New("negative balance")
)

type User struct {
	Username     string
	PasswordHash string
}

type Item struct {
	Name  string
	Price int64
}

type Balance struct {
	Username string
	Amount   int64
}

type Inventory struct {
	Username string
	ItemName string
	Quantity int64
}

type Transaction struct {
	FromUser string
	ToUser   string
	Amount   int64
}
