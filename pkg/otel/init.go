package otel

import (
	"context"
	"fmt"
	"globe-and-citizen/layer8/auth-server/pkg/utils"
	"net"
	"net/url"
	"strings"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace/noop"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

// InitTracer initializes OpenTelemetry with configurable exporter
// Supports: stdout (development) and OTLP (Jaeger, Datadog, generic OTLP collectors)
func InitTracer(serviceName string, cfg Config) (func(context.Context) error, error) {
	if !cfg.Enabled {
		return func(ctx context.Context) error { return nil }, nil
	}

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
			cfg.Protocol = "grpc" // Default to gRPC
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

	res, err := resource.New(context.Background(), resource.WithAttributes(
		attribute.String("service.name", serviceName),
		attribute.String("environment", utils.GetEnvironment()),
	))
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

// initOTLPExporterHTTP creates an HTTP/JSON OTLP exporter
func initOTLPExporterHTTP(endpoint string, authHeader string) (trace.SpanExporter, error) {
	if endpoint == "" {
		endpoint = "http://localhost:4318"
	}

	endpoint = sanitizeEndpoint(endpoint)
	if err := validateOTLPReachability("tcp", endpoint, 3*time.Second); err != nil {
		return nil, fmt.Errorf("OTLP HTTP exporter unreachable at %s: %w", endpoint, err)
	}

	opts := []otlptracehttp.Option{
		otlptracehttp.WithEndpoint(endpoint),
		otlptracehttp.WithInsecure(),
	}

	if authHeader != "" {
		opts = append(opts, otlptracehttp.WithHeaders(map[string]string{
			"Authorization": authHeader,
		}))
	}

	exporter, err := otlptracehttp.New(context.Background(), opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP HTTP exporter (endpoint: %s): %w", endpoint, err)
	}

	return exporter, nil
}

// initOTLPExporterGRPC creates a gRPC OTLP exporter
func initOTLPExporterGRPC(endpoint string, authHeader string) (trace.SpanExporter, error) {
	if endpoint == "" {
		endpoint = "localhost:4317"
	}

	endpoint = sanitizeEndpoint(endpoint)
	if err := validateOTLPReachability("tcp", endpoint, 3*time.Second); err != nil {
		return nil, fmt.Errorf("OTLP gRPC exporter unreachable at %s: %w", endpoint, err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection to %s: %w", endpoint, err)
	}

	opts := []otlptracegrpc.Option{
		otlptracegrpc.WithGRPCConn(conn),
	}

	exporter, err := otlptracegrpc.New(context.Background(), opts...)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to create OTLP gRPC exporter (endpoint: %s): %w", endpoint, err)
	}

	return exporter, nil
}

func validateOTLPReachability(network, address string, timeout time.Duration) error {
	if address == "" {
		return fmt.Errorf("empty OTLP endpoint")
	}

	conn, err := net.DialTimeout(network, address, timeout)
	if err != nil {
		return err
	}
	defer conn.Close()
	return nil
}

// sanitizeEndpoint validates and cleans up the OTLP endpoint
// Removes protocol, paths, and trailing slashes to get just host:port
func sanitizeEndpoint(endpoint string) string {
	endpoint = strings.TrimSpace(endpoint)

	// Parse URL to extract host:port
	if !strings.Contains(endpoint, "://") {
		// No scheme provided, assume http
		endpoint = "http://" + endpoint
	}

	u, err := url.Parse(endpoint)
	if err != nil {
		// If parsing fails, return as-is and let the exporter handle the error
		return endpoint
	}

	// Extract just host:port (remove path, query, fragment)
	if u.Host == "" {
		return endpoint
	}

	return u.Host
}

// GetTracer returns a global tracer from the OpenTelemetry provider
func GetTracer(name string) interface{} {
	return otel.GetTracerProvider().Tracer(name)
}
