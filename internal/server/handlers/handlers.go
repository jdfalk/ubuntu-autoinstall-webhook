package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/jdfalk/ubuntu-autoinstall-webhook/internal/server/logger"
)

// Report structure
type Report struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// WebhookHandler handles incoming webhook requests
func WebhookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var report Report
	err := json.NewDecoder(r.Body).Decode(&report)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Log report
	logEntry := fmt.Sprintf("%s - Received: %+v\n", time.Now().Format(time.RFC3339), report)
	logger.AppendToFile(logEntry)

	// Respond
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, `{"status": "success", "message": "Data received"}`)
}
