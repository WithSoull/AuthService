package env

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/WithSoull/AuthService/internal/config"
)

const (
	securityMaxLoginAttemptsEnvName    = "SECURITY_MAX_LOGIN_ATTEMPTS"
	securityLoginAttemptsWindowEnvName = "SECURITY_LOGIN_ATTEMPTS_WINDOW"
)

const (
	defaultMaxLoginAttempts    = 5
	defaultLoginAttemptsWindow = 15 * time.Minute
)

type SecurityConfig struct {
	maxLoginAttempts    int8
	loginAttemptsWindow time.Duration
}

func NewSecurityConfig() config.SecurityConfig {
	maxLoginAttempts := defaultMaxLoginAttempts
	loginAttemptsWindow := defaultLoginAttemptsWindow

	// Parse max login attempts
	maxLoginAttemptsStr := os.Getenv(securityMaxLoginAttemptsEnvName)
	if len(maxLoginAttemptsStr) == 0 {
		log.Printf("[Config] %s not found in environment variables, using default value: %d",
			securityMaxLoginAttemptsEnvName, defaultMaxLoginAttempts)
	} else {
		if parsed, err := strconv.ParseInt(maxLoginAttemptsStr, 10, 8); err != nil {
			log.Printf("[Config] Invalid format for %s (%s), using default value: %d",
				securityMaxLoginAttemptsEnvName, maxLoginAttemptsStr, defaultMaxLoginAttempts)
		} else {
			maxLoginAttempts = int(parsed)
			log.Printf("[Config] Using %s from environment: %d",
				securityMaxLoginAttemptsEnvName, maxLoginAttempts)
		}
	}

	// Parse login attempts window
	loginAttemptsWindowStr := os.Getenv(securityLoginAttemptsWindowEnvName)
	if len(loginAttemptsWindowStr) == 0 {
		log.Printf("[Config] %s not found in environment variables, using default value: %s",
			securityLoginAttemptsWindowEnvName, defaultLoginAttemptsWindow)
	} else {
		if parsed, err := time.ParseDuration(loginAttemptsWindowStr); err != nil {
			log.Printf("[Config] Invalid format for %s (%s), using default value: %s",
				securityLoginAttemptsWindowEnvName, loginAttemptsWindowStr, defaultLoginAttemptsWindow)
		} else {
			loginAttemptsWindow = parsed
			log.Printf("[Config] Using %s from environment: %s",
				securityLoginAttemptsWindowEnvName, loginAttemptsWindow)
		}
	}

	log.Printf("[Config] Security config loaded successfully - MaxAttempts: %d, Window: %s",
		maxLoginAttempts, loginAttemptsWindow)

	return &SecurityConfig{
		maxLoginAttempts:    int8(maxLoginAttempts),
		loginAttemptsWindow: loginAttemptsWindow,
	}
}

func (cfg *SecurityConfig) MaxLoginAttempts() int8 {
	return cfg.maxLoginAttempts
}

func (cfg *SecurityConfig) LoginAttemptsWindow() time.Duration {
	return cfg.loginAttemptsWindow
}
