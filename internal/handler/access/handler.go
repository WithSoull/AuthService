package access

import (
	"github.com/WithSoull/AuthService/internal/service"
	desc "github.com/WithSoull/AuthService/pkg/access/v1"
)

type accessHandler struct {
	desc.UnimplementedAccessV1Server
	service service.AccessService
}

func NewHandler(service service.AccessService) desc.AccessV1Server {
	return &accessHandler{
		service: service,
	}
}
