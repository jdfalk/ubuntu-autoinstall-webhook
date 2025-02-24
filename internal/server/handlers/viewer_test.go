package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jdfalk/ubuntu-autoinstall-webhook/internal/db"
	"github.com/jdfalk/ubuntu-autoinstall-webhook/internal/server/handlers"
	"github.com/jdfalk/ubuntu-autoinstall-webhook/internal/testutils"
)

func TestViewerHandler(t *testing.T) {
	tdb := testutils.NewTestDB(t)
	db.DB = tdb.DB
	mock := tdb.Mock
	defer db.DB.Close()

	// Prepare expected rows for GetClientLogs.
	now := time.Now()
	rows := sqlmock.NewRows([]string{
		"id", "client_id", "timestamp", "origin", "description", "name", "result", "event_type", "files", "created_at",
	}).AddRow("log1", "client1", now, "origin1", "desc1", "name1", "result1", "event1", "{}", now)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, client_id, timestamp, origin, description, name, result, event_type, files, created_at FROM client_logs ORDER BY created_at DESC")).
		WillReturnRows(rows)

	req, err := http.NewRequest("GET", "/viewer", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.ViewerHandler)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %d", rr.Code)
	}
	var logs []db.ClientLog
	if err := json.NewDecoder(rr.Body).Decode(&logs); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	if len(logs) != 1 {
		t.Errorf("Expected 1 log, got %d", len(logs))
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unmet SQL mock expectations: %v", err)
	}
}

func TestViewerDetailHandler(t *testing.T) {
	tdb := testutils.NewTestDB(t)
	db.DB = tdb.DB
	mock := tdb.Mock
	defer db.DB.Close()

	// Prepare expected row for GetClientLogDetail.
	now := time.Now()
	rows := sqlmock.NewRows([]string{
		"id", "client_id", "timestamp", "origin", "description", "name", "result", "event_type", "files", "created_at",
	}).AddRow("log1", "client1", now, "origin1", "desc1", "name1", "result1", "event1", "{}", now)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, client_id, timestamp, origin, description, name, result, event_type, files, created_at FROM client_logs WHERE id = $1")).
		WithArgs("log1").
		WillReturnRows(rows)

	req, err := http.NewRequest("GET", "/viewer?id=log1", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.ViewerDetailHandler)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %d", rr.Code)
	}

	var logDetail db.ClientLog
	if err := json.NewDecoder(rr.Body).Decode(&logDetail); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	if logDetail.ID != "log1" {
		t.Errorf("Expected log ID 'log1', got '%s'", logDetail.ID)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unmet SQL mock expectations: %v", err)
	}
}
