// Package handlers provides HTTP handlers for processing webhook events and viewing log data.
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

	"github.com/jdfalk/ubuntu-autoinstall-webhook/internal/db"
	"github.com/jdfalk/ubuntu-autoinstall-webhook/internal/ipxe"
	"github.com/jdfalk/ubuntu-autoinstall-webhook/internal/logger"
	"github.com/spf13/viper"
	"github.com/jdfalk/ubuntu-autoinstall-webhook/internal/types"
)

// LogWriter defines an interface for logging events.
type LogWriter interface {
	Write(event types.Event) error
}

// FileLogWriter is the production implementation that writes JSON logs to a file.
type FileLogWriter struct {
	LogDir string
}

// Write writes the event as JSON to a file determined by its SourceIP.
func (w *FileLogWriter) Write(event types.Event) error {
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

// WebhookHandler processes incoming webhook events.
func WebhookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	clientIP := getClientIP(r)

	var event types.Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	event.SourceIP = clientIP

	timestamp := time.Unix(int64(event.Timestamp), 0).Format(time.RFC3339)
	logger.Infof("%s - types.Event: %+v", timestamp, event)

	if err := FileLogger.Write(event); err != nil {
		logger.Errorf("Error logging to file: %v", err)
	}

	// Delegate saving client log and status to the db package.
	db.SaveClientLog(event)
	db.SaveClientStatus(event)

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

// ViewerHandler returns a list of logs in JSON format.
func ViewerHandler(w http.ResponseWriter, r *http.Request) {
	logger.Infof("ViewerHandler called")
	w.Header().Set("Content-Type", "application/json")
	if db.DB == nil {
		logger.Infof("ViewerHandler: No database configured, returning empty log list.")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("[]"))
		return
	}

	logs, err := db.GetClientLogs()
	if err != nil {
		logger.Errorf("ViewerHandler error: %s", err.Error())
		http.Error(w, `{"error": "Internal Server Error"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(logs)
}

// ViewerDetailHandler returns detailed information for a specific log entry.
// Reserved ids (like "status", "logs", "report") are handled by deferring to ViewerHandler.
func ViewerDetailHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if db.DB == nil {
		logger.Infof("ViewerDetailHandler: No database configured, returning empty detail.")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{}"))
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/api/viewer/")
	reserved := map[string]bool{
		"status": true,
		"logs":   true,
		"report": true,
	}
	if id == "" || reserved[id] {
		logger.Infof("ViewerDetailHandler: Reserved id '%s' detected, deferring to ViewerHandler", id)
		ViewerHandler(w, r)
		return
	}

	logger.Infof("ViewerDetailHandler called with id: %s", id)
	logDetail, err := db.GetClientLogDetail(id)
	if err != nil {
		logger.Errorf("ViewerDetailHandler error: %s", err.Error())
		http.Error(w, `{"error": "Internal Server Error"}`, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(logDetail)
}

// ClientLogsHandler returns a list of client logs in JSON format.
func ClientLogsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if db.DB == nil {
		logger.Infof("ClientLogsHandler: No database configured, returning empty log list.")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("[]"))
		return
	}

	logs, err := db.GetClientLogs()
	if err != nil {
		logger.Errorf("ClientLogsHandler error: %s", err.Error())
		http.Error(w, `{"error": "Internal Server Error"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(logs)
}

// ServerLogsHandler returns a list of server logs in JSON format.
func ServerLogsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if db.DB == nil {
		logger.Infof("ServerLogsHandler: No database configured, returning empty log list.")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("[]"))
		return
	}

	logs, err := db.GetServerLogs()
	if err != nil {
		logger.Errorf("ServerLogsHandler error: %s", err.Error())
		http.Error(w, `{"error": "Internal Server Error"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(logs)
}

// IpxeConfigsHandler returns the latest iPXE configurations in JSON format.
func IpxeConfigsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if db.DB == nil {
		logger.Infof("IpxeConfigsHandler: No database configured, returning empty config list.")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("[]"))
		return
	}

	configs, err := db.GetIpxeConfigs()
	if err != nil {
		logger.Errorf("IpxeConfigsHandler error: %s", err.Error())
		http.Error(w, `{"error": "Internal Server Error"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(configs)
}

// HistoricalIpxeConfigsHandler returns historical iPXE configurations in JSON format.
func HistoricalIpxeConfigsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if db.DB == nil {
		logger.Infof("HistoricalIpxeConfigsHandler: No database configured, returning empty config list.")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("[]"))
		return
	}

	configs, err := db.GetHistoricalIpxeConfigs()
	if err != nil {
		logger.Errorf("HistoricalIpxeConfigsHandler error: %s", err.Error())
		http.Error(w, `{"error": "Internal Server Error"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(configs)
}

// CloudInitConfigsHandler returns the current cloud-init configurations in JSON format.
func CloudInitConfigsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if db.DB == nil {
		logger.Infof("CloudInitConfigsHandler: No database configured, returning empty config list.")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("[]"))
		return
	}

	configs, err := db.GetCloudInitConfigs()
	if err != nil {
		logger.Errorf("CloudInitConfigsHandler error: %s", err.Error())
		http.Error(w, `{"error": "Internal Server Error"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(configs)
}

// HistoricalCloudInitConfigsHandler returns historical cloud-init configurations in JSON format.
func HistoricalCloudInitConfigsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if db.DB == nil {
		logger.Infof("HistoricalCloudInitConfigsHandler: No database configured, returning empty config list.")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("[]"))
		return
	}

	configs, err := db.GetHistoricalCloudInitConfigs()
	if err != nil {
		logger.Errorf("HistoricalCloudInitConfigsHandler error: %s", err.Error())
		http.Error(w, `{"error": "Internal Server Error"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(configs)
}

// HardwareInfoHandler processes client hardware submissions.
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

// CloudInitUpdateHandler processes cloud-init updates.
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

// UpdateIPXEOnProgress updates the iPXE file when a client reaches a certain completion threshold.
// It determines the correct phase ("install" or "post-install") based on progress, and passes it to ipxe.UpdateIPXEFile.
func UpdateIPXEOnProgress(clientID string, progress int, macAddress string) {
	var phase string
	if progress >= 25 {
		phase = "post-install"
	} else {
		phase = "install"
	}
	if err := ipxe.UpdateIPXEFile(macAddress, phase); err != nil {
		logger.Errorf("Failed to update iPXE for %s: %v", macAddress, err)
	} else {
		logger.Infof("Updated iPXE file for MAC: %s", macAddress)
	}
}

// formatIPFilename formats the IP address into a valid filename.
func formatIPFilename(ip string) string {
	safeIP := strings.ReplaceAll(ip, ".", "_")
	return fmt.Sprintf("%s.json", safeIP)
}
