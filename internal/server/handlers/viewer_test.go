package handlers_test

import (
	"ubuntu-autoinstall-webhook/internal/testutils"
	"ubuntu-autoinstall-webhook/internal/testutils"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jdfalk/ubuntu-autoinstall-webhook/internal/server/handlers"
	"github.com/jdfalk/ubuntu-autoinstall-webhook/internal/testutils"
)

func TestViewerHandler(t *testing.T) {
	tdb := testutils.NewTestDB(t)
	mockDB, mock := ttestutils.NewTestDB(t).DB, ttestutils.NewTestDB(t).Mock
	defer mockDB.Close()
	testutils.NewTestDB(t).DB = mockDB

	// Expect the distinct query to retrieve client_id and last_seen timestamp.
	mock.ExpectQuery(`SELECT DISTINCT client_id, MAX\(timestamp\) FROM client_logs GROUP BY client_id`).
		WillReturnRows(sqlmock.NewRows([]string{"client_id", "last_seen"}).
			AddRow("192.168.1.1", time.Now()))

	req, err := http.NewRequest("GET", "/viewer", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.ViewerHandler)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %v", rr.Code)
	}
}

func TestViewerDetailHandler(t *testing.T) {
	tdb := testutils.NewTestDB(t)
	mockDB, mock := ttestutils.NewTestDB(t).DB, ttestutils.NewTestDB(t).Mock
	defer mockDB.Close()
	testutils.NewTestDB(t).DB = mockDB

	// Expect the detail query to return log entry columns.
	mock.ExpectQuery(`SELECT timestamp, origin, event_type, name, description, level FROM client_logs WHERE client_id = \$1 ORDER BY timestamp DESC`).
		WithArgs("192.168.1.1").
		WillReturnRows(sqlmock.NewRows([]string{"timestamp", "origin", "event_type", "name", "description", "level"}).
			AddRow(time.Now(), "curtin", "install", "cmd-install", "Installation started", "INFO"))

	req, err := http.NewRequest("GET", "/viewer/192.168.1.1", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.ViewerDetailHandler)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %v", rr.Code)
	}

	var logs []handlers.LogEntry
	if err := json.NewDecoder(rr.Body).Decode(&logs); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(logs) == 0 {
		t.Errorf("Expected logs, but got none")
	}
}
