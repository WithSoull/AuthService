package auth

import "github.com/WithSoull/AuthService/internal/service"

type authService struct {
}

func (s *authService) Login(email, password string) string {
	return ""
}

func (s *authService) GetRefreshToken(refresh_token string) string {
	return ""
}

func (s *authService) GetAccessToken(refresh_token string) string {
	return ""
}

func NewService() service.AuthService {
	return &authService{}
}
