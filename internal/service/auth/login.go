package auth

import (
	"context"
	"log"

	domainerrors "github.com/WithSoull/AuthService/internal/errors/domain"
	"github.com/WithSoull/AuthService/internal/model"
	"github.com/WithSoull/AuthService/internal/utils"
	desc_user "github.com/WithSoull/UserServer/pkg/user/v1"
)

func (s *authService) Login(ctx context.Context, email, password string) (string, error) {
	// rate limiting validation
	attempts, err := s.repository.GetLoginAttempts(ctx, email)
	if err != nil {
		isNeedToLog, grpcErr := domainerrors.ToGRPCStatus(err)
		if isNeedToLog {
			log.Printf("failed to get attempts: %v", err)
		}
		return "", grpcErr
	}
	if attempts >= s.securityConfig.MaxLoginAttempts() {
		_, grpcErr := domainerrors.ToGRPCStatus(domainerrors.ErrTooManyAttempts)
		return "", grpcErr
	}

	// Validate user credentials
	if !utils.IsValidEmail(email) || password == "" {
		s.repository.IncrementLoginAttempts(ctx, email)
		_, grpcErr := domainerrors.ToGRPCStatus(domainerrors.ErrInvalidCredentials)
		return "", grpcErr
	}

	// Check credentials
	res, err := s.userClient.ValidateCredentials(ctx, &desc_user.ValidateCredentialsRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		s.repository.IncrementLoginAttempts(ctx, email)
		isLogNeeded, grpcErr := domainerrors.ToGRPCStatus(err)
		if isLogNeeded {
			log.Printf("[Service Layer] failed to get user credentials: %v", err)
		}
		return "", grpcErr
	}

	if !res.Valid {
		s.repository.IncrementLoginAttempts(ctx, email)
		_, grpcErr := domainerrors.ToGRPCStatus(domainerrors.ErrInvalidCredentials)
		return "", grpcErr
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
		_, grpcErr := domainerrors.ToGRPCStatus(err)
		return "", grpcErr
	}

	return refresh_token, nil
}
