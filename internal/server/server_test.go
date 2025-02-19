package server_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"os"
)

func TestViewerAppServing(t *testing.T) {
	distDir := "viewer-app/dist/viewer-app"
	indexPath := distDir + "/index.html"

	// Ensure index.html exists for the test
	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		t.Fatalf("index.html not found in %s", distDir)
	}

	req := httptest.NewRequest("GET", "/viewer", nil)
	w := httptest.NewRecorder()
	http.DefaultServeMux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestViewerAppStaticFiles(t *testing.T) {
	req := httptest.NewRequest("GET", "/viewer-app/assets/somefile.js", nil)
	w := httptest.NewRecorder()
	http.DefaultServeMux.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound && w.Code != http.StatusOK {
		t.Errorf("Expected status 200 or 404, got %d", w.Code)
	}
}
