// // filepath: /Users/jdfalk/repos/github.com/jdfalk/ubuntu-autoinstall-webhook/internal/db/db_test.go
package db

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestSaveNetworkInterface(t *testing.T) {
	// Create a new sqlmock database and assign it to the global DB variable.
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating sqlmock database: %v", err)
	}
	defer mockDB.Close()
	DB = mockDB

	clientID := "192.168.1.1"
	macAddress := "aa:bb:cc:dd:ee:ff"
	interfaceName := "eth0"
	chipset := "chipset1"
	driver := "driver1"
	networkID := "1"

	// Expect the INSERT ... RETURNING id query.
	mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO network_interfaces (client_id, mac_address, interface_name)
        VALUES ($1, $2, $3)
        ON CONFLICT (mac_address) DO UPDATE SET interface_name = EXCLUDED.interface_name
        RETURNING id;`)).
		WithArgs(clientID, macAddress, interfaceName).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(networkID))

	// Expect the INSERT for chipset info.
	mock.ExpectExec(regexp.QuoteMeta(
		`INSERT INTO network_chipsets (network_interface_id, chipset)
        VALUES ($1, $2)
        ON CONFLICT (network_interface_id) DO UPDATE SET chipset = EXCLUDED.chipset;`)).
		WithArgs(networkID, chipset).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Expect the UPDATE for driver info.
	mock.ExpectExec(regexp.QuoteMeta(
		`UPDATE network_interfaces SET driver = $1 WHERE id = $2;`)).
		WithArgs(driver, networkID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = SaveNetworkInterface(clientID, macAddress, interfaceName, chipset, driver)
	if err != nil {
		t.Errorf("SaveNetworkInterface returned error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unmet expectations: %v", err)
	}
}

func TestSaveCloudInitVersion(t *testing.T) {
	// Create a new sqlmock database and assign it to the global DB variable.
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating sqlmock database: %v", err)
	}
	defer mockDB.Close()
	DB = mockDB

	clientID := "192.168.1.1"
	macAddress := "aa:bb:cc:dd:ee:ff"
	userData := "user_data_test"

	// Expect insertion into the cloud_init_history table.
	mock.ExpectExec(regexp.QuoteMeta(
		`INSERT INTO cloud_init_history (client_id, mac_address, user_data)
        VALUES ($1, $2, $3);`)).
		WithArgs(clientID, macAddress, userData).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Expect deletion query for pruning old versions.
	mock.ExpectExec(regexp.QuoteMeta(
		`DELETE FROM cloud_init_history
        WHERE client_id = $1 AND mac_address = $2
        AND id NOT IN (
            SELECT id FROM cloud_init_history
            WHERE client_id = $1 AND mac_address = $2
            ORDER BY created_at DESC
            LIMIT 5
        );`)).
		WithArgs(clientID, macAddress).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = SaveCloudInitVersion(clientID, macAddress, userData)
	if err != nil {
		t.Errorf("SaveCloudInitVersion returned error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unmet expectations: %v", err)
	}
}

// TestDBPing tests pinging the database, enabling ping monitoring with sqlmock.
func TestDBPing(t *testing.T) {
	// Create a new sqlmock database with MonitorPingsOption enabled and assign it to the global DB variable.
	mockDB, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("Error creating sqlmock database: %v", err)
	}
	defer mockDB.Close()
	DB = mockDB

	// Expect ping to succeed.
	mock.ExpectPing().WillReturnError(nil)

	err = DB.Ping()
	if err != nil {
		t.Errorf("Expected no error on ping, got: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unmet expectations: %v", err)
	}
}
