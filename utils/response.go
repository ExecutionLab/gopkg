//nolint
package utils

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/ExecutionLab/gopkg/apperror"
	"github.com/ExecutionLab/gopkg/logger"
	"github.com/ExecutionLab/gopkg/model"
	"github.com/ExecutionLab/gopkg/sentry"
)

type response struct{}

var Response response

func (response) Success(c echo.Context, result interface{}) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"code":    http.StatusOK,
		"message": "OK",
		"result":  result,
	})
}

func (response) Error(c echo.Context, err apperror.AppError) error {
	logger.NewLogger().Log(err.Raw)

	if err.IsSentry {
		mySentry := sentry.NewSentry()
		claims := c.Get(string(model.KeyContextToken))

		if claims != nil {
			userClaims := claims.(*model.UserClaims)

			mySentry.Option(
				sentry.WithFields("user_id", strconv.FormatInt(userClaims.UserID, 10)),
			)
		}

		mySentry.Log(err)
	}

	return c.JSON(err.HTTPCode, map[string]interface{}{
		"code":    err.ErrorCode,
		"info":    err.Info,
		"message": err.Message,
	})
}
