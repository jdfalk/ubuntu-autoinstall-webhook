// internal/fileeditor/service.go
package fileeditor

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/jdfalk/ubuntu-autoinstall-webhook/internal/observability"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// FileEditor defines the operations for file management.
type FileEditor interface {
	ValidateIpxeFile(content []byte) error
	WriteIpxeFile(ctx context.Context, macAddress string, content []byte) error
	CreateCloudInitDirs(ctx context.Context, macAddress, hostname string) error
	ValidateCloudInitFiles(files map[string][]byte) error
	WriteCloudInitFile(ctx context.Context, macAddress string, fileType string, content []byte) error
	ListFiles(ctx context.Context, fileType string) ([]string, error)
	ReadFile(ctx context.Context, fileType string, filename string) ([]byte, error)
	DeleteFile(ctx context.Context, fileType string, filename string) error
	DeleteCloudInitDir(ctx context.Context, macOrHostname string) error
	CleanupRecycleBin(ctx context.Context) error
	Start(ctx context.Context) error
}

// Service implements the FileEditor interface.
type Service struct {
	fs            afero.Fs
	osFs          *afero.OsFs // Specific OS filesystem for symlink operations
	ipxeDir       string
	cloudInitDir  string
	leaderMutex   sync.Mutex
	isLeader      bool
	tracer        trace.Tracer
	leaderLockKey string
}

// NewService creates a new instance of the file editor service.
func NewService() FileEditor {
	// Create a tracer specific to this service.
	tracer := observability.GetTracer("fileeditor-service")
	fs := afero.NewOsFs()
	// Type assertion to get the concrete *OsFs type
	osFs, ok := fs.(*afero.OsFs)
	if !ok {
		panic("failed to cast to *afero.OsFs")
	}
	return &Service{
		fs:            fs,
		osFs:          osFs,
		ipxeDir:       viper.GetString("fileeditor.ipxe_dir"),
		cloudInitDir:  viper.GetString("fileeditor.cloudinit_dir"),
		isLeader:      false,
		tracer:        tracer,
		leaderLockKey: "fileeditor/leader",
	}
}

func (s *Service) ValidateIpxeFile(content []byte) error {
	// Start a new span for the validation operation.
	ctx, span := s.tracer.Start(context.Background(), "ValidateIpxeFile")
	defer span.End()

	fmt.Println("Validating ipxe file content")
	// TODO: implement actual validation logic.
	_ = ctx // Use context if needed for more advanced scenarios.
	return nil
}

func (s *Service) WriteIpxeFile(ctx context.Context, macAddress string, content []byte) error {
	ctx, span := s.tracer.Start(ctx, "WriteIpxeFile")
	defer span.End()

	span.SetAttributes(attribute.String("mac_address", macAddress))

	if !s.isLeader {
		acquired, err := s.AcquireLeadership(ctx)
		if err != nil {
			span.RecordError(err)
			return fmt.Errorf("failed to acquire leadership: %w", err)
		}

		if !acquired {
			return fmt.Errorf("not the leader, cannot write files")
		}
	}

	// Ensure the iPXE directory exists
	if err := s.ensureDirectory(ctx, s.ipxeDir); err != nil {
		span.RecordError(err)
		return err
	}

	// Normalize the MAC address and create the filename
	normalizedMac := s.normalizeMacAddress(macAddress)
	filename := fmt.Sprintf("mac-%s.ipxe", normalizedMac)
	filepath := filepath.Join(s.ipxeDir, filename)

	// Validate the content
	if err := s.validateIpxeContent(ctx, content); err != nil {
		span.RecordError(err)
		return fmt.Errorf("invalid iPXE content: %w", err)
	}

	// Write the file
	err := afero.WriteFile(s.fs, filepath, content, 0644)
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("failed to write iPXE file: %w", err)
	}

	span.SetAttributes(attribute.String("filepath", filepath))
	span.AddEvent("iPXE file written successfully")
	return nil
}

