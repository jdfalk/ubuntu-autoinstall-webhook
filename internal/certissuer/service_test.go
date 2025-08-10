// internal/certissuer/service_test.go
package certissuer_test

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"os"
	"testing"

	"github.com/jdfalk/ubuntu-autoinstall-webhook/internal/certissuer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Helper function to create a test service with a temporary directory
func setupTestService(t *testing.T) (certissuer.CertIssuer, string) {
	// Create a temporary directory for certificate storage
	tempDir, err := os.MkdirTemp("", "certissuer-test-*")
	require.NoError(t, err, "Failed to create temp directory")

	// Create the service
	svc := certissuer.NewService(tempDir)
	require.NotNil(t, svc, "Service should not be nil")

	return svc, tempDir
}

// Helper function to create a test CSR
func createTestCSR(t *testing.T, commonName string) ([]byte, *rsa.PrivateKey) {
	// Generate a private key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err, "Failed to generate private key")

	// Create a CSR template
	template := x509.CertificateRequest{
		Subject: pkix.Name{
			CommonName:   commonName,
			Organization: []string{"Test Organization"},
		},
		SignatureAlgorithm: x509.SHA256WithRSA,
	}

	// Create the CSR
	csrDER, err := x509.CreateCertificateRequest(rand.Reader, &template, privateKey)
	require.NoError(t, err, "Failed to create CSR")

	// Encode CSR to PEM
	csrPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE REQUEST",
		Bytes: csrDER,
	})

	return csrPEM, privateKey
}

func TestIssueCertificate(t *testing.T) {
	// Setup
	svc, tempDir := setupTestService(t)
	defer os.RemoveAll(tempDir)

	// Create a test CSR
	commonName := "test.example.com"
	csrPEM, _ := createTestCSR(t, commonName)

	// Test issuing a certificate
	clientInfo := map[string]string{
		"common_name":  commonName,
		"organization": "Test Organization",
	}

	// Issue certificate
	ctx := context.Background()
	certPEM, err := svc.IssueCertificate(ctx, csrPEM, clientInfo)
	require.NoError(t, err, "IssueCertificate should not return an error")
	require.NotEmpty(t, certPEM, "Certificate should not be empty")

	// Validate certificate
	block, _ := pem.Decode(certPEM)
	require.NotNil(t, block, "Certificate PEM should decode successfully")
	require.Equal(t, "CERTIFICATE", block.Type, "PEM block type should be CERTIFICATE")

	cert, err := x509.ParseCertificate(block.Bytes)
	require.NoError(t, err, "Should parse certificate successfully")
	require.Equal(t, commonName, cert.Subject.CommonName, "Certificate should have the requested common name")

	// Get Root CA
	caCertPEM, err := svc.GetRootCA(ctx)
	require.NoError(t, err, "GetRootCA should not return an error")
	require.NotEmpty(t, caCertPEM, "CA Certificate should not be empty")

	// Verify certificate was issued by the CA
	caBlock, _ := pem.Decode(caCertPEM)
	require.NotNil(t, caBlock, "CA Certificate PEM should decode successfully")

	caCert, err := x509.ParseCertificate(caBlock.Bytes)
	require.NoError(t, err, "Should parse CA certificate successfully")

	// Verify certificate is signed by CA
	err = cert.CheckSignatureFrom(caCert)
	assert.NoError(t, err, "Certificate should be signed by the CA")
}

