package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type LogEntry struct {
	Timestamp   float64 `json:"timestamp"`
	Origin      string  `json:"origin"`
	EventType   string  `json:"event_type"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Level       string  `json:"level"`
}

type SystemReport struct {
	IPAddress string    `json:"ip_address"`
	LastSeen  time.Time `json:"last_seen"`
}

func GetReportedSystems(logDir string) ([]SystemReport, error) {
	var reports []SystemReport
	files, err := os.ReadDir(logDir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".json") {
			ip := strings.TrimSuffix(file.Name(), ".json")
			filePath := filepath.Join(logDir, file.Name())
			data, err := os.ReadFile(filePath)
			if err != nil {
				continue
			}

			var logs []LogEntry
			if err := json.Unmarshal(data, &logs); err != nil {
				continue
			}

			if len(logs) > 0 {
				sort.Slice(logs, func(i, j int) bool {
					return logs[i].Timestamp > logs[j].Timestamp
				})
				reports = append(reports, SystemReport{
					IPAddress: strings.ReplaceAll(ip, "_", "."),
					LastSeen:  time.Unix(int64(logs[0].Timestamp), 0),
				})
			}
		}
	}
	return reports, nil
}

func ViewerHandler(w http.ResponseWriter, r *http.Request) {
	logDir := viper.GetString("logDir") // Use config or flag value
	// Ensure logDir is set, fallback if empty
	if logDir == "" {
		logDir = "/var/log/autoinstall-webhook"
	}
	reports, err := GetReportedSystems(logDir)
	if err != nil {
		http.Error(w, "Failed to retrieve systems", http.StatusInternalServerError)
		return
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

	logDir := "./logs"
	filePath := filepath.Join(logDir, strings.ReplaceAll(ip, ".", "_")+".json")
	data, err := os.ReadFile(filePath)
	if err != nil {
		http.Error(w, "Log file not found", http.StatusNotFound)
		return
	}

	var logs []LogEntry
	if err := json.Unmarshal(data, &logs); err != nil {
		http.Error(w, "Failed to parse log file", http.StatusInternalServerError)
		return
	}

	sort.Slice(logs, func(i, j int) bool {
		return logs[i].Timestamp > logs[j].Timestamp
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
}
