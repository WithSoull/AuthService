package auth

import (
	"context"

	"github.com/WithSoull/AuthService/internal/config"
	domainerrors "github.com/WithSoull/AuthService/internal/errors/domain"
	"github.com/WithSoull/AuthService/internal/model"
	conditions "github.com/WithSoull/AuthService/internal/validator"
	"github.com/WithSoull/platform_common/pkg/contextx/claimsctx"
	"github.com/WithSoull/platform_common/pkg/logger"
	"github.com/WithSoull/platform_common/pkg/sys/validate"
	"github.com/WithSoull/platform_common/pkg/tracing"
	"go.opentelemetry.io/otel/attribute"
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
	spanName := "user.ValidateCredentials"
	ctx, userSpan := tracing.StartSpan(ctx, spanName)
	res, err := s.userClient.ValidateCredentials(ctx, email, password)
	if err != nil {
		userSpan.RecordError(err)
		userSpan.End()
		s.repository.IncrementLoginAttempts(ctx, email)
		return "", err
	}
	userSpan.SetAttributes(
		attribute.String("email", email),
		attribute.Int("attempts", int(attempts)),
	)

	userSpan.End()

	ctx = claimsctx.InjectUserID(ctx, res.UserID)

	if !res.Valid {
		s.repository.IncrementLoginAttempts(ctx, email)
		return "", domainerrors.ErrInvalidEmailOrPassword
	}

	// Reset attempts counter
	if err := s.repository.ResetLoginAttempts(ctx, email); err != nil {
		logger.Error(ctx, "failed to resert login attempts", zap.Error(err))
	}

	// Create refresh_token
	refresh_token, err := s.tokenService.GenerateRefreshToken(ctx, model.UserInfo{
		UserId: res.UserID,
		Email:  email,
	})
	if err != nil {
		logger.Error(ctx, "failed to generate refresh token", zap.Error(err))
		return "", err
	}

	return refresh_token, nil
}
