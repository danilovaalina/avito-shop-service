package api

import (
	"context"

	"github.com/cockroachdb/errors"

	"avito-shop-service/internal/model"
)

type Service interface {
	Login(ctx context.Context, username, password string) (string, error)
	BuyItem(ctx context.Context, username, itemName string) error
	SendCoin(ctx context.Context, fromUser, toUser string, amount int) error
	Balance(ctx context.Context, username string) (int64, error)
	Inventory(ctx context.Context, username string) ([]model.Inventory, error)
	Transaction(ctx context.Context, username string) ([]model.Transaction, error)
}

type API struct {
	service Service
}

func New(service Service) *API {
	return &API{
		service: service,
	}
}

func (a *API) PostApiAuth(ctx context.Context, request PostApiAuthRequestObject) (PostApiAuthResponseObject, error) {
	token, err := a.service.Login(ctx, request.Body.Username, request.Body.Password)
	if err != nil {
		if errors.Is(err, model.ErrUserUnauthorized) {
			return PostApiAuth401JSONResponse{Errors: "Неавторизован."}, nil
		}
		return nil, err
	}

	return PostApiAuth200JSONResponse{Token: &token}, nil
}

func (a *API) PostApiBuy(ctx context.Context, request PostApiBuyRequestObject) (PostApiBuyResponseObject, error) {
	v := ctx.Value("username")

	username, ok := v.(string)
	if !ok {
		return nil, errors.New("username not found in context")
	}

	err := a.service.BuyItem(ctx, username, request.Body.Item)
	if err != nil {
		if errors.IsAny(err, model.ErrItemNotFound, model.ErrNegativeBalance) {
			return PostApiBuy400JSONResponse{Errors: "Неверный запрос."}, nil
		}
		return nil, err
	}

	return PostApiBuy200Response{}, nil
}

func (a *API) GetApiInfo(ctx context.Context, _ GetApiInfoRequestObject) (GetApiInfoResponseObject, error) {
	v := ctx.Value("username")

	username, ok := v.(string)
	if !ok {
		return nil, errors.New("username not found in context")
	}

	balance, err := a.service.Balance(ctx, username)
	if err != nil {
		return nil, err
	}

	inventory, err := a.service.Inventory(ctx, username)
	if err != nil {
		return nil, err
	}

	transactions, err := a.service.Transaction(ctx, username)
	if err != nil {
		return nil, err
	}

	resp := buildInfoResponse(username, balance, transactions, inventory)

	return GetApiInfo200JSONResponse(resp), nil
}

func (a *API) PostApiSendCoin(ctx context.Context, request PostApiSendCoinRequestObject) (PostApiSendCoinResponseObject, error) {
	v := ctx.Value("username")

	fromUser, ok := v.(string)
	if !ok {
		return nil, errors.New("username not found in context")
	}

	err := a.service.SendCoin(ctx, fromUser, request.Body.ToUser, request.Body.Amount)
	if err != nil {
		if errors.Is(err, model.ErrNegativeBalance) {
			return PostApiSendCoin400JSONResponse{Errors: "Неверный запрос."}, nil
		}
		return nil, err
	}

	return PostApiSendCoin200Response{}, nil
}

func buildInfoResponse(username string, coins int64, transactions []model.Transaction, inventory []model.Inventory) InfoResponse {
	received := make([]struct {
		Amount   int    `json:"amount"`
		FromUser string `json:"fromUser"`
	}, 0)

	sent := make([]struct {
		Amount int    `json:"amount"`
		ToUser string `json:"toUser"`
	}, 0)

	for _, transaction := range transactions {
		if transaction.ToUser == username {
			received = append(received, struct {
				Amount   int    `json:"amount"`
				FromUser string `json:"fromUser"`
			}{
				Amount:   int(transaction.Amount),
				FromUser: transaction.FromUser,
			})
		} else if transaction.FromUser == username {
			sent = append(sent, struct {
				Amount int    `json:"amount"`
				ToUser string `json:"toUser"`
			}{
				Amount: int(transaction.Amount),
				ToUser: transaction.ToUser,
			})
		}
	}

	coinHistory := &struct {
		Received *[]struct {
			Amount   int    `json:"amount"`
			FromUser string `json:"fromUser"`
		} `json:"received,omitempty"`
		Sent *[]struct {
			Amount int    `json:"amount"`
			ToUser string `json:"toUser"`
		} `json:"sent,omitempty"`
	}{
		Received: &received,
		Sent:     &sent,
	}

	inventoryResponse := make([]struct {
		Quantity int    `json:"quantity"`
		Type     string `json:"type"`
	}, len(inventory))

	for i, item := range inventory {
		inventoryResponse[i] = struct {
			Quantity int    `json:"quantity"`
			Type     string `json:"type"`
		}{
			Quantity: int(item.Quantity),
			Type:     item.ItemName,
		}
	}

	return InfoResponse{
		CoinHistory: coinHistory,
		Coins:       int(coins),
		Inventory:   &inventoryResponse,
	}
}
