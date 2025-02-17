package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/jdfalk/ubuntu-autoinstall-webhook/internal/server/handlers"
)

// StartServer initializes and starts the HTTP server
func StartServer(port string) error {
	http.HandleFunc("/webhook", handlers.WebhookHandler)

	addr := fmt.Sprintf("0.0.0.0:%s", port)
	log.Printf("Starting webhook server on %s\n", addr)
	return http.ListenAndServe(addr, nil)
}
