package cmd

import (
	"bytes"
	"testing"

	"github.com/jdfalk/ubuntu-autoinstall-webhook/internal/server"
	"github.com/jdfalk/ubuntu-autoinstall-webhook/internal/logger"
	"github.com/spf13/viper"
)

// Mock function to replace server.StartServer
var mockStartServer = func(port string) error {
	return nil
}

// Helper function to replace server.StartServer and ensure it is restored after the test
func withMockStartServer(mock func(port string) error, testFunc func()) {
	originalStartServer := server.StartServer
	server.StartServer = mock
	defer func() { server.StartServer = originalStartServer }()
	testFunc()
}

func TestServeCommand(t *testing.T) {
	withMockStartServer(mockStartServer, func() {
		t.Run("Serve command exists", func(t *testing.T) {
			if serveCmd.Use != "serve" {
				t.Errorf("Expected 'serve' command, got %s", serveCmd.Use)
			}
		})

		t.Run("Serve command execution", func(t *testing.T) {
			// Set the port in the viper configuration
			viper.Set("port", "8080")

			// Capture the log output from the logger package
			var logOutput bytes.Buffer
			logger.SetOutput(&logOutput)
			defer logger.SetOutput(nil)

			// Execute the serve command
			serveCmd.SetArgs([]string{}) // Ensure no additional arguments are passed
			err := serveCmd.Execute()
			if err != nil {
				t.Errorf("Serve command execution failed: %v", err)
			}

			// Check the log output
			expectedLog := "Webhook server running on port 8080"
			if !bytes.Contains(logOutput.Bytes(), []byte(expectedLog)) {
				t.Errorf("Expected log output %q, got %q", expectedLog, logOutput.String())
			}
		})
	})
}
