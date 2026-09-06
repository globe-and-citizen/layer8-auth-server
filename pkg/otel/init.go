package otel

import (
	"context"
	"fmt"
	"globe-and-citizen/layer8/auth-server/pkg/utils"
	"net/url"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace/noop"
)

// Config holds OpenTelemetry configuration
type Config struct {
	Enabled      bool    `env:"OTEL_ENABLED" env-default:"true"`                   // Enable tracing
	ExporterType string  `env:"OTEL_EXPORTER_TYPE" env-default:"stdout"`           // "stdout" or "otlp"
	Protocol     string  `env:"OTEL_PROTOCOL" env-default:"grpc"`                  // "http" or "grpc" (for otlp)
	Endpoint     string  `env:"OTEL_ENDPOINT" env-default:"http://localhost:4317"` // works with Jaeger, Datadog, etc.
	AuthHeader   string  `env:"OTEL_AUTH_HEADER" env-default:""`                   // Optional auth header
	SamplingRate float64 `env:"OTEL_SAMPLING_RATE" env-default:"1.0"`              // 0.0 to 1.0
}

// InitTracer initializes OpenTelemetry with configurable exporter.
// Supports: stdout (development) and OTLP (Jaeger, Datadog, generic OTLP collectors).
func InitTracer(serviceName string, cfg Config) (func(context.Context) error, error) {
	if !cfg.Enabled {
		return func(ctx context.Context) error { return nil }, nil
	}

	// Configure W3C Trace Context propagation.
	// This is used to extract traceparent/tracestate from incoming requests
	// and inject them into outgoing requests.
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	var exporter trace.SpanExporter
	var err error

	switch strings.ToLower(cfg.ExporterType) {
	case "stdout":
		exporter, err = stdouttrace.New(stdouttrace.WithPrettyPrint())
		if err != nil {
			fmt.Printf("warning: failed to create stdout exporter: %v; tracing disabled\n", err)
			otel.SetTracerProvider(noop.NewTracerProvider())

			return func(ctx context.Context) error { return nil }, nil
		}

	case "otlp":
		if cfg.Protocol == "" {
			cfg.Protocol = "grpc"
		}

		switch strings.ToLower(cfg.Protocol) {
		case "http":
			exporter, err = initOTLPExporterHTTP(cfg.Endpoint, cfg.AuthHeader)
		case "grpc":
			exporter, err = initOTLPExporterGRPC(cfg.Endpoint, cfg.AuthHeader)
		default:
			fmt.Printf("warning: unknown OTLP protocol %q; tracing disabled\n", cfg.Protocol)
			otel.SetTracerProvider(noop.NewTracerProvider())

			return func(ctx context.Context) error { return nil }, nil
		}

		if err != nil {
			fmt.Printf("warning: OTLP exporter initialization failed: %v; tracing disabled\n", err)
			otel.SetTracerProvider(noop.NewTracerProvider())

			return func(ctx context.Context) error { return nil }, nil
		}

	default:
		fmt.Printf("warning: unknown exporter type %q; tracing disabled\n", cfg.ExporterType)
		otel.SetTracerProvider(noop.NewTracerProvider())

		return func(ctx context.Context) error { return nil }, nil
	}

	res, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			attribute.String("service.name", serviceName),
			attribute.String("environment", utils.GetEnvironment()),
		),
	)
	if err != nil {
		fmt.Printf("warning: failed to create resource: %v; tracing disabled\n", err)
		otel.SetTracerProvider(noop.NewTracerProvider())

		return func(ctx context.Context) error { return nil }, nil
	}

	sampler := trace.TraceIDRatioBased(cfg.SamplingRate)

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(res),
		trace.WithSampler(sampler),
	)

	otel.SetTracerProvider(tp)

	return tp.Shutdown, nil
}

// initOTLPExporterHTTP creates an HTTP/JSON OTLP exporter.
//
// The exporter is created without checking whether the collector is currently
// reachable. Temporary collector outages should not prevent the application
// from starting. The exporter will attempt to export when the collector is
// available.
func initOTLPExporterHTTP(
	endpoint string,
	authHeader string,
) (trace.SpanExporter, error) {
	if endpoint == "" {
		endpoint = "http://localhost:4318"
	}

	endpoint = sanitizeEndpoint(endpoint)

	opts := []otlptracehttp.Option{
		otlptracehttp.WithEndpoint(endpoint),
		otlptracehttp.WithInsecure(),
	}

	if authHeader != "" {
		opts = append(opts, otlptracehttp.WithHeaders(
			map[string]string{
				"Authorization": authHeader,
			},
		))
	}

	exporter, err := otlptracehttp.New(context.Background(), opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP HTTP exporter (endpoint: %s): %w", endpoint, err)
	}

	return exporter, nil
}

// initOTLPExporterGRPC creates a gRPC OTLP exporter.
//
// IMPORTANT:
// Do not use grpc.WithBlock() here. The application should not block startup
// waiting for the collector. gRPC ClientConn can reconnect when the remote
// endpoint becomes available again.
func initOTLPExporterGRPC(
	endpoint string,
	authHeader string,
) (trace.SpanExporter, error) {
	if endpoint == "" {
		endpoint = "localhost:4317"
	}

	endpoint = sanitizeEndpoint(endpoint)

	opts := []otlptracegrpc.Option{
		otlptracegrpc.WithEndpoint(endpoint),
		otlptracegrpc.WithInsecure(),
	}

	// NOTE:
	// otlptracegrpc does not accept an arbitrary Authorization header directly
	// through WithHeaders in the same way as the HTTP exporter.
	//
	// If your collector requires Authorization, configure an appropriate
	// gRPC credentials/per-RPC credentials mechanism here.

	_ = authHeader

	exporter, err := otlptracegrpc.New(
		context.Background(),
		opts...,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP gRPC exporter (endpoint: %s): %w", endpoint, err)
	}

	return exporter, nil
}

// sanitizeEndpoint validates and cleans up an OTLP endpoint.
//
// For gRPC this returns host:port.
// For HTTP this returns host:port.
//
// Example:
//
//	localhost:4317
//	http://localhost:4317
//	https://collector.example.com:4317
//
// become:
//
//	localhost:4317
//	collector.example.com:4317
//	collector.example.com:4317
func sanitizeEndpoint(endpoint string) string {
	endpoint = strings.TrimSpace(endpoint)

	if endpoint == "" {
		return endpoint
	}

	// Add a scheme temporarily so url.Parse can correctly identify the host.
	parseEndpoint := endpoint
	if !strings.Contains(parseEndpoint, "://") {
		parseEndpoint = "http://" + parseEndpoint
	}

	u, err := url.Parse(parseEndpoint)
	if err != nil {
		return endpoint
	}

	if u.Host == "" {
		return endpoint
	}

	return u.Host
}

// GetTracer returns a global tracer from the OpenTelemetry provider
func GetTracer(name string) interface{} {
	return otel.GetTracerProvider().Tracer(name)
}
