// internal/certadmin/server.go
package certadmin

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
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
	"github.com/jdfalk/ubuntu-autoinstall-webhook/internal/proto/certadmin"
	pb "github.com/jdfalk/ubuntu-autoinstall-webhook/pkg/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	jwtSecret        = []byte("replace-with-secure-secret-in-production")
	tokenExpiryHours = 24
)

// Server implements the gRPC CertAdmin service
type Server struct {
	certadmin.UnimplementedCertAdminServiceServer
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
	certadmin.RegisterCertAdminServiceServer(grpcServer, s)

	// Enable reflection for tools like grpcurl
	reflection.Register(grpcServer)

	log.Printf("CertAdmin gRPC server listening on %s", listenAddr)
	return grpcServer.Serve(listener)
}

// Authentication interceptor for gRPC
func (s *Server) authInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// Skip authentication for GetServerInfo method
	if info.FullMethod == "/certadmin.CertAdminService/GetServerInfo" {
		return handler(ctx, req)
	}

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

// GetServerInfo returns information about the server
func (s *Server) GetServerInfo(ctx context.Context, _ *emptypb.Empty) (*certadmin.ServerInfoResponse, error) {
	// Get information about the client
	clientInfo := "unknown"
	if p, ok := peer.FromContext(ctx); ok {
		clientInfo = p.Addr.String()
	}

	// Build server info
	serverInfo := &certadmin.ServerInfoResponse{
		ServerVersion: "1.0.0",
		ApiVersion:    "v1",
		ClientAddress: clientInfo,
		ServerTime:    timestamppb.Now(),
		RequiresAuth:  true,
	}

	return serverInfo, nil
}

// GetCACertificate retrieves the CA certificate
func (s *Server) GetCACertificate(ctx context.Context, _ *emptypb.Empty) (*certadmin.CACertificateResponse, error) {
	caCert, err := s.certIssuer.GetRootCA(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get CA certificate: %v", err)
	}

	return &certadmin.CACertificateResponse{
		Certificate: caCert,
	}, nil
}

// IssueCertificate issues a new certificate from a CSR
func (s *Server) IssueCertificate(ctx context.Context, req *certadmin.IssueCertificateRequest) (*certadmin.CertificateResponse, error) {
	if req.Csr == nil || len(req.Csr) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "CSR is required")
	}

	clientInfo := make(map[string]string)
	if req.CommonName != "" {
		clientInfo["common_name"] = req.CommonName
	}
	if req.Organization != "" {
		clientInfo["organization"] = req.Organization
	}
	if len(req.Sans) > 0 {
		clientInfo["sans"] = strings.Join(req.Sans, ",")
	}

	cert, err := s.certIssuer.IssueCertificate(ctx, req.Csr, clientInfo)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to issue certificate: %v", err)
	}

	return &certadmin.CertificateResponse{
		Certificate: cert,
	}, nil
}

// RenewCertificate renews an existing certificate
func (s *Server) RenewCertificate(ctx context.Context, req *certadmin.RenewCertificateRequest) (*certadmin.CertificateResponse, error) {
	if req.Certificate == nil || len(req.Certificate) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "certificate is required")
	}

	renewedCert, err := s.certIssuer.RenewCertificate(ctx, req.Certificate)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to renew certificate: %v", err)
	}

	return &certadmin.CertificateResponse{
		Certificate: renewedCert,
	}, nil
}

// ListCertificates lists all issued certificates
func (s *Server) ListCertificates(ctx context.Context, _ *emptypb.Empty) (*certadmin.ListCertificatesResponse, error) {
	// This would require additional functionality in the certIssuer interface
	// For now, return a placeholder
	return &certadmin.ListCertificatesResponse{
		Certificates: []*certadmin.CertificateInfo{
			{
				SerialNumber: "placeholder",
				Subject:      "placeholder",
				NotBefore:    timestamppb.Now(),
				NotAfter:     timestamppb.Now(),
				IsRevoked:    false,
			},
		},
	}, nil
}

// RevokeCertificate revokes a certificate
func (s *Server) RevokeCertificate(ctx context.Context, req *certadmin.RevokeCertificateRequest) (*emptypb.Empty, error) {
	// This would require additional functionality in the certIssuer interface
	// For now, return a placeholder
	return &emptypb.Empty{}, nil
}

