package server

import (
	"bytes"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"strings"

	"github.com/jdfalk/ubuntu-autoinstall-webhook/assets"
	"github.com/jdfalk/ubuntu-autoinstall-webhook/internal/server/handlers"
)

func StartServer(port string) error {
	a := assets.AssetsFS

	// Check if index.html exists
	entries, err := fs.ReadDir(a, ".")
	if err != nil {
		log.Printf("Error reading directory: %v", err)
		return fmt.Errorf("error reading directory: %v", err)
	}

	for _, entry := range entries {
		log.Printf("Found file: %s", entry.Name())
	}

	_, err = fs.Stat(a, "index.html")
	if err != nil {
		log.Printf("index.html not found: %v", err)
		return fmt.Errorf("index.html not found: %v", err)
	}

	fileServer := http.FileServer(http.FS(a))

	// Serve Angular static files
	http.Handle("/viewer-app/", http.StripPrefix("/viewer-app/", fileServer))

	// Serve API routes
	http.HandleFunc("/api/webhook", handlers.WebhookHandler)
	http.HandleFunc("/api/viewer", handlers.ViewerHandler)
	http.HandleFunc("/api/viewer/", handlers.ViewerDetailHandler)

	// Serve index.html for deep linking
	http.HandleFunc("/viewer-app/", func(w http.ResponseWriter, r *http.Request) {
		// If the requested file does not exist, return index.html
		if !strings.Contains(r.URL.Path, ".") {
			data, err := fs.ReadFile(a, "index.html")
			if err != nil {
				http.Error(w, "index.html not found", http.StatusNotFound)
				return
			}

			fileInfo, err := fs.Stat(a, "index.html")
			if err != nil {
				http.Error(w, "could not get file info", http.StatusInternalServerError)
				return
			}

			http.ServeContent(w, r, "index.html", fileInfo.ModTime(), bytes.NewReader(data))
			return
		}

		// Otherwise, serve the requested file normally
		fileServer.ServeHTTP(w, r)
	})

	addr := fmt.Sprintf("0.0.0.0:%s", port)
	log.Printf("Starting webhook server on %s\n", addr)
	return http.ListenAndServe(addr, nil)
}
