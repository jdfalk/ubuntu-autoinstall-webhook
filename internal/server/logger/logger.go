package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/jdfalk/ubuntu-autoinstall-webhook/internal/db"
	"github.com/spf13/viper"
)

// AppendToFile logs the formatted message into a file.
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

// AppendToSQL logs the message to an SQL database table.
// It uses the global db.DB instance and expects a "system_logs" table with columns: timestamp, level, message.
func AppendToSQL(level, message string) {
	query := `INSERT INTO system_logs (timestamp, level, message) VALUES ($1, $2, $3)`
	_, err := db.DB.Exec(query, time.Now(), level, message)
	if err != nil {
		log.Printf("Error logging to SQL database: %v", err)
	}
}

// Log writes a log entry with the given level and message to both the file and the database.
func Log(level, format string, a ...interface{}) {
	entry := fmt.Sprintf(format, a...)
	// Log to file with level prefix.
	AppendToFile("[%s] %s", level, entry)
	// Log to SQL database.
	AppendToSQL(level, entry)
}

// Debug logs a debug-level message.
func Debug(format string, a ...interface{}) {
	Log("DEBUG", format, a...)
}

// Info logs an info-level message.
func Info(format string, a ...interface{}) {
	Log("INFO", format, a...)
}

// Warning logs a warning-level message.
func Warning(format string, a ...interface{}) {
	Log("WARNING", format, a...)
}

// Error logs an error-level message.
func Error(format string, a ...interface{}) {
	Log("ERROR", format, a...)
}
