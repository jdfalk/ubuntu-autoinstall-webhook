// internal/certissuer/service_test.go
package certissuer_test

import (
	"context"
	"testing"

	"github.com/jdfalk/ubuntu-autoinstall-webhook/internal/certissuer"
)

func TestIssueCertificate(t *testing.T) {
	svc := certissuer.NewService()
	csr := []byte("dummy csr")
	clientInfo := map[string]string{"client": "test"}
	cert, err := svc.IssueCertificate(context.Background(), csr, clientInfo)
	if err != nil {
		t.Fatalf("IssueCertificate returned error: %v", err)
	}
	if len(cert) == 0 {
		t.Errorf("Expected a non-empty certificate")
	}
}

func TestRenewCertificate(t *testing.T) {
	svc := certissuer.NewService()
	cert := []byte("dummy certificate")
	newCert, err := svc.RenewCertificate(context.Background(), cert)
	if err != nil {
		t.Fatalf("RenewCertificate returned error: %v", err)
	}
	if len(newCert) == 0 {
		t.Errorf("Expected a non-empty renewed certificate")
	}
}

func TestGetRootCA(t *testing.T) {
	svc := certissuer.NewService()
	rootCA, err := svc.GetRootCA(context.Background())
	if err != nil {
		t.Fatalf("GetRootCA returned error: %v", err)
	}
	if len(rootCA) == 0 {
		t.Errorf("Expected a non-empty root CA certificate")
	}
}
