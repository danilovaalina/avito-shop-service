package test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	pgcontainer "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"avito-shop-service/internal/api"
	"avito-shop-service/internal/postgres"
	"avito-shop-service/internal/repository"
	"avito-shop-service/internal/service"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestBuyItem(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	url, cleanup, err := setupServer(ctx)
	require.NoError(t, err)
	defer cleanup(ctx)

	client, err := api.NewClientWithResponses(url)
	require.NoError(t, err)

	authBody := api.PostApiAuthJSONRequestBody{
		Username: "testuser",
		Password: "testpassword",
	}

	authResp, err := client.PostApiAuthWithResponse(ctx, authBody)
	require.NoError(t, err)

	require.Equal(t, http.StatusOK, authResp.StatusCode())
	require.NotNil(t, authResp.JSON200.Token)

	buyBody := api.PostApiBuyJSONRequestBody{
		Item: "cup",
	}

	buyResp, err := client.PostApiBuyWithResponse(ctx, buyBody, withToken(*authResp.JSON200.Token))
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, buyResp.StatusCode())

	infoResp, err := client.GetApiInfoWithResponse(ctx, withToken(*authResp.JSON200.Token))
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, infoResp.StatusCode())

	require.NotNil(t, infoResp.JSON200.CoinHistory)
	require.NotNil(t, infoResp.JSON200.CoinHistory.Received)
	require.NotNil(t, infoResp.JSON200.CoinHistory.Sent)
	require.NotNil(t, infoResp.JSON200.Coins)
	require.NotNil(t, infoResp.JSON200.Inventory)

	require.Equal(t, 1, len(*infoResp.JSON200.CoinHistory.Received))
	require.Equal(t, 1, len(*infoResp.JSON200.CoinHistory.Sent))
	require.Equal(t, 1, len(*infoResp.JSON200.Inventory))

	require.Equal(t, 1000, (*infoResp.JSON200.CoinHistory.Received)[0].Amount)
	require.Equal(t, "registration_gift", (*infoResp.JSON200.CoinHistory.Received)[0].FromUser)

	require.Equal(t, 20, (*infoResp.JSON200.CoinHistory.Sent)[0].Amount)
	require.Equal(t, "buy_item", (*infoResp.JSON200.CoinHistory.Sent)[0].ToUser)

	require.Equal(t, 1, (*infoResp.JSON200.Inventory)[0].Quantity)
	require.Equal(t, "cup", (*infoResp.JSON200.Inventory)[0].Type)
}

func TestSendCoin(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	url, cleanup, err := setupServer(ctx)
	require.NoError(t, err)
	defer cleanup(ctx)

	client, err := api.NewClientWithResponses(url)
	require.NoError(t, err)

	senderBody := api.PostApiAuthJSONRequestBody{
		Username: "sender",
		Password: "sender",
	}

	senderResp, err := client.PostApiAuthWithResponse(ctx, senderBody)
	require.NoError(t, err)

	require.Equal(t, http.StatusOK, senderResp.StatusCode())
	require.NotNil(t, senderResp.JSON200.Token)

	receiverBody := api.PostApiAuthJSONRequestBody{
		Username: "receiver",
		Password: "receiver",
	}

	receiverResp, err := client.PostApiAuthWithResponse(ctx, receiverBody)
	require.NoError(t, err)

	require.Equal(t, http.StatusOK, receiverResp.StatusCode())
	require.NotNil(t, receiverResp.JSON200.Token)

	sendCoinBody := api.PostApiSendCoinJSONRequestBody{
		ToUser: "receiver",
		Amount: 500,
	}

	sendCoinResp, err := client.PostApiSendCoinWithResponse(ctx, sendCoinBody, withToken(*senderResp.JSON200.Token))
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, sendCoinResp.StatusCode())

	infoResp, err := client.GetApiInfoWithResponse(ctx, withToken(*senderResp.JSON200.Token))
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, infoResp.StatusCode())

	require.NotNil(t, infoResp.JSON200.CoinHistory)
	require.NotNil(t, infoResp.JSON200.CoinHistory.Received)
	require.NotNil(t, infoResp.JSON200.CoinHistory.Sent)
	require.NotNil(t, infoResp.JSON200.Coins)
	require.NotNil(t, infoResp.JSON200.Inventory)

	require.Equal(t, 1, len(*infoResp.JSON200.CoinHistory.Received))
	require.Equal(t, 1, len(*infoResp.JSON200.CoinHistory.Sent))
	require.Equal(t, 0, len(*infoResp.JSON200.Inventory))

	require.Equal(t, 1000, (*infoResp.JSON200.CoinHistory.Received)[0].Amount)
	require.Equal(t, "registration_gift", (*infoResp.JSON200.CoinHistory.Received)[0].FromUser)

	require.Equal(t, 500, (*infoResp.JSON200.CoinHistory.Sent)[0].Amount)
	require.Equal(t, "receiver", (*infoResp.JSON200.CoinHistory.Sent)[0].ToUser)

	require.Equal(t, 500, infoResp.JSON200.Coins)

	infoResp, err = client.GetApiInfoWithResponse(ctx, withToken(*receiverResp.JSON200.Token))
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, infoResp.StatusCode())

	require.NotNil(t, infoResp.JSON200.CoinHistory)
	require.NotNil(t, infoResp.JSON200.CoinHistory.Received)
	require.NotNil(t, infoResp.JSON200.CoinHistory.Sent)
	require.NotNil(t, infoResp.JSON200.Coins)
	require.NotNil(t, infoResp.JSON200.Inventory)

	require.Equal(t, 2, len(*infoResp.JSON200.CoinHistory.Received))
	require.Equal(t, 0, len(*infoResp.JSON200.CoinHistory.Sent))
	require.Equal(t, 0, len(*infoResp.JSON200.Inventory))

	require.Equal(t, 1000, (*infoResp.JSON200.CoinHistory.Received)[0].Amount)
	require.Equal(t, "registration_gift", (*infoResp.JSON200.CoinHistory.Received)[0].FromUser)

	require.Equal(t, 500, (*infoResp.JSON200.CoinHistory.Received)[1].Amount)
	require.Equal(t, "sender", (*infoResp.JSON200.CoinHistory.Received)[1].FromUser)

	require.Equal(t, 1500, infoResp.JSON200.Coins)
}

