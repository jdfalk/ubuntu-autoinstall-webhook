package cmd

import (
	"database/sql"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jdfalk/ubuntu-autoinstall-webhook/internal/db"
	"github.com/jdfalk/ubuntu-autoinstall-webhook/internal/logger"
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

// WriteFile writes data to a file named by filename with the given permissions.
func (a AferoFs) WriteFile(filename string, data []byte, perm os.FileMode) error {
	return afero.WriteFile(a.Fs, filename, data, perm)
}

// ReadFile reads the contents of the file.
func (a AferoFs) ReadFile(filename string) ([]byte, error) {
	return afero.ReadFile(a.Fs, filename)
}

// UserHomeDir returns a fixed fake home directory for testing.
func (a AferoFs) UserHomeDir() (string, error) {
	return "fake-home", nil
}

// Remove removes the file with the given name.
func (a AferoFs) Remove(name string) error {
	return a.Fs.Remove(name)
}

func TestMain(m *testing.M) {
	// Initialize a dummy DB to avoid "DB not configured" errors.
	db.DB = new(sql.DB)
	logger.SetOutput(io.Discard)
	os.Exit(m.Run())
}

func TestInitConfigDefaults(t *testing.T) {
	// Reset viper for test isolation.
	viper.Reset()

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

	fs := AferoFs{afero.NewMemMapFs()}
	viper.SetFs(fs)

	// Create a fake config file with our desired fallback values.
	// Note: keys here match our defaults exactly.
	fakeConfig := `
port: "5000"
logDir: "test-log-dir"
logFile: "autoinstall_report.log"
logLevel: "INFO"
database:
  enabled: false
  type: postgres
  host: localhost
  port: 5432
  user: user
  password: password
  dbname: autoinstall
  sslmode: require
  max_open_conns: 100
  max_idle_conns: 10
  conn_max_lifetime: 3600
ipxe_folder: "ubuntu-autoinstall-webhook/ipxe"
boot_customization_folder: "ubuntu-autoinstall-webhook/boot"
cloud_init_folder: "ubuntu-autoinstall-webhook/cloud-init"
`
	configPath := filepath.Join(tempWorkDir, "config.yaml")
	// Set the global configFile variable so initConfig uses this file.
	configFile = configPath
	viper.SetConfigFile(configPath)
	if err := fs.WriteFile(configPath, []byte(fakeConfig), 0644); err != nil {
		t.Fatalf("Failed to write fake config file: %v", err)
	}

	initConfig(fs)

	if viper.GetString("port") != "5000" {
		t.Errorf("Expected port 5000, got %s", viper.GetString("port"))
	}
	if filepath.Clean(viper.GetString("logDir")) != filepath.Clean("test-log-dir") {
		t.Errorf("Expected logDir 'test-log-dir', got %s", viper.GetString("logDir"))
	}
	if viper.GetString("logFile") != "autoinstall_report.log" {
		t.Errorf("Expected logFile 'autoinstall_report.log', got %s", viper.GetString("logFile"))
	}
	if viper.GetString("database.host") != "localhost" {
		t.Errorf("Expected database.host 'localhost', got %s", viper.GetString("database.host"))
	}
	if viper.GetInt("database.port") != 5432 {
		t.Errorf("Expected database.port 5432, got %d", viper.GetInt("database.port"))
	}
	if viper.GetString("database.user") != "user" {
		t.Errorf("Expected database.user 'user', got %s", viper.GetString("database.user"))
	}
	if viper.GetString("database.password") != "password" {
		t.Errorf("Expected database.password 'password', got %s", viper.GetString("database.password"))
	}
	if viper.GetString("database.dbname") != "autoinstall" {
		t.Errorf("Expected database.dbname 'autoinstall', got %s", viper.GetString("database.dbname"))
	}
	if viper.GetString("database.sslmode") != "require" {
		t.Errorf("Expected database.sslmode 'require', got %s", viper.GetString("database.sslmode"))
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

	expectedIpxe := "ubuntu-autoinstall-webhook/ipxe"
	expectedBoot := "ubuntu-autoinstall-webhook/boot"
	expectedCloudInit := "ubuntu-autoinstall-webhook/cloud-init"
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

func TestValidatePaths(t *testing.T) {
	// Reset and set FS so that all operations occur in virtual memory.
	viper.Reset()

	fs := AferoFs{afero.NewMemMapFs()}
	viper.SetFs(fs)

	logPath := "fake-logDir"
	ipxePath := "fake-ipxe"
	bootPath := "fake-boot"
	cloudInitPath := "fake-cloud-init"

	viper.Set("logDir", logPath)
	viper.Set("ipxe_folder", ipxePath)
	viper.Set("boot_customization_folder", bootPath)
	viper.Set("cloud_init_folder", cloudInitPath)

	if err := validatePaths(fs); err != nil {
		t.Errorf("validatePaths failed: %v", err)
	}

	paths := []string{
		viper.GetString("logDir"),
		viper.GetString("ipxe_folder"),
		viper.GetString("boot_customization_folder"),
		viper.GetString("cloud_init_folder"),
	}
	for _, p := range paths {
		exists, err := afero.Exists(fs, p)
		if err != nil {
			t.Errorf("Error checking existence of path %s: %v", p, err)
		}
		if !exists {
			t.Errorf("Expected path %s to exist", p)
		}
	}
}

func TestProcessConfigFile(t *testing.T) {
	// Reset viper and set FS for virtual file system.
	viper.Reset()

	fs := AferoFs{afero.NewMemMapFs()}
	viper.SetFs(fs)

	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")
	// Set the global configFile variable.
	configFile = configPath
	viper.SetConfigFile(configPath)

	if err := afero.WriteFile(fs, configPath, []byte(""), 0644); err != nil {
		t.Fatalf("Failed to create temp config file: %v", err)
	}

	if err := processConfigFile(fs); err != nil {
		t.Errorf("processConfigFile failed: %v", err)
	}

	candidatePath := filepath.Join(tempDir, filepath.Base(configPath))
	content, err := afero.ReadFile(fs, candidatePath)
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

func TestWriteFileToCandidateLocations(t *testing.T) {
	// Reset viper and set FS for virtual file system.
	viper.Reset()

	fs := AferoFs{afero.NewMemMapFs()}
	viper.SetFs(fs)

	tempDir := t.TempDir()
	dummyFile := "tmp_test_candidate_file.txt"
	testContent := "Test content for candidate write"
	configPath := filepath.Join(tempDir, "dummyconfig.yaml")
	// Set the global configFile variable.
	configFile = configPath
	viper.SetConfigFile(configPath)

	if err := afero.WriteFile(fs, configPath, []byte("dummy"), 0644); err != nil {
		t.Fatalf("Failed to write dummy config file: %v", err)
	}

	if err := writeFileToCandidateLocations(fs, dummyFile, []byte(testContent)); err != nil {
		t.Errorf("writeFileToCandidateLocations failed: %v", err)
	}

	var candidateFound bool
	candidates := []string{
		filepath.Join(filepath.Dir(configPath), dummyFile),
		filepath.Join(tempDir, "ubuntu-autoinstall-webhook", dummyFile),
		filepath.Join(os.TempDir(), "ubuntu-autoinstall-webhook", dummyFile),
	}
	for _, candidatePath := range candidates {
		if exists, _ := afero.Exists(fs, candidatePath); exists {
			candidateFound = true
			data, err := afero.ReadFile(fs, candidatePath)
			if err != nil {
				t.Errorf("Failed to read candidate file %s: %v", candidatePath, err)
			}
			if string(data) != testContent {
				t.Errorf("Expected content %q in candidate %s, got %q", testContent, candidatePath, string(data))
			}
			break
		}
	}
	if !candidateFound {
		t.Errorf("No candidate file was created")
	}
}
