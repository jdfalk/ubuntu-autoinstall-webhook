package handlers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
)

// Test extracting IP from headers
func TestGetClientIP(t *testing.T) {
	tests := []struct {
		name       string
		headers    map[string]string
		remoteAddr string
		expectedIP string
	}{
		{"X-Forwarded-For", map[string]string{"X-Forwarded-For": "192.168.1.100"}, "", "192.168.1.100"},
		{"X-Real-IP", map[string]string{"X-Real-IP": "10.0.0.2"}, "", "10.0.0.2"},
		{"RemoteAddr", nil, "172.16.2.1:56789", "172.16.2.1"},
	}

	for _, test := range tests {
		req, _ := http.NewRequest("GET", "/", nil)
		for key, value := range test.headers {
			req.Header.Set(key, value)
		}
		if test.remoteAddr != "" {
			req.RemoteAddr = test.remoteAddr
		}

		ip := getClientIP(req)
		if ip != test.expectedIP {
			t.Errorf("%s: expected %s, got %s", test.name, test.expectedIP, ip)
		}
	}
}

// Test valid JSON event with IP logging
func TestWebhookHandlerValidEventWithIPLogging(t *testing.T) {
	testIP := "192.168.1.5"
	testLogDir := "/tmp/autoinstall-webhook-tests"
	testLogFile := filepath.Join(testLogDir, formatIPFilename(testIP))

	// Set up test environment
	_ = os.Remove(testLogFile)
	_ = os.MkdirAll(testLogDir, 0755)

	// Set log directory via Viper (mimicking CLI flag or config file)
	viper.Set("logDir", testLogDir)

	body := []byte(`{
		"origin": "curtin",
		"timestamp": 1440688425.6038516,
		"event_type": "finish",
		"name": "cmd-install",
		"description": "curtin command install",
		"result": "SUCCESS"
	}`)
	req, err := http.NewRequest("POST", "/webhook", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	// Simulate request from testIP
	req.Header.Set("X-Forwarded-For", testIP)

	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(WebhookHandler)
	handler.ServeHTTP(rec, req)

	// Verify response status
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}

	// Verify log file was created
	data, err := ioutil.ReadFile(testLogFile)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	// Verify logged JSON structure
	var loggedEvent Event
	err = json.Unmarshal(data, &loggedEvent)
	if err != nil {
		t.Fatalf("Failed to parse logged JSON: %v", err)
	}

	if loggedEvent.SourceIP != testIP {
		t.Errorf("Expected source_ip %s, got %s", testIP, loggedEvent.SourceIP)
	}

	// Cleanup
	_ = os.Remove(testLogFile)
}

// Test invalid JSON payload
func TestWebhookHandlerInvalidJSON(t *testing.T) {
	body := []byte(`{"invalid_json"}`)
	req, err := http.NewRequest("POST", "/webhook", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(WebhookHandler)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", rec.Code)
	}
}

// Test invalid request method (GET instead of POST)
func TestWebhookHandlerInvalidMethod(t *testing.T) {
	req, err := http.NewRequest("GET", "/webhook", nil)
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(WebhookHandler)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405, got %d", rec.Code)
	}
}
