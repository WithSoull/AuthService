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
	// TODO: REDIS

	// Validate user credentials
	if !utils.IsValidEmail(email) || password == "" {
		_, grpcErr := domainerrors.ToGRPCStatus(domainerrors.ErrInvalidCredentials)
		return "", grpcErr
	}

	res, err := s.userClient.ValidateCredentials(ctx, &desc_user.ValidateCredentialsRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		isLogNeeded, grpcErr := domainerrors.ToGRPCStatus(err)
		if isLogNeeded {
			log.Printf("[Service Layer] failed to get user credentials: %v", err)
		}
		return "", grpcErr
	}

	if !res.Valid {
		_, grpcErr := domainerrors.ToGRPCStatus(domainerrors.ErrInvalidCredentials)
		return "", grpcErr
	}

	// Create refresh_token
	refresh_token, err := s.tokenGenerator.GenerateRefreshToken(model.UserInfo{
		UserId: res.GetUserId(),
		Email:  email,
	})

	return refresh_token, nil
}
