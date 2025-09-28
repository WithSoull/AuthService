package auth

import (
	"github.com/WithSoull/AuthService/internal/service"
	desc "github.com/WithSoull/AuthService/pkg/auth/v1"
)

type authHandler struct {
	desc.UnimplementedAuthV1Server
	service service.AuthService
}

func NewHandler(service service.AuthService) desc.AuthV1Server {
	return &authHandler{
		service: service,
	}
}
