// internal/fileeditor/service.go
package fileeditor

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

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

// DeleteCloudInitDir deletes a cloud-init directory and its associated symlinks
func (s *Service) DeleteCloudInitDir(ctx context.Context, macOrHostname string) error {
	ctx, span := s.tracer.Start(ctx, "DeleteCloudInitDir")
	defer span.End()

	span.SetAttributes(attribute.String("mac_or_hostname", macOrHostname))

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

	// Check if this is a MAC address or hostname
	normalizedInput := s.normalizeMacAddress(macOrHostname)
	macDir := filepath.Join(s.cloudInitDir, normalizedInput)
	macInstallDir := filepath.Join(s.cloudInitDir, normalizedInput+"_install")

	// First try as a MAC address
	if _, err := s.fs.Stat(macDir); err == nil {
		// It's a MAC address directory
		// Find any symlinks pointing to this MAC
		dirs, err := afero.ReadDir(s.fs, s.cloudInitDir)
		if err != nil {
			span.RecordError(err)
			return fmt.Errorf("failed to read cloud-init directory: %w", err)
		}

		// Check for symlinks to this MAC directory
		for _, dir := range dirs {
			if dir.Mode()&os.ModeSymlink != 0 {
				linkPath := filepath.Join(s.cloudInitDir, dir.Name())

				// Use the Symlinker interface to read the link
				if linkReader, ok := s.fs.(afero.LinkReader); ok {
					target, err := linkReader.ReadlinkIfPossible(linkPath)
					if err != nil {
						continue
					}

					// If this symlink points to our MAC directory, remove it
					if target == macDir || target == macInstallDir {
						if err := s.fs.Remove(linkPath); err != nil {
							span.RecordError(err)
							return fmt.Errorf("failed to remove symlink %s: %w", linkPath, err)
						}
					}
				}
			}
		}

		// Now remove the MAC directories
		if err := s.fs.RemoveAll(macDir); err != nil {
			span.RecordError(err)
			return fmt.Errorf("failed to remove directory %s: %w", macDir, err)
		}

		if err := s.fs.RemoveAll(macInstallDir); err != nil {
			span.RecordError(err)
			return fmt.Errorf("failed to remove directory %s: %w", macInstallDir, err)
		}

		span.AddEvent("MAC address directories removed successfully")
		return nil
	}

	// Try as a hostname (symlink)
	hostnameLink := filepath.Join(s.cloudInitDir, macOrHostname)
	hostnameInstallLink := filepath.Join(s.cloudInitDir, macOrHostname+"_install")

	// Check if hostname symlink exists using Lstater interface
	var linkInfo os.FileInfo
	var err error
	if lstater, ok := s.fs.(afero.Lstater); ok {
		linkInfo, _, err = lstater.LstatIfPossible(hostnameLink)
	} else {
		// Fall back to regular Stat if Lstat is not available
		linkInfo, err = s.fs.Stat(hostnameLink)
	}
	if err == nil && linkInfo.Mode()&os.ModeSymlink != 0 {
		// Read the target of the symlink
		if linkReader, ok := s.fs.(afero.LinkReader); ok {
			target, err := linkReader.ReadlinkIfPossible(hostnameLink)
			if err != nil {
				span.RecordError(err)
				return fmt.Errorf("failed to read symlink %s: %w", hostnameLink, err)
			}

			// Remove the hostname symlinks
			if err := s.fs.Remove(hostnameLink); err != nil && !os.IsNotExist(err) {
				span.RecordError(err)
				return fmt.Errorf("failed to remove symlink %s: %w", hostnameLink, err)
			}

			if err := s.fs.Remove(hostnameInstallLink); err != nil && !os.IsNotExist(err) {
				span.RecordError(err)
				// Just log this error as the install symlink may not exist
			}

			// Check if there are other symlinks to the same MAC
			dirs, err := afero.ReadDir(s.fs, s.cloudInitDir)
			if err != nil {
				span.RecordError(err)
				return fmt.Errorf("failed to read cloud-init directory: %w", err)
			}

			// Count symlinks pointing to the same MAC directory
			var linkCount int
			for _, dir := range dirs {
				if dir.Mode()&os.ModeSymlink != 0 {
					// Skip the hostname we're currently deleting
					if dir.Name() == macOrHostname ||
						dir.Name() == macOrHostname+"_install" {
						continue
					}

					linkPath := filepath.Join(s.cloudInitDir, dir.Name())
					linkTarget, err := linkReader.ReadlinkIfPossible(linkPath)
					if err != nil {
						continue
					}

					if linkTarget == target {
						linkCount++
					}
				}
			}

			// If no other symlinks point to this MAC, remove the MAC directories too
			if linkCount == 0 {
				// Extract the last component of the target path which is the MAC directory
				macDirName := filepath.Base(target)
				macDir := filepath.Join(s.cloudInitDir, macDirName)
				macInstallDir := filepath.Join(s.cloudInitDir, macDirName+"_install")

				if err := s.fs.RemoveAll(macDir); err != nil && !os.IsNotExist(err) {
					span.RecordError(err)
					return fmt.Errorf("failed to remove directory %s: %w", macDir, err)
				}

				if err := s.fs.RemoveAll(macInstallDir); err != nil && !os.IsNotExist(err) {
					span.RecordError(err)
					return fmt.Errorf("failed to remove directory %s: %w", macInstallDir, err)
				}
			}

			span.AddEvent("Hostname symlinks removed successfully")
			return nil
		}
	}

	err = fmt.Errorf("no cloud-init directory or symlink found for %s", macOrHostname)
	span.RecordError(err)
	return err
}
