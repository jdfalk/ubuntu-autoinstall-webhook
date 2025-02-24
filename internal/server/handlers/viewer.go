// Package handlers provides HTTP handlers for viewing log data.
package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/jdfalk/ubuntu-autoinstall-webhook/internal/db"
	"github.com/jdfalk/ubuntu-autoinstall-webhook/internal/logger"
)

// ViewerHandler returns a list of logs in JSON format.
// If the database is not configured (db.DB is nil), it returns an empty JSON array.
func ViewerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// If no database is configured, return an empty array.
	if db.DB == nil {
		logger.Infof("ViewerHandler: No database configured, returning empty log list.")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("[]"))
		return
	}

	// Attempt to retrieve logs from the database.
	logs, err := db.GetClientLogs() // Assumes this function exists in your db package.
	if err != nil {
		logger.Errorf("%s", err.Error())
		http.Error(w, `{"error": "Internal Server Error"}`, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(logs)
}

// ViewerDetailHandler returns detailed information for a specific log entry.
// If the database is not configured (db.DB is nil), it returns an empty JSON object.
func ViewerDetailHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// If no database is configured, return an empty object.
	if db.DB == nil {
		logger.Infof("ViewerDetailHandler: No database configured, returning empty detail.")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{}"))
		return
	}

	// Extract the "id" parameter.
	id := r.URL.Query().Get("id")
	if id == "" {
		logger.Errorf("ViewerDetailHandler: Missing id parameter.")
		http.Error(w, `{"error": "Missing id parameter"}`, http.StatusBadRequest)
		return
	}

	// Attempt to retrieve detailed log info.
	logDetail, err := db.GetClientLogDetail(id) // Assumes this function exists in your db package.
	if err != nil {
		logger.Errorf("%s", err.Error())
		http.Error(w, `{"error": "Internal Server Error"}`, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(logDetail)
}

// ClientLogsHandler returns a list of client logs in JSON format.
// If the database is not configured (db.DB is nil), it returns an empty JSON array.
func ClientLogsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if db.DB == nil {
		logger.Infof("ClientLogsHandler: No database configured, returning empty log list.")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("[]"))
		return
	}

	logs, err := db.GetClientLogs() // Implement this in your db package.
	if err != nil {
		logger.Errorf("%s", err.Error())
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

	logs, err := db.GetServerLogs() // Implement this in your db package.
	if err != nil {
		logger.Errorf("%s", err.Error())
		http.Error(w, `{"error": "Internal Server Error"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(logs)
}

// IpxeConfigsHandler returns the current iPXE configurations in JSON format.
func IpxeConfigsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if db.DB == nil {
		logger.Infof("IpxeConfigsHandler: No database configured, returning empty config list.")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("[]"))
		return
	}

	configs, err := db.GetIpxeConfigs() // Implement this in your db package.
	if err != nil {
		logger.Errorf("%s", err.Error())
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

	configs, err := db.GetHistoricalIpxeConfigs() // Implement this in your db package.
	if err != nil {
		logger.Errorf("%s", err.Error())
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

	configs, err := db.GetCloudInitConfigs() // Implement this in your db package.
	if err != nil {
		logger.Errorf("%s", err.Error())
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

	configs, err := db.GetHistoricalCloudInitConfigs() // Implement this in your db package.
	if err != nil {
		logger.Errorf("%s", err.Error())
		http.Error(w, `{"error": "Internal Server Error"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(configs)
}
