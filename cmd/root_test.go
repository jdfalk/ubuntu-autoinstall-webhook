package cmd

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"github.com/jdfalk/ubuntu-autoinstall-webhook/internal/logger"
)

// TestMain configures the new logger for tests.
func TestMain(m *testing.M) {
	// Silence logger output during tests.
	logger.SetOutput(io.Discard)
	os.Exit(m.Run())
}

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

// TestInitConfigDefaults verifies that initConfig correctly sets default configuration values.
// Note: In this test environment the default critical folder paths (which are unwritable)
// fallback to the candidate under "./ubuntu-autoinstall-webhook".
// Therefore, the expected fallback values are now:
//
//	ipxe_folder: "./ubuntu-autoinstall-webhook/ipxe"
//	boot_customization_folder: "./ubuntu-autoinstall-webhook/boot"
//	cloud_init_folder: "./ubuntu-autoinstall-webhook/cloud-init"
func TestInitConfigDefaults(t *testing.T) {
	// Set a temporary working directory to avoid polluting the actual code directory.
	tempWorkDir := t.TempDir()
	oldWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}
	if err := os.Chdir(tempWorkDir); err != nil {
		t.Fatalf("Failed to change working directory: %v", err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(oldWD); err != nil {
			t.Fatalf("Failed to restore working directory: %v", err)
		}
	})

	// Use an in-memory filesystem for the config file.
	fs := afero.NewMemMapFs()
	viper.SetFs(fs)

	// Set a dummy logDir (which we'll let pass through without fallback).
	// For the other keys, the defaults in initConfig (e.g. "/var/www/html/ipxe") are unwritable,
	// so fallback candidates will be used.
	tempDir := "test-log-dir"
	viper.Set("logDir", tempDir)

	// Use a dummy config file in the in-memory FS.
	configPath := filepath.Join(tempWorkDir, "config.yaml")
	viper.SetConfigFile(configPath)
	// Write an empty config file.
	afero.WriteFile(fs, configPath, []byte(""), 0644)

	initConfig(AferoFs{fs})

	// Basic configuration defaults.
	if viper.GetString("port") != "5000" {
		t.Errorf("Expected port 5000, got %s", viper.GetString("port"))
	}
	if filepath.Clean(viper.GetString("logDir")) != filepath.Clean(tempDir) {
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
	// Update expected fallback values to reflect the new candidate.
	expectedIpxe := filepath.Join("ubuntu-autoinstall-webhook", "ipxe")
	expectedBoot := filepath.Join("ubuntu-autoinstall-webhook", "boot")
	expectedCloudInit := filepath.Join("ubuntu-autoinstall-webhook", "cloud-init")
	if filepath.Clean(viper.GetString("ipxe_folder")) != filepath.Clean(expectedIpxe) {
		t.Errorf("Expected ipxe_folder %q, got %s", expectedIpxe, viper.GetString("ipxe_folder"))
	}
	if filepath.Clean(viper.GetString("boot_customization_folder")) != filepath.Clean(expectedBoot) {
		t.Errorf("Expected boot_customization_folder %q, got %s", expectedBoot, viper.GetString("boot_customization_folder"))
	}
	if filepath.Clean(viper.GetString("cloud_init_folder")) != filepath.Clean(expectedCloudInit) {
		t.Errorf("Expected cloud_init_folder %q, got %s", expectedCloudInit, viper.GetString("cloud_init_folder"))
	}
}

// TestValidatePaths verifies that validatePaths creates and validates critical directories using fallback logic.
// Because ensureFolderExists uses OS calls, we use t.TempDir() (a real disk location) here.
func TestValidatePaths(t *testing.T) {
	baseTempDir := t.TempDir()
	logPath := filepath.Join(baseTempDir, "logDir")
	ipxePath := filepath.Join(baseTempDir, "ipxe")
	bootPath := filepath.Join(baseTempDir, "boot")
	cloudInitPath := filepath.Join(baseTempDir, "cloud-init")

	viper.Set("logDir", logPath)
	viper.Set("ipxe_folder", ipxePath)
	viper.Set("boot_customization_folder", bootPath)
	viper.Set("cloud_init_folder", cloudInitPath)

	if err := validatePaths(OsFs{}); err != nil {
		t.Errorf("validatePaths failed: %v", err)
	}

	paths := []string{
		viper.GetString("logDir"),
		viper.GetString("ipxe_folder"),
		viper.GetString("boot_customization_folder"),
		viper.GetString("cloud_init_folder"),
	}
	for _, p := range paths {
		info, err := os.Stat(p)
		if err != nil {
			t.Errorf("Error accessing path %s: %v", p, err)
		} else if !info.IsDir() {
			t.Errorf("Path %s is not a directory", p)
		}
	}
}

