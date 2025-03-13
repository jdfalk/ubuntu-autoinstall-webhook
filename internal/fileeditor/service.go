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
	normalizedMac := s.normalizeMacAddress(macOrHostname)
	macDir := filepath.Join(s.cloudInitDir, normalizedMac)
	macInstallDir := filepath.Join(s.cloudInitDir, normalizedMac+"_install")

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
			// Skip if not a symlink
			if lstater, ok := s.fs.(afero.Lstater); ok {
				info, _, err := lstater.LstatIfPossible(filepath.Join(s.cloudInitDir, dir.Name()))
				if err != nil || info.Mode()&os.ModeSymlink == 0 {
					continue
				}

				// It's a symlink, read its target
				if linkReader, ok := s.fs.(afero.LinkReader); ok {
					target, err := linkReader.ReadlinkIfPossible(filepath.Join(s.cloudInitDir, dir.Name()))
					if err != nil {
						continue
					}

					// If this symlink points to our MAC directory, remove it
					// Handle that target might contain full path in case of BasePathFs
					if strings.HasSuffix(target, macDir) || strings.HasSuffix(target, macInstallDir) {
						if err := s.fs.Remove(filepath.Join(s.cloudInitDir, dir.Name())); err != nil {
							span.RecordError(err)
							return fmt.Errorf("failed to remove symlink: %w", err)
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

		if err := s.fs.RemoveAll(macInstallDir); err != nil && !os.IsNotExist(err) {
			span.RecordError(err)
			return fmt.Errorf("failed to remove directory %s: %w", macInstallDir, err)
		}

		span.AddEvent("MAC address directories removed successfully")
		return nil
	}

	// Try as a hostname (symlink)
	hostnameLink := filepath.Join(s.cloudInitDir, macOrHostname)
	hostnameInstallLink := filepath.Join(s.cloudInitDir, macOrHostname+"_install")

	// Check if hostname symlink exists
	if lstater, ok := s.fs.(afero.Lstater); ok {
		linkInfo, _, err := lstater.LstatIfPossible(hostnameLink)
		if err == nil && linkInfo.Mode()&os.ModeSymlink != 0 {
			// Found a symlink - get the MAC directory it points to
			if linkReader, ok := s.fs.(afero.LinkReader); ok {
				targetPath, err := linkReader.ReadlinkIfPossible(hostnameLink)
				if err != nil {
					span.RecordError(err)
					return fmt.Errorf("failed to read symlink %s: %w", hostnameLink, err)
				}

				// Extract the MAC directory name from the target path
				macDirName := filepath.Base(targetPath)
				macDir = filepath.Join(s.cloudInitDir, macDirName)
				macInstallDir = filepath.Join(s.cloudInitDir, macDirName+"_install")

				// Remove hostname symlinks
				if err := s.fs.Remove(hostnameLink); err != nil && !os.IsNotExist(err) {
					span.RecordError(err)
					return fmt.Errorf("failed to remove symlink %s: %w", hostnameLink, err)
				}

				if err := s.fs.Remove(hostnameInstallLink); err != nil && !os.IsNotExist(err) {
					span.RecordError(err)
					// Just log this error as the install symlink may not exist
				}

				// Check if any other symlinks point to this MAC directory
				otherSymlinksExist := false

				dirs, err := afero.ReadDir(s.fs, s.cloudInitDir)
				if err != nil {
					span.RecordError(err)
					return fmt.Errorf("failed to read cloud-init directory: %w", err)
				}

				for _, entry := range dirs {
					// Skip non-symlinks or the ones we're deleting
					if entry.Name() == macOrHostname || entry.Name() == macOrHostname+"_install" {
						continue
					}

					entryInfo, _, err := lstater.LstatIfPossible(filepath.Join(s.cloudInitDir, entry.Name()))
					if err != nil || entryInfo.Mode()&os.ModeSymlink == 0 {
						continue
					}

					// It's a symlink, check if it points to our MAC directory
					entryTarget, err := linkReader.ReadlinkIfPossible(filepath.Join(s.cloudInitDir, entry.Name()))
					if err != nil {
						continue
					}

					// Use base name comparison to handle path prefix differences
					if filepath.Base(entryTarget) == macDirName {
						otherSymlinksExist = true
						break
					}
				}

				// If no other symlinks point to the MAC directory, delete it
				if !otherSymlinksExist {
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
	}

	err := fmt.Errorf("no cloud-init directory or symlink found for %s", macOrHostname)
	span.RecordError(err)
	return err
}
