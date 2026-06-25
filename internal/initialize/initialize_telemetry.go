package initialize

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	domain_config "github.com/youknow2509/temp-go-ddd/internal/domain/config"
	domain_logger "github.com/youknow2509/temp-go-ddd/internal/domain/logger"
	infra_telemetry_logger "github.com/youknow2509/temp-go-ddd/internal/infrastructure/telemetry/logger"
	infra_telemetry_metric "github.com/youknow2509/temp-go-ddd/internal/infrastructure/telemetry/metric"
	infra_telemetry_tracing "github.com/youknow2509/temp-go-ddd/internal/infrastructure/telemetry/tracing"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// ===
// private function to initialize telemetry
// ===
func initializeTelemetry(ctx context.Context, wg *sync.WaitGroup, config *domain_config.SystemConfig) (domain_logger.ILogger, *sdktrace.TracerProvider, *sdkmetric.MeterProvider, *http.Server, error) {
	// Initialize logger
	logger, err := initializeLogger(config)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	// Initialize metrics
	mp, server, err := initializeMetrics(ctx, wg, config, logger)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	// Initialize tracing
	tp, err := initializeTracing(ctx, config, logger)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	return logger, tp, mp, server, nil
}

// ===
// private function to initialize the logger
// ===

func initializeLogger(config *domain_config.SystemConfig) (domain_logger.ILogger, error) {
	// Create logger configuration
	logger_config := domain_logger.NewLoggerConfigPkg(
		config.Telemetry.Logger,
	)
	// Create logger instance
	logger, err := infra_telemetry_logger.NewZapLogger(logger_config)
	if err != nil {
		return nil, err
	}
	logger.Debug("Logger initialized ...")

	return logger, nil
}

// ===
// private function to initialize metrics
// ===

func initializeMetrics(ctx context.Context, wg *sync.WaitGroup, config *domain_config.SystemConfig, logger domain_logger.ILogger) (*sdkmetric.MeterProvider, *http.Server, error) {
	mp, server, err := infra_telemetry_metric.InitMetrics(ctx, config.Telemetry.Metrics)
	if err != nil {
		return nil, nil, err
	}
	if mp != nil {
		// Track metrics HTTP server routine in wait group
		wg.Add(1)
		// Start HTTP server in goroutine
		go func() {
			defer wg.Done()
			logger.Info(fmt.Sprintf("Metrics exporter HTTP server listening on %s", server.Addr))
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				logger.Error(fmt.Sprintf("Metrics exporter HTTP server error: %v", err))
			}
		}()
	}
	return mp, server, nil
}

// ===
// private function to initialize tracing
// ===

func initializeTracing(ctx context.Context, config *domain_config.SystemConfig, logger domain_logger.ILogger) (*sdktrace.TracerProvider, error) {
	tp, err := infra_telemetry_tracing.InitTracing(ctx, config.Telemetry.Tracing)
	if err != nil {
		return nil, err
	}
	if tp != nil {
		logger.Debug("Tracing (OTLP/Jaeger) initialized ...")
	}

	return tp, nil
}
