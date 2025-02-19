package server

import (
	"bytes"
	"fmt"
	"io/fs"
	"log"
	"net/http"

	"github.com/jdfalk/ubuntu-autoinstall-webhook/internal/server/handlers"
	"github.com/jdfalk/ubuntu-autoinstall-webhook/assets"
)

func StartServer(port string) error {
	a := &assets.Assets
	fileServer := http.FileServer(http.FS(a))

	// Serve Angular app
	http.Handle("/viewer-app/", http.StripPrefix("/viewer-app/", fileServer))

	// Serve API routes
	http.HandleFunc("/api/webhook", handlers.WebhookHandler)
	http.HandleFunc("/api/viewer", handlers.ViewerHandler)
	http.HandleFunc("/api/viewer/", handlers.ViewerDetailHandler)

	// Serve index.html for deep linking
	http.HandleFunc("/viewer", func(w http.ResponseWriter, r *http.Request) {
		// Get file info using fs.Stat
		fileInfo, err := fs.Stat(a, "index.html")
		if err != nil {
			http.Error(w, "index.html not found", http.StatusNotFound)
			return
		}

		// Read file into memory
		data, err := a.ReadFile("index.html")
		if err != nil {
			http.Error(w, "could not read index.html", http.StatusInternalServerError)
			return
		}

		// Serve content using an in-memory reader
		http.ServeContent(w, r, "index.html", fileInfo.ModTime(), bytes.NewReader(data))
	})

	addr := fmt.Sprintf("0.0.0.0:%s", port)
	log.Printf("Starting webhook server on %s\n", addr)
	return http.ListenAndServe(addr, nil)
}
