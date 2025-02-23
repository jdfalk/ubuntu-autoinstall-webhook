package ipxe

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
)

// IPXEHistory stores the last five IPXE configuration versions per client.
type IPXEHistory struct {
	ID        string    `json:"id"`
	ClientID  string    `json:"client_id"`
	Config    string    `json:"config"`
	CreatedAt time.Time `json:"created_at"`
}

// storeFileHistory saves the current version of the iPXE file before modification.
func storeFileHistory(macAddress, ipxeFilePath string) error {
	content, err := ioutil.ReadFile(ipxeFilePath)
	if err != nil {
		return fmt.Errorf("failed to read iPXE file: %w", err)
	}

	historyFile := fmt.Sprintf("%s.history.%d", ipxeFilePath, time.Now().Unix())
	err = ioutil.WriteFile(historyFile, content, 0644)
	if err != nil {
		return fmt.Errorf("failed to store iPXE history: %w", err)
	}

	return nil
}

// UpdateIPXEFile modifies an iPXE file when a client reaches 25% progress
func UpdateIPXEFile(macAddress string) error {
	ipxeFolder := viper.GetString("ipxe_folder")
	bootCustomizationFolder := viper.GetString("boot_customization_folder")

	if ipxeFolder == "" || bootCustomizationFolder == "" {
		return fmt.Errorf("iPXE folder paths are not set in configuration")
	}

	// Construct file paths based on MAC address
	macFile := filepath.Join(bootCustomizationFolder, fmt.Sprintf("mac-%s.ipxe", macAddress))
	defaultFile := filepath.Join(ipxeFolder, fmt.Sprintf("%s.ipxe", macAddress))

	// Determine which file to update
	var ipxeFilePath string
	if _, err := os.Stat(macFile); err == nil {
		ipxeFilePath = macFile
	} else if _, err := os.Stat(defaultFile); err == nil {
		ipxeFilePath = defaultFile
	} else {
		return fmt.Errorf("no iPXE file found for MAC: %s", macAddress)
	}

	// Store file history before modification
	err := storeFileHistory(macAddress, ipxeFilePath)
	if err != nil {
		return err
	}

	// Update the iPXE file to instruct normal boot
	newContent := "#!ipxe\nexit\n"
	err = os.WriteFile(ipxeFilePath, []byte(newContent), 0644)
	if err != nil {
		return fmt.Errorf("failed to update iPXE file: %w", err)
	}

	return nil
}
