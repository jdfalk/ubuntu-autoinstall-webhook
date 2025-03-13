// internal/certissuer/certissuer.go
package certissuer

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	// Default key size for RSA keys
	DefaultKeySize = 2048

	// Certificate validity periods
	CAValidityPeriod   = 10 * 365 * 24 * time.Hour // 10 years
	CertValidityPeriod = 365 * 24 * time.Hour      // 1 year
	RenewBeforeExpiry  = 30 * 24 * time.Hour       // 30 days
)

// CertificateAuthority represents a certificate authority for issuing and managing certificates
type CertificateAuthority struct {
	// CA certificate and private key
	caCert    *x509.Certificate
	caPrivKey *rsa.PrivateKey

	// PEM-encoded versions for storage/return
	caCertPEM []byte
	caKeyPEM  []byte

	// Certificate storage
	certStore map[string]*Certificate
	storePath string

	// Concurrent access
	mu sync.RWMutex
}

// Certificate holds information about an issued certificate
type Certificate struct {
	Serial    *big.Int
	Cert      *x509.Certificate
	CertPEM   []byte
	ExpiresAt time.Time
	IssuedTo  string
	IsRevoked bool
}

// NewCertificateAuthority creates a new certificate authority
func NewCertificateAuthority(storePath string) (*CertificateAuthority, error) {
	ca := &CertificateAuthority{
		certStore: make(map[string]*Certificate),
		storePath: storePath,
	}

	// Create storage directory if it doesn't exist
	if storePath != "" {
		if err := os.MkdirAll(storePath, 0755); err != nil {
			return nil, fmt.Errorf("failed to create certificate storage directory: %w", err)
		}
	}

	// Check if CA already exists in storage
	if err := ca.loadCA(); err != nil {
		// If not, create a new CA
		if err := ca.createCA(); err != nil {
			return nil, fmt.Errorf("failed to create certificate authority: %w", err)
		}
	}

	return ca, nil
}

// createCA generates a new self-signed CA certificate
func (ca *CertificateAuthority) createCA() error {
	// Generate CA private key
	privateKey, err := rsa.GenerateKey(rand.Reader, DefaultKeySize)
	if err != nil {
		return fmt.Errorf("failed to generate CA private key: %w", err)
	}

	// Prepare CA certificate template
	serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return fmt.Errorf("failed to generate serial number: %w", err)
	}

	now := time.Now()
	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization:       []string{"Ubuntu Autoinstall Webhook CA"},
			OrganizationalUnit: []string{"Certificate Authority"},
			CommonName:         "Ubuntu Autoinstall Webhook Root CA",
		},
		NotBefore:             now,
		NotAfter:              now.Add(CAValidityPeriod),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
		IsCA:                  true,
		MaxPathLen:            1,
	}

	// Self-sign the CA certificate
	certBytes, err := x509.CreateCertificate(
		rand.Reader,
		&template,
		&template,
		&privateKey.PublicKey,
		privateKey,
	)
	if err != nil {
		return fmt.Errorf("failed to create CA certificate: %w", err)
	}

	// Parse the certificate bytes
	cert, err := x509.ParseCertificate(certBytes)
	if err != nil {
		return fmt.Errorf("failed to parse CA certificate: %w", err)
	}

	// Store the CA certificate and private key
	ca.caCert = cert
	ca.caPrivKey = privateKey

	// Encode to PEM
	caCertPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})

	caKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	ca.caCertPEM = caCertPEM
	ca.caKeyPEM = caKeyPEM

	// Save to storage if path is provided
	if ca.storePath != "" {
		if err := ca.saveCA(); err != nil {
			return fmt.Errorf("failed to save CA certificate: %w", err)
		}
	}

	return nil
}