// CreateAPIKey creates a new API key
func (s *Server) CreateAPIKey(ctx context.Context, req *certadmin.CreateAPIKeyRequest) (*certadmin.APIKeyResponse, error) {
	if req.Name == "" {
		return nil, status.Errorf(codes.InvalidArgument, "name is required")
	}

	apiKey, err := generateAPIKey()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate API key: %v", err)
	}

	keyInfo := apiKeyInfo{
		Name:        req.Name,
		Key:         apiKey,
		CreatedAt:   time.Now(),
		Description: req.Description,
	}

	s.apiKeysMutex.Lock()
	s.apiKeys[apiKey] = keyInfo
	s.apiKeysMutex.Unlock()

	// Save the updated API keys
	if err := s.saveAPIKeys(); err != nil {
		log.Printf("Warning: Failed to save API keys: %v", err)
	}

	return &certadmin.APIKeyResponse{
		Name:        keyInfo.Name,
		Key:         apiKey,
		CreatedAt:   timestamppb.New(keyInfo.CreatedAt),
		Description: keyInfo.Description,
	}, nil
}

// ListAPIKeys lists all API keys
func (s *Server) ListAPIKeys(ctx context.Context, _ *emptypb.Empty) (*certadmin.ListAPIKeysResponse, error) {
	s.apiKeysMutex.RLock()
	defer s.apiKeysMutex.RUnlock()

	keys := make([]*certadmin.APIKeyInfo, 0, len(s.apiKeys))
	for _, keyInfo := range s.apiKeys {
		keys = append(keys, &certadmin.APIKeyInfo{
			Name:        keyInfo.Name,
			CreatedAt:   timestamppb.New(keyInfo.CreatedAt),
			LastUsedAt:  timestamppb.New(keyInfo.LastUsedAt),
			Description: keyInfo.Description,
		})
	}

	return &certadmin.ListAPIKeysResponse{
		Keys: keys,
	}, nil
}

// RevokeAPIKey revokes an API key
func (s *Server) RevokeAPIKey(ctx context.Context, req *certadmin.RevokeAPIKeyRequest) (*emptypb.Empty, error) {
	if req.Name == "" {
		return nil, status.Errorf(codes.InvalidArgument, "name is required")
	}

	s.apiKeysMutex.Lock()
	defer s.apiKeysMutex.Unlock()

	// Find and remove the API key with the given name
	var found bool
	for key, info := range s.apiKeys {
		if info.Name == req.Name {
			delete(s.apiKeys, key)
			found = true
			break
		}
	}

	if !found {
		return nil, status.Errorf(codes.NotFound, "API key not found")
	}

	// Save the updated API keys
	if err := s.saveAPIKeys(); err != nil {
		log.Printf("Warning: Failed to save API keys: %v", err)
	}

	return &emptypb.Empty{}, nil
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

// validateAPIKey checks if the provided API key is valid
func (s *Server) validateAPIKey(apiKey string) bool {
	s.apiKeysMutex.RLock()
	defer s.apiKeysMutex.RUnlock()

	_, valid := s.apiKeys[apiKey]
	return valid
}

// compareAPIKeys compares two API keys in constant time to prevent timing attacks
func compareAPIKeys(a, b string) bool {
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
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

	// In a real implementation, you'd extract the serial number from the certificate
	// For now, we'll use a placeholder
	serialNumber := "placeholder-serial-number"

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

	// We'll need to add a RevokeCertificate method to the CertIssuer interface
	// For now, we'll simulate success
	return &pb.RevokeCertificateResponse{
		Success: true,
	}, nil
}

// ListCertificates lists all issued certificates
func (s *Server) ListCertificates(ctx context.Context, req *pb.ListCertificatesRequest) (*pb.ListCertificatesResponse, error) {
	// In a real implementation, this would query the certificate store
	// For now, return an empty list
	return &pb.ListCertificatesResponse{
		Certificates: []*pb.CertificateInfo{},
	}, nil
}

// GetCertificateInfo gets detailed information about a certificate
func (s *Server) GetCertificateInfo(ctx context.Context, req *pb.GetCertificateInfoRequest) (*pb.GetCertificateInfoResponse, error) {
	// Validate request
	if req.SerialNumber == "" {
		return nil, status.Error(codes.InvalidArgument, "Serial number is required")
	}

	// In a real implementation, this would query the certificate store
	// For now, return a placeholder certificate
	now := time.Now()
	return &pb.GetCertificateInfoResponse{
		SerialNumber: req.SerialNumber,
		Subject:      "placeholder-subject",
		Issuer:       "placeholder-issuer",
		NotBefore:    timestamppb.New(now),
		NotAfter:     timestamppb.New(now.Add(365 * 24 * time.Hour)),
	}, nil
}
