package access

import (
	"context"

	"github.com/WithSoull/AuthService/internal/service"
)

type accessService struct {
}

func (s *accessService) Check(ctx context.Context, endpoint_address string) error {
	return nil
}

func NewService() service.AccessService {
	return &accessService{}
}
