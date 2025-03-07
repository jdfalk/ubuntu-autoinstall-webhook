// internal/fileeditor/service.go
package fileeditor

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

// FileEditor defines the operations for file management.
type FileEditor interface {
	ValidateIpxeFile(content []byte) error
	WriteIpxeFile(macAddress string, content []byte) error
	CreateCloudInitDirs(macAddress, hostname string) error
	ValidateCloudInitFiles(files map[string][]byte) error
}

// Service implements the FileEditor interface.
type Service struct {
	tracer trace.Tracer
}

// NewService creates a new instance of the file editor service.
func NewService() FileEditor {
	// Create a tracer specific to this service.
	tracer := otel.Tracer("fileeditor")
	return &Service{
		tracer: tracer,
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

func (s *Service) WriteIpxeFile(macAddress string, content []byte) error {
	ctx, span := s.tracer.Start(context.Background(), "WriteIpxeFile")
	defer span.End()

	if err := s.ValidateIpxeFile(content); err != nil {
		return err
	}
	fmt.Printf("Writing ipxe file for MAC %s\n", macAddress)
	// TODO: write the file to /var/www/html/ipxe/boot/mac-{MAC_ADDRESS}.ipxe.
	_ = ctx
	return nil
}

func (s *Service) CreateCloudInitDirs(macAddress, hostname string) error {
	ctx, span := s.tracer.Start(context.Background(), "CreateCloudInitDirs")
	defer span.End()

	fmt.Printf("Creating cloud-init directories for MAC %s and hostname %s\n", macAddress, hostname)
	// TODO: create directories and symlinks, and perform duplicate hostname checks.
	_ = ctx
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
