package main

import (
	"encoding/json"
	"log"
	"net/http"

	// Importing the handlers and configuration packages which use the centralized types.
	"github.com/jdfalk/ubuntu-autoinstall-webhook/internal/handlers"
	"github.com/jdfalk/ubuntu-autoinstall-webhook/internal/config"
)

// main is the entry point for the webhook server.
func main() {
	// Load the database configuration from environment variables.
	dbConfig := config.LoadDBConfig()
	log.Printf("Database Config: %+v", dbConfig)

	// Setup the webhook endpoint using the centralized handler.
	http.HandleFunc("/webhook", handlers.WebhookHandler)
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
