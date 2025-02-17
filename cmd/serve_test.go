package cmd

import (
	"testing"
)

// Test if the serve command exists
func TestServeCommand(t *testing.T) {
	if serveCmd.Use != "serve" {
		t.Errorf("Expected 'serve' command, got %s", serveCmd.Use)
	}
}
