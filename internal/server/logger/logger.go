// // filepath: /Users/jdfalk/repos/github.com/jdfalk/ubuntu-autoinstall-webhook/internal/server/logger/logger.go
package logger

import (
	_ "database/sql"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/jdfalk/ubuntu-autoinstall-webhook/internal/db"
	"github.com/spf13/viper"
)

// Event represents a log event passed from handlers.
type Event struct {
	Origin      string  `json:"origin"`
	Timestamp   float64 `json:"timestamp"`
	EventType   string  `json:"event_type"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Result      string  `json:"result"`
	Status      string  `json:"status"`
	Progress    int     `json:"progress"`
	Message     string  `json:"message"`
	SourceIP    string  `json:"source_ip"`
}

/*
AppendToFile constructs a formatted log entry using the provided format and arguments,
ensures that the log directory exists (creating it if necessary), and then appends the
entry to the designated log file.

The log file path is determined by combining the "logDir" and "logFile" values obtained
from viper configuration.
*/
func AppendToFile(format string, a ...interface{}) {
	entry := fmt.Sprintf(format, a...)
	logDir := viper.GetString("logDir")
	logFile := filepath.Join(logDir, viper.GetString("logFile"))

	// Ensure the log directory exists. Create it if it does not exist.
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		if err := os.MkdirAll(logDir, 0755); err != nil {
			log.Printf("Error creating log directory: %v", err)
			return
		}
	}

	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Error opening log file: %v", err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("%s\n", entry))
	if err != nil {
		log.Printf("Error writing to log file: %v", err)
	}
}

/*
AppendToSystemSQL writes a log entry to the SQL database for system-level logs.
It uses the global db.DB instance and expects a table named "system_logs" with
at least columns: timestamp, level, and message.
*/
func AppendToSystemSQL(level, message string) {
	query := `INSERT INTO system_logs (timestamp, level, message) VALUES ($1, $2, $3)`
	_, err := db.DB.Exec(query, time.Now(), level, message)
	if err != nil {
		log.Printf("Error logging to system_logs: %v", err)
	}
}

/*
AppendToClientLogsSQL writes an event to the client_logs table.
*/
func AppendToClientLogsSQL(event Event) {
	query := `
        INSERT INTO client_logs (timestamp, origin, event_type, name, description, result, status, progress, message, source_ip)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);
    `
	_, err := db.DB.Exec(
		query,
		time.Now(),
		event.Origin,
		event.EventType,
		event.Name,
		event.Description,
		event.Result,
		event.Status,
		event.Progress,
		event.Message,
		event.SourceIP,
	)
	if err != nil {
		log.Printf("Error inserting into client_logs: %v", err)
	}
}

/*
AppendToClientStatusSQL writes or updates client status in the client_status table.
It uses event.Origin as the client identifier.
*/
func AppendToClientStatusSQL(event Event) {
	query := `
        INSERT INTO client_status (client_id, status, progress, message, updated_at)
        VALUES ((SELECT id FROM client_identification WHERE id = $1), $2, $3, $4, NOW())
        ON CONFLICT (client_id) DO UPDATE SET status = $2, progress = $3, message = $4, updated_at = NOW();
    `
	_, err := db.DB.Exec(query, event.Origin, event.Status, event.Progress, event.Message)
	if err != nil {
		log.Printf("Error inserting/updating client_status: %v", err)
	}
}

/*
LogSystem logs a message to both the log file and the system_logs table.
Use this function for system-level events.
*/
func LogSystem(level, format string, a ...interface{}) {
	entry := fmt.Sprintf(format, a...)
	// Log to file with a level prefix.
	AppendToFile("[%s] %s", level, entry)
	// Log to SQL system_logs.
	AppendToSystemSQL(level, entry)
	// Output to standard out.
	fmt.Printf("[%s] %s\n", level, entry)
}

/*
LogClient logs a client event to both the log file and client-specific tables.
It calls both AppendToClientLogsSQL and AppendToClientStatusSQL.
*/
func LogClient(eventType, format string, a ...interface{}) {
	event := Event{
		EventType: eventType,
		Message:   fmt.Sprintf(format, a...),
	}
	// Log to file with event type and message.
	AppendToFile("[%s] %s", event.EventType, event.Message)
	// Send event to client_logs table.
	AppendToClientLogsSQL(event)
	// Send event to client_status table.
	AppendToClientStatusSQL(event)
}

// Existing convenience wrappers for system logging (backwards–compatibility)
func Debug(format string, a ...interface{}) {
	LogSystem("DEBUG", format, a...)
}

func Info(format string, a ...interface{}) {
	LogSystem("INFO", format, a...)
}

func Warning(format string, a ...interface{}) {
	LogSystem("WARNING", format, a...)
}

func Error(format string, a ...interface{}) {
	LogSystem("ERROR", format, a...)
}

/*
Debugf logs a formatted message with the DEBUG level.
*/
func Debugf(format string, a ...interface{}) {
	LogSystem("DEBUG", format, a...)
}

/*
Infof logs a formatted message with the INFO level.
*/
func Infof(format string, a ...interface{}) {
	LogSystem("INFO", format, a...)
}

/*
Warningf logs a formatted message with the WARNING level.
*/
func Warningf(format string, a ...interface{}) {
	LogSystem("WARNING", format, a...)
}

/*
Errorf logs a formatted message with the ERROR level.
*/
func Errorf(format string, a ...interface{}) {
	LogSystem("ERROR", format, a...)
}

// New Wrappers for client logging
/*
Debugf logs a formatted message with the DEBUG level.
*/
func ClientDebugf(format string, a ...interface{}) {
	LogClient("DEBUG", format, a...)
}

/*
Infof logs a formatted message with the INFO level.
*/
func ClientInfof(format string, a ...interface{}) {
	LogClient("INFO", format, a...)
}

/*
Warningf logs a formatted message with the WARNING level.
*/
func ClientWarningf(format string, a ...interface{}) {
	LogClient("WARNING", format, a...)
}

/*
Errorf logs a formatted message with the ERROR level.
*/
func ClientErrorf(format string, a ...interface{}) {
	LogClient("ERROR", format, a...)
}

// SetOutput sets the output destination for the logger.
func SetOutput(w io.Writer) {
	log.SetOutput(w)
}
