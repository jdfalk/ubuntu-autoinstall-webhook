// cmd/cert-issuer.go
package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/jdfalk/ubuntu-autoinstall-webhook/internal/certissuer"
	"github.com/spf13/cobra"
)

var (
	certStoragePath string
	listenAddr      string
)

var certIssuerCmd = &cobra.Command{
	Use:   "cert-issuer",
	Short: "Starts the cert-issuer microservice",
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

		// Create an HTTP server to handle certificate requests
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
		serverErr := make(chan error, 1)
		go func() {
			fmt.Printf("Certificate Issuer server listening on %s\n", listenAddr)
			serverErr <- server.ListenAndServe()
		}()

		// Wait for interrupt signal or server error
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

		select {
		case err := <-serverErr:
			if err != http.ErrServerClosed {
				fmt.Printf("Error starting server: %v\n", err)
				return err
			}
		case <-stop:
			fmt.Println("Shutting down server...")
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			if err := server.Shutdown(ctx); err != nil {
				fmt.Printf("Error shutting down server: %v\n", err)
				return err
			}
		}

		fmt.Println("Server stopped gracefully")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(certIssuerCmd)

	// Add command-line flags
	certIssuerCmd.Flags().StringVar(&certStoragePath, "cert-path", "", "Path to store certificates (default: ~/.autoinstall-webhook/certificates)")
	certIssuerCmd.Flags().StringVar(&listenAddr, "listen", ":8443", "Address to listen on for certificate requests")
}
