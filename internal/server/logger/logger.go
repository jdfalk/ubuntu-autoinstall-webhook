package logger

import (
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// log.Printf logs data to a file
func log.Printf(entry string) {
	logDir := viper.GetString("logDir")
	logFile := filepath.Join(logDir, viper.GetString("logFile"))

	// Ensure log directory exists
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		os.MkdirAll(logDir, 0755)
	}

	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Error opening log file: %v", err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(entry)
	if err != nil {
		log.Printf("Error writing to log file: %v", err)
	}
}
