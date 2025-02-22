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

	"github.com/jdfalk/ubuntu-autoinstall-webhook/internal/ipxe"
	"github.com/jdfalk/ubuntu-autoinstall-webhook/internal/server/logger"

	"github.com/jdfalk/ubuntu-autoinstall-webhook/internal/db"
	"github.com/spf13/viper"
)

// LogWriter defines an interface for logging events.
type LogWriter interface {
	Write(event Event) error
}

// FileLogWriter is the production implementation that writes JSON logs to a file.
type FileLogWriter struct {
	LogDir string
}

// Write writes the event as JSON to a file determined by its SourceIP.
func (w *FileLogWriter) Write(event Event) error {
	ipFilename := formatIPFilename(event.SourceIP)
	logFilePath := filepath.Join(w.LogDir, ipFilename)

	// Ensure the log directory exists.
	if _, err := os.Stat(w.LogDir); os.IsNotExist(err) {
		if err := os.MkdirAll(w.LogDir, 0755); err != nil {
			return fmt.Errorf("error creating log directory: %v", err)
		}
	}

	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error opening log file for %s: %v", event.SourceIP, err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(event); err != nil {
		return fmt.Errorf("error writing JSON log for %s: %v", event.SourceIP, err)
	}
	return nil
}

// FileLogger is a package-level variable used by the handler.
// In production it will be assigned to a FileLogWriter instance.
// In tests we can override it with a mock.
var FileLogger LogWriter = &FileLogWriter{
	LogDir: func() string {
		dir := viper.GetString("logDir")
		if dir == "" {
			return "/var/log/autoinstall-webhook"
		}
		return dir
	}(),
}

// File represents an optional file entry in the event payload.
type File struct {
	Content  string `json:"content"`  // Base64 encoded content
	Path     string `json:"path"`     // File path
	Encoding string `json:"encoding"` // Encoding format (e.g., base64)
}

// Event represents the webhook payload structure.
type Event struct {
	Origin      string  `json:"origin"`
	Timestamp   float64 `json:"timestamp"`
	EventType   string  `json:"event_type"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Result      string  `json:"result,omitempty"` // Some events may have a "result"
	Files       []File  `json:"files,omitempty"`  // Optional files attribute
	SourceIP    string  `json:"source_ip"`        // IP of the client
	Status      string  `json:"status,omitempty"`
	Progress    int     `json:"progress,omitempty"`
	Message     string  `json:"message,omitempty"`
}

// WebhookHandler processes incoming webhook events.
func WebhookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Extract client IP.
	clientIP := getClientIP(r)

	// Decode request body.
	var event Event
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	event.SourceIP = clientIP

	// Convert timestamp to human-readable format.
	timestamp := time.Unix(int64(event.Timestamp), 0).Format(time.RFC3339)

	// Log event to standard log.
	logEntry := fmt.Sprintf("%s - Event: %+v\n", timestamp, event)
	logger.Info(logEntry)

	// Log event per source IP in JSON format using FileLogger.
	if err := FileLogger.Write(event); err != nil {
		fmt.Println(err)
	}

	// Save client log to the database.
	saveClientLogToDB(event)

	// If the event includes a status update, save it to the database.
	if event.Status != "" {
		saveClientStatus(event)
	}

	// Respond.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, `{"status": "success", "message": "Event received"}`)
}

// getClientIP extracts the real IP of the client, handling proxies if needed.
func getClientIP(r *http.Request) string {
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		ips := strings.Split(xff, ",")
		return strings.TrimSpace(ips[0])
	}
	xRealIP := r.Header.Get("X-Real-IP")
	if xRealIP != "" {
		return xRealIP
	}
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "unknown"
	}
	return ip
}

// saveClientLogToDB saves the client log to the database.
func saveClientLogToDB(event Event) {
	query := `INSERT INTO client_logs (client_id, timestamp, origin, description, name, result, event_type, files, created_at)
		VALUES ((SELECT id FROM client_identification WHERE id = $1), $2, $3, $4, $5, $6, $7, $8, NOW())`
	filesJSON, _ := json.Marshal(event.Files)
	_, err := db.DB.Exec(query, event.SourceIP, time.Unix(int64(event.Timestamp), 0), event.Origin,
		event.Description, event.Name, event.Result, event.EventType, string(filesJSON))
	if err != nil {
		logger.AppendToFile("Error saving client log: %v\n", err)
	}
}

// saveClientStatus saves client status update to the database.
func saveClientStatus(event Event) {
	query := `INSERT INTO client_status (client_id, status, progress, message, updated_at)
		VALUES ((SELECT id FROM client_identification WHERE id = $1), $2, $3, $4, NOW())
		ON CONFLICT (client_id) DO UPDATE
		SET status = $2, progress = $3, message = $4, updated_at = NOW();`
	_, err := db.DB.Exec(query, event.SourceIP, event.Status, event.Progress, event.Message)
	if err != nil {
		logger.AppendToFile("Error saving client status: %v\n", err)
	}
}

// formatIPFilename formats the IP address into a valid filename.
func formatIPFilename(ip string) string {
	safeIP := strings.ReplaceAll(ip, ".", "_")
	return fmt.Sprintf("%s.json", safeIP)
}

// HardwareInfoHandler processes client hardware submissions
func HardwareInfoHandler(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		ClientID      string `json:"client_id"`
		MacAddress    string `json:"mac_address"`
		InterfaceName string `json:"interface_name"`
		Chipset       string `json:"chipset"`
		Driver        string `json:"driver"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err := db.SaveNetworkInterface(payload.ClientID, payload.MacAddress, payload.InterfaceName, payload.Chipset, payload.Driver)
	if err != nil {
		http.Error(w, "Failed to save network interface", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hardware info saved successfully"))
}

// CloudInitUpdateHandler processes cloud-init updates
func CloudInitUpdateHandler(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		ClientID   string `json:"client_id"`
		MacAddress string `json:"mac_address"`
		UserData   string `json:"user_data"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err := db.SaveCloudInitVersion(payload.ClientID, payload.MacAddress, payload.UserData)
	if err != nil {
		http.Error(w, "Failed to save cloud-init configuration", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Cloud-init data updated successfully"))
}

// UpdateIPXEOnProgress updates the iPXE file when a client reaches 25% completion
func UpdateIPXEOnProgress(clientID string, progress int, macAddress string) {
	if progress >= 25 {
		err := ipxe.UpdateIPXEFile(macAddress)
		if err != nil {
			logger.Errorf("Failed to update iPXE for %s: %v", macAddress, err)
		} else {
			logger.Infof("Updated iPXE file for MAC: %s", macAddress)
		}
	}
}
