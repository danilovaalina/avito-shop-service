package service

import (
	"context"
	"time"

	"github.com/cockroachdb/errors"

	"avito-shop-service/internal/model"
)

const defaultTokenLifetime = 24 * time.Hour

type Repository interface {
	GetUser(ctx context.Context, username string) (model.User, error)
	CreateUser(ctx context.Context, username string, passwordHash string) (model.User, error)
	UpdateBalance(ctx context.Context, username string, itemName string) error
	SwapBalance(ctx context.Context, fromUser, toUser string, amount int) error
	Balance(ctx context.Context, username string) (int64, error)
	Inventory(ctx context.Context, username string) ([]model.Inventory, error)
	Transaction(ctx context.Context, username string) ([]model.Transaction, error)
}

type Tokenizer interface {
	CreateToken(username string) (string, error)
}

type Encryptor interface {
	Encrypt(password string) (string, error)
	Verify(hash string, password string) error
}

type Service struct {
	repository Repository
	tokenizer  Tokenizer
	encryptor  Encryptor
}

func New(repository Repository, tokenizer Tokenizer, encryptor Encryptor) *Service {
	return &Service{
		repository: repository,
		tokenizer:  tokenizer,
		encryptor:  encryptor,
	}
}

func (s *Service) Login(ctx context.Context, username, password string) (string, error) {
	user, err := s.repository.GetUser(ctx, username)
	if err != nil && !errors.Is(model.ErrUserNotFound, err) {
		return "", err
	}

	if errors.Is(err, model.ErrUserNotFound) {
		user, err = s.createUser(ctx, username, password)
		if err != nil {
			return "", err
		}
	} else {
		err = s.encryptor.Verify(user.PasswordHash, password)
		if err != nil {
			return "", errors.Join(err, model.ErrUserUnauthorized)
		}
	}

	token, err := s.tokenizer.CreateToken(username)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *Service) BuyItem(ctx context.Context, username, itemName string) error {
	err := s.repository.UpdateBalance(ctx, username, itemName)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) SendCoin(ctx context.Context, fromUser, toUser string, amount int) error {
	err := s.repository.SwapBalance(ctx, fromUser, toUser, amount)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) Balance(ctx context.Context, username string) (int64, error) {
	balance, err := s.repository.Balance(ctx, username)
	if err != nil {
		return balance, err
	}

	return balance, nil
}

func (s *Service) Inventory(ctx context.Context, username string) ([]model.Inventory, error) {
	inventory, err := s.repository.Inventory(ctx, username)
	if err != nil {
		return nil, err
	}

	return inventory, nil
}

func (s *Service) Transaction(ctx context.Context, username string) ([]model.Transaction, error) {
	transactions, err := s.repository.Transaction(ctx, username)
	if err != nil {
		return nil, err
	}

	return transactions, nil
}

func (s *Service) createUser(ctx context.Context, username string, password string) (model.User, error) {
	hash, err := s.encryptor.Encrypt(password)
	if err != nil {
		return model.User{}, err
	}

	user, err := s.repository.CreateUser(ctx, username, hash)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}
