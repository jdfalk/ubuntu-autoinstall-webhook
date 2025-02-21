package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/jdfalk/ubuntu-autoinstall-webhook/internal/db"
)

type LogEntry struct {
	Timestamp   time.Time `json:"timestamp"`
	Origin      string    `json:"origin"`
	EventType   string    `json:"event_type"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Level       string    `json:"level,omitempty"`
}

type SystemReport struct {
	IPAddress string    `json:"ip_address"`
	LastSeen  time.Time `json:"last_seen"`
}

func ViewerHandler(w http.ResponseWriter, r *http.Request) {
	query := `SELECT DISTINCT client_id, MAX(timestamp) FROM client_logs GROUP BY client_id`
	rows, err := db.DB.Query(query)
	if err != nil {
		http.Error(w, "Failed to retrieve systems", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var reports []SystemReport
	for rows.Next() {
		var report SystemReport
		var lastSeen time.Time
		if err := rows.Scan(&report.IPAddress, &lastSeen); err != nil {
			http.Error(w, "Error scanning row", http.StatusInternalServerError)
			return
		}
		report.LastSeen = lastSeen
		reports = append(reports, report)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reports)
}

func ViewerDetailHandler(w http.ResponseWriter, r *http.Request) {
	ip := strings.TrimPrefix(r.URL.Path, "/viewer/")
	if ip == "" {
		http.Error(w, "Missing IP address", http.StatusBadRequest)
		return
	}

	query := `SELECT timestamp, origin, event_type, name, description, level FROM client_logs WHERE client_id = $1 ORDER BY timestamp DESC`
	rows, err := db.DB.Query(query, ip)
	if err != nil {
		http.Error(w, "Failed to retrieve logs", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var logs []LogEntry
	for rows.Next() {
		var log LogEntry
		if err := rows.Scan(&log.Timestamp, &log.Origin, &log.EventType, &log.Name, &log.Description, &log.Level); err != nil {
			http.Error(w, "Error scanning row", http.StatusInternalServerError)
			return
		}
		logs = append(logs, log)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
}
