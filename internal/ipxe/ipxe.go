package ipxe

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/viper"

	"github.com/jdfalk/ubuntu-autoinstall-webhook/internal/logger"
)

// IPXEHistory stores the last five IPXE configuration versions per client.
type IPXEHistory struct {
	ID        string    `json:"id"`
	ClientID  string    `json:"client_id"`
	Config    string    `json:"config"`
	CreatedAt time.Time `json:"created_at"`
}

// generateIPXEFilePath constructs the iPXE file path based solely on the MAC address.
// It uses boot_customization_folder if provided; otherwise, it falls back to ipxe_folder.
// The MAC address is normalized to lowercase without colons.
func generateIPXEFilePath(macAddress string) (string, error) {
	baseDir := viper.GetString("boot_customization_folder")
	if baseDir == "" {
		baseDir = viper.GetString("ipxe_folder")
		if baseDir == "" {
			err := fmt.Errorf("no iPXE folder configured")
			logger.Errorf("%s", err.Error())
			return "", err
		}
	}
	normalizedMac := strings.ToLower(macAddress)
	normalizedMac = strings.ReplaceAll(normalizedMac, ":", "")
	fileName := fmt.Sprintf("mac-%s.ipxe", normalizedMac)
	ipxePath := filepath.Join(baseDir, fileName)
	logger.Infof("Generated iPXE file path: %s", ipxePath)
	return ipxePath, nil
}

// storeFileHistory saves the current version of the iPXE file before modification.
func storeFileHistory(ipxeFilePath string) error {
	logger.Infof("Storing history for iPXE file: %s", ipxeFilePath)
	// Read the current file contents using os.ReadFile.
	content, err := os.ReadFile(ipxeFilePath)
	if err != nil {
		errMsg := fmt.Sprintf("failed to read iPXE file %s: %v", ipxeFilePath, err)
		logger.Errorf("%s", errMsg)
		return fmt.Errorf("%s", errMsg)
	}

	historyFile := fmt.Sprintf("%s.history.%d", ipxeFilePath, time.Now().Unix())
	// Write the file contents to the history file using os.WriteFile.
	err = os.WriteFile(historyFile, content, 0644)
	if err != nil {
		errMsg := fmt.Sprintf("failed to store iPXE history to %s: %v", historyFile, err)
		logger.Errorf("%s", errMsg)
		return fmt.Errorf("%s", errMsg)
	}

	logger.Infof("Successfully stored iPXE history to %s", historyFile)
	return nil
}

// UpdateIPXEFile modifies an iPXE file when a client reaches 25% progress.
// It generates the file path based solely on the provided MAC address.
func UpdateIPXEFile(macAddress string) error {
	logger.Infof("Starting update of iPXE file for MAC address: %s", macAddress)
	ipxeFilePath, err := generateIPXEFilePath(macAddress)
	if err != nil {
		return err
	}

	// If the file exists, store its history before updating.
	if _, err := os.Stat(ipxeFilePath); err == nil {
		logger.Infof("iPXE file %s exists; storing history before update.", ipxeFilePath)
		if err := storeFileHistory(ipxeFilePath); err != nil {
			return err
		}
	} else {
		logger.Infof("iPXE file %s does not exist; it will be created.", ipxeFilePath)
	}

	// Update the iPXE file to instruct normal boot.
	newContent := "#!ipxe\nexit\n"
	err = os.WriteFile(ipxeFilePath, []byte(newContent), 0644)
	if err != nil {
		errMsg := fmt.Sprintf("failed to update iPXE file %s: %v", ipxeFilePath, err)
		logger.Errorf("%s", errMsg)
		return fmt.Errorf("%s", errMsg)
	}

	logger.Infof("Successfully updated iPXE file %s", ipxeFilePath)
	return nil
}