func (s *Service) CreateCloudInitDirs(ctx context.Context, macAddress, hostname string) error {
	ctx, span := s.tracer.Start(ctx, "CreateCloudInitDirs")
	defer span.End()

	span.SetAttributes(
		attribute.String("mac_address", macAddress),
		attribute.String("hostname", hostname),
	)

	if !s.isLeader {
		acquired, err := s.AcquireLeadership(ctx)
		if err != nil {
			span.RecordError(err)
			return fmt.Errorf("failed to acquire leadership: %w", err)
		}

		if !acquired {
			return fmt.Errorf("not the leader, cannot create directories")
		}
	}

	// Ensure the cloud-init directory exists
	if err := s.ensureDirectory(ctx, s.cloudInitDir); err != nil {
		span.RecordError(err)
		return err
	}

	// Normalize the MAC address
	normalizedMac := s.normalizeMacAddress(macAddress)

	// Create MAC address directories
	macDir := filepath.Join(s.cloudInitDir, normalizedMac)
	macInstallDir := filepath.Join(s.cloudInitDir, normalizedMac+"_install")

	if err := s.ensureDirectory(ctx, macDir); err != nil {
		span.RecordError(err)
		return err
	}

	if err := s.ensureDirectory(ctx, macInstallDir); err != nil {
		span.RecordError(err)
		return err
	}

	// Create hostname symlinks
	hostnameDir := filepath.Join(s.cloudInitDir, hostname)
	hostnameInstallDir := filepath.Join(s.cloudInitDir, hostname+"_install")

	// Remove existing symlinks if they exist
	_ = s.fs.Remove(hostnameDir)
	_ = s.fs.Remove(hostnameInstallDir)

	// Create the symlinks
	if symlinker, ok := s.fs.(afero.Symlinker); ok {
		if err := symlinker.SymlinkIfPossible(macDir, hostnameDir); err != nil {
			span.RecordError(err)
			return fmt.Errorf("failed to create hostname symlink: %w", err)
		}

		if err := symlinker.SymlinkIfPossible(macInstallDir, hostnameInstallDir); err != nil {
			span.RecordError(err)
			return fmt.Errorf("failed to create hostname_install symlink: %w", err)
		}
	} else {
		err := fmt.Errorf("filesystem does not support symlinks")
		span.RecordError(err)
		return err
	}

	span.AddEvent("Cloud-init directories and symlinks created successfully")
	return nil
}

func (s *Service) ValidateCloudInitFiles(files map[string][]byte) error {
	ctx, span := s.tracer.Start(context.Background(), "ValidateCloudInitFiles")
	defer span.End()

	fmt.Println("Validating cloud-init files")
	// TODO: implement validation using cloud-init libraries.
	_ = ctx
	return nil
}

// AcquireLeadership attempts to acquire leadership for file operations
// In a real implementation, this would use a distributed lock mechanism
func (s *Service) AcquireLeadership(ctx context.Context) (bool, error) {
	ctx, span := s.tracer.Start(ctx, "AcquireLeadership")
	defer span.End()

	s.leaderMutex.Lock()
	defer s.leaderMutex.Unlock()

	// TODO: Implement proper distributed locking using database or k8s
	// For now, just simulate being the leader
	s.isLeader = true

	span.SetAttributes(attribute.Bool("is_leader", s.isLeader))
	return s.isLeader, nil
}

// ReleaseLeadership releases the leadership lock
func (s *Service) ReleaseLeadership(ctx context.Context) error {
	ctx, span := s.tracer.Start(ctx, "ReleaseLeadership")
	defer span.End()

	s.leaderMutex.Lock()
	defer s.leaderMutex.Unlock()

	// Reset leader status
	s.isLeader = false

	span.SetAttributes(attribute.Bool("is_leader", s.isLeader))
	return nil
}

// ensureDirectory ensures a directory exists, creating it if necessary
func (s *Service) ensureDirectory(ctx context.Context, path string) error {
	_, span := s.tracer.Start(ctx, "ensureDirectory")
	defer span.End()

	span.SetAttributes(attribute.String("path", path))

	info, err := s.fs.Stat(path)
	if os.IsNotExist(err) {
		err = s.fs.MkdirAll(path, 0755)
		if err != nil {
			span.RecordError(err)
			return fmt.Errorf("failed to create directory %s: %w", path, err)
		}
		return nil
	}

	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("failed to check directory %s: %w", path, err)
	}

	if !info.IsDir() {
		err = fmt.Errorf("%s exists but is not a directory", path)
		span.RecordError(err)
		return err
	}

	return nil
}

// normalizeMacAddress converts a MAC address to the format used in filenames
func (s *Service) normalizeMacAddress(mac string) string {
	// Convert to lowercase and replace colons with hyphens
	return strings.ToLower(strings.ReplaceAll(mac, ":", "-"))
}

