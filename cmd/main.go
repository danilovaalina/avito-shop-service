package main

import (
	"context"
	"net"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"

	"avito-shop-service/internal/api"
	"avito-shop-service/internal/logutil"
	"avito-shop-service/internal/postgres"
	"avito-shop-service/internal/repository"
	"avito-shop-service/internal/service"
)

const (
	defaultPort = "8080"
)

func main() {
	logutil.Setup()

	pool, err := postgres.Pool(context.Background(), "")
	if err != nil {
		log.Fatal().Stack().Err(err).Send()
	}
	defer pool.Close()

	a := api.New(service.New(repository.New(pool), service.NewJWT(), service.NewBcrypt()))

	e := echo.New()
	api.RegisterHandlers(e, api.NewStrictHandler(a, nil))

	e.Use(api.LoggerMiddleware())
	e.Use(api.JWTParser())

	v, err := api.ValidatorMiddleware()
	if err != nil {
		log.Fatal().Stack().Err(err).Send()
	}

	e.Use(v)

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = defaultPort
	}

	addr := net.JoinHostPort("", port)

	err = e.Start(addr)
	if err != nil {
		log.Fatal().Stack().Err(err).Send()
	}
}
