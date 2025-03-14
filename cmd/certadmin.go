// cmd/certadmin.go
package cmd

import (
	"context"
	"crypto/tls"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jdfalk/ubuntu-autoinstall-webhook/internal/proto/certadmin"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	serverAddr   string
	apiKeyFile   string
	apiKey       string
	insecureConn bool
	timeout      time.Duration
)

// Main certadmin command
var certAdminCmd = &cobra.Command{
	Use:   "certadmin",
	Short: "Certificate administration client",
	Long: `Certificate administration client for interacting with the certificate issuer server.
Example usage:
  certadmin get-ca -o ca.crt
  certadmin issue --csr client.csr --common-name example.com -o client.crt
  certadmin renew --cert client.crt -o renewed.crt
  certadmin list`,
	PersistentPreRunE: setupClient,
}

// Get CA certificate command
var getCACmd = &cobra.Command{
	Use:   "get-ca",
	Short: "Get the CA certificate",
	RunE:  getCACertificate,
}

// Issue certificate command
var issueCmd = &cobra.Command{
	Use:   "issue",
	Short: "Issue a new certificate",
	RunE:  issueCertificate,
}

// Renew certificate command
var renewCmd = &cobra.Command{
	Use:   "renew",
	Short: "Renew an existing certificate",
	RunE:  renewCertificate,
}

// List certificates command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all certificates",
	RunE:  listCertificates,
}

// API key management commands
var apiKeyCmd = &cobra.Command{
	Use:   "apikey",
	Short: "API key management",
}

var createAPIKeyCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new API key",
	RunE:  createAPIKey,
}

var listAPIKeysCmd = &cobra.Command{
	Use:   "list",
	Short: "List all API keys",
	RunE:  listAPIKeys,
}

var revokeAPIKeyCmd = &cobra.Command{
	Use:   "revoke",
	Short: "Revoke an API key",
	RunE:  revokeAPIKey,
}

func init() {
	rootCmd.AddCommand(certAdminCmd)

	// Add common flags
	certAdminCmd.PersistentFlags().StringVar(&serverAddr, "server", "localhost:8443", "Certificate server address")
	certAdminCmd.PersistentFlags().StringVar(&apiKeyFile, "api-key-file", "", "Path to API key file")
	certAdminCmd.PersistentFlags().StringVar(&apiKey, "api-key", "", "API key (if not using file)")
	certAdminCmd.PersistentFlags().BoolVar(&insecureConn, "insecure", false, "Use insecure connection")
	certAdminCmd.PersistentFlags().DurationVar(&timeout, "timeout", 30*time.Second, "Request timeout")

	// Add commands
	certAdminCmd.AddCommand(getCACmd)
	certAdminCmd.AddCommand(issueCmd)
	certAdminCmd.AddCommand(renewCmd)
	certAdminCmd.AddCommand(listCmd)
	certAdminCmd.AddCommand(apiKeyCmd)

	// API key commands
	apiKeyCmd.AddCommand(createAPIKeyCmd)
	apiKeyCmd.AddCommand(listAPIKeysCmd)
	apiKeyCmd.AddCommand(revokeAPIKeyCmd)

	// Get CA certificate flags
	getCACmd.Flags().StringP("output", "o", "ca.crt", "Output file for CA certificate")

	// Issue certificate flags
	issueCmd.Flags().String("csr", "", "Path to CSR file")
	issueCmd.Flags().String("common-name", "", "Common name for the certificate")
	issueCmd.Flags().String("org", "", "Organization for the certificate")
	issueCmd.Flags().StringSlice("sans", []string{}, "Subject Alternative Names")
	issueCmd.Flags().StringP("output", "o", "cert.crt", "Output file for certificate")
	issueCmd.MarkFlagRequired("csr")

	// Renew certificate flags
	renewCmd.Flags().String("cert", "", "Path to certificate file to renew")
	renewCmd.Flags().StringP("output", "o", "renewed.crt", "Output file for renewed certificate")
	renewCmd.MarkFlagRequired("cert")

	// Create API key flags
	createAPIKeyCmd.Flags().String("name", "", "Name for the API key")
	createAPIKeyCmd.Flags().String("description", "", "Description for the API key")
	createAPIKeyCmd.MarkFlagRequired("name")

	// Revoke API key flags
	revokeAPIKeyCmd.Flags().String("name", "", "Name of the API key to revoke")
	revokeAPIKeyCmd.MarkFlagRequired("name")
}

// setupClient sets up the gRPC client connection
func setupClient(cmd *cobra.Command, args []string) error {
	// Check if we're just getting server info (doesn't need auth)
	if cmd.CalledAs() == "server-info" {
		return nil
	}

	// Load API key if specified from file
	if apiKey == "" && apiKeyFile != "" {
		data, err := os.ReadFile(apiKeyFile)
		if err != nil {
			return fmt.Errorf("failed to read API key file: %w", err)
		}
		apiKey = strings.TrimSpace(string(data))
	}

	// Skip API key check for certain commands that don't require auth
	if cmd.CalledAs() != "get-ca" && apiKey == "" {
		return fmt.Errorf("API key is required")
	}

	return nil
}

// getGRPCClient creates a new gRPC client connection
func getGRPCClient() (*grpc.ClientConn, certadmin.CertAdminServiceClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var opts []grpc.DialOption

	if insecureConn {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		// Use TLS but skip verification for now (in a real system, you'd verify the server cert)
		tlsConfig := &tls.Config{
			InsecureSkipVerify: true, // Not recommended for production
		}
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))
	}

	conn, err := grpc.DialContext(ctx, serverAddr, opts...)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to server: %w", err)
	}

	client := certadmin.NewCertAdminServiceClient(conn)
	return conn, client, nil
}

