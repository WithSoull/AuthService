package auth

import (
	"context"
	"log"

	"github.com/WithSoull/AuthService/internal/model"
	"github.com/WithSoull/platform_common/pkg/sys"
	"github.com/WithSoull/platform_common/pkg/sys/codes"
)

func (s *authService) GetRefreshToken(ctx context.Context, old_refresh_token string) (string, error) {
	claims, err := s.tokenGenerator.VerifyRefreshToken(old_refresh_token)
	if err != nil {
		return "", sys.NewCommonError("ivalid refresh token", codes.Unauthenticated)
	}

	new_refresh_token, err := s.tokenGenerator.GenerateRefreshToken(model.UserInfo{
		UserId: claims.UserId,
		Email:  claims.Email,
	})
	if err != nil {
		log.Printf("[Service Layer] failed to generate refresh token with error: %v", err)
		return "", err
	}

	return new_refresh_token, nil
}
