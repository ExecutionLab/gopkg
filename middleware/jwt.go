//nolint
package middleware

import (
	"context"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/ExecutionLab/gopkg/model"
)

type TokenType string

const (
	TokenTypeAuth    TokenType = "auth"
	TokenTypeRefresh TokenType = "refresh"
)

type TokenService struct {
	Key string
}

func NewTokenSvc(key string) TokenService {
	return TokenService{
		Key: key,
	}
}

func (t TokenService) EncodeRefreshToken(
	userID int64,
	email string,
	issuer string,
	expire time.Duration,
) (string, error) {
	if issuer == "" {
		issuer = "Husol"
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, model.UserClaims{
		UserID: userID,
		Email:  email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(expire).Unix(),
			Issuer:    issuer,
			Subject:   string(TokenTypeRefresh),
		},
	})

	return token.SignedString([]byte(t.Key))
}

func (t TokenService) Encode(userID int64, email string, issuer string, expire time.Duration) (string, error) {
	if issuer == "" {
		issuer = "Husol"
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, model.UserClaims{
		UserID: userID,
		Email:  email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(expire).Unix(),
			Issuer:    issuer,
			Subject:   string(TokenTypeAuth),
		},
	})

	return token.SignedString([]byte(t.Key))
}

// Decode ValidateToken.
func (t TokenService) Decode(tokenString string) (*model.UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &model.UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(t.Key), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*model.UserClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, nil
}

func (t TokenService) DecodeAuthToken(tokenString string) (*model.UserClaims, error) {
	claims, err := t.Decode(tokenString)
	if err != nil {
		return nil, err
	}
	if claims.Subject != string(TokenTypeAuth) {
		return nil, nil
	}

	return claims, nil
}

func (t TokenService) DecodeRefreshToken(tokenString string) (*model.UserClaims, error) {
	claims, err := t.Decode(tokenString)
	if err != nil {
		return nil, err
	}
	if claims.Subject != string(TokenTypeRefresh) {
		return nil, nil
	}

	return claims, nil
}

func SetClaim(ctx context.Context, data *model.UserClaims) context.Context {
	newCtx := context.WithValue(ctx, model.KeyContextToken, data)

	return newCtx
}

func GetClaim(c context.Context) *model.UserClaims {
	tokenPayloadInterface := c.Value(model.KeyContextToken)

	tokenPayload, ok := tokenPayloadInterface.(*model.UserClaims)
	if !ok {
		return nil
	}

	return tokenPayload
}
