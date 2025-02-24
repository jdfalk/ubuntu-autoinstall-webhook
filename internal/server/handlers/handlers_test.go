package handlers_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jdfalk/ubuntu-autoinstall-webhook/internal/db"
	"github.com/jdfalk/ubuntu-autoinstall-webhook/internal/logger"
	"github.com/jdfalk/ubuntu-autoinstall-webhook/internal/server/handlers"
	"github.com/jdfalk/ubuntu-autoinstall-webhook/internal/testutils"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/mock"
)

var mockDB *sql.DB
var sqlMock sqlmock.Sqlmock
var appFs afero.Fs

// MockLogWriter implements the LogWriter interface for testing.
type MockLogWriter struct {
	mock.Mock
}

func (m *MockLogWriter) Write(event handlers.Event) error {
	args := m.Called(event)
	return args.Error(0)
}

func TestMain(m *testing.M) {
	// Create a test DB using our helper.
	tdb := testutils.NewTestDB(&testing.T{})
	mockDB, sqlMock = tdb.DB, tdb.Mock
	// Assign the global DB variable used in production code.
	db.DB = mockDB

	// Use an in-memory filesystem for testing.
	appFs = afero.NewMemMapFs()

	// Define the logs subdirectory.
	tempBaseDir := "/tmp/autoinstall-webhook-test"
	logDir := tempBaseDir + "/logs"

	// Ensure the base and logs directories exist.
	if err := appFs.MkdirAll(logDir, 0755); err != nil {
		logger.Errorf("%s", "Failed to create logs directory")
		os.Exit(1)
	}

	// Define the log file path.
	logFilePath := logDir + "/log.json"
	// Remove any pre-existing file/directory at that path.
	_ = appFs.RemoveAll(logFilePath)

	// Create the log file in the in-memory filesystem.
	file, err := appFs.Create(logFilePath)
	if err != nil {
		logger.Errorf("%s", "Failed to create log file")
		os.Exit(1)
	}
	file.Close()

	// Set the log directory and log file name.
	viper.Set("logDir", logDir)
	viper.Set("logFile", "log.json")

	// Silence logger output during tests.
	logger.SetOutput(io.Discard)

	code := m.Run()
	mockDB.Close()
	os.Exit(code)
}

func TestWebhookHandler_Success(t *testing.T) {
	// Set up expected database calls.
	sqlMock.ExpectExec("INSERT INTO client_logs").WillReturnResult(sqlmock.NewResult(1, 1))
	sqlMock.ExpectExec(`INSERT INTO client_status $begin:math:text$client_id, status, progress, message, updated_at$end:math:text$ VALUES $begin:math:text$\\(SELECT id FROM client_identification WHERE id = \\$1$end:math:text$, \$2, \$3, \$4, NOW$begin:math:text$$end:math:text$\) ON CONFLICT $begin:math:text$client_id$end:math:text$ DO UPDATE SET status = \$2, progress = \$3, message = \$4, updated_at = NOW$begin:math:text$$end:math:text$`).WillReturnResult(sqlmock.NewResult(1, 1))

	// Override the package-level FileLogger with our mock.
	mockLogger := new(MockLogWriter)
	handlers.FileLogger = mockLogger

	// Create an event. (Note: SourceIP will be derived from r.RemoteAddr.)
	event := handlers.Event{
		Origin:      "curtin",
		Timestamp:   float64(time.Now().Unix()),
		EventType:   "finish",
		Name:        "cmd-install",
		Description: "curtin command install",
		Result:      "SUCCESS",
		Status:      "installing",
		Progress:    50,
		Message:     "Installation is halfway done",
	}

	// Prepare the request. Set RemoteAddr so getClientIP returns "192.168.1.1".
	body, _ := json.Marshal(event)
	req, err := http.NewRequest("POST", "/webhook", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.RemoteAddr = "192.168.1.1:1234"
	req.Header.Set("Content-Type", "application/json")

	// Set expectation on the mock logger for the "finish" event.
	expected := event
	expected.SourceIP = "192.168.1.1"
	mockLogger.
		On("Write", mock.MatchedBy(func(e handlers.Event) bool {
			return e.Origin == expected.Origin &&
				e.EventType == expected.EventType &&
				e.Name == expected.Name &&
				e.SourceIP == expected.SourceIP
		})).
		Return(nil)

	rr := httptest.NewRecorder()
	http.HandlerFunc(handlers.WebhookHandler).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %v", rr.Code)
	}

	mockLogger.AssertExpectations(t)
}

func TestWebhookHandler_InvalidMethod(t *testing.T) {
	req, err := http.NewRequest("GET", "/webhook", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	rr := httptest.NewRecorder()
	http.HandlerFunc(handlers.WebhookHandler).ServeHTTP(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status Method Not Allowed, got %v", rr.Code)
	}
}

func TestWebhookHandler_InvalidJSON(t *testing.T) {
	body := []byte("{invalid json}")
	req, err := http.NewRequest("POST", "/webhook", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	http.HandlerFunc(handlers.WebhookHandler).ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status Bad Request, got %v", rr.Code)
	}
}

func TestWebhookHandler_StatusUpdate(t *testing.T) {
	// Setup database expectation for client_logs insert.
	sqlMock.ExpectExec("INSERT INTO client_logs").WillReturnResult(sqlmock.NewResult(1, 1))
	// Setup database expectation for client_status update.
	sqlMock.ExpectExec(`INSERT INTO client_status $begin:math:text$client_id, status, progress, message, updated_at$end:math:text$ VALUES $begin:math:text$\\(SELECT id FROM client_identification WHERE id = \\$1$end:math:text$, \$2, \$3, \$4, NOW$begin:math:text$$end:math:text$\) ON CONFLICT $begin:math:text$client_id$end:math:text$ DO UPDATE SET status = \$2, progress = \$3, message = \$4, updated_at = NOW$begin:math:text$$end:math:text$`).WillReturnResult(sqlmock.NewResult(1, 1))

	// Override the package-level FileLogger with our mock for status event.
	mockLogger := new(MockLogWriter)
	handlers.FileLogger = mockLogger

	event := handlers.Event{
		Origin:      "curtin",
		Timestamp:   float64(time.Now().Unix()),
		EventType:   "status",
		Name:        "install-progress",
		Description: "Updating install progress",
		Status:      "in-progress",
		Progress:    75,
		Message:     "Installation nearing completion",
	}

	// Prepare request; setting RemoteAddr ensures the SourceIP is derived correctly.
	body, _ := json.Marshal(event)
	req, err := http.NewRequest("POST", "/webhook", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.RemoteAddr = "192.168.1.2:1234"
	req.Header.Set("Content-Type", "application/json")

	// Set expectation on the mock logger.
	expected := event
	expected.SourceIP = "192.168.1.2"
	mockLogger.
		On("Write", mock.MatchedBy(func(e handlers.Event) bool {
			return e.Origin == expected.Origin &&
				e.EventType == expected.EventType &&
				e.Name == expected.Name &&
				e.SourceIP == expected.SourceIP
		})).
		Return(nil)

	rr := httptest.NewRecorder()
	http.HandlerFunc(handlers.WebhookHandler).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %v", rr.Code)
	}

	mockLogger.AssertExpectations(t)
}
