package api_test

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"avito-shop-service/internal/api"
	"avito-shop-service/internal/logutil"
	mockapi "avito-shop-service/internal/mocks/api"
	"avito-shop-service/internal/model"
)

func TestMain(m *testing.M) {
	logutil.Setup()
	os.Exit(m.Run())
}

func TestAPI_PostApiAuth(t *testing.T) {
	s := mockapi.NewService(t)
	ctx := context.Background()
	s.EXPECT().Login(ctx, "user", "password").Return("token", nil).Once()

	req := api.PostApiAuthRequestObject{
		Body: &api.PostApiAuthJSONRequestBody{
			Username: "user",
			Password: "password",
		},
	}
	a := api.New(s)

	resp, err := a.PostApiAuth(ctx, req)

	token := "token"
	require.NoError(t, err)
	require.Equal(t, api.PostApiAuth200JSONResponse{Token: &token}, resp)
}

func TestAPI_PostApiBuy(t *testing.T) {
	s := mockapi.NewService(t)
	ctx := context.WithValue(context.Background(), "username", "user")
	s.EXPECT().BuyItem(ctx, "user", "book").Return(nil).Once()

	req := api.PostApiBuyRequestObject{
		Body: &api.PostApiBuyJSONRequestBody{
			Item: "book",
		},
	}
	a := api.New(s)

	resp, err := a.PostApiBuy(ctx, req)

	require.NoError(t, err)
	require.Equal(t, api.PostApiBuy200Response{}, resp)
}

func TestAPI_PostGetApiInfo(t *testing.T) {
	s := mockapi.NewService(t)
	ctx := context.WithValue(context.Background(), "username", "user")
	s.EXPECT().Balance(ctx, "user").Return(700, nil).Once()
	s.EXPECT().Inventory(ctx, "user").Return([]model.Inventory{
		{
			Username: "user",
			ItemName: "hoody",
			Quantity: 1,
		},
	}, nil).Once()
	s.EXPECT().Transaction(ctx, "user").Return([]model.Transaction{
		{
			FromUser: "registration_gift",
			ToUser:   "user",
			Amount:   1000,
		},
		{
			FromUser: "user",
			ToUser:   "buy_item",
			Amount:   300,
		},
	}, nil).Once()

	req := api.GetApiInfoRequestObject{}
	a := api.New(s)

	_, err := a.GetApiInfo(ctx, req)

	require.NoError(t, err)
}

func TestAPI_PostApiSendCoin(t *testing.T) {
	s := mockapi.NewService(t)
	ctx := context.WithValue(context.Background(), "username", "user")
	s.EXPECT().SendCoin(ctx, "user", "testuser", 500).Return(nil).Once()

	req := api.PostApiSendCoinRequestObject{
		Body: &api.PostApiSendCoinJSONRequestBody{
			Amount: 500,
			ToUser: "testuser",
		},
	}
	a := api.New(s)

	resp, err := a.PostApiSendCoin(ctx, req)

	require.NoError(t, err)
	require.Equal(t, api.PostApiSendCoin200Response{}, resp)
}