// validateIpxeContent validates the content of an iPXE file
func (s *Service) validateIpxeContent(ctx context.Context, content []byte) error {
	ctx, span := s.tracer.Start(ctx, "validateIpxeContent")
	defer span.End()

	// Check if content is empty
	if len(content) == 0 {
		err := fmt.Errorf("iPXE content cannot be empty")
		span.RecordError(err)
		return err
	}

	// Check if the content starts with #!ipxe
	if !strings.HasPrefix(string(content), "#!ipxe") {
		err := fmt.Errorf("iPXE content must start with #!ipxe")
		span.RecordError(err)
		return err
	}

	// Additional validation logic could be added here

	return nil
}

// WriteCloudInitFile writes a cloud-init file for a given MAC address and file type
func (s *Service) WriteCloudInitFile(ctx context.Context, macAddress string, fileType string, content []byte) error {
	ctx, span := s.tracer.Start(ctx, "WriteCloudInitFile")
	defer span.End()

	span.SetAttributes(
		attribute.String("mac_address", macAddress),
		attribute.String("file_type", fileType),
	)

	if !s.isLeader {
		acquired, err := s.AcquireLeadership(ctx)
		if err != nil {
			span.RecordError(err)
			return fmt.Errorf("failed to acquire leadership: %w", err)
		}

		if !acquired {
			return fmt.Errorf("not the leader, cannot write files")
		}
	}

	// Normalize the MAC address
	normalizedMac := s.normalizeMacAddress(macAddress)

	// Determine the directory based on file type
	var dirPath string
	if strings.HasSuffix(fileType, "_install") {
		baseType := strings.TrimSuffix(fileType, "_install")
		dirPath = filepath.Join(s.cloudInitDir, normalizedMac+"_install")
		fileType = baseType
	} else {
		dirPath = filepath.Join(s.cloudInitDir, normalizedMac)
	}

	// Validate that the directory exists
	if err := s.ensureDirectory(ctx, dirPath); err != nil {
		span.RecordError(err)
		return err
	}

	// Validate the content based on file type
	if err := s.validateCloudInitContent(ctx, fileType, content); err != nil {
		span.RecordError(err)
		return fmt.Errorf("invalid cloud-init content: %w", err)
	}

	// Write the file
	filePath := filepath.Join(dirPath, fileType)
	err := afero.WriteFile(s.fs, filePath, content, 0644)
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("failed to write cloud-init file: %w", err)
	}

	span.SetAttributes(attribute.String("filepath", filePath))
	span.AddEvent("Cloud-init file written successfully")
	return nil
}

// validateCloudInitContent validates the content of a cloud-init file based on its type
func (s *Service) validateCloudInitContent(ctx context.Context, fileType string, content []byte) error {
	ctx, span := s.tracer.Start(ctx, "validateCloudInitContent")
	defer span.End()

	span.SetAttributes(attribute.String("file_type", fileType))

	// Basic validation - check if content is empty
	if len(content) == 0 {
		err := fmt.Errorf("cloud-init content cannot be empty")
		span.RecordError(err)
		return err
	}

	switch fileType {
	case "meta-data":
		// Basic validation for meta-data
		// Could add more specific validation here
		return nil

	case "user-data":
		// Basic validation for user-data (should be valid YAML)
		// TODO: Add more comprehensive validation using cloud-init libraries
		if !strings.HasPrefix(string(content), "#cloud-config") {
			err := fmt.Errorf("user-data must start with #cloud-config")
			span.RecordError(err)
			return err
		}
		return nil

	case "network-config":
		// Network config validation
		return nil

	case "variables.sh":
		// Shell script validation
		if !strings.HasPrefix(string(content), "#!/bin/bash") &&
			!strings.HasPrefix(string(content), "#!/bin/sh") {
			err := fmt.Errorf("shell scripts must start with a valid shebang (#!/bin/bash or #!/bin/sh)")
			span.RecordError(err)
			return err
		}
		return nil

	default:
		// Unknown file type
		err := fmt.Errorf("unknown cloud-init file type: %s", fileType)
		span.RecordError(err)
		return err
	}
}

// ListFiles lists all files of a specific type in the given directory
func (s *Service) ListFiles(ctx context.Context, fileType string) ([]string, error) {
	ctx, span := s.tracer.Start(ctx, "ListFiles")
	defer span.End()

	span.SetAttributes(attribute.String("file_type", fileType))

	var dir string
	var pattern string

	switch fileType {
	case "ipxe":
		dir = s.ipxeDir
		pattern = "mac-*.ipxe"
	case "cloudinit":
		dir = s.cloudInitDir
		pattern = "*"
	default:
		err := fmt.Errorf("unknown file type: %s", fileType)
		span.RecordError(err)
		return nil, err
	}

	files, err := afero.Glob(s.fs, filepath.Join(dir, pattern))
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("failed to list files: %w", err)
	}

	// Extract just the filename part
	var result []string
	for _, file := range files {
		result = append(result, filepath.Base(file))
	}

	span.SetAttributes(attribute.Int("file_count", len(result)))
	return result, nil
}

