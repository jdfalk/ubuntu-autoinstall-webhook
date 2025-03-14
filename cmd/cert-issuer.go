// cmd/cert-issuer.go
package cmd

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/jdfalk/ubuntu-autoinstall-webhook/internal/certadmin"
	"github.com/jdfalk/ubuntu-autoinstall-webhook/internal/certissuer"
	pb "github.com/jdfalk/ubuntu-autoinstall-webhook/pkg/proto"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	certStoragePath string
	httpListenAddr  string
	grpcListenAddr  string
	apiKey          string
	generateApiKey  bool
)

var certIssuerCmd = &cobra.Command{
	Use:   "cert-issuer",
	Short: "Starts the cert-issuer microservice",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// Generate API key if requested
		if generateApiKey {
			key, err := generateRandomAPIKey(32)
			if err != nil {
				return fmt.Errorf("failed to generate API key: %w", err)
			}
			fmt.Printf("Generated API key: %s\n", key)
			// Store this key for use
			apiKey = key
		}

		// Validate that we have an API key
		if apiKey == "" {
			return fmt.Errorf("API key is required, either provide one with --api-key or generate one with --generate-api-key")
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Starting cert-issuer microservice...")

		// Set default storage path if not provided
		if certStoragePath == "" {
			homeDir, err := os.UserHomeDir()
			if err == nil {
				certStoragePath = filepath.Join(homeDir, ".autoinstall-webhook", "certificates")
			} else {
				certStoragePath = "/var/lib/autoinstall-webhook/certificates"
			}
		}

		fmt.Printf("Using certificate storage path: %s\n", certStoragePath)

		// Create certificate service
		certService := certissuer.NewService(certStoragePath)

		// Setup API keys
		apiKeys := map[string]string{
			apiKey: "admin", // Use provided or generated API key
		}

		// Start HTTP server
		httpServer := startHTTPServer(certService, httpListenAddr)

		// Start gRPC server
		grpcServer, grpcListener := startGRPCServer(certService, apiKeys, grpcListenAddr)

		// Wait for interrupt or termination signal
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
		<-stop

		fmt.Println("Shutting down servers...")

		// Graceful shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Stop HTTP server
		if err := httpServer.Shutdown(ctx); err != nil {
			fmt.Printf("Error shutting down HTTP server: %v\n", err)
		}

		// Stop gRPC server
		grpcServer.GracefulStop()
		grpcListener.Close()

		fmt.Println("Servers stopped gracefully")
		return nil
	},
}

// startHTTPServer starts the HTTP API server
func startHTTPServer(certService certissuer.CertIssuer, addr string) *http.Server {
	// Create HTTP router
	mux := http.NewServeMux()

	// Handler for CA certificate requests
	mux.HandleFunc("/api/v1/ca", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		caCert, err := certService.GetRootCA(r.Context())
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to get CA certificate: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/x-pem-file")
		w.Write(caCert)
	})

	// Handler for certificate issuance
	mux.HandleFunc("/api/v1/issue", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Parse the request body
		var req struct {
			CSR        string            `json:"csr"`
			ClientInfo map[string]string `json:"client_info"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, fmt.Sprintf("Failed to parse request: %v", err), http.StatusBadRequest)
			return
		}

		if req.CSR == "" {
			http.Error(w, "CSR is required", http.StatusBadRequest)
			return
		}

		if req.ClientInfo == nil {
			req.ClientInfo = make(map[string]string)
		}

		// Issue certificate
		cert, err := certService.IssueCertificate(r.Context(), []byte(req.CSR), req.ClientInfo)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to issue certificate: %v", err), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"certificate": string(cert),
		})
	})

	// Handler for certificate renewal
	mux.HandleFunc("/api/v1/renew", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Parse the request body
		var req struct {
			Certificate string `json:"certificate"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, fmt.Sprintf("Failed to parse request: %v", err), http.StatusBadRequest)
			return
		}

		if req.Certificate == "" {
			http.Error(w, "Certificate is required", http.StatusBadRequest)
			return
		}

		// Renew certificate
		renewedCert, err := certService.RenewCertificate(r.Context(), []byte(req.Certificate))
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to renew certificate: %v", err), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"certificate": string(renewedCert),
		})
	})

	// Simple health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	// Create and start HTTP server
	server := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		fmt.Printf("HTTP server listening on %s\n", addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("HTTP server error: %v\n", err)
		}
	}()

	return server
}

// startGRPCServer starts the gRPC server
func startGRPCServer(certService certissuer.CertIssuer, apiKeys map[string]string, addr string) (*grpc.Server, net.Listener) {
	// Create authentication interceptor
	authInterceptor := certadmin.NewAuthInterceptor(apiKeys)

	// Create gRPC server with authentication
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(authInterceptor.Unary()),
	)

	// Create config directory for API keys
	configDir := filepath.Join(certStoragePath, "../admin")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		panic(fmt.Sprintf("failed to create config directory: %v", err))
	}

	// Register the CertAdmin service - update to match your actual server constructor
	certAdminServer, err := certadmin.NewServer(certService, configDir, "", "")
	if err != nil {
		panic(fmt.Sprintf("failed to create cert admin server: %v", err))
	}
	pb.RegisterCertAdminServer(grpcServer, certAdminServer)

	// Enable reflection for easier client development
	reflection.Register(grpcServer)

	// Listen on the specified port
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		panic(fmt.Sprintf("failed to listen on %s: %v", addr, err))
	}

	// Start gRPC server
	go func() {
		fmt.Printf("gRPC server listening on %s\n", addr)
		if err := grpcServer.Serve(lis); err != nil {
			fmt.Printf("gRPC server error: %v\n", err)
		}
	}()

	return grpcServer, lis
}

// generateRandomAPIKey creates a secure random API key
func generateRandomAPIKey(length int) (string, error) {
	// Generate random bytes
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	// Encode to base64
	key := base64.URLEncoding.EncodeToString(bytes)
	// Trim to the requested length
	if len(key) > length {
		key = key[:length]
	}

	return key, nil
}

func init() {
	rootCmd.AddCommand(certIssuerCmd)

	// Add command-line flags
	certIssuerCmd.Flags().StringVar(&certStoragePath, "cert-path", "", "Path to store certificates (default: ~/.autoinstall-webhook/certificates)")
	certIssuerCmd.Flags().StringVar(&httpListenAddr, "http-listen", ":8443", "Address to listen on for HTTP requests")
	certIssuerCmd.Flags().StringVar(&grpcListenAddr, "grpc-listen", ":9443", "Address to listen on for gRPC requests")
	certIssuerCmd.Flags().StringVar(&apiKey, "api-key", "", "API key for authenticating admin requests")
	certIssuerCmd.Flags().BoolVar(&generateApiKey, "generate-api-key", false, "Generate and print a random API key")
}
