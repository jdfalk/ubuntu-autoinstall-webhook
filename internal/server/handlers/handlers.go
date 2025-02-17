package handlers

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jdfalk/ubuntu-autoinstall-webhook/internal/server/logger"
	"github.com/spf13/viper"
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
	SourceIP    string  `json:"source_ip"`        // IP of the client
}

// WebhookHandler processes incoming webhook events
func WebhookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Extract client IP
	clientIP := getClientIP(r)

	// Decode request body
	var event Event
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Convert timestamp to human-readable format
	timestamp := time.Unix(int64(event.Timestamp), 0).Format(time.RFC3339)
	event.SourceIP = clientIP

	// Log event to standard log
	logEntry := fmt.Sprintf("%s - Event: %+v\n", timestamp, event)
	logger.AppendToFile(logEntry)

	// Log event per source IP in JSON format
	logEventByIP(event)

	// Respond
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, `{"status": "success", "message": "Event received"}`)
}

// Extract the real IP of the client, handling proxies if needed
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header first
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		ips := strings.Split(xff, ",")
		return strings.TrimSpace(ips[0]) // First IP in the list is the client
	}

	// Check X-Real-IP header
	xRealIP := r.Header.Get("X-Real-IP")
	if xRealIP != "" {
		return xRealIP
	}

	// Default: get remote address
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "unknown"
	}
	return ip
}

// Log the event per source IP in JSON format
func logEventByIP(event Event) {
	logDir := viper.GetString("logDir") // Use config or flag value

	// Ensure logDir is set, fallback if empty
	if logDir == "" {
		logDir = "/var/log/autoinstall-webhook"
	}

	ipFilename := formatIPFilename(event.SourceIP)
	logFilePath := filepath.Join(logDir, ipFilename)

	// Ensure log directory exists
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		os.MkdirAll(logDir, 0755)
	}

	// Open or create the log file
	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Error opening log file for %s: %v\n", event.SourceIP, err)
		return
	}
	defer file.Close()

	// Encode event as JSON and append to file
	encoder := json.NewEncoder(file)
	err = encoder.Encode(event)
	if err != nil {
		fmt.Printf("Error writing JSON log for %s: %v\n", event.SourceIP, err)
	}
}

// Format the IP address into a valid filename
func formatIPFilename(ip string) string {
	safeIP := strings.ReplaceAll(ip, ".", "_")
	return fmt.Sprintf("%s.json", safeIP)
}
