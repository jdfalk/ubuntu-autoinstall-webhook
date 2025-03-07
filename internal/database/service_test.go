// internal/database/service_test.go
package database_test

import (
	"context"
	"testing"

	"github.com/jdfalk/ubuntu-autoinstall-webhook/internal/database"
)

func TestConnectAndMigrate(t *testing.T) {
	dbService := database.NewService()

	// Test connection
	if err := dbService.Connect(context.Background()); err != nil {
		t.Fatalf("Connect() returned an error: %v", err)
	}

	// Test schema migration
	if err := dbService.MigrateSchema(context.Background()); err != nil {
		t.Fatalf("MigrateSchema() returned an error: %v", err)
	}
}

func TestInsertRecord(t *testing.T) {
	dbService := database.NewService()

	record := map[string]interface{}{
		"id":   1,
		"name": "Test Record",
	}

	if err := dbService.InsertRecord(context.Background(), record); err != nil {
		t.Errorf("InsertRecord() returned an error: %v", err)
	}
}

func TestQuery(t *testing.T) {
	dbService := database.NewService()

	_, err := dbService.Query(context.Background(), "SELECT * FROM records WHERE id = ?", 1)
	if err != nil {
		t.Errorf("Query() returned an error: %v", err)
	}
}
