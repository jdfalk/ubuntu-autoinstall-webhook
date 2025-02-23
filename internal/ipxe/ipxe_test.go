package ipxe

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
)

func TestStoreFileHistory(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "ipxe_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a temporary iPXE file
	macAddress := "00:11:22:33:44:55"
	ipxeFilePath := filepath.Join(tempDir, "test.ipxe")
	initialContent := []byte("#!ipxe\necho Test Config\n")
	if err := os.WriteFile(ipxeFilePath, initialContent, 0644); err != nil {
		t.Fatalf("Failed to write temp iPXE file: %v", err)
	}

	// Store file history
	err = storeFileHistory(macAddress, ipxeFilePath)
	if err != nil {
		t.Fatalf("storeFileHistory failed: %v", err)
	}

	// Check if history file was created
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

	// Create a temporary iPXE file
	macAddress := "00:11:22:33:44:55"
	ipxeFilePath := filepath.Join(tempDir, macAddress+".ipxe")
	initialContent := []byte("#!ipxe\necho Initial Config\n")
	if err := os.WriteFile(ipxeFilePath, initialContent, 0644); err != nil {
		t.Fatalf("Failed to write temp iPXE file: %v", err)
	}

	// Mock viper configuration
	viper.Set("ipxe_folder", tempDir)
	viper.Set("boot_customization_folder", tempDir+"/boot")
	viper.Set("cloud_init_folder", tempDir+"/cloud-init")

	// Test updating the file
	err = UpdateIPXEFile(macAddress)
	if err != nil {
		t.Fatalf("UpdateIPXEFile failed: %v", err)
	}
}
