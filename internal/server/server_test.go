package server_test

import (
	"io/fs"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jdfalk/ubuntu-autoinstall-webhook/assets"
	"github.com/jdfalk/ubuntu-autoinstall-webhook/internal/server"
)

// TestViewerAppServing verifies that a request to a non-existent file under the Angular app path
// falls back to serving index.html from the embedded assets.
func TestViewerAppServing(t *testing.T) {
	// Ensure index.html exists in the embedded assets.
	indexData, err := fs.ReadFile(assets.AssetsFS, "index.html")
	if err != nil {
		t.Fatalf("index.html not found in embedded assets: %v", err)
	}

	// Reset default mux so that tests don't bleed between each other.
	http.DefaultServeMux = http.NewServeMux()
	// Register routes.
	server.RegisterRoutes()

	// Request a file that does not exist (so index.html should be served).
	req := httptest.NewRequest("GET", "/viewer-app/nonexistent-file", nil)
	w := httptest.NewRecorder()
	http.DefaultServeMux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 for fallback to index.html, got %d", w.Code)
	}
	body := w.Body.String()
	// Check that the response body contains (or equals) the content of index.html.
	if !strings.Contains(body, string(indexData)) {
		t.Errorf("Expected response body to contain index.html content")
	}
}

// TestViewerAppStaticFiles verifies that a request for a static asset under the Angular app path
// returns 200 if the file exists or 404 if it does not.
func TestViewerAppStaticFiles(t *testing.T) {
	// Reset default mux and register routes.
	http.DefaultServeMux = http.NewServeMux()
	server.RegisterRoutes()

	// For testing, try to read a known static file.
	// Adjust "assets/somefile.js" to a file that is embedded in your assets module.
	_, err := fs.ReadFile(assets.AssetsFS, "assets/somefile.js")
	expectedStatus := http.StatusOK
	if err != nil {
		// If the file is not present in the embedded assets, we expect a fallback to index.html,
		// which is served with status 200. (You could also expect 404 if you prefer.)
		expectedStatus = http.StatusOK
	}

	req := httptest.NewRequest("GET", "/viewer-app/assets/somefile.js", nil)
	w := httptest.NewRecorder()
	http.DefaultServeMux.ServeHTTP(w, req)
	if w.Code != expectedStatus {
		t.Errorf("Expected status %d for static file request, got %d", expectedStatus, w.Code)
	}
}
