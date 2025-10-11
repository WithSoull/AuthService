package auth

import (
	"github.com/WithSoull/AuthService/internal/config"
	"github.com/WithSoull/AuthService/internal/repository"
	"github.com/WithSoull/AuthService/internal/service"
	"github.com/WithSoull/AuthService/internal/tokens"
	desc_user "github.com/WithSoull/UserServer/pkg/user/v1"
)

type authService struct {
	userClient     desc_user.UserV1Client
	tokenGenerator tokens.TokenGenerator
	repository     repository.AuthRepository
	securityConfig config.SecurityConfig
}

func NewService(userClient desc_user.UserV1Client, tokenGenerator tokens.TokenGenerator, repository repository.AuthRepository) service.AuthService {
	return &authService{
		userClient:     userClient,
		tokenGenerator: tokenGenerator,
		repository:     repository,
	}
}
