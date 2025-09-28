package app

import (
	"context"
	"log"

	"github.com/WithSoull/AuthService/internal/config"
	"github.com/WithSoull/AuthService/internal/config/env"
	handler_access "github.com/WithSoull/AuthService/internal/handler/access"
	handler_auth "github.com/WithSoull/AuthService/internal/handler/auth"
	"github.com/WithSoull/AuthService/internal/service"
	service_access "github.com/WithSoull/AuthService/internal/service/access"
	service_auth "github.com/WithSoull/AuthService/internal/service/auth"
	access_v1 "github.com/WithSoull/AuthService/pkg/access/v1"
	auth_v1 "github.com/WithSoull/AuthService/pkg/auth/v1"
)

type serviceProvider struct {
	grpcConfig config.GRPCConfig

	authHandler   auth_v1.AuthV1Server
	accessHandler access_v1.AccessV1Server

	authService   service.AuthService
	accessService service.AccessService
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) GRPCConfig() config.GRPCConfig {
	if s.grpcConfig == nil {
		cfg, err := env.NewGRPCConfig()
		if err != nil {
			log.Fatalf("failed to get grpc config: %s", err.Error())
		}

		s.grpcConfig = cfg
	}

	return s.grpcConfig
}

func (s *serviceProvider) AuthService(ctx context.Context) service.AuthService {
	if s.authService == nil {
		s.authService = service_auth.NewService()
	}
	return s.authService
}

func (s *serviceProvider) AccessService(ctx context.Context) service.AccessService {
	if s.accessService == nil {
		s.accessService = service_access.NewService()
	}
	return s.accessService
}

func (s *serviceProvider) AuthHandler(ctx context.Context) auth_v1.AuthV1Server {
	if s.authHandler == nil {
		s.authHandler = handler_auth.NewHandler(s.AuthService(ctx))
	}
	return s.authHandler
}

func (s *serviceProvider) AccessHandler(ctx context.Context) access_v1.AccessV1Server {
	if s.accessHandler == nil {
		s.accessHandler = handler_access.NewHandler(s.AccessService(ctx))
	}
	return s.accessHandler
}
