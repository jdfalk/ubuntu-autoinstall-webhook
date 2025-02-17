package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/jdfalk/ubuntu-autoinstall-webhook/internal/server/logger"
)

// File represents an optional file entry in the event payload
type File struct {
	Content  string `json:"content"`  // Base64 encoded content
	Path     string `json:"path"`     // File path
	Encoding string `json:"encoding"` // Encoding format (e.g., base64)
}

// Event represents the webhook payload structure
type Event struct {
	Origin      string  `json:"origin"`
	Timestamp   float64 `json:"timestamp"`
	EventType   string  `json:"event_type"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Result      string  `json:"result,omitempty"` // Some events may have a "result"
	Files       []File  `json:"files,omitempty"`  // Optional files attribute
}

// WebhookHandler processes incoming webhook events
func WebhookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var event Event
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Convert timestamp to human-readable format
	timestamp := time.Unix(int64(event.Timestamp), 0).Format(time.RFC3339)

	// Log event metadata
	logEntry := fmt.Sprintf("%s - Event: %+v\n", timestamp, event)
	logger.AppendToFile(logEntry)

	// Process files if provided
	if len(event.Files) > 0 {
		for _, file := range event.Files {
			fileLog := fmt.Sprintf("Received file: %s (Encoding: %s)", file.Path, file.Encoding)
			logger.AppendToFile(fileLog)
		}
	}

	// Respond
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, `{"status": "success", "message": "Event received"}`)
}
