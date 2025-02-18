package api

import (
	"context"
	"net/http"

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
			AuthenticationFunc: processJWT,
		},
		SilenceServersWarning: true,
		Skipper:               skipper,
	})

	return validator, nil
}

func ErrHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	var he *echo.HTTPError
	if errors.As(err, &he) {
		_ = c.JSON(he.Code, echo.Map{"errors": he.Message})
		return
	}

	_ = c.JSON(http.StatusInternalServerError, echo.Map{"errors": "Внутренняя ошибка сервера."})
}

func processJWT(ctx context.Context, input *openapi3filter.AuthenticationInput) error {
	c := middleware.GetEchoContext(ctx)

	user := c.Get("user").(*jwt.Token)

	claims, _ := user.Claims.(jwt.MapClaims)

	username, _ := claims["username"].(string)

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