// ReadFile reads a file from the appropriate directory based on type
func (s *Service) ReadFile(ctx context.Context, fileType string, filename string) ([]byte, error) {
	ctx, span := s.tracer.Start(ctx, "ReadFile")
	defer span.End()

	span.SetAttributes(
		attribute.String("file_type", fileType),
		attribute.String("filename", filename),
	)

	var filePath string

	switch fileType {
	case "ipxe":
		filePath = filepath.Join(s.ipxeDir, filename)
	case "cloudinit":
		parts := strings.Split(filename, "/")
		if len(parts) != 2 {
			err := fmt.Errorf("invalid cloudinit filename format, expected 'dir/file'")
			span.RecordError(err)
			return nil, err
		}
		filePath = filepath.Join(s.cloudInitDir, parts[0], parts[1])
	default:
		err := fmt.Errorf("unknown file type: %s", fileType)
		span.RecordError(err)
		return nil, err
	}

	content, err := afero.ReadFile(s.fs, filePath)
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	span.SetAttributes(attribute.Int("content_size", len(content)))
	return content, nil
}

// DeleteFile deletes a file from the appropriate directory
func (s *Service) DeleteFile(ctx context.Context, fileType string, filename string) error {
	ctx, span := s.tracer.Start(ctx, "DeleteFile")
	defer span.End()

	span.SetAttributes(
		attribute.String("file_type", fileType),
		attribute.String("filename", filename),
	)

	if !s.isLeader {
		acquired, err := s.AcquireLeadership(ctx)
		if err != nil {
			span.RecordError(err)
			return fmt.Errorf("failed to acquire leadership: %w", err)
		}

		if !acquired {
			return fmt.Errorf("not the leader, cannot delete files")
		}
	}

	var filePath string

	switch fileType {
	case "ipxe":
		filePath = filepath.Join(s.ipxeDir, filename)
	case "cloudinit":
		parts := strings.Split(filename, "/")
		if len(parts) != 2 {
			err := fmt.Errorf("invalid cloudinit filename format, expected 'dir/file'")
			span.RecordError(err)
			return err
		}
		filePath = filepath.Join(s.cloudInitDir, parts[0], parts[1])
	default:
		err := fmt.Errorf("unknown file type: %s", fileType)
		span.RecordError(err)
		return err
	}

	err := s.fs.Remove(filePath)
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("failed to delete file %s: %w", filePath, err)
	}

	span.AddEvent("File deleted successfully")
	return nil
}

