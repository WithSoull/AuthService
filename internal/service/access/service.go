package access

import "github.com/WithSoull/AuthService/internal/service"

type accessService struct {
}

func (s *accessService) Check(endpoint_address string) error {
	return nil
}

func NewService() service.AccessService {
	return &accessService{}
}
