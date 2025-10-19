package app

import (
	"context"
	"flag"
	"net"
	"syscall"
	"time"

	"github.com/WithSoull/AuthService/internal/config"
	desc_auth "github.com/WithSoull/AuthService/pkg/auth/v1"
	"github.com/WithSoull/platform_common/pkg/closer"
	"github.com/WithSoull/platform_common/pkg/logger"
	"github.com/WithSoull/platform_common/pkg/metric"
	metricsInterceptor "github.com/WithSoull/platform_common/pkg/middleware/metrics"
	validationInterceptor "github.com/WithSoull/platform_common/pkg/middleware/validation"
	"github.com/WithSoull/platform_common/pkg/tracing"
	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

const (
	shutdownTimeout = 5 * time.Second
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

type App struct {
	serviceProvider *serviceProvider
	grpcServer      *grpc.Server
}

func NewApp(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		a.initLogger,
		a.initCloser,
		a.initServiceProvider,
		a.initMetrics,
		a.initGRPCServer,
		a.initTracing,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initConfig(_ context.Context) error {
	err := config.Load(configPath)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) initLogger(_ context.Context) error {
	return logger.Init(config.AppConfig().Logger)
}

func (a *App) initCloser(_ context.Context) error {
	closer.Configure(logger.Logger(), shutdownTimeout, syscall.SIGINT, syscall.SIGTERM)
	return nil
}

func (a *App) initServiceProvider(_ context.Context) error {
	a.serviceProvider = newServiceProvider()
	return nil
}

func (a *App) initMetrics(ctx context.Context) error {
	meterProvider, err := metric.InitOTELMetrics(config.AppConfig().Metrics)
	if err != nil {
		logger.Error(ctx, "failed to create meter provider", zap.Error(err))
	}

	closer.AddNamed("OTEL Metrics", meterProvider.Shutdown)

	if err := metric.Init(ctx, config.AppConfig().Metrics); err != nil {
		logger.Error(ctx, "failed to init metrics", zap.Error(err))
		return err
	}
	return nil
}

func (a *App) initGRPCServer(ctx context.Context) error {
	a.grpcServer = grpc.NewServer(
		grpc.Creds(insecure.NewCredentials()),
		grpc.UnaryInterceptor(
			grpcMiddleware.ChainUnaryServer(
				metricsInterceptor.MetricsInterceptor,
				validationInterceptor.ErrorCodesInterceptor(logger.Logger()),
				tracing.UnaryServerInterceptor(config.AppConfig().Tracing.ServiceName()),
			),
		),
	)

	closer.AddNamed("GRPC server", func(ctx context.Context) error {
		a.grpcServer.GracefulStop()
		return nil
	})

	reflection.Register(a.grpcServer)

	desc_auth.RegisterAuthV1Server(a.grpcServer, a.serviceProvider.AuthHandler(ctx))

	return nil
}

func (a *App) initTracing(ctx context.Context) error {
	err := tracing.InitTracer(ctx, config.AppConfig().Tracing)
	if err != nil {
		return err
	}

	closer.AddNamed("tracer", tracing.ShutdownTracer)

	return nil
}

func (a *App) runGRPCServer() error {
	lis, err := net.Listen("tcp", config.AppConfig().GRPC.Address())
	if err != nil {
		return err
	}

	logger.Info(context.Background(), "GRPC server listening", zap.String("address", config.AppConfig().GRPC.Address()))

	err = a.grpcServer.Serve(lis)
	if err != nil {
		return err
	}

	logger.Info(context.Background(), "GRPC server stopped gracefully")
	return nil
}

func (a *App) Run() error {
	return a.runGRPCServer()
}
