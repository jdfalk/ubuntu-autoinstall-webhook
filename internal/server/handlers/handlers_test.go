package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Test valid JSON event with files
func TestWebhookHandlerValidEventWithFiles(t *testing.T) {
	body := []byte(`{
		"origin": "curtin",
		"timestamp": 1440688425.6038516,
		"event_type": "finish",
		"name": "cmd-install",
		"description": "curtin command install",
		"result": "SUCCESS",
		"files": [
			{
				"content": "fCBzZmRpc2s....gLS1uby1yZX",
				"path": "/var/log/curtin/install.log",
				"encoding": "base64"
			},
			{
				"content": "fCBzZmRpc2s....gLS1uby1yZX",
				"path": "/var/log/syslog",
				"encoding": "base64"
			}
		]
	}`)
	req, err := http.NewRequest("POST", "/webhook", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(WebhookHandler)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}
}

// Test valid JSON event without files
func TestWebhookHandlerValidEventWithoutFiles(t *testing.T) {
	body := []byte(`{
		"origin": "curtin",
		"timestamp": 1440688425.6038516,
		"event_type": "start",
		"name": "cmd-install",
		"description": "curtin command install"
	}`)
	req, err := http.NewRequest("POST", "/webhook", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(WebhookHandler)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}
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
	req, err := http.NewRequest("GET", "/webhook", nil) // Using GET instead of POST
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
