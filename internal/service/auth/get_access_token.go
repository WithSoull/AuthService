package auth

import (
	"context"

	domainerrors "github.com/WithSoull/AuthService/internal/errors/domain"
	"github.com/WithSoull/AuthService/internal/model"
)

func (s *authService) GetAccessToken(ctx context.Context, refresh_token string) (string, error) {
	claims, err := s.tokenGenerator.VerifyRefreshToken(ctx, refresh_token)
	if err != nil {
		return "", domainerrors.ErrInvalidRefreshToken
	}

	new_access_token, err := s.tokenGenerator.GenerateAccessToken(ctx, model.UserInfo{
		UserId: claims.UserId,
		Email:  claims.Email,
	})
	if err != nil {
		return "", err
	}

	return new_access_token, nil
}
