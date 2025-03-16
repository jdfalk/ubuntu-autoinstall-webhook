// internal/fileeditor/service_test.go
package fileeditor

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace/noop"
)

// setupTestService creates a test service with a filesystem that supports symlinks
func setupTestService(t *testing.T) (*Service, afero.Fs, string) {
	t.Helper()

	// Set up test configuration
	viper.Set("fileeditor.ipxe_dir", "/var/www/html/ipxe/boot")
	viper.Set("fileeditor.cloudinit_dir", "/var/lib/cloud-init")

	// Create a temporary directory for testing (using os instead of ioutil)
	tempDir, err := os.MkdirTemp("", "fileeditor-test")
	require.NoError(t, err)

	// Clean up the temporary directory after the test
	t.Cleanup(func() {
		os.RemoveAll(tempDir)
	})

	// Create a filesystem using the temporary directory as base path
	// This gives us real symlink support while isolating from the real filesystem
	fs := afero.NewBasePathFs(afero.NewOsFs(), tempDir)

	// Get the underlying OsFs for symlink operations
	osFs := afero.NewOsFs().(*afero.OsFs)

	// Ensure filesystem supports symlinks
	if _, ok := fs.(afero.Symlinker); !ok {
		t.Fatal("filesystem does not support symlinks")
	}

	if _, ok := fs.(afero.LinkReader); !ok {
		t.Fatal("filesystem does not support reading symlinks")
	}

	// Create a mock tracer
	mockTracer := noop.NewTracerProvider().Tracer("test")

	// Create the service with the filesystem
	service := &Service{
		fs:            fs,
		osFs:          osFs,
		ipxeDir:       viper.GetString("fileeditor.ipxe_dir"),
		cloudInitDir:  viper.GetString("fileeditor.cloudinit_dir"),
		isLeader:      true, // Always a leader in tests
		tracer:        mockTracer,
		leaderLockKey: "test-key",
	}

	return service, fs, tempDir
}

// Helper function to check symlinks using afero interfaces
func checkSymlink(t *testing.T, fs afero.Fs, linkPath string, expectedTarget string) {
	// Check if path exists
	exists, err := afero.Exists(fs, linkPath)
	require.NoError(t, err)
	require.True(t, exists, "Symlink does not exist: %s", linkPath)

	// Check if it's a symlink
	lstater, ok := fs.(afero.Lstater)
	require.True(t, ok, "filesystem does not support Lstater")

	info, _, err := lstater.LstatIfPossible(linkPath)
	require.NoError(t, err)
	require.True(t, info.Mode()&os.ModeSymlink != 0, "Not a symlink: %s", linkPath)

	// Check the target
	linkReader, ok := fs.(afero.LinkReader)
	require.True(t, ok, "filesystem does not support LinkReader")

	target, err := linkReader.ReadlinkIfPossible(linkPath)
	require.NoError(t, err)

	// When using BasePathFs, the symlink target will include the temp directory
	// So we'll just check if the target ends with our expected path
	assert.True(t, strings.HasSuffix(target, expectedTarget),
		"Expected symlink target to end with %s, got %s", expectedTarget, target)
}