// TestInvalidPaths verifies that invalid paths are correctly rejected.
// This test temporarily changes the working directory so that fallback file/directory creation occurs in a temp folder.
func TestInvalidPaths(t *testing.T) {
	tempWorkDir := t.TempDir()
	oldWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}
	if err := os.Chdir(tempWorkDir); err != nil {
		t.Fatalf("Failed to change working directory: %v", err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(oldWD); err != nil {
			t.Fatalf("Failed to restore working directory: %v", err)
		}
	})

	invalidPaths := []string{
		"../invalid",
		"~/invalid",
		"//invalid",
	}
	for _, p := range invalidPaths {
		viper.Set("logDir", p)
		if err := validatePaths(OsFs{}); err == nil {
			t.Errorf("Expected error for invalid path %s, but did not get one", p)
		}
	}
}

// TestProcessConfigFile tests that processConfigFile correctly organizes the configuration file.
// It creates a temporary config file, invokes processConfigFile, and checks for expected section headers.
func TestProcessConfigFile(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")
	if err := os.WriteFile(configPath, []byte(""), 0644); err != nil {
		t.Fatalf("Failed to create temp config file: %v", err)
	}
	viper.SetConfigFile(configPath)
	t.Cleanup(func() { os.Remove(configPath) })

	if err := processConfigFile(); err != nil {
		t.Errorf("processConfigFile failed: %v", err)
	}

	candidatePath := filepath.Join(tempDir, filepath.Base(configPath))
	content, err := os.ReadFile(candidatePath)
	if err != nil {
		t.Fatalf("Failed to read candidate file %s: %v", candidatePath, err)
	}
	expectedSections := []string{
		"# Logging Configuration",
		"# Database Configuration",
		"# iPXE Settings",
		"# Cloud-Init Settings",
	}
	for _, section := range expectedSections {
		if !strings.Contains(string(content), section) {
			t.Errorf("Expected section header %q in processed config file, got: %s", section, string(content))
		}
	}
}

// TestWriteFileToCandidateLocations tests that writeFileToCandidateLocations writes data to the first valid candidate.
// It also cleans up any files written to avoid polluting the code directory.
func TestWriteFileToCandidateLocations(t *testing.T) {
	tempDir := t.TempDir()
	dummyFile := "tmp_test_candidate_file.txt"
	testContent := "Test content for candidate write"
	configPath := filepath.Join(tempDir, "dummyconfig.yaml")
	viper.SetConfigFile(configPath)

	t.Cleanup(func() {
		os.Remove(filepath.Join(filepath.Dir(configPath), dummyFile))
		os.Remove(filepath.Join(tempDir, "ubuntu-autoinstall-webhook", dummyFile))
		os.Remove(filepath.Join(os.TempDir(), "ubuntu-autoinstall-webhook", dummyFile))
	})

	if err := writeFileToCandidateLocations(dummyFile, []byte(testContent)); err != nil {
		t.Errorf("writeFileToCandidateLocations failed: %v", err)
	}
	candidate1 := filepath.Join(filepath.Dir(configPath), dummyFile)
	candidate2 := filepath.Join(tempDir, "ubuntu-autoinstall-webhook", dummyFile)
	var usedCandidate string
	if _, err := os.Stat(candidate1); err == nil {
		usedCandidate = candidate1
	} else if _, err := os.Stat(candidate2); err == nil {
		usedCandidate = candidate2
	}
	if usedCandidate == "" {
		t.Errorf("No candidate file was created")
		return
	}
	data, err := os.ReadFile(usedCandidate)
	if err != nil {
		t.Errorf("Failed to read file from candidate %s: %v", usedCandidate, err)
	}
	if string(data) != testContent {
		t.Errorf("Expected content %q, got %q", testContent, string(data))
	}
}
