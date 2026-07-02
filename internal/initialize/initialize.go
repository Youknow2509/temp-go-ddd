package initialize

import (
	"context"
	"net/http"
	"sync"

	domain_config "github.com/youknow2509/temp-go-ddd/internal/domain/config"
	domain_logger "github.com/youknow2509/temp-go-ddd/internal/domain/logger"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// AppResources contains resource instances that need explicit lifecycle management (e.g. shutdown)
type AppResources struct {
	Config         *domain_config.SystemConfig
	Logger         domain_logger.ILogger
	TracerProvider *sdktrace.TracerProvider
	MeterProvider  *sdkmetric.MeterProvider
	MetricServer   *http.Server
	Connections    *Connection
	Interfaces     *AppInterface
}

// ===
// Initialize the system
// ===
func Initialize(ctx context.Context, wg *sync.WaitGroup) (*AppResources, error) {
	// Initialize configurations system
	config, err := initializeConfig(ctx)
	if err != nil {
		return nil, err
	}
	// Initialize telemetry (logger, metric, tracing)
	logger, tp, mp, server, err := initializeTelemetry(ctx, wg, config)
	if err != nil {
		return nil, err
	}

	// Initialize connection
	connections, err := initializeConnection(ctx, config, logger)
	if err != nil {
		return nil, err
	}

	// Initialize services
	services, err := initializeAppService(config, logger, connections)
	if err != nil {
		return nil, err
	}

	// Initialize use cases
	if err := initializeUseCase(config, logger, services); err != nil {
		return nil, err
	}

	// Initialize pool
	if err := initializeAppPool(config, logger); err != nil {
		return nil, err
	}

	// Initialize interfaces
	interfaces, err := initializeInterfaces(config, logger, wg)
	if err != nil {
		return nil, err
	}

	resources := &AppResources{
		Config:         config,
		Logger:         logger,
		TracerProvider: tp,
		MeterProvider:  mp,
		MetricServer:   server,
		Connections:    connections,
		Interfaces:     interfaces,
	}

	return resources, nil
}