// DeleteCloudInitDir moves a cloud-init directory and its associated files to the recycle bin
func (s *Service) DeleteCloudInitDir(ctx context.Context, macOrHostname string) error {
	fmt.Printf("\n===== DELETE CLOUD-INIT DIR START: %s =====\n", macOrHostname)
	ctx, span := s.tracer.Start(ctx, "DeleteCloudInitDir")
	defer span.End()

	span.SetAttributes(attribute.String("mac_or_hostname", macOrHostname))

	// Leadership check
	if !s.isLeader {
		acquired, err := s.AcquireLeadership(ctx)
		if err != nil {
			span.RecordError(err)
			return fmt.Errorf("failed to acquire leadership: %w", err)
		}
		if !acquired {
			return fmt.Errorf("not the leader, cannot delete directories")
		}
	}

	// Ensure recycle bin directory exists
	recycleBinPath := filepath.Join(s.cloudInitDir, "recycle_bin")
	fmt.Printf("DEBUG: Recycle bin path: %s\n", recycleBinPath)

	if err := s.ensureDirectory(ctx, recycleBinPath); err != nil {
		span.RecordError(err)
		return fmt.Errorf("failed to create recycle bin directory: %w", err)
	}
	fmt.Printf("DEBUG: Recycle bin exists\n")

	// We need to determine if this is a MAC address or hostname
	// Get Lstater and LinkReader for symlink operations
	lstater, ok := s.fs.(afero.Lstater)
	if !ok {
		return fmt.Errorf("filesystem doesn't support required operations")
	}

	linkReader, ok := s.fs.(afero.LinkReader)
	if !ok {
		return fmt.Errorf("filesystem doesn't support required operations")
	}

	var targetMacName string
	var foundResource bool
	var isSymlink bool

	// First, check if it's a valid MAC address format
	var normalizedMac string
	if isMacAddress(macOrHostname) {
		normalizedMac = s.normalizeMacAddress(macOrHostname)
		fmt.Printf("DEBUG: Validated as MAC address: %s -> %s\n", macOrHostname, normalizedMac)

		macDir := filepath.Join(s.cloudInitDir, normalizedMac)
		macDirExists, _ := afero.Exists(s.fs, macDir)

		if macDirExists {
			foundResource = true
			targetMacName = normalizedMac
			fmt.Printf("DEBUG: Found MAC directory: %s\n", macDir)
		}
	}

	// If not found as MAC or not a MAC format, try as a hostname (symlink)
	hostnameDir := filepath.Join(s.cloudInitDir, macOrHostname)
	if !foundResource {
		fmt.Printf("DEBUG: Checking if %s is a hostname (symlink)\n", macOrHostname)

		info, _, err := lstater.LstatIfPossible(hostnameDir)
		if err == nil && info.Mode()&os.ModeSymlink != 0 {
			fmt.Printf("DEBUG: Confirmed as symlink\n")
			isSymlink = true
			foundResource = true

			// Get the MAC directory it points to
			targetPath, err := linkReader.ReadlinkIfPossible(hostnameDir)
			if err != nil {
				span.RecordError(err)
				return fmt.Errorf("failed to read symlink: %w", err)
			}
			fmt.Printf("DEBUG: Symlink target: %s\n", targetPath)

			// Extract the MAC name from the target path
			targetMacName = filepath.Base(targetPath)
			fmt.Printf("DEBUG: Extracted MAC name: %s\n", targetMacName)
		}
	}

	if !foundResource {
		err := fmt.Errorf("no cloud-init directory or symlink found for %s", macOrHostname)
		span.RecordError(err)
		return err
	}

	// If it's a symlink (hostname), we need to check for other symlinks and remove hostname symlinks
	if isSymlink {
		// Remove the hostname symlinks first
		hostnameInstallDir := filepath.Join(s.cloudInitDir, macOrHostname+"_install")

		fmt.Printf("DEBUG: Removing hostname symlinks\n")
		if err := s.fs.Remove(hostnameDir); err != nil && !os.IsNotExist(err) {
			span.RecordError(err)
			return fmt.Errorf("failed to remove hostname symlink: %w", err)
		}

		if err := s.fs.Remove(hostnameInstallDir); err != nil && !os.IsNotExist(err) {
			// Just log this error as the install symlink might not exist
			span.RecordError(err)
		}

		// Check if any other symlinks point to this MAC directory
		fmt.Printf("DEBUG: Checking for other symlinks to %s\n", targetMacName)
		dirs, err := afero.ReadDir(s.fs, s.cloudInitDir)
		if err != nil {
			span.RecordError(err)
			return fmt.Errorf("failed to read cloud-init directory: %w", err)
		}

		// Look for other symlinks pointing to the same MAC directory
		otherSymlinksExist := false
		for _, dir := range dirs {
			// Skip non-symlinks and the symlinks we just removed
			if dir.Name() == macOrHostname || dir.Name() == macOrHostname+"_install" {
				continue
			}

			entryPath := filepath.Join(s.cloudInitDir, dir.Name())
			entryInfo, _, err := lstater.LstatIfPossible(entryPath)
			if err != nil || entryInfo.Mode()&os.ModeSymlink == 0 {
				continue
			}

			// Read symlink target
			linkTarget, err := linkReader.ReadlinkIfPossible(entryPath)
			if err != nil {
				continue
			}

			// Check if this symlink points to our MAC directory
			if filepath.Base(linkTarget) == targetMacName {
				fmt.Printf("DEBUG: Found other symlink %s pointing to %s\n", dir.Name(), targetMacName)
				otherSymlinksExist = true
				break
			}
		}

		// If other symlinks exist, don't delete the MAC directories
		if otherSymlinksExist {
			fmt.Printf("DEBUG: Other symlinks exist, not deleting MAC directories\n")
			return nil
		}

		fmt.Printf("DEBUG: No other symlinks found, will delete MAC directories\n")
	}

	// Move MAC directories to recycle bin with timestamp
	timestamp := time.Now().Format("20060102_150405")
	fmt.Printf("DEBUG: Moving directories to recycle bin with timestamp %s\n", timestamp)

	// Move the main directory
	macDir := filepath.Join(s.cloudInitDir, targetMacName)
	macDirExists, _ := afero.Exists(s.fs, macDir)
	fmt.Printf("DEBUG: MAC dir %s exists? %v\n", macDir, macDirExists)

	if macDirExists {
		recycleDestPath := filepath.Join(recycleBinPath,
			fmt.Sprintf("%s_delete_me_%s", targetMacName, timestamp))
		fmt.Printf("DEBUG: Moving %s -> %s\n", macDir, recycleDestPath)

		if err := s.fs.Rename(macDir, recycleDestPath); err != nil {
			span.RecordError(err)
			return fmt.Errorf("failed to move directory to recycle bin: %w", err)
		}
		fmt.Printf("DEBUG: Successfully moved directory to recycle bin\n")
	}

	// Move the install directory
	macInstallDir := filepath.Join(s.cloudInitDir, targetMacName+"_install")
	installDirExists, _ := afero.Exists(s.fs, macInstallDir)
	fmt.Printf("DEBUG: MAC install dir %s exists? %v\n", macInstallDir, installDirExists)

	if installDirExists {
		recycleDestPath := filepath.Join(recycleBinPath,
			fmt.Sprintf("%s_install_delete_me_%s", targetMacName, timestamp))
		fmt.Printf("DEBUG: Moving %s -> %s\n", macInstallDir, recycleDestPath)

		if err := s.fs.Rename(macInstallDir, recycleDestPath); err != nil {
			span.RecordError(err)
			return fmt.Errorf("failed to move install directory to recycle bin: %w", err)
		}
		fmt.Printf("DEBUG: Successfully moved install directory to recycle bin\n")
	}

	fmt.Printf("===== DELETE CLOUD-INIT DIR COMPLETE =====\n\n")
	span.AddEvent("Directories moved to recycle bin")
	return nil
}

