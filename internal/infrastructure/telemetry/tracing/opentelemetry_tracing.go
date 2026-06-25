package tracing

import (
	"context"
	"fmt"
	"strings"
	"time"

	domain_config "github.com/youknow2509/temp-go-ddd/internal/domain/config"
	"go.opentelemetry.io/contrib/propagators/jaeger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

// InitTracing initializes the OpenTelemetry TracerProvider with Jaeger OTLP exporter
func InitTracing(ctx context.Context, setting domain_config.TelemetryTracingSetting) (*sdktrace.TracerProvider, error) {
	if !setting.Enabled {
		return nil, nil
	}

	var client otlptrace.Client
	endpoint := setting.OtlpEndpoint

	// Parse options
	headers := setting.Headers
	// Determine if TLS is used based on scheme (standard OTLP default behavior)
	insecure := true
	if strings.HasPrefix(endpoint, "https://") {
		insecure = false
	}

	// Remove schema prefixes for OTLP gRPC/HTTP clients
	endpoint = strings.TrimPrefix(endpoint, "http://")
	endpoint = strings.TrimPrefix(endpoint, "https://")

	if setting.Protocol == "grpc" {
		grpcOpts := []otlptracegrpc.Option{
			otlptracegrpc.WithEndpoint(endpoint),
			otlptracegrpc.WithHeaders(headers),
		}
		if insecure {
			grpcOpts = append(grpcOpts, otlptracegrpc.WithInsecure())
		}
		client = otlptracegrpc.NewClient(grpcOpts...)
	} else {
		httpOpts := []otlptracehttp.Option{
			otlptracehttp.WithEndpoint(endpoint),
			otlptracehttp.WithHeaders(headers),
		}
		if insecure {
			httpOpts = append(httpOpts, otlptracehttp.WithInsecure())
		}
		client = otlptracehttp.NewClient(httpOpts...)
	}

	exporter, err := otlptrace.New(ctx, client)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP trace exporter: %w", err)
	}

	// Build Resource
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(setting.ServiceName),
			semconv.ServiceVersionKey.String(setting.ServiceVersion),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// Configure BatchSpanProcessor
	batchOpts := []sdktrace.BatchSpanProcessorOption{}
	if setting.Batch.MaxQueueSize > 0 {
		batchOpts = append(batchOpts, sdktrace.WithMaxQueueSize(setting.Batch.MaxQueueSize))
	}
	if setting.Batch.MaxExportBatchSize > 0 {
		batchOpts = append(batchOpts, sdktrace.WithMaxExportBatchSize(setting.Batch.MaxExportBatchSize))
	}
	if setting.Batch.ScheduleDelayMs > 0 {
		batchOpts = append(batchOpts, sdktrace.WithBatchTimeout(time.Duration(setting.Batch.ScheduleDelayMs)*time.Millisecond))
	}
	if setting.Batch.ExportTimeoutMs > 0 {
		batchOpts = append(batchOpts, sdktrace.WithExportTimeout(time.Duration(setting.Batch.ExportTimeoutMs)*time.Millisecond))
	}

	bsp := sdktrace.NewBatchSpanProcessor(exporter, batchOpts...)

	// Configure Sampler
	sampler := sdktrace.ParentBased(sdktrace.TraceIDRatioBased(setting.SamplingRatio))

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sampler),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)

	// Register the TracerProvider with the OpenTelemetry SDK
	otel.SetTracerProvider(tp)

	// Configure Propagators
	var propagators []propagation.TextMapPropagator
	for _, prop := range setting.Propagation {
		switch prop {
		case "tracecontext":
			propagators = append(propagators, propagation.TraceContext{})
		case "baggage":
			propagators = append(propagators, propagation.Baggage{})
		case "jaeger":
			propagators = append(propagators, jaeger.Jaeger{})
		}
	}

	if len(propagators) > 0 {
		otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagators...))
	} else {
		otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	}

	return tp, nil
}
