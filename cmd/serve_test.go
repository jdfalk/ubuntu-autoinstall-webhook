package cmd

import (
    "bytes"
    "testing"


)

func TestServeCommand(t *testing.T) {
    t.Run("Serve command exists", func(t *testing.T) {
        if serveCmd.Use != "serve" {
            t.Errorf("Expected 'serve' command, got %s", serveCmd.Use)
        }
    })

    t.Run("Serve command execution", func(t *testing.T) {
        buf := new(bytes.Buffer)
        serveCmd.SetOut(buf)

        err := serveCmd.Execute()
        if err != nil {
            t.Errorf("Serve command execution failed: %v", err)
        }

        output := buf.String()
        if output == "" {
            t.Errorf("Expected output from serve command, got empty string")
        }
    })
}
