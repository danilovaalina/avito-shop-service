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
	return &API{service: service}
}

func (a *API) PostApiAuth(ctx context.Context, request PostApiAuthRequestObject) (PostApiAuthResponseObject, error) {
	token, err := a.service.Login(ctx, request.Body.Username, request.Body.Password)

	if err != nil {
		if errors.Is(err, model.ErrUserUnauthorized) {
			return PostApiAuth401JSONResponse{Errors: ptrString(err.Error())}, nil
		}
		return PostApiAuth400JSONResponse{Errors: ptrString(err.Error())}, nil
	}

	return PostApiAuth200JSONResponse{Token: &token}, nil
}

func (a *API) PostApiBuy(ctx context.Context, request PostApiBuyRequestObject) (PostApiBuyResponseObject, error) {
	v := ctx.Value("username")

	username, ok := v.(string)
	if !ok {
		return PostApiBuy500JSONResponse{Errors: ptrString("invalid username format")}, nil
	}

	err := a.service.BuyItem(ctx, username, request.Body.Item)
	if err != nil {
		return PostApiBuy400JSONResponse{Errors: ptrString(err.Error())}, nil
	}
	return PostApiBuy200Response{}, nil
}

func (a *API) GetApiInfo(ctx context.Context, request GetApiInfoRequestObject) (GetApiInfoResponseObject, error) {
	v := ctx.Value("username")

	username, ok := v.(string)
	if !ok {
		return GetApiInfo500JSONResponse{Errors: ptrString("invalid username format")}, nil
	}

	balance, err := a.service.Balance(ctx, username)
	if err != nil {
		return GetApiInfo500JSONResponse{Errors: ptrString(err.Error())}, nil
	}

	inventory, err := a.service.Inventory(ctx, username)
	if err != nil {
		return GetApiInfo500JSONResponse{Errors: ptrString(err.Error())}, nil
	}

	transactions, err := a.service.Transaction(ctx, username)
	if err != nil {
		return GetApiInfo500JSONResponse{Errors: ptrString(err.Error())}, nil
	}

	resp := buildInfoResponse(username, balance, transactions, inventory)

	return GetApiInfo200JSONResponse(resp), nil
}

func (a *API) PostApiSendCoin(ctx context.Context, request PostApiSendCoinRequestObject) (PostApiSendCoinResponseObject, error) {
	v := ctx.Value("username")

	fromUser, ok := v.(string)
	if !ok {
		return PostApiSendCoin500JSONResponse{Errors: ptrString("invalid username format")}, nil
	}

	err := a.service.SendCoin(ctx, fromUser, request.Body.ToUser, request.Body.Amount)
	if err != nil {
		return PostApiSendCoin400JSONResponse{Errors: ptrString(err.Error())}, nil
	}

	return PostApiSendCoin200Response{}, nil
}

func ptrString(s string) *string {
	return &s
}

func ptrInt(i int) *int {
	return &i
}

func buildInfoResponse(username string, coins int64, transactions []model.Transaction, inventory []model.Inventory) InfoResponse {
	received := make([]struct {
		Amount   *int    `json:"amount,omitempty"`
		FromUser *string `json:"fromUser,omitempty"`
	}, 0)

	sent := make([]struct {
		Amount *int    `json:"amount,omitempty"`
		ToUser *string `json:"toUser,omitempty"`
	}, 0)

	for _, transaction := range transactions {
		if transaction.ToUser == username {
			received = append(received, struct {
				Amount   *int    `json:"amount,omitempty"`
				FromUser *string `json:"fromUser,omitempty"`
			}{
				Amount:   ptrInt(int(transaction.Amount)),
				FromUser: ptrString(transaction.FromUser),
			})
		} else if transaction.FromUser == username {
			sent = append(sent, struct {
				Amount *int    `json:"amount,omitempty"`
				ToUser *string `json:"toUser,omitempty"`
			}{
				Amount: ptrInt(int(transaction.Amount)),
				ToUser: ptrString(transaction.ToUser),
			})
		}
	}

	coinHistory := &struct {
		Received *[]struct {
			Amount   *int    `json:"amount,omitempty"`
			FromUser *string `json:"fromUser,omitempty"`
		} `json:"received,omitempty"`
		Sent *[]struct {
			Amount *int    `json:"amount,omitempty"`
			ToUser *string `json:"toUser,omitempty"`
		} `json:"sent,omitempty"`
	}{
		Received: &received,
		Sent:     &sent,
	}

	inventoryResponse := make([]struct {
		Quantity *int    `json:"quantity,omitempty"`
		Type     *string `json:"type,omitempty"`
	}, len(inventory))

	for i, item := range inventory {
		inventoryResponse[i] = struct {
			Quantity *int    `json:"quantity,omitempty"`
			Type     *string `json:"type,omitempty"`
		}{
			Quantity: ptrInt(int(item.Quantity)),
			Type:     ptrString(item.ItemName),
		}
	}

	return InfoResponse{
		CoinHistory: coinHistory,
		Coins:       ptrInt(int(coins)),
		Inventory:   &inventoryResponse,
	}
}
