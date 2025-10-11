package jwt

import (
	"fmt"
	"log"
	"time"

	"github.com/WithSoull/AuthService/internal/config"
	"github.com/WithSoull/AuthService/internal/model"
	"github.com/golang-jwt/jwt/v5"
)

type TokenType string

const (
	AccessToken  TokenType = "access"
	RefreshToken TokenType = "refresh"
)

type JWTService struct {
}

func NewJWTService() *JWTService {
	return &JWTService{}
}

func (j *JWTService) GenerateAccessToken(info model.UserInfo) (string, error) {
	secretKey := []byte(config.AppConfig().JWT.AccessTokenSecretKey())
	duration := config.AppConfig().JWT.AccessTokenExpiration()
	return j.generateToken(info, duration, AccessToken, secretKey)
}

func (j *JWTService) GenerateRefreshToken(info model.UserInfo) (string, error) {
	secretKey := []byte(config.AppConfig().JWT.RefreshTokenSecretKey())
	duration := config.AppConfig().JWT.RefreshTokenExpiration()
	return j.generateToken(info, duration, RefreshToken, secretKey)
}

func (j *JWTService) generateToken(info model.UserInfo, duration time.Duration, tokenType TokenType, secretKey []byte) (string, error) {
	if len(secretKey) == 0 {
		return "", fmt.Errorf("%s secret key is empty", tokenType)
	}

	now := time.Now()
	claims := model.UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(duration)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
		UserId:    info.UserId,
		Email:     info.Email,
		TokenType: string(tokenType),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		log.Printf("failed to sign %s token: %w", tokenType, err)
		return "", fmt.Errorf("failed to sign %s token: %w", tokenType, err)
	}

	return signedToken, nil
}

func (j *JWTService) VerifyAccessToken(tokenStr string) (*model.UserClaims, error) {
	secretKey := []byte(config.AppConfig().JWT.AccessTokenSecretKey())
	claims, err := j.verifyToken(tokenStr, secretKey, AccessToken)
	if err != nil {
		return nil, fmt.Errorf("access token verification failed: %w", err)
	}
	return claims, nil
}

func (j *JWTService) VerifyRefreshToken(tokenStr string) (*model.UserClaims, error) {
	secretKey := []byte(config.AppConfig().JWT.RefreshTokenSecretKey())
	claims, err := j.verifyToken(tokenStr, secretKey, RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("refresh token verification failed: %w", err)
	}
	return claims, nil
}

func (j *JWTService) verifyToken(tokenStr string, secretKey []byte, expectedType TokenType) (*model.UserClaims, error) {
	if tokenStr == "" {
		return nil, fmt.Errorf("token string is empty")
	}

	token, err := jwt.ParseWithClaims(
		tokenStr,
		&model.UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				err := fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				log.Printf("[JWTService] failed to verifyToken: %v", err)
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return secretKey, nil
		},
	)
	if err != nil {
		switch {
		case err == jwt.ErrTokenExpired:
			return nil, fmt.Errorf("token has expired")
		case err == jwt.ErrTokenNotValidYet:
			return nil, fmt.Errorf("token is not valid yet")
		case err == jwt.ErrTokenMalformed:
			return nil, fmt.Errorf("token is malformed")
		default:
			return nil, fmt.Errorf("invalid token: %w", err)
		}
	}

	if !token.Valid {
		return nil, fmt.Errorf("token is not valid")
	}

	claims, ok := token.Claims.(*model.UserClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims type")
	}

	if claims.TokenType != string(expectedType) {
		return nil, fmt.Errorf("invalid token type: expected %s, got %s", expectedType, claims.TokenType)
	}

	if claims.ExpiresAt != nil && time.Now().After(claims.ExpiresAt.Time) {
		return nil, fmt.Errorf("token has expired")
	}

	return claims, nil
}