// Helper function to determine if a string is a MAC address
func isMacAddress(s string) bool {
	// Simple regex to match MAC address formats like 00:11:22:33:44:55 or 00-11-22-33-44-55
	matched, _ := regexp.MatchString(`^([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})$`, s)
	return matched
}

// CleanupRecycleBin removes all files from the recycle bin directory
func (s *Service) CleanupRecycleBin(ctx context.Context) error {
	ctx, span := s.tracer.Start(ctx, "CleanupRecycleBin")
	defer span.End()

	recycleBinPath := filepath.Join(s.cloudInitDir, "recycle_bin")

	// Check if recycle bin exists
	exists, err := afero.Exists(s.fs, recycleBinPath)
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("failed to check if recycle bin exists: %w", err)
	}

	// If it doesn't exist, create it
	if !exists {
		if err := s.fs.MkdirAll(recycleBinPath, 0755); err != nil {
			span.RecordError(err)
			return fmt.Errorf("failed to create recycle bin: %w", err)
		}
		return nil
	}

	// Get all entries in the recycle bin
	entries, err := afero.ReadDir(s.fs, recycleBinPath)
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("failed to read recycle bin: %w", err)
	}

	// Delete each entry
	for _, entry := range entries {
		entryPath := filepath.Join(recycleBinPath, entry.Name())
		if err := s.fs.RemoveAll(entryPath); err != nil {
			span.RecordError(err)
			return fmt.Errorf("failed to remove %s from recycle bin: %w", entry.Name(), err)
		}
	}

	span.AddEvent("Recycle bin cleaned up")
	return nil
}

// Add this method to the Service struct in service.go
func (s *Service) Start(ctx context.Context) error {
	ctx, span := s.tracer.Start(ctx, "Start")
	defer span.End()

	span.AddEvent("Starting file editor service")

	// Ensure the recycle bin exists
	recycleBinPath := filepath.Join(s.cloudInitDir, "recycle_bin")
	if err := s.ensureDirectory(ctx, recycleBinPath); err != nil {
		span.RecordError(err)
		return fmt.Errorf("failed to create recycle bin directory: %w", err)
	}

	// Start a goroutine for periodic cleanup
	go func() {
		ticker := time.NewTicker(24 * time.Hour) // Configure cleanup period as needed
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if s.isLeader {
					cleanCtx := context.Background()
					if err := s.CleanupRecycleBin(cleanCtx); err != nil {
						// Log error but continue
						fmt.Printf("Failed to clean recycle bin: %v\n", err)
					}
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}
