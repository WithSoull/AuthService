package config

import (
	"time"
)

type GRPCConfig interface {
	Address() string
}

type LoggerConfig interface {
	LogLevel() string
	AsJSON() bool
	EnableOLTP() bool
	ServiceName() string
	OTLPEndpoint() string
	ServiceEnvironment() string
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

type TracingConfig interface {
	CollectorEndpoint() string
	ServiceName() string
	Environment() string
	ServiceVersion() string
}
