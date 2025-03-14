// cmd/cert-issuer.go
package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/jdfalk/ubuntu-autoinstall-webhook/internal/certadmin"
	"github.com/jdfalk/ubuntu-autoinstall-webhook/internal/certissuer"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

var (
	certStoragePath   string
	httpListenAddr    string
	grpcListenAddr    string
	adminAPIKey       string
	tlsCertFile       string
	tlsKeyFile        string
	enableGRPCReflect bool
)

var certIssuerCmd = &cobra.Command{
	Use:   "cert-issuer",
	Short: "Starts the cert-issuer microservice with HTTP API and gRPC admin interface",
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

		// Store the admin API key if provided
		if adminAPIKey == "" {
			adminAPIKey = os.Getenv("CERT_ADMIN_API_KEY")
			if adminAPIKey == "" {
				// Generate a random API key if not provided
				adminAPIKey = generateRandomAPIKey(32)
				fmt.Printf("Generated admin API key: %s\n", adminAPIKey)
				fmt.Println("Store this key securely for future admin operations!")
			}
		}

		// Start HTTP server
		httpServer, httpErrCh := startHTTPServer(certService, httpListenAddr)

		// Start gRPC server
		grpcServer, grpcErrCh := startGRPCServer(certService, grpcListenAddr)

		// Wait for interrupt signal or server errors
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

		select {
		case err := <-httpErrCh:
			fmt.Printf("HTTP server error: %v\n", err)
			return err
		case err := <-grpcErrCh:
			fmt.Printf("gRPC server error: %v\n", err)
			return err
		case <-stop:
			fmt.Println("Shutting down servers...")

			// Gracefully stop the HTTP server
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			if err := httpServer.Shutdown(ctx); err != nil {
				fmt.Printf("Error shutting down HTTP server: %v\n", err)
			}

			// Gracefully stop the gRPC server
			grpcServer.GracefulStop()
		}

		fmt.Println("Servers stopped gracefully")
		return nil
	},
}

func startHTTPServer(certService certissuer.CertIssuer, listenAddr string) (*http.Server, chan error) {
	// Create HTTP mux and handlers (same code as before)
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

		// Parse the request body to get CSR and client info
		var req struct {
			CSR        string            `json:"csr"`
			ClientInfo map[string]string `json:"client_info"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, fmt.Sprintf("Failed to parse request: %v", err), http.StatusBadRequest)
			return
		}

		// Validate input
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

		// Return the certificate
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

		// Parse the request body to get the certificate
		var req struct {
			Certificate string `json:"certificate"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, fmt.Sprintf("Failed to parse request: %v", err), http.StatusBadRequest)
			return
		}

		// Validate input
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

		// Return the renewed certificate
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

	// Create an HTTP server
	server := &http.Server{
		Addr:         listenAddr,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Start the server in a goroutine
	errCh := make(chan error, 1)
	go func() {
		fmt.Printf("HTTP server listening on %s\n", listenAddr)
		errCh <- server.ListenAndServe()
	}()

	return server, errCh
}

func startGRPCServer(certService certissuer.CertIssuer, listenAddr string) (*grpc.Server, chan error) {
	errCh := make(chan error, 1)

	// Create a listener for the gRPC server
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		errCh <- fmt.Errorf("failed to listen on %s: %w", listenAddr, err)
		return nil, errCh
	}

	var opts []grpc.ServerOption

	// Add TLS if certificate and key files are provided
	if tlsCertFile != "" && tlsKeyFile != "" {
		fmt.Printf("Using TLS for gRPC with cert=%s, key=%s\n", tlsCertFile, tlsKeyFile)
		creds, err := credentials.NewServerTLSFromFile(tlsCertFile, tlsKeyFile)
		if err != nil {
			errCh <- fmt.Errorf("failed to load TLS credentials: %w", err)
			return nil, errCh
		}
		opts = append(opts, grpc.Creds(creds))
	} else {
		fmt.Println("Warning: gRPC server running without TLS")
	}

	// Add authentication interceptor
	opts = append(opts, grpc.UnaryInterceptor(authInterceptor))

	// Create a gRPC server with the options
	grpcServer := grpc.NewServer(opts...)

	// Create and register the admin service
	adminService := certadmin.NewServer(certService)
	certadmin.RegisterCertAdminServer(grpcServer, adminService)

	// Enable reflection if requested
	if enableGRPCReflect {
		fmt.Println("Enabling gRPC reflection (useful for debugging)")
		reflection.Register(grpcServer)
	}

	// Start the gRPC server in a goroutine
	go func() {
		fmt.Printf("gRPC server listening on %s\n", listenAddr)
		if err := grpcServer.Serve(listener); err != nil {
			errCh <- fmt.Errorf("gRPC server error: %w", err)
		}
	}()

	return grpcServer, errCh
}

// Authentication interceptor for gRPC
func authInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// Extract authentication from request
	var authMsg *certadmin.Auth

	// Check for auth field in various request types
	switch v := req.(type) {
	case *certadmin.GetCARequest:
		authMsg = v.Auth
	case *certadmin.IssueCertificateRequest:
		authMsg = v.Auth
	case *certadmin.RenewCertificateRequest:
		authMsg = v.Auth
	case *certadmin.RevokeCertificateRequest:
		authMsg = v.Auth
	case *certadmin.ListCertificatesRequest:
		authMsg = v.Auth
	}

	// Validate the API key
	if authMsg == nil || authMsg.ApiKey != adminAPIKey {
		return nil, grpc.Errorf(grpc.Code(401), "unauthorized: invalid API key")
	}

	// Proceed with the request
	return handler(ctx, req)
}

// Helper function to generate a random API key
func generateRandomAPIKey(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func init() {
	rootCmd.AddCommand(certIssuerCmd)

	// Add command-line flags
	certIssuerCmd.Flags().StringVar(&certStoragePath, "cert-path", "", "Path to store certificates (default: ~/.autoinstall-webhook/certificates)")
	certIssuerCmd.Flags().StringVar(&httpListenAddr, "http-listen", ":8443", "Address to listen on for HTTP certificate requests")
	certIssuerCmd.Flags().StringVar(&grpcListenAddr, "grpc-listen", ":8444", "Address to listen on for gRPC admin requests")
	certIssuerCmd.Flags().StringVar(&adminAPIKey, "admin-api-key", "", "API key for admin operations (will generate one if not provided)")
	certIssuerCmd.Flags().StringVar(&tlsCertFile, "tls-cert", "", "TLS certificate file for gRPC server")
	certIssuerCmd.Flags().StringVar(&tlsKeyFile, "tls-key", "", "TLS key file for gRPC server")
	certIssuerCmd.Flags().BoolVar(&enableGRPCReflect, "grpc-reflect", false, "Enable gRPC reflection (useful for debugging)")
}
