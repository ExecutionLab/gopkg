package model

import "github.com/dgrijalva/jwt-go"

// UserClaims .
type UserClaims struct {
	jwt.StandardClaims

	UserID int64  `json:"user_id"`
	Email  string `json:"email"`
}
