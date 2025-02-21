package testutils

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

// TestDB provides a shared mock database setup for tests.
type TestDB struct {
	DB   *sql.DB
	Mock sqlmock.Sqlmock
}

// NewTestDB initializes and returns a mock database instance.
func NewTestDB(t *testing.T) *TestDB {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to initialize mock database: %v", err)
	}

	return &TestDB{
		DB:   db,
		Mock: mock,
	}
}
