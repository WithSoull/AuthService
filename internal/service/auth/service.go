package auth

import (
	"github.com/WithSoull/AuthService/internal/client/grpc"
	"github.com/WithSoull/AuthService/internal/config"
	"github.com/WithSoull/AuthService/internal/repository"
	"github.com/WithSoull/AuthService/internal/service"
	"github.com/WithSoull/platform_common/pkg/tokens"
)

type authService struct {
	userClient     grpc.UserClient
	tokenService   tokens.TokenService
	repository     repository.AuthRepository
	securityConfig config.SecurityConfig
}

func NewService(userClient grpc.UserClient, tokenService tokens.TokenService, repository repository.AuthRepository) service.AuthService {
	return &authService{
		userClient:   userClient,
		tokenService: tokenService,
		repository:   repository,
	}
}
