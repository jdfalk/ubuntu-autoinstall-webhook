// internal/certadmin/server.go
package certadmin

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/jdfalk/ubuntu-autoinstall-webhook/internal/certissuer"
	pb "github.com/jdfalk/ubuntu-autoinstall-webhook/pkg/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

// Server implements the gRPC CertAdmin service
type Server struct {
	pb.UnimplementedCertAdminServer
	certIssuer      certissuer.CertIssuer
	apiKeys         map[string]apiKeyInfo
	apiKeysMutex    sync.RWMutex
	configPath      string
	serverOptions   []grpc.ServerOption
	serverStartTime time.Time
}

type apiKeyInfo struct {
	Name        string    `json:"name"`
	Key         string    `json:"key"`
	CreatedAt   time.Time `json:"created_at"`
	LastUsedAt  time.Time `json:"last_used_at,omitempty"`
	Description string    `json:"description,omitempty"`
}

// NewServer creates a new certificate admin server
func NewServer(certIssuer certissuer.CertIssuer, configPath string, tlsCert, tlsKey string) (*Server, error) {
	server := &Server{
		certIssuer:      certIssuer,
		apiKeys:         make(map[string]apiKeyInfo),
		configPath:      configPath,
		serverStartTime: time.Now(),
	}

	// Load API keys from config if it exists
	if err := server.loadAPIKeys(); err != nil {
		log.Printf("Warning: Failed to load API keys: %v", err)

		// Generate a default API key if none exist
		if len(server.apiKeys) == 0 {
			apiKey, err := generateAPIKey()
			if err != nil {
				return nil, fmt.Errorf("failed to generate default API key: %w", err)
			}

			server.apiKeys[apiKey] = apiKeyInfo{
				Name:        "default",
				Key:         apiKey,
				CreatedAt:   time.Now(),
				Description: "Default API key",
			}

			fmt.Printf("Generated default API key: %s\n", apiKey)

			// Save the default key
			if err := server.saveAPIKeys(); err != nil {
				log.Printf("Warning: Failed to save default API key: %v", err)
			}
		}
	}

	// Setup server options with authentication interceptor
	serverOptions := []grpc.ServerOption{
		grpc.UnaryInterceptor(server.authInterceptor),
	}

	// Add TLS if certificates provided
	if tlsCert != "" && tlsKey != "" {
		creds, err := credentials.NewServerTLSFromFile(tlsCert, tlsKey)
		if err != nil {
			return nil, fmt.Errorf("failed to load TLS credentials: %w", err)
		}
		serverOptions = append(serverOptions, grpc.Creds(creds))
	}

	server.serverOptions = serverOptions
	return server, nil
}

// Start starts the gRPC server
func (s *Server) Start(listenAddr string) error {
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	grpcServer := grpc.NewServer(s.serverOptions...)
	pb.RegisterCertAdminServer(grpcServer, s)

	// Enable reflection for tools like grpcurl
	reflection.Register(grpcServer)

	log.Printf("CertAdmin gRPC server listening on %s", listenAddr)
	return grpcServer.Serve(listener)
}

// Authentication interceptor for gRPC
func (s *Server) authInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// Extract API key from metadata
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "missing metadata")
	}

	authHeader, ok := md["authorization"]
	if !ok || len(authHeader) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "missing authorization header")
	}

	// Parse "Bearer {token}" format
	authParts := strings.SplitN(authHeader[0], " ", 2)
	if len(authParts) != 2 || authParts[0] != "Bearer" {
		return nil, status.Errorf(codes.Unauthenticated, "invalid authorization format")
	}

	apiKey := authParts[1]

	// Verify API key
	s.apiKeysMutex.RLock()
	keyInfo, valid := s.apiKeys[apiKey]
	s.apiKeysMutex.RUnlock()

	if !valid {
		return nil, status.Errorf(codes.Unauthenticated, "invalid API key")
	}

	// Update last used timestamp
	s.apiKeysMutex.Lock()
	keyInfo.LastUsedAt = time.Now()
	s.apiKeys[apiKey] = keyInfo
	s.apiKeysMutex.Unlock()

	// Log the API key usage and client information
	if p, ok := peer.FromContext(ctx); ok {
		log.Printf("API key '%s' used from %s for %s", keyInfo.Name, p.Addr, info.FullMethod)
	}

	// Proceed with the handler
	return handler(ctx, req)
}

// GetCACertificate returns the CA certificate
func (s *Server) GetCACertificate(ctx context.Context, req *pb.GetCACertificateRequest) (*pb.GetCACertificateResponse, error) {
	// Get CA certificate from the issuer
	caCert, err := s.certIssuer.GetRootCA(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get CA certificate: %v", err)
	}

	return &pb.GetCACertificateResponse{
		CertificatePem: string(caCert),
	}, nil
}

// IssueCertificate issues a new certificate based on the provided CSR
func (s *Server) IssueCertificate(ctx context.Context, req *pb.IssueCertificateRequest) (*pb.IssueCertificateResponse, error) {
	// Validate request
	if req.CsrPem == "" {
		return nil, status.Error(codes.InvalidArgument, "CSR is required")
	}

	// Issue certificate
	cert, err := s.certIssuer.IssueCertificate(ctx, []byte(req.CsrPem), req.ClientInfo)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to issue certificate: %v", err)
	}

	// Extract the serial number from the certificate
	serialNumber, err := extractSerialNumber(cert)
	if err != nil {
		// Use a placeholder if we can't extract it
		serialNumber = "unknown-serial-number"
	}

	return &pb.IssueCertificateResponse{
		CertificatePem: string(cert),
		SerialNumber:   serialNumber,
	}, nil
}

