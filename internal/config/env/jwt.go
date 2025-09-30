package env

import (
	"errors"
	"os"
	"time"

	"github.com/WithSoull/AuthService/internal/config"
)

const (
	refreshTokenSecretKeyEnvName = "REFRESH_TOKEN_SECRET"
	accessTokenSecretKeyEnvName  = "ACCESS_TOKEN_SECRET"

	refreshTokenTTLEnvName = "REFRESH_TOKEN_TTL"
	accessTokenTTLEnvName  = "ACCESS_TOKEN_TTL"
)

type jwtConfig struct {
	refreshTokenSecretKey string
	accessTokenSecretKey  string

	refreshTokenTTL time.Duration
	accessTokenTTL  time.Duration
}

func NewJWTConfig() (config.JWTConfig, error) {
	refreshSecret := os.Getenv(refreshTokenSecretKeyEnvName)
	if len(refreshSecret) == 0 {
		return nil, errors.New("refresh secret key not found")
	}

	accessSecret := os.Getenv(accessTokenSecretKeyEnvName)
	if len(accessSecret) == 0 {
		return nil, errors.New("access secret key not found")
	}

	refreshTTLStr := os.Getenv(refreshTokenTTLEnvName)
	if len(refreshTTLStr) == 0 {
		return nil, errors.New("refresh token ttl not found")
	}

	accessTTLStr := os.Getenv(accessTokenTTLEnvName)
	if len(accessTTLStr) == 0 {
		return nil, errors.New("access token ttl not found")
	}

	refreshTTL, err := time.ParseDuration(refreshTTLStr)
	if err != nil {
		return nil, err
	}

	accessTTL, err := time.ParseDuration(accessTTLStr)
	if err != nil {
		return nil, err
	}

	return &jwtConfig{
		refreshTokenSecretKey: refreshSecret,
		accessTokenSecretKey:  accessSecret,
		refreshTokenTTL:       refreshTTL,
		accessTokenTTL:        accessTTL,
	}, nil
}

func (c *jwtConfig) RefreshTokenSecretKey() string {
	return c.refreshTokenSecretKey
}

func (c *jwtConfig) AccessTokenSecretKey() string {
	return c.accessTokenSecretKey
}

func (c *jwtConfig) RefreshTokenExpiration() time.Duration {
	return c.refreshTokenTTL
}

func (c *jwtConfig) AccessTokenExpiration() time.Duration {
	return c.accessTokenTTL
}

