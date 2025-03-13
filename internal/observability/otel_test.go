// internal/observability/otel.go
package observability

import (
	"context"
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	// Setup test configuration
	viper.SetDefault("observability.tracer.mode", "noop")
	viper.SetDefault("version", "test")
	viper.SetDefault("environment", "test")

	os.Exit(m.Run())
}

func TestGetTracerWithoutInit(t *testing.T) {
	// Reset the singleton state
	tracerProvider = nil
	isInitialized = false

	// GetTracer should return a no-op tracer when not initialized
	tracer := GetTracer("test-service")
	assert.NotNil(t, tracer)

	// Create a span and verify it's a no-op span
	_, span := tracer.Start(context.Background(), "test-span")
	assert.NotNil(t, span)
	span.End()
}

func TestInitTracing(t *testing.T) {
	// Reset the singleton state
	tracerProvider = nil
	isInitialized = false

	// Test with noop mode
	viper.Set("observability.tracer.mode", "noop")
	err := InitTracing(context.Background(), "test-service")
	assert.NoError(t, err)
	assert.True(t, isInitialized)

	// Get tracer should return a non-nil tracer
	tracer := GetTracer("test-component")
	assert.NotNil(t, tracer)

	// Should be able to create spans
	_, span := tracer.Start(context.Background(), "test-operation")
	assert.NotNil(t, span)
	span.End()

	// Test shutdown
	err = Shutdown(context.Background())
	assert.NoError(t, err)
}

func TestUnsupportedTracerMode(t *testing.T) {
	// Reset the singleton state
	tracerProvider = nil
	isInitialized = false

	// Test with invalid mode
	viper.Set("observability.tracer.mode", "invalid-mode")
	err := InitTracing(context.Background(), "test-service")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown tracer mode")
}

func TestTracerIsSingleton(t *testing.T) {
	// Reset the singleton state
	tracerProvider = nil
	isInitialized = false

	// Initialize tracing
	viper.Set("observability.tracer.mode", "noop")
	err := InitTracing(context.Background(), "test-service")
	assert.NoError(t, err)

	// Get the first provider
	firstProvider := tracerProvider

	// Try to initialize again - should use the singleton
	err = InitTracing(context.Background(), "another-service")
	assert.NoError(t, err)

	// Should be the same instance
	assert.Same(t, firstProvider, tracerProvider)
}

func TestTracerWithStdout(t *testing.T) {
	// Skip in CI environments
	if os.Getenv("CI") != "" {
		t.Skip("Skipping test in CI environment")
	}

	// Reset the singleton state
	tracerProvider = nil
	isInitialized = false

	// Test with stdout mode
	viper.Set("observability.tracer.mode", "stdout")
	err := InitTracing(context.Background(), "test-service")
	assert.NoError(t, err)

	// Create a span - this would be printed to stdout
	tracer := GetTracer("test-component")
	ctx, span := tracer.Start(context.Background(), "test-operation")
	span.End()

	// Verify we can create child spans
	_, child := tracer.Start(ctx, "child-operation")
	child.End()

	// Clean up
	err = Shutdown(context.Background())
	assert.NoError(t, err)
}

func TestGetTracerWithDefaultParams(t *testing.T) {
	// Reset the singleton state
	tracerProvider = nil
	isInitialized = false

	// Initialize tracing
	viper.Set("observability.tracer.mode", "noop")
	err := InitTracing(context.Background(), "test-service")
	assert.NoError(t, err)

	// Get a tracer with default parameters
	tracer := GetTracer("test-component")
	assert.NotNil(t, tracer)

	// Ensure we can use it
	_, span := tracer.Start(context.Background(), "test-operation")
	span.End()
}
