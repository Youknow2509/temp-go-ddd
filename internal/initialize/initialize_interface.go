package initialize

import (
	"sync"

	domain_config "github.com/youknow2509/temp-go-ddd/internal/domain/config"
	domain_logger "github.com/youknow2509/temp-go-ddd/internal/domain/logger"
	interface_grpc "github.com/youknow2509/temp-go-ddd/internal/interface/grpc"
	interface_http "github.com/youknow2509/temp-go-ddd/internal/interface/http"
	interface_http_router "github.com/youknow2509/temp-go-ddd/internal/interface/http/router"
	"google.golang.org/grpc"
)

// AppInterface defines all interfaces for system: http - gin, grpc - grpc-go, ...
type AppInterface struct {
	GrpcServer *interface_grpc.Server
	HttpServer *interface_http.Server
}

// InitializeInterfaces initializes the application interfaces and returns an instance of AppInterface.
func initializeInterfaces(
	config *domain_config.SystemConfig,
	logger domain_logger.ILogger,
	waitGroup *sync.WaitGroup,
) (*AppInterface, error) {
	// Initialize Grpc server
	grpcServer, err := initializeGrpcServer(logger, config, waitGroup)
	if err != nil {
		return nil, err
	}

	// Initialize Http server
	httpServer, err := initializeHttpServer(logger, config, waitGroup)
	if err != nil {
		return nil, err
	}

	logger.Info("All interfaces initialized successfully")
	return &AppInterface{
		GrpcServer: grpcServer,
		HttpServer: httpServer,
	}, nil
}

// InitializeHttpServer initializes the HTTP server and returns an instance of HttpServer.
func initializeHttpServer(
	logger domain_logger.ILogger,
	config *domain_config.SystemConfig,
	waitGroup *sync.WaitGroup,
) (*interface_http.Server, error) {
	srv, err := interface_http.New(
		config.HttpServer,
		logger,
		interface_http_router.RegisterRoutes,
	)
	if err != nil {
		logger.Error("create http server: ", err)
		return nil, err
	}
	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()
		if err := srv.Start(); err != nil {
			logger.Error("http server error: ", err)
		}
	}()
	logger.Info("HTTP server initialized successfully")
	return srv, nil
}

// InitializeGrpcServer initializes the gRPC server and returns an instance of GrpcServer.
func initializeGrpcServer(
	logger domain_logger.ILogger,
	config *domain_config.SystemConfig,
	waitGroup *sync.WaitGroup,
) (*interface_grpc.Server, error) {
	srv, err := interface_grpc.New(
		config.GrpcServer,
		logger,
		func(s *grpc.Server) {
			// Register custom services here if any
		},
	)
	if err != nil {
		logger.Error("create grpc server: ", err)
		return nil, err
	}
	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()
		if err := srv.Start(); err != nil {
			logger.Error("grpc server error: ", err)
		}
	}()
	logger.Info("gRPC server initialized successfully")
	return srv, nil
}