// createAuthContext creates an authenticated context
func createAuthContext(ctx context.Context) context.Context {
	if apiKey != "" {
		return metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+apiKey)
	}
	return ctx
}

// getCACertificate gets the CA certificate
func getCACertificate(cmd *cobra.Command, args []string) error {
	outputFile, _ := cmd.Flags().GetString("output")

	conn, client, err := getGRPCClient()
	if err != nil {
		return err
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// No authentication needed for CA certificate
	res, err := client.GetCACertificate(ctx, &emptypb.Empty{})
	if err != nil {
		return fmt.Errorf("failed to get CA certificate: %w", err)
	}

	if err := os.WriteFile(outputFile, res.Certificate, 0644); err != nil {
		return fmt.Errorf("failed to write CA certificate: %w", err)
	}

	fmt.Printf("CA certificate saved to %s\n", outputFile)
	return nil
}

// issueCertificate issues a new certificate
func issueCertificate(cmd *cobra.Command, args []string) error {
	csrFile, _ := cmd.Flags().GetString("csr")
	commonName, _ := cmd.Flags().GetString("common-name")
	org, _ := cmd.Flags().GetString("org")
	sans, _ := cmd.Flags().GetStringSlice("sans")
	outputFile, _ := cmd.Flags().GetString("output")

	// Read CSR file
	csrData, err := os.ReadFile(csrFile)
	if err != nil {
		return fmt.Errorf("failed to read CSR file: %w", err)
	}

	conn, client, err := getGRPCClient()
	if err != nil {
		return err
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	req := &certadmin.IssueCertificateRequest{
		Csr:        csrData,
		CommonName: commonName,
		Org:        org,
		Sans:       sans,
	}

	ctx = createAuthContext(ctx)
	res, err := client.IssueCertificate(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to issue certificate: %w", err)
	}

	if err := os.WriteFile(outputFile, res.Certificate, 0644); err != nil {
		return fmt.Errorf("failed to write certificate: %w", err)
	}

	fmt.Printf("Certificate saved to %s\n", outputFile)
	return nil
}

// renewCertificate renews an existing certificate
func renewCertificate(cmd *cobra.Command, args []string) error {
	certFile, _ := cmd.Flags().GetString("cert")
	outputFile, _ := cmd.Flags().GetString("output")

	// Read certificate file
	certData, err := os.ReadFile(certFile)
	if err != nil {
		return fmt.Errorf("failed to read certificate file: %w", err)
	}

	conn, client, err := getGRPCClient()
	if err != nil {
		return err
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	req := &certadmin.RenewCertificateRequest{
		Certificate: certData,
	}

	ctx = createAuthContext(ctx)
	res, err := client.RenewCertificate(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to renew certificate: %w", err)
	}

	if err := os.WriteFile(outputFile, res.Certificate, 0644); err != nil {
		return fmt.Errorf("failed to write renewed certificate: %w", err)
	}

	fmt.Printf("Renewed certificate saved to %s\n", outputFile)
	return nil
}

// listCertificates lists all certificates
func listCertificates(cmd *cobra.Command, args []string) error {
	conn, client, err := getGRPCClient()
	if err != nil {
		return err
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	ctx = createAuthContext(ctx)
	res, err := client.ListCertificates(ctx, &emptypb.Empty{})
	if err != nil {
		return fmt.Errorf("failed to list certificates: %w", err)
	}

	for _, cert := range res.Certificates {
		fmt.Printf("Serial: %s, Common Name: %s, Org: %s, Expiry: %s\n",
			cert.SerialNumber, cert.CommonName, cert.Org, cert.Expiry)
	}

	return nil
}

// createAPIKey creates a new API key
func createAPIKey(cmd *cobra.Command, args []string) error {
	name, _ := cmd.Flags().GetString("name")
	description, _ := cmd.Flags().GetString("description")

	conn, client, err := getGRPCClient()
	if err != nil {
		return err
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	req := &certadmin.CreateAPIKeyRequest{
		Name:        name,
		Description: description,
	}

	ctx = createAuthContext(ctx)
	res, err := client.CreateAPIKey(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to create API key: %w", err)
	}

	fmt.Printf("API key created: %s\n", res.ApiKey)
	return nil
}

// listAPIKeys lists all API keys
func listAPIKeys(cmd *cobra.Command, args []string) error {
	conn, client, err := getGRPCClient()
	if err != nil {
		return err
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	ctx = createAuthContext(ctx)
	res, err := client.ListAPIKeys(ctx, &emptypb.Empty{})
	if err != nil {
		return fmt.Errorf("failed to list API keys: %w", err)
	}

	for _, key := range res.ApiKeys {
		fmt.Printf("Name: %s, Description: %s, Created: %s\n",
			key.Name, key.Description, key.CreatedAt)
	}

	return nil
}

// revokeAPIKey revokes an API key
func revokeAPIKey(cmd *cobra.Command, args []string) error {
	name, _ := cmd.Flags().GetString("name")

	conn, client, err := getGRPCClient()
	if err != nil {
		return err
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	req := &certadmin.RevokeAPIKeyRequest{
		Name: name,
	}

	ctx = createAuthContext(ctx)
	_, err = client.RevokeAPIKey(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to revoke API key: %w", err)
	}

	fmt.Printf("API key %s revoked\n", name)
	return nil
}