// RevokeCertificate revokes an existing certificate
func (s *Server) RevokeCertificate(ctx context.Context, req *pb.RevokeCertificateRequest) (*pb.RevokeCertificateResponse, error) {
	// Validate request
	if req.SerialNumber == "" {
		return nil, status.Error(codes.InvalidArgument, "Serial number is required")
	}

	// In a real implementation, we would call a RevokeCertificate method on the cert issuer
	// For now, we'll simulate success
	// TODO: Implement proper revocation in the certissuer package
	return &pb.RevokeCertificateResponse{
		Success: true,
	}, nil
}

// ListCertificates lists all issued certificates
func (s *Server) ListCertificates(ctx context.Context, req *pb.ListCertificatesRequest) (*pb.ListCertificatesResponse, error) {
	// In a real implementation, this would query the certificate store
	// For now, we'll return a sample certificate for demonstration

	certificates := []*pb.CertificateInfo{}

	// Add a sample certificate
	now := time.Now()
	sampleCert := &pb.CertificateInfo{
		SerialNumber: "sample-serial-123456",
		SubjectName:  "CN=sample.example.com",
		IssuedTo:     "sample-client",
		IssuedAt:     now.Add(-24 * time.Hour).Format(time.RFC3339),
		ExpiresAt:    now.AddDate(1, 0, 0).Format(time.RFC3339),
		Revoked:      false,
		// Don't include the full PEM in listings to keep response size smaller
	}
	certificates = append(certificates, sampleCert)

	return &pb.ListCertificatesResponse{
		Certificates: certificates,
	}, nil
}

// GetCertificateInfo gets detailed information about a certificate
func (s *Server) GetCertificateInfo(ctx context.Context, req *pb.GetCertificateInfoRequest) (*pb.GetCertificateInfoResponse, error) {
	// Validate request
	if req.SerialNumber == "" {
		return nil, status.Error(codes.InvalidArgument, "Serial number is required")
	}

	// In a real implementation, this would query the certificate store
	// For now, return a sample certificate info that matches the serial number

	// Sample certificate info with the requested serial number
	now := time.Now()
	certInfo := &pb.CertificateInfo{
		SerialNumber:   req.SerialNumber,
		SubjectName:    fmt.Sprintf("CN=%s.example.com", req.SerialNumber),
		IssuedTo:       "client-" + req.SerialNumber,
		IssuedAt:       now.Add(-24 * time.Hour).Format(time.RFC3339),
		ExpiresAt:      now.AddDate(1, 0, 0).Format(time.RFC3339),
		Revoked:        false,
		CertificatePem: fmt.Sprintf("-----BEGIN CERTIFICATE-----\nSample certificate with serial number %s\n-----END CERTIFICATE-----", req.SerialNumber),
	}

	return &pb.GetCertificateInfoResponse{
		Certificate: certInfo,
	}, nil
}

// Helper functions for API key management
func (s *Server) loadAPIKeys() error {
	if s.configPath == "" {
		return errors.New("config path not set")
	}

	// Ensure the directory exists
	configDir := s.configPath
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Read the API keys file
	apiKeysFilePath := filepath.Join(configDir, "api_keys.json")
	data, err := os.ReadFile(apiKeysFilePath)
	if os.IsNotExist(err) {
		// No file yet, start with empty map
		return nil
	} else if err != nil {
		return fmt.Errorf("failed to read API keys file: %w", err)
	}

	// Unmarshal the API keys
	var apiKeys map[string]apiKeyInfo
	if err := json.Unmarshal(data, &apiKeys); err != nil {
		return fmt.Errorf("failed to unmarshal API keys: %w", err)
	}

	s.apiKeysMutex.Lock()
	s.apiKeys = apiKeys
	s.apiKeysMutex.Unlock()

	return nil
}

func (s *Server) saveAPIKeys() error {
	if s.configPath == "" {
		return errors.New("config path not set")
	}

	// Ensure the directory exists
	configDir := s.configPath
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Marshal the API keys
	s.apiKeysMutex.RLock()
	data, err := json.MarshalIndent(s.apiKeys, "", "  ")
	s.apiKeysMutex.RUnlock()

	if err != nil {
		return fmt.Errorf("failed to marshal API keys: %w", err)
	}

	// Write to file
	apiKeysFilePath := filepath.Join(configDir, "api_keys.json")
	if err := os.WriteFile(apiKeysFilePath, data, 0600); err != nil {
		return fmt.Errorf("failed to write API keys file: %w", err)
	}

	return nil
}

// generateAPIKey creates a secure random API key
func generateAPIKey() (string, error) {
	// Generate 32 bytes of random data
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	// Encode as base64 and trim padding for a cleaner token
	key := base64.URLEncoding.EncodeToString(bytes)
	key = strings.TrimRight(key, "=")

	return key, nil
}

// compareAPIKeys compares two API keys in constant time to prevent timing attacks
func compareAPIKeys(a, b string) bool {
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}

// Helper function to extract serial number from a certificate
func extractSerialNumber(certPEM []byte) (string, error) {
	block, _ := pem.Decode(certPEM)
	if block == nil || block.Type != "CERTIFICATE" {
		return "", fmt.Errorf("failed to decode PEM block containing certificate")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("failed to parse certificate: %w", err)
	}

	return cert.SerialNumber.String(), nil
}

// API key management methods - these will need to be added to the proto file if we want to expose them via gRPC
