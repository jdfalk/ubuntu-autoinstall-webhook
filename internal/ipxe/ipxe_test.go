package ipxe

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/viper"
)

func TestStoreFileHistory(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "ipxe_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a temporary iPXE file.
	ipxeFilePath := filepath.Join(tempDir, "test.ipxe")
	initialContent := []byte("#!ipxe\necho Test Config\n")
	if err := os.WriteFile(ipxeFilePath, initialContent, 0644); err != nil {
		t.Fatalf("Failed to write temp iPXE file: %v", err)
	}

	// Store file history using the updated signature (only ipxeFilePath is needed).
	err = storeFileHistory(ipxeFilePath)
	if err != nil {
		t.Fatalf("storeFileHistory failed: %v", err)
	}

	// Check if a history file was created.
	historyFiles, err := filepath.Glob(ipxeFilePath + ".history.*")
	if err != nil {
		t.Fatalf("Failed to check history files: %v", err)
	}

	if len(historyFiles) == 0 {
		t.Fatal("No history files were created")
	}
}

func TestUpdateIPXEFile(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "ipxe_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create the boot folder where the file is expected.
	bootFolder := filepath.Join(tempDir, "boot")
	if err := os.MkdirAll(bootFolder, 0755); err != nil {
		t.Fatalf("Failed to create boot folder: %v", err)
	}

	macAddress := "00:11:22:33:44:55"
	// Normalize MAC address: all lowercase and without colons.
	normalizedMac := strings.ToLower(macAddress)
	normalizedMac = strings.ReplaceAll(normalizedMac, ":", "")
	expectedFileName := fmt.Sprintf("mac-%s.ipxe", normalizedMac)
	ipxeFilePath := filepath.Join(bootFolder, expectedFileName)

	initialContent := []byte("#!ipxe\necho Initial Config\n")
	if err := os.WriteFile(ipxeFilePath, initialContent, 0644); err != nil {
		t.Fatalf("Failed to write temp iPXE file: %v", err)
	}

	// Configure viper to use the temporary directories.
	viper.Set("ipxe_folder", tempDir)
	viper.Set("boot_customization_folder", bootFolder)
	viper.Set("cloud_init_folder", filepath.Join(tempDir, "cloud-init"))

	// Test updating the file.
	err = UpdateIPXEFile(macAddress)
	if err != nil {
		t.Fatalf("UpdateIPXEFile failed: %v", err)
	}

	// Verify that the file content has been updated.
	updatedContent, err := os.ReadFile(ipxeFilePath)
	if err != nil {
		t.Fatalf("Failed to read updated ipxe file: %v", err)
	}
	expectedContent := "#!ipxe\nexit\n"
	if string(updatedContent) != expectedContent {
		t.Errorf("Unexpected updated content:\nGot: %s\nExpected: %s", string(updatedContent), expectedContent)
	}
}
