package metric

import (
	"context"
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	domain_config "github.com/youknow2509/temp-go-ddd/internal/domain/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"
)

// InitMetrics initializes the OpenTelemetry MeterProvider with Prometheus exporter and starts HTTP server
func InitMetrics(ctx context.Context, setting domain_config.TelemetryMetricsSetting) (*metric.MeterProvider, *http.Server, error) {
	if !setting.Enabled {
		return nil, nil, nil
	}

	// Create Prometheus exporter
	exporter, err := prometheus.New()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create prometheus exporter: %w", err)
	}

	// Create MeterProvider with Prometheus reader
	mp := metric.NewMeterProvider(
		metric.WithReader(exporter),
	)

	// Register the MeterProvider with the OpenTelemetry SDK
	otel.SetMeterProvider(mp)

	// Start Prometheus HTTP server
	mux := http.NewServeMux()
	mux.Handle(setting.Path, promhttp.Handler())

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", setting.Host, setting.Port),
		Handler: mux,
	}

	return mp, server, nil
}
