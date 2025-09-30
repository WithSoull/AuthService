package model

import "github.com/golang-jwt/jwt/v5"

type UserClaims struct {
	jwt.RegisteredClaims
	UserId    int64  `json:"user_id"`
	Email     string `json:"email"`
	TokenType string `json:"token_type"` // refresh or access
}
