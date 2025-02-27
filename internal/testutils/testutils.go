// testutils.go
//
// Package testutils contains shared testing utilities and helpers
// that can be used across multiple test files to keep our code DRY.
package testutils

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

// TestDB provides a shared mock database setup for tests.
type TestDB struct {
	DB   *sql.DB         // DB is the mocked database connection.
	Mock sqlmock.Sqlmock // Mock is the sqlmock instance used for setting expectations.
}

// NewTestDB initializes and returns a mock database instance.
// It registers a cleanup function to close the database connection when the test finishes.
func NewTestDB(t *testing.T) *TestDB {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to initialize mock database: %v", err)
	}
	// Register cleanup to close the DB after the test.
	t.Cleanup(func() {
		db.Close()
	})
	return &TestDB{
		DB:   db,
		Mock: mock,
	}
}

// MockDBInit returns a function that simulates successful database initialization.
// This is useful for overriding real database initialization during tests.
func MockDBInit() func() error {
	return func() error {
		return nil
	}
}
