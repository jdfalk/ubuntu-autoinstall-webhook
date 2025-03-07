// cmd/cert-issuer.go
package cmd

import (
    "fmt"

    "github.com/spf13/cobra"
    "github.com/jdfalk/ubuntu-autoinstall-webhook/internal/certissuer"
)

var certIssuerCmd = &cobra.Command{
    Use:   "cert-issuer",
    Short: "Starts the cert-issuer microservice",
    RunE: func(cmd *cobra.Command, args []string) error {
        fmt.Println("Starting cert-issuer microservice...")
        certService := certissuer.NewService()
        // Issue a dummy certificate for demonstration.
        dummyCSR := []byte("dummy csr")
        cert, err := certService.IssueCertificate(dummyCSR, map[string]string{"client": "dummy"})
        if err != nil {
            return err
        }
        fmt.Println("Issued certificate:", string(cert))
        fmt.Println("Cert-issuer microservice started successfully.")
        return nil
    },
}

func init() {
    rootCmd.AddCommand(certIssuerCmd)
}
