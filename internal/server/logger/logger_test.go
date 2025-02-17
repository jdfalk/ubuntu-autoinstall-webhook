package logger

import (
	"os"
	"testing"

	"github.com/spf13/viper"
)

// Test log file creation
func TestAppendToFile(t *testing.T) {
	testLog := "/tmp/test_log.log"
	os.Remove(testLog) // Ensure clean start

	viper.Set("logDir", "/tmp")
	viper.Set("logFile", "test_log.log")

	AppendToFile("Test log entry")

	// Check if file exists
	if _, err := os.Stat(testLog); os.IsNotExist(err) {
		t.Errorf("Log file was not created")
	}
}
