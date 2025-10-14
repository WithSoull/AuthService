package auth

import (
	"context"

	"github.com/WithSoull/AuthService/internal/config"
	domainerrors "github.com/WithSoull/AuthService/internal/errors/domain"
	"github.com/WithSoull/AuthService/internal/model"
	conditions "github.com/WithSoull/AuthService/internal/validator"
	desc_user "github.com/WithSoull/UserServer/pkg/user/v1"
	"github.com/WithSoull/platform_common/pkg/contextx/claimsctx"
	"github.com/WithSoull/platform_common/pkg/logger"
	"github.com/WithSoull/platform_common/pkg/sys/validate"
	"go.uber.org/zap"
)

func (s *authService) Login(ctx context.Context, email, password string) (string, error) {
	ctx = claimsctx.InjectUserEmail(ctx, email)

	err := validate.Validate(
		ctx,
		conditions.ValidateNotEmptyEmailAndPassword(email, password),
	)
	if err != nil {
		return "", err
	}

	// rate limiting validation
	attempts, err := s.repository.GetLoginAttempts(ctx, email)
	if err != nil {
		return "", err
	}
	if attempts >= config.AppConfig().Security.MaxLoginAttempts() {
		return "", domainerrors.ErrTooManyAttempts
	}

	// Check credentials
	res, err := s.userClient.ValidateCredentials(ctx, &desc_user.ValidateCredentialsRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		s.repository.IncrementLoginAttempts(ctx, email)
		return "", err
	}

	ctx = claimsctx.InjectUserID(ctx, res.UserId)

	if !res.Valid {
		s.repository.IncrementLoginAttempts(ctx, email)
		return "", domainerrors.ErrInvalidEmailOrPassword
	}

	// Reset attempts counter
	if err := s.repository.ResetLoginAttempts(ctx, email); err != nil {
		logger.Error(ctx, "failed to resert login attempts", zap.Error(err))
	}

	// Create refresh_token
	refresh_token, err := s.tokenGenerator.GenerateRefreshToken(ctx, model.UserInfo{
		UserId: res.GetUserId(),
		Email:  email,
	})
	if err != nil {
		logger.Error(ctx, "failed to generate refresh token", zap.Error(err))
		return "", err
	}

	return refresh_token, nil
}
