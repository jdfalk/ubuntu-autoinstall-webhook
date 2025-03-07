// internal/fileeditor/service_test.go
package fileeditor_test

import (
	"testing"

	"github.com/jdfalk/ubuntu-autoinstall-webhook/internal/fileeditor"
)

func TestWriteIpxeFile(t *testing.T) {
	// Create a new instance of the FileEditor service.
	fe := fileeditor.NewService()

	// Call WriteIpxeFile with dummy data.
	err := fe.WriteIpxeFile("AA:BB:CC:DD:EE:FF", []byte("#!ipxe\necho Booting..."))
	if err != nil {
		t.Errorf("WriteIpxeFile returned an error: %v", err)
	}
}
