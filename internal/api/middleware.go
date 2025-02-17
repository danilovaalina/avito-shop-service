package api

import (
	"context"
	"net/http"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	middleware "github.com/oapi-codegen/echo-middleware"
	"github.com/rs/zerolog/log"
)

func LoggerMiddleware() echo.MiddlewareFunc {
	return echomiddleware.RequestLoggerWithConfig(echomiddleware.RequestLoggerConfig{
		LogValuesFunc: func(c echo.Context, v echomiddleware.RequestLoggerValues) error {
			if v.Error != nil {
				log.Error().Stack().Err(v.Error).Send()
			}

			return nil
		},
		HandleError: true,
		LogError:    true,
	})
}

func JWTParser() echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte("secret"),
		Skipper:    skipper,
	})
}

func ValidatorMiddleware() (echo.MiddlewareFunc, error) {
	spec, err := GetSwagger()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	spec.Servers = nil

	validator := middleware.OapiRequestValidatorWithOptions(spec, &middleware.Options{
		Options: openapi3filter.Options{
			AuthenticationFunc: jwtValidator,
		},
		SilenceServersWarning: true,
		Skipper:               skipper,
	})

	return validator, nil
}

func ErrHandler(_ error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	_ = c.JSON(http.StatusInternalServerError, echo.Map{"errors": "Внутренняя ошибка сервера."})
}

func jwtValidator(ctx context.Context, input *openapi3filter.AuthenticationInput) error {
	c := middleware.GetEchoContext(ctx)

	user := c.Get("user").(*jwt.Token)
	if user == nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Токен отсутствует")
	}

	claims, ok := user.Claims.(jwt.MapClaims)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "Неверный формат токена")
	}

	username, ok := claims["username"].(string)
	if !ok || username == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "Отсутствует или некорректное имя пользователя в токене")
	}

	exp, ok := claims["exp"].(float64)
	if !ok || exp < float64(time.Now().Unix()) {
		return echo.NewHTTPError(http.StatusUnauthorized, "Токен истек")
	}

	ctx = context.WithValue(c.Request().Context(), "username", username)
	input.RequestValidationInput.Request = c.Request().WithContext(ctx)
	c.SetRequest(input.RequestValidationInput.Request)

	return nil
}

func skipper(c echo.Context) bool {
	if c.Request().URL.Path == "/api/auth" {
		return true
	}
	return false
}
