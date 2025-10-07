package auth

import (
	"context"
	"fmt"

	"github.com/WithSoull/AuthService/internal/client/cache"
	"github.com/WithSoull/AuthService/internal/config"
	"github.com/WithSoull/AuthService/internal/repository"
	"github.com/gomodule/redigo/redis"
)

type redisRepository struct {
	client         cache.CacheClient
	securityConfig config.SecurityConfig
}

func NewRedisRepository(client cache.CacheClient, cfg config.SecurityConfig) repository.AuthRepository {
	return &redisRepository{client: client, securityConfig: cfg}
}

const (
	prefixKey            = "login_attempts"
	maxLoginAttemptLimit = 127 // int8 limitations
)

func (r *redisRepository) getKey(text string) string {
	return fmt.Sprintf("%s:%s", prefixKey, text)
}

func (r *redisRepository) replyToInt8(reply any) (int8, error) {
	count64, err := redis.Int64(reply, nil)
	if err != nil {
		return 0, err
	}

	if count64 > maxLoginAttemptLimit {
		count64 = maxLoginAttemptLimit
	}

	return int8(count64), nil
}

func (r *redisRepository) IncrementLoginAttempts(ctx context.Context, email string) (int8, error) {
	reply, err := r.client.Incr(ctx, r.getKey(email))
	if err != nil {
		return 0, fmt.Errorf("failed to increment login attempts: %w", err)
	}

	count, err := r.replyToInt8(reply)
	if err != nil {
		return 0, err
	}

	if count == 1 {
		r.client.Expire(ctx, r.getKey(email), r.securityConfig.LoginAttemptsWindow())
	}

	return count, nil
}

func (r *redisRepository) ResetLoginAttempts(ctx context.Context, email string) error {
	err := r.client.Del(ctx, r.getKey(email))
	if err != nil {
		return fmt.Errorf("failed to reset login attempts: %w", err)
	}

	return nil
}

func (r *redisRepository) GetLoginAttempts(ctx context.Context, email string) (int8, error) {
	reply, err := r.client.Get(ctx, r.getKey(email))

	if err != nil {
		return 0, fmt.Errorf("failed to get login attempts: %w", err)
	}

	if reply == nil { // key not found
		return 0, nil
	}

	count, err := r.replyToInt8(reply)
	if err != nil {
		return 0, err
	}

	return count, nil
}
