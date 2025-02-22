package logger

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jdfalk/ubuntu-autoinstall-webhook/internal/db"
	"github.com/spf13/viper"
)

// setupFileLogging configures the log file location in /tmp and returns its full path.
func setupFileLogging(t *testing.T, filename string) string {
	t.Helper()
	tempDir := os.TempDir() // e.g. "/tmp"
	viper.Set("logDir", tempDir)
	viper.Set("logFile", filename)
	fullPath := filepath.Join(tempDir, filename)
	// Remove any previous log file.
	os.Remove(fullPath)
	return fullPath
}

// readFileContent is a helper to read a file's content.
func readFileContent(t *testing.T, path string) string {
	t.Helper()
	data, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read file %s: %v", path, err)
	}
	return string(data)
}

func TestAppendToFile(t *testing.T) {
	logFile := setupFileLogging(t, "test_log.log")
	AppendToFile("Test log entry")

	content := readFileContent(t, logFile)
	if !strings.Contains(content, "Test log entry") {
		t.Errorf("Expected log file to contain 'Test log entry', got: %s", content)
	}
}

func TestAppendToSQL(t *testing.T) {
	// Create sqlmock DB and set the global db.DB.
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to open sqlmock database: %v", err)
	}
	defer sqlDB.Close()
	db.DB = sqlDB

	// Expect the SQL log to be written. Use sqlmock.AnyArg() for the timestamp.
	mock.ExpectExec("INSERT INTO system_logs").
		WithArgs(sqlmock.AnyArg(), "INFO", "Test SQL log").
		WillReturnResult(sqlmock.NewResult(1, 1))

	AppendToSQL("INFO", "Test SQL log")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("SQL expectations were not met: %v", err)
	}
}

func TestLog(t *testing.T) {
	logFile := setupFileLogging(t, "test_log_log.log")

	// Setup sqlmock for SQL logging.
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to open sqlmock: %v", err)
	}
	defer sqlDB.Close()
	db.DB = sqlDB

	// Expect SQL logging with INFO level.
	mock.ExpectExec("INSERT INTO system_logs").
		WithArgs(sqlmock.AnyArg(), "INFO", "Test log message").
		WillReturnResult(sqlmock.NewResult(1, 1))

	Log("INFO", "Test log message")

	// Verify file contains the log entry.
	content := readFileContent(t, logFile)
	if !strings.Contains(content, "[INFO]") || !strings.Contains(content, "Test log message") {
		t.Errorf("Expected file to contain '[INFO] Test log message', got: %s", content)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("SQL expectations were not met: %v", err)
	}
}

func TestDebug(t *testing.T) {
	logFile := setupFileLogging(t, "test_log_debug.log")
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to open sqlmock: %v", err)
	}
	defer sqlDB.Close()
	db.DB = sqlDB

	mock.ExpectExec("INSERT INTO system_logs").
		WithArgs(sqlmock.AnyArg(), "DEBUG", "Debug message").
		WillReturnResult(sqlmock.NewResult(1, 1))

	Debug("Debug message")

	content := readFileContent(t, logFile)
	if !strings.Contains(content, "[DEBUG]") || !strings.Contains(content, "Debug message") {
		t.Errorf("Expected file to contain '[DEBUG] Debug message', got: %s", content)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("SQL expectations were not met: %v", err)
	}
}

func TestInfo(t *testing.T) {
	logFile := setupFileLogging(t, "test_log_info.log")
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to open sqlmock: %v", err)
	}
	defer sqlDB.Close()
	db.DB = sqlDB

	mock.ExpectExec("INSERT INTO system_logs").
		WithArgs(sqlmock.AnyArg(), "INFO", "Info message").
		WillReturnResult(sqlmock.NewResult(1, 1))

	Info("Info message")

	content := readFileContent(t, logFile)
	if !strings.Contains(content, "[INFO]") || !strings.Contains(content, "Info message") {
		t.Errorf("Expected file to contain '[INFO] Info message', got: %s", content)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("SQL expectations were not met: %v", err)
	}
}

func TestWarning(t *testing.T) {
	logFile := setupFileLogging(t, "test_log_warning.log")
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to open sqlmock: %v", err)
	}
	defer sqlDB.Close()
	db.DB = sqlDB

	mock.ExpectExec("INSERT INTO system_logs").
		WithArgs(sqlmock.AnyArg(), "WARNING", "Warning message").
		WillReturnResult(sqlmock.NewResult(1, 1))

	Warning("Warning message")

	content := readFileContent(t, logFile)
	if !strings.Contains(content, "[WARNING]") || !strings.Contains(content, "Warning message") {
		t.Errorf("Expected file to contain '[WARNING] Warning message', got: %s", content)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("SQL expectations were not met: %v", err)
	}
}

func TestError(t *testing.T) {
	logFile := setupFileLogging(t, "test_log_error.log")
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to open sqlmock: %v", err)
	}
	defer sqlDB.Close()
	db.DB = sqlDB

	mock.ExpectExec("INSERT INTO system_logs").
		WithArgs(sqlmock.AnyArg(), "ERROR", "Error message").
		WillReturnResult(sqlmock.NewResult(1, 1))

	Error("Error message")

	content := readFileContent(t, logFile)
	if !strings.Contains(content, "[ERROR]") || !strings.Contains(content, "Error message") {
		t.Errorf("Expected file to contain '[ERROR] Error message', got: %s", content)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("SQL expectations were not met: %v", err)
	}
}
