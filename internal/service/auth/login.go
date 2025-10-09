package auth

import (
	"context"
	"log"

	// VI = Validator Interceptor

	"github.com/WithSoull/AuthService/internal/model"
	conditions "github.com/WithSoull/AuthService/internal/validator"
	desc_user "github.com/WithSoull/UserServer/pkg/user/v1"
	"github.com/WithSoull/platform_common/pkg/sys"
	"github.com/WithSoull/platform_common/pkg/sys/codes"
	"github.com/WithSoull/platform_common/pkg/sys/validate"
)

func (s *authService) Login(ctx context.Context, email, password string) (string, error) {
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
	if attempts >= s.securityConfig.MaxLoginAttempts() {
		return "", sys.NewCommonError("Too many attempts", codes.ResourceExhausted)
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

	if !res.Valid {
		s.repository.IncrementLoginAttempts(ctx, email)
		return "", sys.NewCommonError("Ivalid email or password", codes.InvalidArgument)
	}

	// Reset attempts counter
	if err := s.repository.ResetLoginAttempts(ctx, email); err != nil {
		log.Printf("[Service Layer] failed to resert login attempts for %s: %v", email, err)
	}

	// Create refresh_token
	refresh_token, err := s.tokenGenerator.GenerateRefreshToken(model.UserInfo{
		UserId: res.GetUserId(),
		Email:  email,
	})
	if err != nil {
		log.Printf("[Service Layer] failed to generate refresh token: %v", err)
		return "", err
	}

	return refresh_token, nil
}
