package cmd

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/jdfalk/ubuntu-autoinstall-webhook/internal/server"
	"github.com/spf13/viper"
)

// mockStartServer is a mock function that replaces server.StartServer in tests.
var mockStartServer = func(port string) error {
	return nil
}

// withMockStartServer temporarily replaces server.StartServer with a mock function.
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
			// Override configuration values for test so no permission errors occur.
			viper.Set("port", "8080")
			viper.Set("logDir", "./test-log")
			viper.Set("ipxe_folder", "./test-ipxe")
			viper.Set("boot_customization_folder", "./test-ipxe/boot")
			viper.Set("cloud_init_folder", "./test-cloud-init")

			// Ensure these directories exist.
			if err := os.MkdirAll("./test-log", 0755); err != nil {
				t.Fatalf("Failed to create test log directory: %v", err)
			}
			if err := os.MkdirAll("./test-ipxe/boot", 0755); err != nil {
				t.Fatalf("Failed to create test ipxe boot directory: %v", err)
			}
			if err := os.MkdirAll("./test-cloud-init", 0755); err != nil {
				t.Fatalf("Failed to create test cloud-init directory: %v", err)
			}

			// Capture stdout (which our logger writes to) by redirecting os.Stdout.
			oldStdout := os.Stdout
			r, w, err := os.Pipe()
			if err != nil {
				t.Fatalf("Failed to create pipe: %v", err)
			}
			os.Stdout = w

			// Execute the serve command using rootCmd.
			rootCmd.SetArgs([]string{"serve"})
			err = rootCmd.Execute()
			if err != nil {
				t.Errorf("Serve command execution failed: %v", err)
			}

			// Close writer and restore stdout.
			w.Close()
			os.Stdout = oldStdout

			var buf bytes.Buffer
			if _, err := io.Copy(&buf, r); err != nil {
				t.Fatalf("Failed to read captured output: %v", err)
			}
			captured := buf.String()

			// Check the captured log output.
			expectedLog := "Webhook server running on port 8080"
			if !bytes.Contains([]byte(captured), []byte(expectedLog)) {
				t.Errorf("Expected log output %q, got %q", expectedLog, captured)
			}
		})
	})
}
