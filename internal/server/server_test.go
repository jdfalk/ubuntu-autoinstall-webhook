package server

import (
	"net/http"
	"testing"
	"time"
)

// Test server starts on correct port
func TestStartServer(t *testing.T) {
	port := "8085"

	// Run server in a goroutine
	go func() {
		err := StartServer(port)
		if err != nil {
			t.Errorf("Failed to start server: %v", err)
		}
	}()

	// Wait for server to start
	time.Sleep(2 * time.Second)

	// Make a request
	resp, err := http.Get("http://localhost:" + port + "/webhook")
	if err != nil {
		t.Errorf("Server did not start correctly: %v", err)
	} else {
		resp.Body.Close()
	}
}
