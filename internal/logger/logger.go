package logger

// Package logger provides centralized logging functionality for the application.
// It handles logging to both files and the database. To enable database logging,
// the application must inject a DB executor (for example, the global DB connection)
// by calling SetDBExecutor. This package also provides a standard logger instance.

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/spf13/viper"
)

// DBExecutor defines the minimal interface required to execute SQL commands.
// Any type that implements Exec(query string, args ...interface{}) (sql.Result, error)
// can be used (for example, a *sql.DB instance).
type DBExecutor interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
}

// dbExecutor holds the injected database executor for logging to the database.
var dbExecutor DBExecutor

// bufferedLog represents a log message that would be written to the system_logs table.
type bufferedLog struct {
	level   string
	message string
}

var (
	// systemSQLBuffer buffers log messages until the DB executor is initialized.
	systemSQLBuffer []bufferedLog
	// bufferMutex protects access to the systemSQLBuffer.
	bufferMutex sync.Mutex
)

// SetDBExecutor injects the DB executor to be used by the logger package.
// Call this during initialization (after establishing a DB connection).
// When set, any buffered log messages are flushed.
func SetDBExecutor(executor DBExecutor) {
	bufferMutex.Lock()
	defer bufferMutex.Unlock()
	dbExecutor = executor
	// Flush any buffered messages to the database.
	for _, bl := range systemSQLBuffer {
		query := `INSERT INTO system_logs (timestamp, level, message) VALUES ($1, $2, $3)`
		_, err := dbExecutor.Exec(query, time.Now(), bl.level, bl.message)
		if err != nil {
			log.Printf("Error logging buffered message to system_logs: %v", err)
		}
	}
	systemSQLBuffer = nil
}

// Logger is the standard logger instance used by the application.
var Logger *log.Logger

func init() {
	// Initialize the logger with a default output to stdout.
	Logger = log.New(os.Stdout, "APP_LOG: ", log.Ldate|log.Ltime|log.Lshortfile)
}

// Event represents a log event with optional fields for client logging.
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
AppendToFile constructs a formatted log entry and appends it to a log file.
It ensures that the log directory exists and creates it if necessary.
The log file path is determined by combining the "logDir" and "logFile" values obtained
from viper configuration.
*/
func AppendToFile(format string, a ...interface{}) {
	entry := fmt.Sprintf(format, a...)
	logDir := viper.GetString("logDir")
	logFile := filepath.Join(logDir, viper.GetString("logFile"))

	// Ensure the log directory exists.
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

// AppendToSystemSQL writes a log entry to the system_logs table in the database.
// If the database executor has not yet been injected, the message is buffered.
func AppendToSystemSQL(level, message string) {
	bufferMutex.Lock()
	defer bufferMutex.Unlock()
	if dbExecutor == nil {
		// Buffer the log message instead of printing "DB not configured..."
		systemSQLBuffer = append(systemSQLBuffer, bufferedLog{level: level, message: message})
		return
	}
	query := `INSERT INTO system_logs (timestamp, level, message) VALUES ($1, $2, $3)`
	_, err := dbExecutor.Exec(query, time.Now(), level, message)
	if err != nil {
		log.Printf("Error logging to system_logs: %v", err)
	}
}

// AppendToClientLogsSQL writes an event to the client_logs table in the database.
func AppendToClientLogsSQL(event Event) {
	if dbExecutor == nil {
		fmt.Printf("DB not configured; skipping client logs: %+v\n", event)
		return
	}
	query := `
        INSERT INTO client_logs (timestamp, origin, event_type, name, description, result, status, progress, message, source_ip)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);
    `
	_, err := dbExecutor.Exec(
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

// AppendToClientStatusSQL writes or updates the client status in the client_status table.
func AppendToClientStatusSQL(event Event) {
	if dbExecutor == nil {
		fmt.Printf("DB not configured; skipping client status update: %+v\n", event)
		return
	}
	query := `
        INSERT INTO client_status (client_id, status, progress, message, updated_at)
        VALUES ($1, $2, $3, $4, NOW())
        ON CONFLICT (client_id) DO UPDATE SET status = $2, progress = $3, message = $4, updated_at = NOW();
    `
	_, err := dbExecutor.Exec(query, event.Origin, event.Status, event.Progress, event.Message)
	if err != nil {
		log.Printf("Error inserting/updating client_status: %v", err)
	}
}

// LogSystem logs a system-level event by writing to the file, database, and stdout.
func LogSystem(level, format string, a ...interface{}) {
	entry := fmt.Sprintf(format, a...)
	// Log to file.
	AppendToFile("[%s] %s", level, entry)
	// Log to the system_logs table.
	AppendToSystemSQL(level, entry)
	// Also output to standard output.
	fmt.Printf("[%s] %s\n", level, entry)
}

// LogClient logs a client event by writing to the file and to the client_logs and client_status tables.
func LogClient(eventType, format string, a ...interface{}) {
	event := Event{
		EventType: eventType,
		Message:   fmt.Sprintf(format, a...),
	}
	// Log to file.
	AppendToFile("[%s] %s", event.EventType, event.Message)
	// Write event to the client_logs table.
	AppendToClientLogsSQL(event)
	// Update client status in the client_status table.
	AppendToClientStatusSQL(event)
}

// Debug logs a debug-level system message.
func Debug(format string, a ...interface{}) {
	LogSystem("DEBUG", format, a...)
}

// Info logs an informational system message.
func Info(format string, a ...interface{}) {
	LogSystem("INFO", format, a...)
}

// Warning logs a warning-level system message.
func Warning(format string, a ...interface{}) {
	LogSystem("WARNING", format, a...)
}

// Error logs an error-level system message.
func Error(format string, a ...interface{}) {
	LogSystem("ERROR", format, a...)
}

// Debugf logs a formatted debug-level system message.
func Debugf(format string, a ...interface{}) {
	LogSystem("DEBUG", format, a...)
}

// Infof logs a formatted informational system message.
func Infof(format string, a ...interface{}) {
	LogSystem("INFO", format, a...)
}

// Warningf logs a formatted warning-level system message.
func Warningf(format string, a ...interface{}) {
	LogSystem("WARNING", format, a...)
}

// Errorf logs a formatted error-level system message.
func Errorf(format string, a ...interface{}) {
	LogSystem("ERROR", format, a...)
}

// ClientDebugf logs a formatted debug message for client events.
func ClientDebugf(format string, a ...interface{}) {
	LogClient("DEBUG", format, a...)
}

// ClientInfof logs a formatted informational message for client events.
func ClientInfof(format string, a ...interface{}) {
	LogClient("INFO", format, a...)
}

// ClientWarningf logs a formatted warning message for client events.
func ClientWarningf(format string, a ...interface{}) {
	LogClient("WARNING", format, a...)
}

// ClientErrorf logs a formatted error message for client events.
func ClientErrorf(format string, a ...interface{}) {
	LogClient("ERROR", format, a...)
}

// SetOutput sets the output destination for the logger.
func SetOutput(w io.Writer) {
	log.SetOutput(w)
}