func TestCheckJWTValidation(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	url, cleanup, err := setupServer(ctx, withJWT(service.NewJWT(service.WithLifetime(time.Second))))
	require.NoError(t, err)
	defer cleanup(ctx)

	client, err := api.NewClientWithResponses(url)
	require.NoError(t, err)

	authBody := api.PostApiAuthJSONRequestBody{
		Username: "testuser",
		Password: "testpassword",
	}

	authResp, err := client.PostApiAuthWithResponse(ctx, authBody)
	require.NoError(t, err)

	require.Equal(t, http.StatusOK, authResp.StatusCode())
	require.NotNil(t, authResp.JSON200.Token)

	time.Sleep(time.Second)

	infoResp, err := client.GetApiInfoWithResponse(ctx, withToken(*authResp.JSON200.Token))
	require.NoError(t, err)
	require.Equal(t, http.StatusUnauthorized, infoResp.StatusCode())
}

type serverOptions struct {
	jwt *service.JWT
}
type serverOption func(opts *serverOptions)

func withJWT(jwt *service.JWT) serverOption {
	return func(o *serverOptions) {
		o.jwt = jwt
	}
}

func setupServer(ctx context.Context, opts ...serverOption) (string, func(context.Context), error) {
	o := new(serverOptions)
	for _, opt := range opts {
		opt(o)
	}
	if o.jwt == nil {
		o.jwt = service.NewJWT()
	}

	container, connStr, err := setupDB(ctx)
	if err != nil {
		return "", nil, err
	}

	pool, err := postgres.Pool(ctx, connStr)
	if err != nil {
		_ = container.Terminate(ctx)
		return "", nil, err
	}

	a := api.New(service.New(repository.New(pool), o.jwt, service.NewBcrypt()))

	e := echo.New()
	e.HTTPErrorHandler = api.ErrHandler
	api.RegisterHandlers(e, api.NewStrictHandler(a, nil))

	e.Use(api.JWTParser())

	v, err := api.ValidatorMiddleware()
	if err != nil {
		pool.Close()
		_ = container.Terminate(ctx)
		return "", nil, err
	}

	e.Use(v)

	s := httptest.NewServer(e)

	return s.URL, func(ctx context.Context) {
		s.Close()
		pool.Close()
		_ = container.Terminate(ctx)
	}, nil
}

func setupDB(ctx context.Context) (testcontainers.Container, string, error) {
	container, err := pgcontainer.Run(ctx,
		"postgres:14",
		pgcontainer.WithUsername("testuser"),
		pgcontainer.WithPassword("testpass"),
		pgcontainer.WithDatabase("testdb"),
		pgcontainer.WithInitScripts("../scripts/init.sql"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		return nil, "", errors.WithStack(err)
	}

	connStr, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return nil, "", errors.WithStack(err)
	}

	return container, connStr, nil
}

func withToken(token string) api.RequestEditorFn {
	return func(ctx context.Context, req *http.Request) error {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
		return nil
	}
}
