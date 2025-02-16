package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	middleware "github.com/oapi-codegen/echo-middleware"
)

func JWTParser() echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte("secret"),
		Skipper:    skipper,
	})
}

func jwtValidator(ctx context.Context, input *openapi3filter.AuthenticationInput) error {
	c := middleware.GetEchoContext(ctx)

	//x := input.RequestValidationInput.Request.Header.Get("Authorization")
	//
	//splitToken := strings.Split(x, "Bearer ")
	//if len(splitToken) != 2 {
	//	// Error: Bearer token not in proper format
	//}
	//
	//token, err := jwt.Parse(splitToken[1], func(token *jwt.Token) (interface{}, error) {
	//	return []byte("secret"), nil
	//})
	//if err != nil {
	//	return err
	//}

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

	//c := middleware.GetEchoContext(ctx)
	ctx = context.WithValue(c.Request().Context(), "username", username)
	input.RequestValidationInput.Request = c.Request().WithContext(ctx)
	c.SetRequest(input.RequestValidationInput.Request)

	return nil
}

func ValidatorMiddleware() (echo.MiddlewareFunc, error) {
	// Загружаем OpenAPI-спецификацию
	spec, err := GetSwagger()
	if err != nil {
		return nil, fmt.Errorf("loading spec: %w", err)
	}

	spec.Servers = nil

	// Создаем middleware с валидацией запросов
	validator := middleware.OapiRequestValidatorWithOptions(spec, &middleware.Options{
		Options: openapi3filter.Options{
			AuthenticationFunc: jwtValidator,
		},
		SilenceServersWarning: true,
		Skipper:               skipper,
	})

	return validator, nil
}

func skipper(c echo.Context) bool {
	if c.Request().URL.Path == "/api/auth" {
		return true
	}
	return false
}