func TestRenewCertificate(t *testing.T) {
	// Setup
	svc, tempDir := setupTestService(t)
	defer os.RemoveAll(tempDir)

	// First issue a certificate that we can renew
	commonName := "renew.example.com"
	csrPEM, _ := createTestCSR(t, commonName)

	clientInfo := map[string]string{
		"common_name":  commonName,
		"organization": "Test Organization",
	}

	ctx := context.Background()
	certPEM, err := svc.IssueCertificate(ctx, csrPEM, clientInfo)
	require.NoError(t, err, "IssueCertificate should not return an error")

	// Now test renewal
	renewedCertPEM, err := svc.RenewCertificate(ctx, certPEM)
	require.NoError(t, err, "RenewCertificate should not return an error")
	require.NotEmpty(t, renewedCertPEM, "Renewed certificate should not be empty")

	// Ensure we got a different certificate
	assert.NotEqual(t, certPEM, renewedCertPEM, "Renewed certificate should be different from original")

	// Decode and check the renewed certificate
	block, _ := pem.Decode(renewedCertPEM)
	require.NotNil(t, block, "Renewed certificate PEM should decode successfully")

	renewedCert, err := x509.ParseCertificate(block.Bytes)
	require.NoError(t, err, "Should parse renewed certificate successfully")

	// Verify it has the same common name
	assert.Equal(t, commonName, renewedCert.Subject.CommonName,
		"Renewed certificate should have the same common name")
}

func TestGetRootCA(t *testing.T) {
	// Setup
	svc, tempDir := setupTestService(t)
	defer os.RemoveAll(tempDir)

	// Get the root CA certificate
	ctx := context.Background()
	caCertPEM, err := svc.GetRootCA(ctx)
	require.NoError(t, err, "GetRootCA should not return an error")
	require.NotEmpty(t, caCertPEM, "CA certificate should not be empty")

	// Decode the CA certificate
	block, _ := pem.Decode(caCertPEM)
	require.NotNil(t, block, "CA Certificate PEM should decode successfully")
	require.Equal(t, "CERTIFICATE", block.Type, "PEM block type should be CERTIFICATE")

	// Parse the certificate
	caCert, err := x509.ParseCertificate(block.Bytes)
	require.NoError(t, err, "Should parse CA certificate successfully")

	// Verify it's a CA certificate
	assert.True(t, caCert.IsCA, "Certificate should be a CA certificate")

	// Check if the certificate file was created in the temp directory
	files, err := os.ReadDir(tempDir)
	require.NoError(t, err, "Should be able to read temp directory")

	foundCAFile := false
	for _, file := range files {
		if file.Name() == "ca.crt" {
			foundCAFile = true
			break
		}
	}

	assert.True(t, foundCAFile, "CA certificate file should exist in storage directory")
}

func TestPersistenceAndReload(t *testing.T) {
	// Create a service and issue a certificate
	svc1, tempDir := setupTestService(t)
	defer os.RemoveAll(tempDir)

	ctx := context.Background()
	caCert1, err := svc1.GetRootCA(ctx)
	require.NoError(t, err)

	// Create a second service using the same directory
	svc2 := certissuer.NewService(tempDir)
	require.NotNil(t, svc2)

	// Get the CA from the second service
	caCert2, err := svc2.GetRootCA(ctx)
	require.NoError(t, err)

	// Compare the two CA certificates - they should be identical since loaded from the same files
	assert.True(t, bytes.Equal(caCert1, caCert2),
		"CA certificates from two service instances should be identical")
}

func TestInvalidCSR(t *testing.T) {
	svc, tempDir := setupTestService(t)
	defer os.RemoveAll(tempDir)

	// Test with invalid CSR data
	invalidCSR := []byte("not a valid CSR")
	clientInfo := map[string]string{"common_name": "test.example.com"}

	ctx := context.Background()
	_, err := svc.IssueCertificate(ctx, invalidCSR, clientInfo)
	assert.Error(t, err, "IssueCertificate should return an error with invalid CSR")
}

func TestInvalidCertificateRenewal(t *testing.T) {
	svc, tempDir := setupTestService(t)
	defer os.RemoveAll(tempDir)

	// Test with invalid certificate
	invalidCert := []byte("not a valid certificate")

	ctx := context.Background()
	_, err := svc.RenewCertificate(ctx, invalidCert)
	assert.Error(t, err, "RenewCertificate should return an error with invalid certificate")
}
