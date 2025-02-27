package ipxe

import (
	"fmt"
	"io"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jdfalk/ubuntu-autoinstall-webhook/internal/db"
	"github.com/jdfalk/ubuntu-autoinstall-webhook/internal/logger"
	"github.com/jdfalk/ubuntu-autoinstall-webhook/internal/testutils"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

func TestStoreFileHistory(t *testing.T) {
	// Suppress logging.
	logger.SetOutput(io.Discard)

	// Override the package-level FS with an in-memory FS.
	origFs := Fs
	Fs = testutils.NewTestFS(t)
	defer func() { Fs = origFs }()

	tempDir := "/tmp/ipxe_test"
	// Create a temporary directory in the in-memory FS.
	if err := Fs.MkdirAll(tempDir, 0755); err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}

	// Create a temporary iPXE file.
	ipxeFilePath := filepath.Join(tempDir, "test.ipxe")
	initialContent := []byte("#!ipxe\necho Test Config\n")
	if err := afero.WriteFile(Fs, ipxeFilePath, initialContent, 0644); err != nil {
		t.Fatalf("Failed to write temp iPXE file: %v", err)
	}

	// Call storeFileHistory.
	err := storeFileHistory(ipxeFilePath)
	if err != nil {
		t.Fatalf("storeFileHistory failed: %v", err)
	}

	// Check if a history file was created.
	historyFiles, err := afero.Glob(Fs, ipxeFilePath+".history.*")
	if err != nil {
		t.Fatalf("Failed to check history files: %v", err)
	}
	if len(historyFiles) == 0 {
		t.Fatal("No history files were created")
	}
}

func TestUpdateIPXEFile(t *testing.T) {
	// Suppress logging.
	logger.SetOutput(io.Discard)

	// Override the package-level FS with an in-memory FS.
	origFs := Fs
	Fs = testutils.NewTestFS(t)
	defer func() { Fs = origFs }()

	// Set up a test DB using testutils.
	testDB := testutils.NewTestDB(t)
	db.DB = testDB.DB
	// Optionally override logger's DBExecutor if needed:
	// logger.SetDBExecutor(testDB.DB)

	tempDir := "/tmp/ipxe_test"
	bootFolder := filepath.Join(tempDir, "boot")
	// Create the boot folder in the in-memory FS.
	if err := Fs.MkdirAll(bootFolder, 0755); err != nil {
		t.Fatalf("Failed to create boot folder: %v", err)
	}

	macAddress := "00:11:22:33:44:55"
	normalizedMac := strings.ToLower(macAddress)
	normalizedMac = strings.ReplaceAll(normalizedMac, ":", "")
	expectedFileName := fmt.Sprintf("mac-%s.ipxe", normalizedMac)
	ipxeFilePath := filepath.Join(bootFolder, expectedFileName)

	// Write an initial iPXE file to simulate an existing configuration.
	initialContent := []byte("#!ipxe\necho Initial Config\n")
	if err := afero.WriteFile(Fs, ipxeFilePath, initialContent, 0644); err != nil {
		t.Fatalf("Failed to write initial iPXE file: %v", err)
	}

	// Configure viper to use our in-memory paths.
	viper.Set("ipxe_folder", tempDir)
	viper.Set("boot_customization_folder", bootFolder)
	viper.Set("cloud_init_folder", filepath.Join(tempDir, "cloud-init"))

	// Prepare expectation for the DB call inside SaveIPXEConfiguration.
	expectedQuery := regexp.QuoteMeta("INSERT INTO ipxe_configuration (client_id, config, phase, created_at) VALUES ($1, $2, $3, $4)")
	newContent := fmt.Sprintf("#!ipxe\n# phase: %s\nexit\n", "install")
	testDB.Mock.ExpectExec(expectedQuery).
		WithArgs(macAddress, newContent, "install", sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Call UpdateIPXEFile with phase "install".
	err := UpdateIPXEFile(macAddress, "install")
	if err != nil {
		t.Fatalf("UpdateIPXEFile failed: %v", err)
	}

	// Verify that the file content has been updated.
	updatedContent, err := afero.ReadFile(Fs, ipxeFilePath)
	if err != nil {
		t.Fatalf("Failed to read updated iPXE file: %v", err)
	}
	expectedContent := "#!ipxe\n# phase: install\nexit\n"
	if string(updatedContent) != expectedContent {
		t.Errorf("Unexpected updated content:\nGot: %s\nExpected: %s", string(updatedContent), expectedContent)
	}

	// Ensure the DB expectations were met.
	if err := testDB.Mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unmet DB expectations: %v", err)
	}
}
