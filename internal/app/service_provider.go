package app

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"os"

	"github.com/WithSoull/AuthService/internal/client/cache"
	redis_client "github.com/WithSoull/AuthService/internal/client/cache/redis"
	"github.com/WithSoull/AuthService/internal/config"
	handler_auth "github.com/WithSoull/AuthService/internal/handler/auth"
	"github.com/WithSoull/AuthService/internal/repository"
	"github.com/WithSoull/AuthService/internal/repository/auth"
	"github.com/WithSoull/AuthService/internal/service"
	service_auth "github.com/WithSoull/AuthService/internal/service/auth"
	"github.com/WithSoull/AuthService/internal/tokens"
	"github.com/WithSoull/AuthService/internal/tokens/jwt"
	access_v1 "github.com/WithSoull/AuthService/pkg/access/v1"
	auth_v1 "github.com/WithSoull/AuthService/pkg/auth/v1"
	desc_user "github.com/WithSoull/UserServer/pkg/user/v1"
	"github.com/WithSoull/platform_common/pkg/closer"
	"github.com/WithSoull/platform_common/pkg/logger"
	"github.com/gomodule/redigo/redis"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	servicePemPath = "service.pem"
)

type serviceProvider struct {
	authHandler   auth_v1.AuthV1Server
	accessHandler access_v1.AccessV1Server

	authService    service.AuthService
	authRepository repository.AuthRepository
	cacheClient    cache.CacheClient

	tokenGenerator tokens.TokenGenerator

	redisPool  *redis.Pool
	userClient desc_user.UserV1Client
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) RedisPool() *redis.Pool {
	if s.redisPool == nil {
		redisPool := &redis.Pool{
			MaxIdle:     int(config.AppConfig().Redis.MaxIdle()),
			IdleTimeout: config.AppConfig().Redis.IdleTimeout(),
			DialContext: func(ctx context.Context) (redis.Conn, error) {
				return redis.DialContext(ctx, "tcp", config.AppConfig().Redis.Address())
			},
		}

		closer.AddNamed("RedisPool", func(ctx context.Context) error {
			return redisPool.Close()
		})

		s.redisPool = redisPool
	}

	return s.redisPool
}

func (s *serviceProvider) CacheClient() cache.CacheClient {
	if s.cacheClient == nil {
		s.cacheClient = redis_client.NewClient(s.RedisPool())
	}

	return s.cacheClient
}

func (s *serviceProvider) AuthRepository() repository.AuthRepository {
	if s.authRepository == nil {
		s.authRepository = auth.NewRedisRepository(s.CacheClient())
	}

	return s.authRepository
}

func (s *serviceProvider) AuthService(ctx context.Context) service.AuthService {
	if s.authService == nil {
		s.authService = service_auth.NewService(s.UserClient(ctx), s.TokenGenerator(ctx), s.AuthRepository())
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
		s.tokenGenerator = jwt.NewJWTService()
	}
	return s.tokenGenerator
}

func (s *serviceProvider) UserClient(ctx context.Context) desc_user.UserV1Client {
	if s.userClient == nil {
		caCert, err := os.ReadFile("ca.cert")
		if err != nil {
			logger.Fatal(ctx, "could not read ca certificate", zap.Error(err))
		}

		certPool := x509.NewCertPool()
		if !certPool.AppendCertsFromPEM(caCert) {
			logger.Fatal(ctx, "failed to append ca certificate")
		}

		tlsConfig := &tls.Config{
			ServerName: "localhost",
			RootCAs:    certPool,
		}

		creds := credentials.NewTLS(tlsConfig)

		conn, err := grpc.DialContext(ctx, config.AppConfig().UserClient.Address(), grpc.WithTransportCredentials(creds))
		if err != nil {
			logger.Fatal(ctx, "failed to dial gRPC server", zap.Error(err))
		}

		closer.AddNamed("UserClientGRPC", func(ctx context.Context) error {
			return conn.Close()
		})

		logger.Debug(ctx, "Succesfully create UserServer client", zap.Any("connection", conn))
		s.userClient = desc_user.NewUserV1Client(conn)
	}

	return s.userClient
}
