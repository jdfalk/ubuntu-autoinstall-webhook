package server

import (
	"bytes"
	"fmt"
	"io/fs"
	"mime"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/jdfalk/ubuntu-autoinstall-webhook/assets"
	"github.com/jdfalk/ubuntu-autoinstall-webhook/internal/logger"
	"github.com/jdfalk/ubuntu-autoinstall-webhook/internal/server/handlers"
)

// RegisterRoutes registers all HTTP routes (for serving the Angular app and API endpoints)
// into the default HTTP mux.
func RegisterRoutes() {
	// Serve Angular app (with deep linking support).
	http.Handle("/viewer-app/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Trim the /viewer-app/ prefix from the requested path.
		filePath := strings.TrimPrefix(r.URL.Path, "/viewer-app/")
		// Try to read the requested file from the embedded assets.
		data, err := fs.ReadFile(assets.AssetsFS, filePath)
		if err == nil {
			w.Header().Set("Content-Type", mime.TypeByExtension(path.Ext(filePath)))
			http.ServeContent(w, r, filePath, fileInfoModTime(filePath), bytes.NewReader(data))
			return
		}

		// If file does not exist, serve index.html for deep linking.
		indexData, err := fs.ReadFile(assets.AssetsFS, "index.html")
		if err != nil {
			logger.Errorf("Error reading index.html: %v", err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		http.ServeContent(w, r, "index.html", fileInfoModTime("index.html"), bytes.NewReader(indexData))
	}))

	// Register API routes.
	http.HandleFunc("/api/webhook", handlers.WebhookHandler)
	http.HandleFunc("/api/hardware-info", handlers.HardwareInfoHandler)
	http.HandleFunc("/api/cloud-init", handlers.CloudInitUpdateHandler)
	http.HandleFunc("/api/viewer", handlers.ViewerHandler)
	http.HandleFunc("/api/viewer/", handlers.ViewerDetailHandler)
}

// StartServer configures HTTP handlers then starts the HTTP server on the specified port.
var StartServer = func(port string) error {
	RegisterRoutes()
	addr := fmt.Sprintf("0.0.0.0:%s", port)
	logger.Infof("Starting webhook server on %s", addr)
	return http.ListenAndServe(addr, nil)
}

// fileInfoModTime returns the modification time for a given file in the embedded assets.
// If an error occurs, it returns the current time.
func fileInfoModTime(filePath string) time.Time {
	file, err := assets.AssetsFS.Open(filePath)
	if err == nil {
		info, err := file.Stat()
		if err == nil {
			return info.ModTime()
		}
	}
	return time.Now()
}
