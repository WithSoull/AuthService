package config

import (
	"time"

	"github.com/joho/godotenv"
)

func Load(path string) error {
	if err := godotenv.Load(path); err != nil {
		return err
	}
	return nil
}

type GRPCConfig interface {
	Address() string
}

type RedisConfig interface {
	Address() string
	MaxIdle() int8
	ConnTimeout() time.Duration
	IdleTimeout() time.Duration
}

type SecurityConfig interface {
	MaxLoginAttempts() int8
	LoginAttemptsWindow() time.Duration
}

type JWTConfig interface {
	RefreshTokenSecretKey() string
	AccessTokenSecretKey() string

	RefreshTokenExpiration() time.Duration
	AccessTokenExpiration() time.Duration
}
