package auth

import (
	"context"
	"log"

	domainerrors "github.com/WithSoull/AuthService/internal/errors/domain"
	"github.com/WithSoull/AuthService/internal/model"
)

func (s *authService) GetAccessToken(ctx context.Context, refresh_token string) (string, error) {
	claims, err := s.tokenGenerator.VerifyRefreshToken(refresh_token)
	if err != nil {
		_, grpcErr := domainerrors.ToGRPCStatus(domainerrors.ErrInvalidCredentials)
		return "", grpcErr
	}

	new_access_token, err := s.tokenGenerator.GenerateAccessToken(model.UserInfo{
		UserId: claims.UserId,
		Email:  claims.Email,
	})
	if err != nil {
		isLogNeeded, grpcErr := domainerrors.ToGRPCStatus(err)
		if isLogNeeded {
			log.Printf("[Service Layer] failed to generate access token with error: %v", err)
		}
		return "", grpcErr
	}

	return new_access_token, nil
}
