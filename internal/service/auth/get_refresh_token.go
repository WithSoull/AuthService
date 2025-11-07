package auth

import (
	"context"

	"github.com/WithSoull/AuthService/internal/model"
	"github.com/WithSoull/platform_common/pkg/contextx/claimsctx"
	"github.com/WithSoull/platform_common/pkg/logger"
	"github.com/WithSoull/platform_common/pkg/sys"
	"github.com/WithSoull/platform_common/pkg/sys/codes"
	"go.uber.org/zap"
)

func (s *authService) GetRefreshToken(ctx context.Context, old_refresh_token string) (string, error) {
	claims, err := s.tokenService.VerifyRefreshToken(ctx, old_refresh_token)
	if err != nil {
		return "", sys.NewCommonError("ivalid refresh token", codes.Unauthenticated)
	}

	ctx = claimsctx.InjectUserEmail(ctx, claims.Email)
	ctx = claimsctx.InjectUserID(ctx, claims.UserId)

	new_refresh_token, err := s.tokenService.GenerateRefreshToken(ctx, model.UserInfo{
		UserId: claims.UserId,
		Email:  claims.Email,
	})
	if err != nil {
		logger.Error(ctx, "failed to generate refresh token", zap.Error(err))
		return "", err
	}

	return new_refresh_token, nil
}
