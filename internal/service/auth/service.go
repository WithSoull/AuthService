package auth

import (
	"github.com/WithSoull/AuthService/internal/client/grpc"
	"github.com/WithSoull/AuthService/internal/config"
	"github.com/WithSoull/AuthService/internal/repository"
	"github.com/WithSoull/AuthService/internal/service"
	"github.com/WithSoull/AuthService/internal/tokens"
)

type authService struct {
	userClient     grpc.UserClient
	tokenGenerator tokens.TokenGenerator
	repository     repository.AuthRepository
	securityConfig config.SecurityConfig
}

func NewService(userClient grpc.UserClient, tokenGenerator tokens.TokenGenerator, repository repository.AuthRepository) service.AuthService {
	return &authService{
		userClient:     userClient,
		tokenGenerator: tokenGenerator,
		repository:     repository,
	}
}
