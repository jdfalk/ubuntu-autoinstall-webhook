// internal/observability/otel.go
package observability

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
)

var (
	tracerProvider trace.TracerProvider
	once           sync.Once
	shutdown       func(context.Context) error
	isInitialized  bool
)

// InitTracing initializes OpenTelemetry tracing
func InitTracing(ctx context.Context, serviceName string) error {
	var err error

	once.Do(func() {
		tracerMode := viper.GetString("observability.tracer.mode")

		// Create resource with service information
		res, resErr := resource.New(ctx,
			resource.WithAttributes(
				semconv.ServiceNameKey.String(serviceName),
				semconv.ServiceVersionKey.String(viper.GetString("version")),
				semconv.DeploymentEnvironmentKey.String(viper.GetString("environment")),
			),
		)

		if resErr != nil {
			err = resErr
			return
		}

		var traceProvider trace.TracerProvider

		switch tracerMode {
		case "otlp":
			// Configure OTLP exporter
			endpoint := viper.GetString("observability.tracer.otlp.endpoint")
			if endpoint == "" {
				endpoint = "localhost:4317"
			}

			client := otlptracegrpc.NewClient(
				otlptracegrpc.WithEndpoint(endpoint),
				otlptracegrpc.WithInsecure(),
			)

			exporter, exportErr := otlptrace.New(ctx, client)
			if exportErr != nil {
				err = exportErr
				return
			}

			bsp := sdktrace.NewBatchSpanProcessor(exporter)
			tp := sdktrace.NewTracerProvider(
				sdktrace.WithSampler(sdktrace.AlwaysSample()),
				sdktrace.WithResource(res),
				sdktrace.WithSpanProcessor(bsp),
			)

			traceProvider = tp
			shutdown = func(ctx context.Context) error {
				return tp.Shutdown(ctx)
			}

		case "stdout":
			// Configure stdout exporter (useful for debugging)
			exporter, exportErr := stdouttrace.New(stdouttrace.WithPrettyPrint())
			if exportErr != nil {
				err = exportErr
				return
			}

			bsp := sdktrace.NewBatchSpanProcessor(exporter)
			tp := sdktrace.NewTracerProvider(
				sdktrace.WithSampler(sdktrace.AlwaysSample()),
				sdktrace.WithResource(res),
				sdktrace.WithSpanProcessor(bsp),
			)

			traceProvider = tp
			shutdown = func(ctx context.Context) error {
				return tp.Shutdown(ctx)
			}

		case "noop", "":
			// No-op tracer for testing or when tracing is disabled
			traceProvider = noop.NewTracerProvider()
			shutdown = func(ctx context.Context) error { return nil }

		default:
			err = fmt.Errorf("unknown tracer mode: %s", tracerMode)
			return
		}

		// Set global tracer provider and propagator
		otel.SetTracerProvider(traceProvider)
		otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		))

		tracerProvider = traceProvider
		isInitialized = true
	})

	return err
}

// GetTracer returns a tracer for the given name
func GetTracer(name string) trace.Tracer {
	// If not initialized, return a no-op tracer
	if !isInitialized {
		return noop.NewTracerProvider().Tracer(name)
	}

	// Get tracer from the provider
	return tracerProvider.Tracer(name,
		trace.WithInstrumentationVersion("1.0.0"),
		trace.WithSchemaURL(semconv.SchemaURL),
	)
}

// Shutdown gracefully shuts down the tracer provider
func Shutdown(ctx context.Context) error {
	if shutdown != nil {
		// Use a timeout to ensure we don't hang forever
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		return shutdown(ctx)
	}
	return nil
}
