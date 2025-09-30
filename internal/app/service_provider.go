package app

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"log"
	"os"

	"github.com/WithSoull/AuthService/internal/config"
	"github.com/WithSoull/AuthService/internal/config/env"
	handler_auth "github.com/WithSoull/AuthService/internal/handler/auth"
	"github.com/WithSoull/AuthService/internal/service"
	service_auth "github.com/WithSoull/AuthService/internal/service/auth"
	"github.com/WithSoull/AuthService/internal/tokens"
	"github.com/WithSoull/AuthService/internal/tokens/jwt"
	access_v1 "github.com/WithSoull/AuthService/pkg/access/v1"
	auth_v1 "github.com/WithSoull/AuthService/pkg/auth/v1"
	desc_user "github.com/WithSoull/UserServer/pkg/user/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	servicePemPath = "service.pem"
)

type serviceProvider struct {
	grpcConfig     config.GRPCConfig
	userGrpcConfig config.GRPCConfig
	jwtConfig      config.JWTConfig

	authHandler   auth_v1.AuthV1Server
	accessHandler access_v1.AccessV1Server

	authService service.AuthService

	tokenGenerator tokens.TokenGenerator

	userClient desc_user.UserV1Client
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) GRPCConfig() config.GRPCConfig {
	if s.grpcConfig == nil {
		cfg, err := env.NewGRPCConfig()
		if err != nil {
			log.Fatalf("failed to get grpc config: %v", err)
		}

		s.grpcConfig = cfg
	}

	return s.grpcConfig
}

func (s *serviceProvider) UserGRPCConfig() config.GRPCConfig {
	if s.userGrpcConfig == nil {
		cfg, err := env.NewUserGRPCConfig()
		if err != nil {
			log.Fatalf("failed to get grpc config: %v", err)
		}

		s.userGrpcConfig = cfg
	}

	return s.userGrpcConfig
}

func (s *serviceProvider) JWTConfig() config.JWTConfig {
	if s.jwtConfig == nil {
		cfg, err := env.NewJWTConfig()
		if err != nil {
			log.Fatalf("failed to get jwt config: %v", err)
		}
		s.jwtConfig = cfg
	}

	return s.jwtConfig
}

func (s *serviceProvider) AuthService(ctx context.Context) service.AuthService {
	if s.authService == nil {
		s.authService = service_auth.NewService(s.UserClient(ctx), s.TokenGenerator(ctx))
	}
	return s.authService
}

func (s *serviceProvider) AuthHandler(ctx context.Context) auth_v1.AuthV1Server {
	if s.authHandler == nil {
		s.authHandler = handler_auth.NewHandler(s.AuthService(ctx))
	}
	return s.authHandler
}

func (s *serviceProvider) TokenGenerator(ctx context.Context) tokens.TokenGenerator {
	if s.tokenGenerator == nil {
		s.tokenGenerator = jwt.NewJWTService(s.JWTConfig())
	}
	return s.tokenGenerator
}

func (s *serviceProvider) UserClient(ctx context.Context) desc_user.UserV1Client {
	if s.userClient == nil {
		caCert, err := os.ReadFile("ca.cert")
		if err != nil {
			log.Fatalf("could not read ca certificate: %v", err)
		}

		certPool := x509.NewCertPool()
		if !certPool.AppendCertsFromPEM(caCert) {
			log.Fatalf("failed to append ca certificate")
		}

		tlsConfig := &tls.Config{
			ServerName: "localhost", // Должно совпадать с CN или SAN в сертификате
			RootCAs:    certPool,
		}

		creds := credentials.NewTLS(tlsConfig)

		conn, err := grpc.DialContext(ctx, s.UserGRPCConfig().Address(),
			grpc.WithTransportCredentials(creds))
		if err != nil {
			log.Fatalf("failed to dial gRPC server: %v", err)
		}

		// conn, err := grpc.DialContext(ctx, s.UserGRPCConfig().Address(), grpc.WithTransportCredentials(insecure.NewCredentials()))
		// if err != nil {
		// 	log.Fatalf("failed to dial gRPC server: %v", err)
		// }

		s.userClient = desc_user.NewUserV1Client(conn)
	}

	return s.userClient
}