// loadCA loads the CA certificate and private key from storage
func (ca *CertificateAuthority) loadCA() error {
	if ca.storePath == "" {
		return fmt.Errorf("no storage path provided")
	}

	// Read CA certificate
	certPath := filepath.Join(ca.storePath, "ca.crt")
	certPEM, err := os.ReadFile(certPath)
	if err != nil {
		return err
	}

	// Read CA private key
	keyPath := filepath.Join(ca.storePath, "ca.key")
	keyPEM, err := os.ReadFile(keyPath)
	if err != nil {
		return err
	}

	// Parse certificate PEM
	certBlock, _ := pem.Decode(certPEM)
	if certBlock == nil {
		return fmt.Errorf("failed to decode CA certificate PEM")
	}

	cert, err := x509.ParseCertificate(certBlock.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse CA certificate: %w", err)
	}

	// Parse private key PEM
	keyBlock, _ := pem.Decode(keyPEM)
	if keyBlock == nil {
		return fmt.Errorf("failed to decode CA private key PEM")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(keyBlock.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse CA private key: %w", err)
	}

	// Store in memory
	ca.caCert = cert
	ca.caPrivKey = privateKey
	ca.caCertPEM = certPEM
	ca.caKeyPEM = keyPEM

	return nil
}

// saveCA saves the CA certificate and private key to storage
func (ca *CertificateAuthority) saveCA() error {
	if ca.storePath == "" {
		return fmt.Errorf("no storage path provided")
	}

	// Save CA certificate
	certPath := filepath.Join(ca.storePath, "ca.crt")
	if err := os.WriteFile(certPath, ca.caCertPEM, 0644); err != nil {
		return fmt.Errorf("failed to save CA certificate: %w", err)
	}

	// Save CA private key (restricted permissions)
	keyPath := filepath.Join(ca.storePath, "ca.key")
	if err := os.WriteFile(keyPath, ca.caKeyPEM, 0600); err != nil {
		return fmt.Errorf("failed to save CA private key: %w", err)
	}

	return nil
}

// GetCACertificate returns the PEM-encoded CA certificate
func (ca *CertificateAuthority) GetCACertificate() []byte {
	ca.mu.RLock()
	defer ca.mu.RUnlock()

	return ca.caCertPEM
}

// IssueCertificateFromCSR issues a certificate from a CSR
func (ca *CertificateAuthority) IssueCertificateFromCSR(csrPEM []byte, clientInfo map[string]string) ([]byte, error) {
	ca.mu.Lock()
	defer ca.mu.Unlock()

	// Decode CSR
	block, _ := pem.Decode(csrPEM)
	if block == nil || block.Type != "CERTIFICATE REQUEST" {
		return nil, fmt.Errorf("failed to decode CSR, invalid PEM format")
	}

	// Parse CSR
	csr, err := x509.ParseCertificateRequest(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse CSR: %w", err)
	}

	// Verify CSR signature
	if err := csr.CheckSignature(); err != nil {
		return nil, fmt.Errorf("invalid CSR signature: %w", err)
	}

	// Generate serial number
	serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return nil, fmt.Errorf("failed to generate serial number: %w", err)
	}

	// Extract client information
	commonName := clientInfo["common_name"]
	if commonName == "" {
		commonName = csr.Subject.CommonName
	}

	now := time.Now()
	// Create certificate template
	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName: commonName,
			Organization: []string{func() string {
				if org, ok := clientInfo["organization"]; ok && org != "" {
					return org
				}
				return "Ubuntu Autoinstall"
			}()},
		},
		NotBefore:             now,
		NotAfter:              now.Add(CertValidityPeriod),
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
		IsCA:                  false,
	}

	// Add Subject Alternative Names if provided
	if sans, ok := clientInfo["sans"]; ok && sans != "" {
		// Parse comma-separated list of SANs
		// Note: In a real implementation, you would parse this properly
		template.DNSNames = []string{sans}
	}

	// Sign the certificate
	certBytes, err := x509.CreateCertificate(
		rand.Reader,
		&template,
		ca.caCert,
		csr.PublicKey,
		ca.caPrivKey,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create certificate: %w", err)
	}

	// Encode to PEM
	certPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})

	// Parse the certificate for storage
	cert, err := x509.ParseCertificate(certBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse issued certificate: %w", err)
	}

	// Store the certificate
	certInfo := &Certificate{
		Serial:    serialNumber,
		Cert:      cert,
		CertPEM:   certPEM,
		ExpiresAt: template.NotAfter,
		IssuedTo:  commonName,
	}

	// Store the certificate using the serial number as a key
	ca.certStore[serialNumber.String()] = certInfo

	// Save to storage if configured
	if ca.storePath != "" {
		certPath := filepath.Join(ca.storePath, fmt.Sprintf("%s.crt", serialNumber.String()))
		if err := os.WriteFile(certPath, certPEM, 0644); err != nil {
			// Log but don't fail
			fmt.Printf("Warning: Failed to save certificate to storage: %v\n", err)
		}
	}

	return certPEM, nil
}

