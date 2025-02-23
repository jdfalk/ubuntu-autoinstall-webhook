package cmd

import (
	"os"
	"testing"

	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

// AferoFs is a wrapper for afero.Fs that implements the FileSystem interface.
type AferoFs struct {
	afero.Fs
}

// Stat returns the FileInfo structure describing the file.
func (a AferoFs) Stat(name string) (os.FileInfo, error) {
	return a.Fs.Stat(name)
}

// MkdirAll creates a directory named path, along with any necessary parents.
func (a AferoFs) MkdirAll(path string, perm os.FileMode) error {
	return a.Fs.MkdirAll(path, perm)
}

// Test default values
func TestInitConfigDefaults(t *testing.T) {
	// Set up a temporary in-memory file system for testing
	fs := afero.NewMemMapFs()
	viper.SetFs(fs)

	tempDir := "/tmp/test"
	viper.Set("logDir", tempDir)

	initConfig()

	if viper.GetString("port") != "5000" {
		t.Errorf("Expected port 5000, got %s", viper.GetString("port"))
	}
	if viper.GetString("logDir") != tempDir {
		t.Errorf("Expected logDir %s, got %s", tempDir, viper.GetString("logDir"))
	}
	if viper.GetString("logFile") != "autoinstall_report.log" {
		t.Errorf("Expected logFile autoinstall_report.log, got %s", viper.GetString("logFile"))
	}
	if viper.GetString("database.host") != "localhost" {
		t.Errorf("Expected database.host localhost, got %s", viper.GetString("database.host"))
	}
	if viper.GetInt("database.port") != 5432 {
		t.Errorf("Expected database.port 5432, got %d", viper.GetInt("database.port"))
	}
	if viper.GetString("database.user") != "user" {
		t.Errorf("Expected database.user user, got %s", viper.GetString("database.user"))
	}
	if viper.GetString("database.password") != "password" {
		t.Errorf("Expected database.password password, got %s", viper.GetString("database.password"))
	}
	if viper.GetString("database.dbname") != "autoinstall" {
		t.Errorf("Expected database.dbname autoinstall, got %s", viper.GetString("database.dbname"))
	}
	if viper.GetString("database.sslmode") != "disable" {
		t.Errorf("Expected database.sslmode disable, got %s", viper.GetString("database.sslmode"))
	}
	if viper.GetInt("database.max_open_conns") != 100 {
		t.Errorf("Expected database.max_open_conns 100, got %d", viper.GetInt("database.max_open_conns"))
	}
	if viper.GetInt("database.max_idle_conns") != 10 {
		t.Errorf("Expected database.max_idle_conns 10, got %d", viper.GetInt("database.max_idle_conns"))
	}
	if viper.GetInt("database.conn_max_lifetime") != 3600 {
		t.Errorf("Expected database.conn_max_lifetime 3600, got %d", viper.GetInt("database.conn_max_lifetime"))
	}
	if viper.GetString("ipxe_folder") != "/var/www/html/ipxe" {
		t.Errorf("Expected ipxe_folder /var/www/html/ipxe, got %s", viper.GetString("ipxe_folder"))
	}
	if viper.GetString("boot_customization_folder") != "/var/www/html/ipxe/boot" {
		t.Errorf("Expected boot_customization_folder /var/www/html/ipxe/boot, got %s", viper.GetString("boot_customization_folder"))
	}
	if viper.GetString("cloud_init_folder") != "/var/www/html/cloud-init" {
		t.Errorf("Expected cloud_init_folder /var/www/html/cloud-init, got %s", viper.GetString("cloud_init_folder"))
	}
}

// Test path validation
func TestValidatePaths(t *testing.T) {
	// Set up a temporary in-memory file system for testing
	fs := afero.NewMemMapFs()
	viper.SetFs(fs)

	tempDir := "/tmp/test"
	viper.Set("logDir", tempDir)
	viper.Set("ipxe_folder", tempDir+"/ipxe")
	viper.Set("boot_customization_folder", tempDir+"/boot")
	viper.Set("cloud_init_folder", tempDir+"/cloud-init")

	if err := validatePaths(AferoFs{fs}); err != nil {
		t.Errorf("validatePaths failed: %v", err)
	}

	// Check if the directories were created
	paths := []string{
		viper.GetString("logDir"),
		viper.GetString("ipxe_folder"),
		viper.GetString("boot_customization_folder"),
		viper.GetString("cloud_init_folder"),
	}

	for _, p := range paths {
		info, err := fs.Stat(p)
		if err != nil {
			t.Errorf("Error accessing path %s: %v", p, err)
		} else if !info.IsDir() {
			t.Errorf("Path %s is not a directory", p)
		}
	}
}

// Test invalid paths
func TestInvalidPaths(t *testing.T) {
	// Set up a temporary in-memory file system for testing
	fs := afero.NewMemMapFs()
	viper.SetFs(fs)

	invalidPaths := []string{
		"../invalid",
		"~/invalid",
		"//invalid",
	}

	for _, p := range invalidPaths {
		viper.Set("logDir", p)
		err := validatePaths(AferoFs{fs})
		if err == nil {
			t.Errorf("Expected error for invalid path %s, but did not get one", p)
		}
	}
}
