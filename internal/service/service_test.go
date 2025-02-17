package service_test

import (
	"context"
	"os"
	"testing"

	"github.com/cockroachdb/errors"
	"github.com/stretchr/testify/require"

	"avito-shop-service/internal/logutil"
	mockservice "avito-shop-service/internal/mocks/service"
	"avito-shop-service/internal/model"
	"avito-shop-service/internal/service"
)

func TestMain(m *testing.M) {
	logutil.Setup()
	os.Exit(m.Run())
}

func TestService_LoginAgain(t *testing.T) {
	r := mockservice.NewRepository(t)

	r.EXPECT().GetUser(context.Background(), "user").Return(model.User{
		Username:     "user",
		PasswordHash: "passwordHash",
	}, nil).Once()

	enc := mockservice.NewEncryptor(t)

	enc.EXPECT().Verify("passwordHash", "password").Return(nil).Once()

	tkz := mockservice.NewTokenizer(t)

	tkz.EXPECT().CreateToken("user").Return("token", nil).Once()

	s := service.New(r, tkz, enc)

	token, err := s.Login(context.Background(), "user", "password")

	require.NoError(t, err)
	require.Equal(t, "token", token)
}

func TestService_LoginFirstTime(t *testing.T) {
	r := mockservice.NewRepository(t)

	r.EXPECT().GetUser(context.Background(), "user").Return(model.User{}, model.ErrUserNotFound).Once()

	enc := mockservice.NewEncryptor(t)

	enc.EXPECT().Encrypt("password").Return("passwordHash", nil).Once()

	r.EXPECT().CreateUser(context.Background(), "user", "passwordHash").Return(model.User{
		Username:     "user",
		PasswordHash: "passwordHash",
	}, nil).Once()

	tkz := mockservice.NewTokenizer(t)

	tkz.EXPECT().CreateToken("user").Return("token", nil).Once()

	s := service.New(r, tkz, enc)

	token, err := s.Login(context.Background(), "user", "password")

	require.NoError(t, err)
	require.Equal(t, "token", token)
}

func TestService_LoginGetUserErr(t *testing.T) {
	r := mockservice.NewRepository(t)

	r.EXPECT().GetUser(context.Background(), "user").Return(model.User{}, errors.New("test")).Once()

	enc := mockservice.NewEncryptor(t)

	tkz := mockservice.NewTokenizer(t)

	s := service.New(r, tkz, enc)

	token, err := s.Login(context.Background(), "user", "password")

	require.Error(t, err)
	require.Empty(t, token)
}

func TestService_BuyItem(t *testing.T) {
	r := mockservice.NewRepository(t)

	r.EXPECT().UpdateBalance(context.Background(), "user", "item").Return(nil).Once()

	enc := mockservice.NewEncryptor(t)
	tkz := mockservice.NewTokenizer(t)
	s := service.New(r, tkz, enc)

	err := s.BuyItem(context.Background(), "user", "item")

	require.NoError(t, err)
}

func TestService_SendCoin(t *testing.T) {
	r := mockservice.NewRepository(t)

	r.EXPECT().SwapBalance(context.Background(), "sender", "receiver", 500).Return(nil).Once()

	enc := mockservice.NewEncryptor(t)
	tkz := mockservice.NewTokenizer(t)
	s := service.New(r, tkz, enc)

	err := s.SendCoin(context.Background(), "sender", "receiver", 500)

	require.NoError(t, err)
}

func TestService_Balance(t *testing.T) {
	r := mockservice.NewRepository(t)

	r.EXPECT().Balance(context.Background(), "user").Return(1000, nil).Once()

	enc := mockservice.NewEncryptor(t)
	tkz := mockservice.NewTokenizer(t)
	s := service.New(r, tkz, enc)

	balance, err := s.Balance(context.Background(), "user")
	require.NoError(t, err)
	require.Equal(t, int64(1000), balance)
}

func TestService_Inventory(t *testing.T) {
	r := mockservice.NewRepository(t)

	r.EXPECT().Inventory(context.Background(), "user").Return([]model.Inventory{
		{
			Username: "user",
			ItemName: "book",
			Quantity: 1,
		},
	}, nil).Once()

	enc := mockservice.NewEncryptor(t)
	tkz := mockservice.NewTokenizer(t)
	s := service.New(r, tkz, enc)

	inventory, err := s.Inventory(context.Background(), "user")
	require.NoError(t, err)

	require.Equal(t, 1, len(inventory))
	require.Equal(t, "user", inventory[0].Username)
	require.Equal(t, "book", inventory[0].ItemName)
	require.Equal(t, int64(1), inventory[0].Quantity)
}

func TestService_Transaction(t *testing.T) {
	r := mockservice.NewRepository(t)

	r.EXPECT().Transaction(context.Background(), "sender").Return([]model.Transaction{
		{
			FromUser: "sender",
			ToUser:   "receiver",
			Amount:   500,
		},
	}, nil).Once()

	enc := mockservice.NewEncryptor(t)
	tkz := mockservice.NewTokenizer(t)
	s := service.New(r, tkz, enc)

	transactions, err := s.Transaction(context.Background(), "sender")
	if err != nil {
		return
	}

	require.Equal(t, 1, len(transactions))
	require.Equal(t, "sender", transactions[0].FromUser)
	require.Equal(t, "receiver", transactions[0].ToUser)
	require.Equal(t, int64(500), transactions[0].Amount)
}
