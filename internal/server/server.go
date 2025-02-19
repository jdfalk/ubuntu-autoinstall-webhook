package server

import (
	"fmt"
	"log"
	"net/http"
	"github.com/jdfalk/ubuntu-autoinstall-webhook/internal/server/handlers"
	"github.com/jdfalk/ubuntu-autoinstall-webhook/internal/assets"
)


func StartServer(port string) error {

	a :=&assets.Assets
	fileServer := http.FileServer(http.FS(a))

	// Serve Angular app
	http.Handle("/viewer-app/", http.StripPrefix("/viewer-app/", fileServer))

	// Serve API routes
	http.HandleFunc("/api/webhook", handlers.WebhookHandler)
	http.HandleFunc("/api/viewer", handlers.ViewerHandler)
	http.HandleFunc("/api/viewer/", handlers.ViewerDetailHandler)

	// Serve index.html for deep linking
	http.HandleFunc("/viewer", func(w http.ResponseWriter, r *http.Request) {
		indexFile, err := a.Open("index.html")
		if err != nil {
			http.Error(w, "index.html not found", http.StatusNotFound)
			return
		}
		fileInfo, err := indexFile.Stat()
		if err != nil {
			http.Error(w, "could not get file info", http.StatusInternalServerError)
			return
		}
		http.ServeContent(w, r, "index.html", fileInfo.ModTime(), indexFile)
	})

	addr := fmt.Sprintf("0.0.0.0:%s", port)
	log.Printf("Starting webhook server on %s\n", addr)
	return http.ListenAndServe(addr, nil)
}
