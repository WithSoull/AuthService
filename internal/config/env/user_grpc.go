package env

import (
	"errors"
	"net"
	"os"

	"github.com/WithSoull/AuthService/internal/config"
)

const (
	userGrpcHostEnvName = "USER_SERVER_GRPC_HOST"
	userGrpcPortEnvName = "USER_SERVER_GRPC_PORT"
)

type userGrpcConfig struct {
	host string
	port string
}

func NewUserGRPCConfig() (config.GRPCConfig, error) {
	host := os.Getenv(userGrpcHostEnvName)
	if len(host) == 0 {
		return nil, errors.New("grpc host not found")
	}

	port := os.Getenv(userGrpcPortEnvName)
	if len(port) == 0 {
		return nil, errors.New("grpc port not found")
	}

	return &grpcConfig{
		host: host,
		port: port,
	}, nil
}

func (cfg *userGrpcConfig) Address() string {
	address := net.JoinHostPort(cfg.host, cfg.port)
	return address
}