// RenewCertificate renews an existing certificate
func (ca *CertificateAuthority) RenewCertificate(certPEM []byte) ([]byte, error) {
	ca.mu.Lock()
	defer ca.mu.Unlock()

	// Decode certificate
	block, _ := pem.Decode(certPEM)
	if block == nil || block.Type != "CERTIFICATE" {
		return nil, fmt.Errorf("failed to decode certificate, invalid PEM format")
	}

	// Parse certificate
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse certificate: %w", err)
	}

	// Verify certificate was issued by this CA
	if err := cert.CheckSignatureFrom(ca.caCert); err != nil {
		return nil, fmt.Errorf("certificate not issued by this CA")
	}

	// Check if certificate is revoked
	if storedCert, ok := ca.certStore[cert.SerialNumber.String()]; ok && storedCert.IsRevoked {
		return nil, fmt.Errorf("certificate has been revoked")
	}

	// Generate new serial number
	serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return nil, fmt.Errorf("failed to generate serial number: %w", err)
	}

	// Create new certificate with same information but new validity period
	now := time.Now()
	template := x509.Certificate{
		SerialNumber:          serialNumber,
		Subject:               cert.Subject,
		NotBefore:             now,
		NotAfter:              now.Add(CertValidityPeriod),
		KeyUsage:              cert.KeyUsage,
		ExtKeyUsage:           cert.ExtKeyUsage,
		BasicConstraintsValid: true,
		IsCA:                  false,
		DNSNames:              cert.DNSNames,
		IPAddresses:           cert.IPAddresses,
	}

	// Sign the certificate with CA
	newCertBytes, err := x509.CreateCertificate(
		rand.Reader,
		&template,
		ca.caCert,
		cert.PublicKey,
		ca.caPrivKey,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create renewed certificate: %w", err)
	}

	// Encode to PEM
	newCertPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: newCertBytes,
	})

	// Parse for storage
	newCert, err := x509.ParseCertificate(newCertBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse renewed certificate: %w", err)
	}

	// Store the renewed certificate
	certInfo := &Certificate{
		Serial:    serialNumber,
		Cert:      newCert,
		CertPEM:   newCertPEM,
		ExpiresAt: template.NotAfter,
		IssuedTo:  cert.Subject.CommonName,
	}

	ca.certStore[serialNumber.String()] = certInfo

	// Save to storage if configured
	if ca.storePath != "" {
		certPath := filepath.Join(ca.storePath, fmt.Sprintf("%s.crt", serialNumber.String()))
		if err := os.WriteFile(certPath, newCertPEM, 0644); err != nil {
			// Log but don't fail
			fmt.Printf("Warning: Failed to save renewed certificate to storage: %v\n", err)
		}
	}

	return newCertPEM, nil
}

// RevokeCertificate marks a certificate as revoked
func (ca *CertificateAuthority) RevokeCertificate(serialNumber string) error {
	ca.mu.Lock()
	defer ca.mu.Unlock()

	cert, ok := ca.certStore[serialNumber]
	if !ok {
		return fmt.Errorf("certificate with serial number %s not found", serialNumber)
	}

	cert.IsRevoked = true
	return nil
}
