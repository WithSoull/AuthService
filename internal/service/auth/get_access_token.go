package auth

import (
	"context"

	"github.com/WithSoull/AuthService/internal/model"
	"github.com/WithSoull/platform_common/pkg/sys"
	"github.com/WithSoull/platform_common/pkg/sys/codes"
)

func (s *authService) GetAccessToken(ctx context.Context, refresh_token string) (string, error) {
	claims, err := s.tokenGenerator.VerifyRefreshToken(refresh_token)
	if err != nil {
		return "", sys.NewCommonError("ivalid refresh token", codes.Unauthenticated)
	}

	new_access_token, err := s.tokenGenerator.GenerateAccessToken(model.UserInfo{
		UserId: claims.UserId,
		Email:  claims.Email,
	})
	if err != nil {
		return "", err
	}

	return new_access_token, nil
}
