package env

import (
	"errors"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/WithSoull/AuthService/internal/config"
)

const (
	// Redis connection constants
	redisHostEnvName              = "CACHE_HOST"
	redisPortEnvName              = "CACHE_PORT"
	redisConnectionTimeoutEnvName = "CACHE_CONNECTION_TIMEOUT"
	redisMaxIdleEnvName           = "CACHE_MAX_IDLE"
	redisIdleTimeoutEnvName       = "CACHE_IDLE_TIMEOUT"
)

// Default values for Redis connection
const (
	defaultRedisPort   = "6379"
	defaultMaxIdle     = 10
	defaultConnTimeout = 5 * time.Second
	defaultIdleTimeout = 240 * time.Second
)

type redisConfig struct {
	host        string
	port        string
	maxIdle     int8
	connTimeout time.Duration
	idleTimeout time.Duration
}

func NewRedisConfig() (config.RedisConfig, error) {
	host := os.Getenv(redisHostEnvName)
	if len(host) == 0 {
		return nil, errors.New("redis host not found - this is required")
	}

	port := defaultRedisPort
	if portEnv := os.Getenv(redisPortEnvName); len(portEnv) == 0 {
		log.Printf("[Config] %s not found in environment variables, using default value: %s",
			redisPortEnvName, defaultRedisPort)
	} else {
		port = portEnv
		log.Printf("[Config] Using %s from environment: %s", redisPortEnvName, port)
	}

	maxIdle := defaultMaxIdle
	maxIdleStr := os.Getenv(redisMaxIdleEnvName)
	if len(maxIdleStr) == 0 {
		log.Printf("[Config] %s not found in environment variables, using default value: %d",
			redisMaxIdleEnvName, defaultMaxIdle)
	} else {
		if parsed, err := strconv.ParseInt(maxIdleStr, 10, 8); err != nil {
			log.Printf("[Config] Invalid format for %s (%s), using default value: %d",
				redisMaxIdleEnvName, maxIdleStr, defaultMaxIdle)
		} else {
			maxIdle = int(parsed)
			log.Printf("[Config] Using %s from environment: %d", redisMaxIdleEnvName, maxIdle)
		}
	}

	connTimeout := defaultConnTimeout
	connTimeoutStr := os.Getenv(redisConnectionTimeoutEnvName)
	if len(connTimeoutStr) == 0 {
		log.Printf("[Config] %s not found in environment variables, using default value: %s",
			redisConnectionTimeoutEnvName, defaultConnTimeout)
	} else {
		if parsed, err := time.ParseDuration(connTimeoutStr); err != nil {
			log.Printf("[Config] Invalid format for %s (%s), using default value: %s",
				redisConnectionTimeoutEnvName, connTimeoutStr, defaultConnTimeout)
		} else {
			connTimeout = parsed
			log.Printf("[Config] Using %s from environment: %s",
				redisConnectionTimeoutEnvName, connTimeout)
		}
	}

	idleTimeout := defaultIdleTimeout
	idleTimeoutStr := os.Getenv(redisIdleTimeoutEnvName)
	if len(idleTimeoutStr) == 0 {
		log.Printf("[Config] %s not found in environment variables, using default value: %s",
			redisIdleTimeoutEnvName, defaultIdleTimeout)
	} else {
		if parsed, err := time.ParseDuration(idleTimeoutStr); err != nil {
			log.Printf("[Config] Invalid format for %s (%s), using default value: %s",
				redisIdleTimeoutEnvName, idleTimeoutStr, defaultIdleTimeout)
		} else {
			idleTimeout = parsed
			log.Printf("[Config] Using %s from environment: %s",
				redisIdleTimeoutEnvName, idleTimeout)
		}
	}

	log.Printf("[Config] Redis connection config loaded successfully - Host: %s, Port: %s", host, port)

	return &redisConfig{
		host:        host,
		port:        port,
		maxIdle:     int8(maxIdle),
		connTimeout: connTimeout,
		idleTimeout: idleTimeout,
	}, nil
}

func (cfg *redisConfig) Address() string {
	address := net.JoinHostPort(cfg.host, cfg.port)
	return address
}

func (cfg *redisConfig) MaxIdle() int8 {
	return cfg.maxIdle
}

func (cfg *redisConfig) ConnTimeout() time.Duration {
	return cfg.connTimeout
}

func (cfg *redisConfig) IdleTimeout() time.Duration {
	return cfg.idleTimeout
}
