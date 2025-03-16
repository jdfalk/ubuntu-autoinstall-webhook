// internal/certissuer/service.go
package certissuer

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
)

// CertIssuer defines the interface for certificate management.
type CertIssuer interface {
	// IssueCertificate issues a certificate based on the provided CSR and client info.
	IssueCertificate(ctx context.Context, csr []byte, clientInfo map[string]string) ([]byte, error)
	// RenewCertificate renews the given certificate.
	RenewCertificate(ctx context.Context, cert []byte) ([]byte, error)
	// GetRootCA returns the root CA certificate.
	GetRootCA(ctx context.Context) ([]byte, error)
}

// Service implements the CertIssuer interface.
type Service struct {
	ca          *CertificateAuthority
	certStorage string
}

// NewService creates a new instance of the CertIssuer service.
func NewService(certStorage string) CertIssuer {
	// Default storage location if none provided
	if certStorage == "" {
		homeDir, err := os.UserHomeDir()
		if err == nil {
			certStorage = filepath.Join(homeDir, ".autoinstall-webhook", "certificates")
		} else {
			certStorage = "/var/lib/autoinstall-webhook/certificates"
		}
	}

	ca, err := NewCertificateAuthority(certStorage)
	if err != nil {
		fmt.Printf("Error initializing certificate authority: %v\n", err)
		// Fallback to in-memory CA
		ca, _ = NewCertificateAuthority("")
	}

	return &Service{
		ca:          ca,
		certStorage: certStorage,
	}
}

// IssueCertificate issues a certificate based on the provided CSR and client info.
func (s *Service) IssueCertificate(ctx context.Context, csr []byte, clientInfo map[string]string) ([]byte, error) {
	fmt.Println("Issuing certificate for client:", clientInfo)
	return s.ca.IssueCertificateFromCSR(csr, clientInfo)
}

// RenewCertificate renews the given certificate.
func (s *Service) RenewCertificate(ctx context.Context, cert []byte) ([]byte, error) {
	fmt.Println("Renewing certificate")
	return s.ca.RenewCertificate(cert)
}

// GetRootCA returns the root CA certificate.
func (s *Service) GetRootCA(ctx context.Context) ([]byte, error) {
	fmt.Println("Retrieving root CA certificate")
	return s.ca.GetCACertificate(), nil
}
