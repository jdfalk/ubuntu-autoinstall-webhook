package ipxe

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/afero"
	"github.com/spf13/viper"

	"github.com/jdfalk/ubuntu-autoinstall-webhook/internal/db"
	"github.com/jdfalk/ubuntu-autoinstall-webhook/internal/logger"
)

// Fs is the file system used for file operations.
// In production it uses the real OS file system; tests may override it.
var Fs afero.Fs = afero.NewOsFs()

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
	content, err := afero.ReadFile(Fs, ipxeFilePath)
	if err != nil {
		errMsg := fmt.Sprintf("failed to read iPXE file %s: %v", ipxeFilePath, err)
		logger.Errorf("%s", errMsg)
		return fmt.Errorf("%s", errMsg)
	}

	historyFile := fmt.Sprintf("%s.history.%d", ipxeFilePath, time.Now().Unix())
	err = afero.WriteFile(Fs, historyFile, content, 0644)
	if err != nil {
		errMsg := fmt.Sprintf("failed to store iPXE history to %s: %v", historyFile, err)
		logger.Errorf("%s", errMsg)
		return fmt.Errorf("%s", errMsg)
	}

	logger.Infof("Successfully stored iPXE history to %s", historyFile)
	return nil
}

// UpdateIPXEFile modifies an iPXE file based on the provided MAC address and phase.
// It writes new content to disk (including the phase in a comment) and saves the new configuration to the database.
func UpdateIPXEFile(macAddress string, phase string) error {
	logger.Infof("Starting update of iPXE file for MAC address: %s with phase: %s", macAddress, phase)
	ipxeFilePath, err := generateIPXEFilePath(macAddress)
	if err != nil {
		return err
	}

	// If the file exists, store its history before updating.
	if _, err := Fs.Stat(ipxeFilePath); err == nil {
		logger.Infof("iPXE file %s exists; storing history before update.", ipxeFilePath)
		if err := storeFileHistory(ipxeFilePath); err != nil {
			return err
		}
	} else {
		logger.Infof("iPXE file %s does not exist; it will be created.", ipxeFilePath)
	}

	// Create new content including the phase.
	newContent := fmt.Sprintf("#!ipxe\n# phase: %s\nexit\n", phase)
	err = afero.WriteFile(Fs, ipxeFilePath, []byte(newContent), 0644)
	if err != nil {
		errMsg := fmt.Sprintf("failed to update iPXE file %s: %v", ipxeFilePath, err)
		logger.Errorf("%s", errMsg)
		return fmt.Errorf("%s", errMsg)
	}

	logger.Infof("Successfully updated iPXE file %s", ipxeFilePath)

	// Save the new configuration to the database.
	if err := db.SaveIPXEConfiguration(macAddress, newContent, phase); err != nil {
		logger.Errorf("Failed to save iPXE configuration for %s: %v", macAddress, err)
		return err
	}

	return nil
}

// UpdateIPXEOnProgress updates the iPXE file when a client reaches a certain progress threshold.
// It determines the new phase ("install" or "post-install"), updates the configuration, and logs the latest config.
func UpdateIPXEOnProgress(clientID string, progress int, macAddress string) {
	var newPhase string
	if progress >= 25 {
		newPhase = "post-install"
	} else {
		newPhase = "install"
	}

	err := UpdateIPXEFile(macAddress, newPhase)
	if err != nil {
		logger.Errorf("Failed to update iPXE for %s: %v", macAddress, err)
		return
	}

	// Retrieve the latest IPXE configuration matching the phase.
	cfg, err := db.GetLatestIPXEConfig(macAddress, newPhase)
	if err != nil {
		logger.Errorf("Failed to retrieve latest IPXE config for %s with phase %s: %v", macAddress, newPhase, err)
	} else {
		logger.Infof("Latest IPXE configuration for MAC %s (phase %s): %+v", macAddress, newPhase, cfg)
	}
}
