//nolint
package middleware

import (
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/ExecutionLab/gopkg/apperror"
	"github.com/ExecutionLab/gopkg/model"
	"github.com/ExecutionLab/gopkg/utils"
)

func Auth(secretKeySign string, skipper middleware.Skipper, isRefresh bool) func(next echo.HandlerFunc) echo.HandlerFunc {
	handlerFunc := func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if skipper(c) {
				return next(c)
			}

			tokenService := &TokenService{
				Key: secretKeySign,
			}

			token := c.Request().Header.Get("Authorization")
			token = strings.TrimSpace(strings.TrimPrefix(token, "Bearer "))

			if token == "" {
				token = c.QueryParam("token")
			}

			var decodeFunc func(tokenString string) (*model.UserClaims, error)

			if isRefresh {
				decodeFunc = tokenService.DecodeRefreshToken
			} else {
				decodeFunc = tokenService.DecodeAuthToken
			}

			claims, err := decodeFunc(token)
			if err != nil {
				if err.(*jwt.ValidationError).Errors == jwt.ValidationErrorExpired {
					return utils.Response.Error(c, apperror.ErrUnauthorizedExpiredToken(err))
				}

				return utils.Response.Error(c, apperror.ErrUnauthorized(err))
			}

			if claims == nil {
				return utils.Response.Error(c, apperror.ErrUnauthorized(nil))
			}

			c.Set(string(model.KeyContextToken), claims)
			return next(c)
		}
	}
	return handlerFunc
}
