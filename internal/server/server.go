package server

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"log"
	"mime"
	"net/http"
	"path"
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

	// Serve Angular app and handle deep linking
	http.Handle("/viewer-app/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		filePath := strings.TrimPrefix(r.URL.Path, "/viewer-app/")
		file, err := a.Open(filePath)
		if err == nil {
			defer file.Close()
			fileInfo, err := file.Stat()
			if err == nil {
				data, _ := io.ReadAll(file)
				w.Header().Set("Content-Type", mime.TypeByExtension(path.Ext(filePath)))
				http.ServeContent(w, r, filePath, fileInfo.ModTime(), bytes.NewReader(data))
				return
			}
		}

		// If file does not exist, serve index.html for deep linking
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

		w.Header().Set("Content-Type", "text/html")
		http.ServeContent(w, r, "index.html", fileInfo.ModTime(), bytes.NewReader(data))
	}))

	// Serve API routes
	http.HandleFunc("/api/webhook", handlers.WebhookHandler)
	http.HandleFunc("/api/viewer", handlers.ViewerHandler)
	http.HandleFunc("/api/viewer/", handlers.ViewerDetailHandler)

	addr := fmt.Sprintf("0.0.0.0:%s", port)
	log.Printf("Starting webhook server on %s\n", addr)
	return http.ListenAndServe(addr, nil)
}