func TestValidateIpxeContent(t *testing.T) {
	service, _, _ := setupTestService(t)
	ctx := context.Background()

	testCases := []struct {
		name     string
		content  []byte
		expectOk bool
	}{
		{
			name:     "Valid iPXE content",
			content:  []byte("#!ipxe\necho Starting iPXE boot\nchain http://server/boot.ipxe"),
			expectOk: true,
		},
		{
			name:     "Empty content",
			content:  []byte{},
			expectOk: false,
		},
		{
			name:     "Missing shebang",
			content:  []byte("echo Starting iPXE boot\nchain http://server/boot.ipxe"),
			expectOk: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := service.validateIpxeContent(ctx, tc.content)
			if tc.expectOk {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestValidateCloudInitContent(t *testing.T) {
	service, _, _ := setupTestService(t)
	ctx := context.Background()

	testCases := []struct {
		name     string
		fileType string
		content  []byte
		expectOk bool
	}{
		{
			name:     "Valid user-data",
			fileType: "user-data",
			content:  []byte("#cloud-config\npackages:\n  - htop"),
			expectOk: true,
		},
		{
			name:     "Invalid user-data",
			fileType: "user-data",
			content:  []byte("packages:\n  - htop"),
			expectOk: false,
		},
		{
			name:     "Valid variables.sh with bash",
			fileType: "variables.sh",
			content:  []byte("#!/bin/bash\nexport VAR=value"),
			expectOk: true,
		},
		{
			name:     "Valid variables.sh with sh",
			fileType: "variables.sh",
			content:  []byte("#!/bin/sh\nexport VAR=value"),
			expectOk: true,
		},
		{
			name:     "Invalid variables.sh",
			fileType: "variables.sh",
			content:  []byte("export VAR=value"),
			expectOk: false,
		},
		{
			name:     "Valid meta-data",
			fileType: "meta-data",
			content:  []byte("instance-id: test-instance"),
			expectOk: true,
		},
		{
			name:     "Unknown filetype",
			fileType: "unknown",
			content:  []byte("content"),
			expectOk: false,
		},
		{
			name:     "Empty content",
			fileType: "user-data",
			content:  []byte{},
			expectOk: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := service.validateCloudInitContent(ctx, tc.fileType, tc.content)
			if tc.expectOk {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestWriteIpxeFile(t *testing.T) {
	service, fs, _ := setupTestService(t)
	ctx := context.Background()

	macAddress := "00:11:22:33:44:55"
	content := []byte("#!ipxe\necho Hello World")

	// Write the file
	err := service.WriteIpxeFile(ctx, macAddress, content)
	require.NoError(t, err)

	// Check if the file was written correctly
	expectedPath := filepath.Join(service.ipxeDir, "mac-00-11-22-33-44-55.ipxe")
	exists, err := afero.Exists(fs, expectedPath)
	require.NoError(t, err)
	assert.True(t, exists)

	// Check the content
	fileContent, err := afero.ReadFile(fs, expectedPath)
	require.NoError(t, err)
	assert.Equal(t, content, fileContent)
}

func TestWriteIpxeFileInvalidContent(t *testing.T) {
	service, _, _ := setupTestService(t)
	ctx := context.Background()

	macAddress := "00:11:22:33:44:55"
	invalidContent := []byte("invalid content")

	// Write the file with invalid content
	err := service.WriteIpxeFile(ctx, macAddress, invalidContent)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid iPXE content")
}

func TestCreateCloudInitDirs(t *testing.T) {
	service, fs, _ := setupTestService(t)
	ctx := context.Background()

	macAddress := "00:11:22:33:44:55"
	hostname := "test-host"

	// Create the directories
	err := service.CreateCloudInitDirs(ctx, macAddress, hostname)
	require.NoError(t, err)

	// Check if MAC directories exist
	normalizedMac := "00-11-22-33-44-55"
	macDir := filepath.Join(service.cloudInitDir, normalizedMac)
	macInstallDir := filepath.Join(service.cloudInitDir, normalizedMac+"_install")

	exists, err := afero.DirExists(fs, macDir)
	require.NoError(t, err)
	assert.True(t, exists)

	exists, err = afero.DirExists(fs, macInstallDir)
	require.NoError(t, err)
	assert.True(t, exists)

	// Check if hostname symlinks exist and point to the correct directories
	hostnameDir := filepath.Join(service.cloudInitDir, hostname)
	hostnameInstallDir := filepath.Join(service.cloudInitDir, hostname+"_install")

	// Use our helper function instead of direct OS calls
	checkSymlink(t, fs, hostnameDir, macDir)
	checkSymlink(t, fs, hostnameInstallDir, macInstallDir)
}

func TestCreateCloudInitDirsExistingHostname(t *testing.T) {
	service, fs, _ := setupTestService(t)
	ctx := context.Background()

	// Create first MAC directory and hostname symlink
	firstMac := "00:11:22:33:44:55"
	hostname := "test-host"
	err := service.CreateCloudInitDirs(ctx, firstMac, hostname)
	require.NoError(t, err)

	// Try to create another MAC directory with same hostname
	secondMac := "aa:bb:cc:dd:ee:ff"
	err = service.CreateCloudInitDirs(ctx, secondMac, hostname)
	require.NoError(t, err) // It should succeed by overwriting the symlink

	// Check if the hostname symlink now points to the second MAC
	normalizedSecondMac := "aa-bb-cc-dd-ee-ff"
	secondMacDir := filepath.Join(service.cloudInitDir, normalizedSecondMac)
	hostnameDir := filepath.Join(service.cloudInitDir, hostname)

	// Use our helper function instead of direct OS calls
	checkSymlink(t, fs, hostnameDir, secondMacDir)
}

func TestWriteCloudInitFile(t *testing.T) {
	service, fs, _ := setupTestService(t)
	ctx := context.Background()

	macAddress := "00:11:22:33:44:55"
	normalizedMac := "00-11-22-33-44-55"

	// Create the directories first
	macDir := filepath.Join(service.cloudInitDir, normalizedMac)
	macInstallDir := filepath.Join(service.cloudInitDir, normalizedMac+"_install")

	err := fs.MkdirAll(macDir, 0755)
	require.NoError(t, err)

	err = fs.MkdirAll(macInstallDir, 0755)
	require.NoError(t, err)

	testCases := []struct {
		name      string
		fileType  string
		content   []byte
		isInstall bool
	}{
		{
			name:      "Write user-data",
			fileType:  "user-data",
			content:   []byte("#cloud-config\npackages:\n  - htop"),
			isInstall: false,
		},
		{
			name:      "Write meta-data",
			fileType:  "meta-data",
			content:   []byte("instance-id: test-instance"),
			isInstall: false,
		},
		{
			name:      "Write variables.sh",
			fileType:  "variables.sh",
			content:   []byte("#!/bin/bash\nexport VAR=value"),
			isInstall: false,
		},
		{
			name:      "Write install user-data",
			fileType:  "user-data_install",
			content:   []byte("#cloud-config\npackages:\n  - htop"),
			isInstall: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Write the file
			err := service.WriteCloudInitFile(ctx, macAddress, tc.fileType, tc.content)
			require.NoError(t, err)

			// Check if the file was written correctly
			var expectedPath string
			fileType := tc.fileType
			if tc.isInstall {
				fileType = strings.TrimSuffix(tc.fileType, "_install")
				expectedPath = filepath.Join(service.cloudInitDir, normalizedMac+"_install", fileType)
			} else {
				expectedPath = filepath.Join(service.cloudInitDir, normalizedMac, fileType)
			}

			exists, err := afero.Exists(fs, expectedPath)
			require.NoError(t, err)
			assert.True(t, exists)

			// Check the content
			fileContent, err := afero.ReadFile(fs, expectedPath)
			require.NoError(t, err)
			assert.Equal(t, tc.content, fileContent)
		})
	}
}

func TestListFiles(t *testing.T) {
	service, fs, _ := setupTestService(t)
	ctx := context.Background()

	// Create some test files
	err := fs.MkdirAll(service.ipxeDir, 0755)
	require.NoError(t, err)

	ipxeFiles := []string{
		"mac-00-11-22-33-44-55.ipxe",
		"mac-aa-bb-cc-dd-ee-ff.ipxe",
	}

	for _, file := range ipxeFiles {
		err := afero.WriteFile(fs, filepath.Join(service.ipxeDir, file), []byte("#!ipxe"), 0644)
		require.NoError(t, err)
	}

	// List the files
	files, err := service.ListFiles(ctx, "ipxe")
	require.NoError(t, err)
	assert.ElementsMatch(t, ipxeFiles, files)

	// Test with unknown type
	_, err = service.ListFiles(ctx, "unknown")
	assert.Error(t, err)
}

func TestReadFile(t *testing.T) {
	service, fs, _ := setupTestService(t)
	ctx := context.Background()

	// Create test files
	ipxeContent := []byte("#!ipxe\necho Hello")
	cloudInitContent := []byte("#cloud-config\npackages:\n  - htop")

	err := fs.MkdirAll(service.ipxeDir, 0755)
	require.NoError(t, err)

	err = fs.MkdirAll(filepath.Join(service.cloudInitDir, "00-11-22-33-44-55"), 0755)
	require.NoError(t, err)

	ipxeFile := filepath.Join(service.ipxeDir, "mac-00-11-22-33-44-55.ipxe")
	cloudInitFile := filepath.Join(service.cloudInitDir, "00-11-22-33-44-55", "user-data")

	err = afero.WriteFile(fs, ipxeFile, ipxeContent, 0644)
	require.NoError(t, err)

	err = afero.WriteFile(fs, cloudInitFile, cloudInitContent, 0644)
	require.NoError(t, err)

	// Read iPXE file
	content, err := service.ReadFile(ctx, "ipxe", "mac-00-11-22-33-44-55.ipxe")
	require.NoError(t, err)
	assert.Equal(t, ipxeContent, content)

	// Read cloud-init file
	content, err = service.ReadFile(ctx, "cloudinit", "00-11-22-33-44-55/user-data")
	require.NoError(t, err)
	assert.Equal(t, cloudInitContent, content)

	// Test with non-existent file
	_, err = service.ReadFile(ctx, "ipxe", "non-existent.ipxe")
	assert.Error(t, err)

	// Test with invalid cloudinit path
	_, err = service.ReadFile(ctx, "cloudinit", "invalid-format")
	assert.Error(t, err)

	// Test with unknown type
	_, err = service.ReadFile(ctx, "unknown", "file")
	assert.Error(t, err)
}

func TestDeleteFile(t *testing.T) {
	service, fs, _ := setupTestService(t)
	ctx := context.Background()

	// Create test files
	err := fs.MkdirAll(service.ipxeDir, 0755)
	require.NoError(t, err)

	err = fs.MkdirAll(filepath.Join(service.cloudInitDir, "00-11-22-33-44-55"), 0755)
	require.NoError(t, err)

	ipxeFile := filepath.Join(service.ipxeDir, "mac-00-11-22-33-44-55.ipxe")
	cloudInitFile := filepath.Join(service.cloudInitDir, "00-11-22-33-44-55", "user-data")

	err = afero.WriteFile(fs, ipxeFile, []byte("#!ipxe"), 0644)
	require.NoError(t, err)

	err = afero.WriteFile(fs, cloudInitFile, []byte("#cloud-config"), 0644)
	require.NoError(t, err)

	// Delete iPXE file
	err = service.DeleteFile(ctx, "ipxe", "mac-00-11-22-33-44-55.ipxe")
	require.NoError(t, err)

	exists, err := afero.Exists(fs, ipxeFile)
	require.NoError(t, err)
	assert.False(t, exists)

	// Delete cloud-init file
	err = service.DeleteFile(ctx, "cloudinit", "00-11-22-33-44-55/user-data")
	require.NoError(t, err)

	exists, err = afero.Exists(fs, cloudInitFile)
	require.NoError(t, err)
	assert.False(t, exists)

	// Test with invalid cloudinit path
	err = service.DeleteFile(ctx, "cloudinit", "invalid-format")
	assert.Error(t, err)

	// Test with unknown type
	err = service.DeleteFile(ctx, "unknown", "file")
	assert.Error(t, err)
}

func TestDeleteCloudInitDir(t *testing.T) {
	service, fs, _ := setupTestService(t)
	ctx := context.Background()

	// Get a Symlinker and Lstater from the fs
	symlinker, ok := fs.(afero.Symlinker)
	require.True(t, ok, "filesystem must support symlinks")

	_, ok = fs.(afero.Lstater)
	require.True(t, ok, "filesystem must support Lstater")

	// Ensure recycle bin exists
	recycleBinPath := filepath.Join(service.cloudInitDir, "recycle_bin")
	err := fs.MkdirAll(recycleBinPath, 0755)
	require.NoError(t, err)

	// Delete by MAC address
	t.Run("Delete by MAC", func(t *testing.T) {
		// Use unique paths for this subtest
		macAddress := "00:11:22:33:44:55"
		hostname := "test-host-mac"
		normalizedMac := "00-11-22-33-44-55"

		// Create fresh directories and symlinks
		macDir := filepath.Join(service.cloudInitDir, normalizedMac)
		macInstallDir := filepath.Join(service.cloudInitDir, normalizedMac+"_install")
		hostnameDir := filepath.Join(service.cloudInitDir, hostname)
		hostnameInstallDir := filepath.Join(service.cloudInitDir, hostname+"_install")

		// Clean up from previous test runs
		_ = fs.RemoveAll(macDir)
		_ = fs.RemoveAll(macInstallDir)
		_ = fs.RemoveAll(hostnameDir)
		_ = fs.RemoveAll(hostnameInstallDir)

		// Clear recycle bin
		entries, _ := afero.ReadDir(fs, recycleBinPath)
		for _, entry := range entries {
			_ = fs.RemoveAll(filepath.Join(recycleBinPath, entry.Name()))
		}

		err := fs.MkdirAll(macDir, 0755)
		require.NoError(t, err)

		err = fs.MkdirAll(macInstallDir, 0755)
		require.NoError(t, err)

		err = symlinker.SymlinkIfPossible(macDir, hostnameDir)
		require.NoError(t, err)

		err = symlinker.SymlinkIfPossible(macInstallDir, hostnameInstallDir)
		require.NoError(t, err)

		// Now run the actual test
		err = service.DeleteCloudInitDir(ctx, macAddress)
		require.NoError(t, err)

		// Verify symlinks are deleted
		exists, err := afero.Exists(fs, hostnameDir)
		require.NoError(t, err)
		assert.False(t, exists)

		exists, err = afero.Exists(fs, hostnameInstallDir)
		require.NoError(t, err)
		assert.False(t, exists)

		// Verify original directories are gone
		exists, err = afero.Exists(fs, macDir)
		require.NoError(t, err)
		assert.False(t, exists)

		exists, err = afero.Exists(fs, macInstallDir)
		require.NoError(t, err)
		assert.False(t, exists)

		// Check that directories were moved to recycle bin
		entries, err = afero.ReadDir(fs, recycleBinPath)
		require.NoError(t, err)

		// Should have 2 entries in recycle bin: macDir and macInstallDir
		assert.Equal(t, 2, len(entries))

		// Check naming pattern (should contain normalized MAC and delete_me)
		macDirMoved := false
		macInstallDirMoved := false

		for _, entry := range entries {
			name := entry.Name()
			if strings.HasPrefix(name, normalizedMac+"_delete_me_") {
				macDirMoved = true
			} else if strings.HasPrefix(name, normalizedMac+"_install_delete_me_") {
				macInstallDirMoved = true
			}
		}

		assert.True(t, macDirMoved, "MAC directory not moved to recycle bin")
		assert.True(t, macInstallDirMoved, "MAC install directory not moved to recycle bin")
	})

	// Complete the "Delete by Hostname" test
	t.Run("Delete by Hostname", func(t *testing.T) {
		// Use completely different paths from the first test
		hostname := "test-host-hostname"
		normalizedMac := "11-22-33-44-55-66"

		// Create fresh directories with different names
		macDir := filepath.Join(service.cloudInitDir, normalizedMac)
		macInstallDir := filepath.Join(service.cloudInitDir, normalizedMac+"_install")
		hostnameDir := filepath.Join(service.cloudInitDir, hostname)
		hostnameInstallDir := filepath.Join(service.cloudInitDir, hostname+"_install")

		// Clean up any existing files from previous test runs
		_ = fs.RemoveAll(macDir)
		_ = fs.RemoveAll(macInstallDir)
		_ = fs.RemoveAll(hostnameDir)
		_ = fs.RemoveAll(hostnameInstallDir)

		// Clear recycle bin
		entries, _ := afero.ReadDir(fs, recycleBinPath)
		for _, entry := range entries {
			_ = fs.RemoveAll(filepath.Join(recycleBinPath, entry.Name()))
		}

		// Create the directories
		err := fs.MkdirAll(macDir, 0755)
		require.NoError(t, err)

		err = fs.MkdirAll(macInstallDir, 0755)
		require.NoError(t, err)

		// Create the symlinks
		err = symlinker.SymlinkIfPossible(macDir, hostnameDir)
		require.NoError(t, err)

		err = symlinker.SymlinkIfPossible(macInstallDir, hostnameInstallDir)
		require.NoError(t, err)

		// Add debugging before deletion
		if linkReader, ok := fs.(afero.LinkReader); ok {
			target, _ := linkReader.ReadlinkIfPossible(hostnameDir)
			fmt.Printf("Hostname symlink target: %s\n", target)
			fmt.Printf("MAC dir: %s\n", macDir)
			fmt.Printf("Basename of target: %s\n", filepath.Base(target))
		}

		// Call DeleteCloudInitDir
		err = service.DeleteCloudInitDir(ctx, hostname)
		require.NoError(t, err)

		// Verify symlinks are deleted
		exists, err := afero.Exists(fs, hostnameDir)
		require.NoError(t, err)
		assert.False(t, exists)

		exists, err = afero.Exists(fs, hostnameInstallDir)
		require.NoError(t, err)
		assert.False(t, exists)

		// Verify original directories are gone
		exists, err = afero.Exists(fs, macDir)
		require.NoError(t, err)
		assert.False(t, exists)

		exists, err = afero.Exists(fs, macInstallDir)
		require.NoError(t, err)
		assert.False(t, exists)

		// Check that directories were moved to recycle bin
		entries, err = afero.ReadDir(fs, recycleBinPath)
		require.NoError(t, err)

		// Should have 2 entries in recycle bin: macDir and macInstallDir
		assert.Equal(t, 2, len(entries))

		// Check naming pattern
		macDirMoved := false
		macInstallDirMoved := false

		for _, entry := range entries {
			name := entry.Name()
			if strings.HasPrefix(name, normalizedMac+"_delete_me_") {
				macDirMoved = true
			} else if strings.HasPrefix(name, normalizedMac+"_install_delete_me_") {
				macInstallDirMoved = true
			}
		}

		assert.True(t, macDirMoved, "MAC directory not moved to recycle bin")
		assert.True(t, macInstallDirMoved, "MAC install directory not moved to recycle bin")
	})

	// Complete the "Multiple Symlinks" test
	t.Run("Multiple Symlinks", func(t *testing.T) {
		// Use another set of unique paths
		macAddress := "aa:bb:cc:dd:ee:ff"
		hostname := "test-host-multiple"
		otherHostname := "other-host-multiple"
		normalizedMac := service.normalizeMacAddress(macAddress)

		// Create fresh directories and symlinks
		macDir := filepath.Join(service.cloudInitDir, normalizedMac)
		macInstallDir := filepath.Join(service.cloudInitDir, normalizedMac+"_install")
		hostnameDir := filepath.Join(service.cloudInitDir, hostname)
		hostnameInstallDir := filepath.Join(service.cloudInitDir, hostname+"_install")
		otherHostnameDir := filepath.Join(service.cloudInitDir, otherHostname)
		otherHostnameInstallDir := filepath.Join(service.cloudInitDir, otherHostname+"_install")

		// Clean up any existing files from previous test runs
		_ = fs.RemoveAll(macDir)
		_ = fs.RemoveAll(macInstallDir)
		_ = fs.RemoveAll(hostnameDir)
		_ = fs.RemoveAll(hostnameInstallDir)
		_ = fs.RemoveAll(otherHostnameDir)
		_ = fs.RemoveAll(otherHostnameInstallDir)

		// Clear recycle bin - Make sure this runs before checking the bin later!
		entries, _ := afero.ReadDir(fs, recycleBinPath)
		for _, entry := range entries {
			_ = fs.RemoveAll(filepath.Join(recycleBinPath, entry.Name()))
		}

		// Verify recycle bin is empty at start
		entries, err = afero.ReadDir(fs, recycleBinPath)
		require.NoError(t, err)
		assert.Equal(t, 0, len(entries), "Recycle bin should be empty at test start")

		// Create the directories
		err = fs.MkdirAll(macDir, 0755)
		require.NoError(t, err)

		err = fs.MkdirAll(macInstallDir, 0755)
		require.NoError(t, err)

		// Create multiple symlinks to the same MAC address
		err = symlinker.SymlinkIfPossible(macDir, hostnameDir)
		require.NoError(t, err)

		err = symlinker.SymlinkIfPossible(macInstallDir, hostnameInstallDir)
		require.NoError(t, err)

		err = symlinker.SymlinkIfPossible(macDir, otherHostnameDir)
		require.NoError(t, err)

		err = symlinker.SymlinkIfPossible(macInstallDir, otherHostnameInstallDir)
		require.NoError(t, err)

		// Delete one hostname only
		err = service.DeleteCloudInitDir(ctx, hostname)
		require.NoError(t, err)

		// First hostname symlinks should be gone
		exists, err := afero.Exists(fs, hostnameDir)
		require.NoError(t, err)
		assert.False(t, exists)

		exists, err = afero.Exists(fs, hostnameInstallDir)
		require.NoError(t, err)
		assert.False(t, exists)

		// But MAC directories should still exist (referenced by other symlink)
		exists, err = afero.Exists(fs, macDir)
		require.NoError(t, err)
		assert.True(t, exists)

		exists, err = afero.Exists(fs, macInstallDir)
		require.NoError(t, err)
		assert.True(t, exists)

		// Other hostname symlinks should still exist
		exists, err = afero.Exists(fs, otherHostnameDir)
		require.NoError(t, err)
		assert.True(t, exists)

		exists, err = afero.Exists(fs, otherHostnameInstallDir)
		require.NoError(t, err)
		assert.True(t, exists)

		// Check recycle bin - should be empty since MAC dirs are still referenced
		entries, err = afero.ReadDir(fs, recycleBinPath)
		require.NoError(t, err)
		assert.Equal(t, 0, len(entries), "Recycle bin should be empty")
	})

	// Add a "Non-existent" test case
	t.Run("Non-existent", func(t *testing.T) {
		err := service.DeleteCloudInitDir(ctx, "non-existent-hostname")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no cloud-init directory or symlink found")
	})

	// Test CleanupRecycleBin function
	t.Run("CleanupRecycleBin", func(t *testing.T) {
		// First put some test files in the recycle bin
		testFile1 := filepath.Join(recycleBinPath, "test-mac_delete_me_20250313")
		testFile2 := filepath.Join(recycleBinPath, "test-mac-install_delete_me_20250313")

		err := afero.WriteFile(fs, testFile1, []byte("test content"), 0644)
		require.NoError(t, err)

		err = afero.WriteFile(fs, testFile2, []byte("test content"), 0644)
		require.NoError(t, err)

		// Verify files exist
		exists, err := afero.Exists(fs, testFile1)
		require.NoError(t, err)
		assert.True(t, exists)

		exists, err = afero.Exists(fs, testFile2)
		require.NoError(t, err)
		assert.True(t, exists)

		// Call CleanupRecycleBin
		err = service.CleanupRecycleBin(ctx)
		require.NoError(t, err)

		// Verify recycle bin is empty
		entries, err := afero.ReadDir(fs, recycleBinPath)
		require.NoError(t, err)
		assert.Equal(t, 0, len(entries), "Recycle bin should be empty after cleanup")

		// Verify recycle bin directory still exists
		exists, err = afero.DirExists(fs, recycleBinPath)
		require.NoError(t, err)
		assert.True(t, exists, "Recycle bin directory should still exist after cleanup")
	})
}

func TestLeadership(t *testing.T) {
	service, _, _ := setupTestService(t)
	ctx := context.Background()

	// Initially not a leader
	service.isLeader = false

	// Acquire leadership
	acquired, err := service.AcquireLeadership(ctx)
	require.NoError(t, err)
	assert.True(t, acquired)
	assert.True(t, service.isLeader)

	// Release leadership
	err = service.ReleaseLeadership(ctx)
	require.NoError(t, err)
	assert.False(t, service.isLeader)
}

func TestEnsureDirectory(t *testing.T) {
	service, fs, _ := setupTestService(t)
	ctx := context.Background()

	// Test creating a new directory
	testDir := "/test/directory"
	err := service.ensureDirectory(ctx, testDir)
	require.NoError(t, err)

	exists, err := afero.Exists(fs, testDir)
	require.NoError(t, err)
	assert.True(t, exists)

	// Test with existing directory
	err = service.ensureDirectory(ctx, testDir)
	require.NoError(t, err)

	// Test with existing file (not a directory)
	filePath := "/test/file"
	err = afero.WriteFile(fs, filePath, []byte("test"), 0644)
	require.NoError(t, err)

	err = service.ensureDirectory(ctx, filePath)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "exists but is not a directory")
}

func TestNormalizeMacAddress(t *testing.T) {
	service, _, _ := setupTestService(t)

	testCases := []struct {
		input    string
		expected string
	}{
		{
			input:    "00:11:22:33:44:55",
			expected: "00-11-22-33-44-55",
		},
		{
			input:    "00-11-22-33-44-55",
			expected: "00-11-22-33-44-55",
		},
		{
			input:    "AA:BB:CC:DD:EE:FF",
			expected: "aa-bb-cc-dd-ee-ff",
		},
		{
			input:    "00:11:22:33:44:55:66:77",
			expected: "00-11-22-33-44-55-66-77",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result := service.normalizeMacAddress(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}
