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

	// Open the log file in append mode, creating it if necessary.
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Error opening log file: %v", err)
		return
	}
	defer file.Close()

	// Append the log entry to the file with a newline.
	_, err = file.WriteString(fmt.Sprintf("%s\n", entry))
	if err != nil {
		log.Printf("Error writing to log file: %v", err)
	}
}

/*
AppendToSQL writes a log entry to an SQL database. It uses the global db.DB instance
and expects a table named "system_logs" with at least three columns: timestamp, level,
and message. The current timestamp is added along with the provided log level and message.
Any error during execution is logged to the standard logger.
*/
func AppendToSQL(level, message string) {
	query := `INSERT INTO system_logs (timestamp, level, message) VALUES ($1, $2, $3)`
	_, err := db.DB.Exec(query, time.Now(), level, message)
	if err != nil {
		log.Printf("Error logging to SQL database: %v", err)
	}
}

/*
Log is the central logging function of the module. It creates a single log entry by
formatting the provided message and then writes that entry to both the log file and
the SQL database. The provided 'level' (e.g., DEBUG, INFO, etc.) is prepended to the
message in the file log.
*/
func Log(level, format string, a ...interface{}) {
	entry := fmt.Sprintf(format, a...)
	// Log to file with a level prefix.
	AppendToFile("[%s] %s", level, entry)
	// Also log to the SQL database.
	AppendToSQL(level, entry)
}

// The following functions are convenience wrappers for the primary Log function.

/*
Debug logs a message with the DEBUG level. It can be used for detailed troubleshooting.
*/
func Debug(format string, a ...interface{}) {
	Log("DEBUG", format, a...)
}

/*
Info logs a message with the INFO level, typically used for general operational information.
*/
func Info(format string, a ...interface{}) {
	Log("INFO", format, a...)
}

/*
Warning logs a message with the WARNING level, typically indicating a potential problem.
*/
func Warning(format string, a ...interface{}) {
	Log("WARNING", format, a...)
}

/*
Error logs a message with the ERROR level, typically indicating that an error has occurred.
*/
func Error(format string, a ...interface{}) {
	Log("ERROR", format, a...)
}

// The following functions are identical in functionality to their respective wrappers
// above; they use an "f" suffix to make explicit that they accept a format string.

/*
Debugf logs a formatted message with the DEBUG level.
*/
func Debugf(format string, a ...interface{}) {
	Log("DEBUG", format, a...)
}

/*
Infof logs a formatted message with the INFO level.
*/
func Infof(format string, a ...interface{}) {
	Log("INFO", format, a...)
}

/*
Warningf logs a formatted message with the WARNING level.
*/
func Warningf(format string, a ...interface{}) {
	Log("WARNING", format, a...)
}

/*
Errorf logs a formatted message with the ERROR level.
*/
func Errorf(format string, a ...interface{}) {
	Log("ERROR", format, a...)
}
