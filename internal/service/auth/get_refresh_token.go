package auth

import (
	"context"
	"log"

	domainerrors "github.com/WithSoull/AuthService/internal/errors/domain"
	"github.com/WithSoull/AuthService/internal/model"
)

func (s *authService) GetRefreshToken(ctx context.Context, old_refresh_token string) (string, error) {
	claims, err := s.tokenGenerator.VerifyRefreshToken(old_refresh_token)
	if err != nil {
		_, grpcErr := domainerrors.ToGRPCStatus(domainerrors.ErrInvalidCredentials)
		return "", grpcErr
	}

	// TODO: Check if refresh token exists in Redis and is not blacklisted
	// exists, err := s.redisClient.RefreshTokenExists(ctx, claims.UserId, refreshToken)
	// if err != nil || !exists {
	//     log.Printf("[Service Layer] refresh token not found in Redis for user: %s", claims.UserId)
	//     _, grpcErr := domainerrors.ToGRPCStatus(domainerrors.ErrInvalidCredentials)
	//     return nil, grpcErr
	// }

	new_refresh_token, err := s.tokenGenerator.GenerateRefreshToken(model.UserInfo{
		UserId: claims.UserId,
		Email:  claims.Email,
	})
	if err != nil {
		isLogNeeded, grpcErr := domainerrors.ToGRPCStatus(err)
		if isLogNeeded {
			log.Printf("[Service Layer] failed to generate refresh token with error: %v", err)
		}
		return "", grpcErr
	}

	// TODO: Redis operations
	// 1. Blacklist old refresh token
	// err = s.redisClient.BlacklistRefreshToken(ctx, claims.UserId, refreshToken)
	// if err != nil {
	//     log.Printf("[Service Layer] failed to blacklist old refresh token: %v", err)
	// }
	//
	// 2. Save new refresh token
	// err = s.redisClient.SaveRefreshToken(ctx, claims.UserId, tokenPair.RefreshToken)
	// if err != nil {
	//     log.Printf("[Service Layer] failed to save new refresh token: %v", err)
	// }

	return new_refresh_token, nil
}
