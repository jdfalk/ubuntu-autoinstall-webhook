// cmd/cert-issuer.go
package cmd

import (
	"context"
	"fmt"

	"github.com/jdfalk/ubuntu-autoinstall-webhook/internal/certissuer"
	"github.com/spf13/cobra"
)

var certIssuerCmd = &cobra.Command{
	Use:   "cert-issuer",
	Short: "Starts the cert-issuer microservice",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Starting cert-issuer microservice...")
		certService := certissuer.NewService()
		// Pass context.Background() along with dummy CSR and client information.
		dummyCSR := []byte("dummy csr")
		cert, err := certService.IssueCertificate(context.Background(), dummyCSR, map[string]string{"client": "dummy"})
		if err != nil {
			return err
		}
		fmt.Println("Issued certificate:", string(cert))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(certIssuerCmd)
}
