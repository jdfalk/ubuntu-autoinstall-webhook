// Package cmd provides the main command-line interface for the Ubuntu Autoinstall Webhook service.
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jdfalk/ubuntu-autoinstall-webhook/internal/logger" // Import the centralized logger
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Global variables for flags.
var configFile string
var logDir string

// ConfigBlock represents a configuration block along with any preceding comments.
type ConfigBlock struct {
	// Key is the configuration key.
	Key string
	// Comments are any preceding comment lines.
	Comments []string
	// Lines include the key line and any indented child lines.
	Lines []string
}

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "webhook",
	Short: "Ubuntu Autoinstall Webhook CLI",
	Long:  "A webhook service for capturing Ubuntu Autoinstall reports",
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logger.Errorf("%s", err.Error())
		os.Exit(1)
	}
}

// init sets up global flags and initializes configuration.
func init() {
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "Path to the config file")
	rootCmd.PersistentFlags().StringVar(&logDir, "logDir", "", "Directory for log storage")

	// OsFs is implemented in filesystem.go.
	cobra.OnInitialize(func() { initConfig(OsFs{}) })
}

// keyExistsInContent checks if a key (with colon) exists anywhere in the provided content.
// It returns true if the key followed by a colon is found.
func keyExistsInContent(content, key string) bool {
	return strings.Contains(content, key+":")
}

// initConfig loads configuration from file/environment, processes the config file,
// and validates critical paths. It also sets default values including a new option for log level.
func initConfig(fs FileSystem) {
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		viper.AddConfigPath("/etc/webhook/")
	}

	// Override logDir if set via CLI flag.
	if logDir != "" {
		viper.Set("logDir", logDir)
	}

	// Set default values.
	viper.SetDefault("port", "5000")
	viper.SetDefault("logDir", "/var/log/autoinstall-webhook")
	viper.SetDefault("logFile", "autoinstall_report.log")
	// New viper default for log level output.
	viper.SetDefault("logLevel", "INFO")
	// Database defaults – note that now the default for "enabled" is false.
	viper.SetDefault("database.enabled", false)
	viper.SetDefault("database.type", "postgres")
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.user", "user")
	viper.SetDefault("database.password", "password")
	viper.SetDefault("database.dbname", "autoinstall")
	viper.SetDefault("database.sslmode", "require")
	viper.SetDefault("database.max_open_conns", 100)
	viper.SetDefault("database.max_idle_conns", 10)
	viper.SetDefault("database.conn_max_lifetime", 3600)
	viper.SetDefault("ipxe_folder", "/var/www/html/ipxe")
	viper.SetDefault("boot_customization_folder", "/var/www/html/ipxe/boot")
	viper.SetDefault("cloud_init_folder", "/var/www/html/cloud-init/")

	// Read config file if available.
	if err := viper.ReadInConfig(); err == nil {
		logger.Infof("Using config file: %s", viper.ConfigFileUsed())
	} else {
		logger.Debugf("No config file found or error reading config: %v", err)
	}

	// Process (ensure missing options are added and file is organized).
	if err := processConfigFile(); err != nil {
		logger.Errorf("Failed to process config file: %v", err)
	}

	// Re-read the config file so that Viper picks up the new defaults and organization.
	if err := viper.ReadInConfig(); err != nil {
		logger.Warningf("Error re-reading config file: %v", err)
	} else {
		logger.Infof("Re-loaded config file with updated settings")
	}

	// DEBUG: Log all settings and the database.enabled value.
	logger.Infof("All config settings: %#v", viper.AllSettings())
	logger.Infof("Raw database.enabled: %#v", viper.Get("database.enabled"))
	logger.Infof("Database enabled (bool): %v", viper.GetBool("database.enabled"))

	// Validate paths with fallback logic.
	if err := validatePaths(); err != nil {
		logger.Errorf("%s", err.Error())
		os.Exit(1)
	}

	// Enable ENV variables (e.g. WEBHOOK_PORT, WEBHOOK_LOGDIR, etc.).
	viper.AutomaticEnv()
}

// deduplicate removes duplicate lines from a slice while preserving order.
func deduplicate(lines []string) []string {
	seen := make(map[string]struct{})
	var result []string
	for _, line := range lines {
		if _, exists := seen[line]; !exists {
			seen[line] = struct{}{}
			result = append(result, line)
		}
	}
	return result
}

// deduplicateConsecutive removes consecutive duplicate lines.
func deduplicateConsecutive(lines []string) []string {
	if len(lines) == 0 {
		return lines
	}
	result := []string{lines[0]}
	for i := 1; i < len(lines); i++ {
		if strings.TrimSpace(lines[i]) != strings.TrimSpace(lines[i-1]) {
			result = append(result, lines[i])
		}
	}
	return result
}

// processConfigFile reads the existing config file (or creates a new one)
// and appends any missing configuration options with default values, then
// organizes the file into fixed sections.
func processConfigFile() error {
	configPath := viper.ConfigFileUsed()
	if configPath == "" {
		configPath = "config.yaml"
	}

	// Read existing config.
	contentBytes, err := os.ReadFile(configPath)
	existingContent := ""
	if err == nil {
		existingContent = string(contentBytes)
	} else {
		logger.Debugf("Error reading config file, will create a new one: %v", err)
	}

	var missingEntries []string
	header := "# Missing configuration options (added automatically)"
	// Check for missing flat keys.
	if !keyExistsInContent(existingContent, "port") {
		missingEntries = append(missingEntries, "# port: 25000")
	}
	if !keyExistsInContent(existingContent, "logDir") {
		missingEntries = append(missingEntries, "# logDir: \"/opt/custom-logs\"")
	}
	if !keyExistsInContent(existingContent, "logFile") {
		missingEntries = append(missingEntries, "# logFile: \"autoinstall_report.log\"")
	}
	// For database, check if the block exists; if not, add a nested block.
	if !strings.Contains(existingContent, "database:") {
		missingEntries = append(missingEntries, "# Database Configuration")
		missingEntries = append(missingEntries, "database:")
		missingEntries = append(missingEntries, "  enabled: false")
		missingEntries = append(missingEntries, "  type: postgres")
		missingEntries = append(missingEntries, "  host: localhost")
		missingEntries = append(missingEntries, "  port: 5432")
		missingEntries = append(missingEntries, "  user: user")
		missingEntries = append(missingEntries, "  password: password")
		missingEntries = append(missingEntries, "  dbname: autoinstall")
		missingEntries = append(missingEntries, "  sslmode: enabled")
		missingEntries = append(missingEntries, "  max_open_conns: 100")
		missingEntries = append(missingEntries, "  max_idle_conns: 10")
		missingEntries = append(missingEntries, "  conn_max_lifetime: 3600")
	}
	// Other keys.
	if !keyExistsInContent(existingContent, "ipxe_folder") {
		missingEntries = append(missingEntries, "# ipxe_folder: \"/var/www/html/ipxe\"")
	}
	if !keyExistsInContent(existingContent, "boot_customization_folder") {
		missingEntries = append(missingEntries, "# boot_customization_folder: \"/var/www/html/ipxe/boot\"")
	}
	if !keyExistsInContent(existingContent, "cloud_init_folder") {
		missingEntries = append(missingEntries, "# cloud_init_folder: \"/var/www/html/cloud-init/\"")
	}

	var combinedContent string
	if len(missingEntries) > 0 {
		missingEntries = append([]string{header}, missingEntries...)
		combinedContent = existingContent + "\n" + strings.Join(missingEntries, "\n") + "\n"
	} else {
		combinedContent = existingContent
	}

	// Organize the config file using fixed sections.
	// Treat "database" as one block.
	knownKeys := map[string]bool{
		"port":                      true,
		"logDir":                    true,
		"logFile":                   true,
		"database":                  true,
		"ipxe_folder":               true,
		"boot_customization_folder": true,
		"cloud_init_folder":         true,
	}
	blocks := make(map[string]ConfigBlock)
	var pendingComments []string
	missingHeader := "# Missing configuration options (added automatically)"
	seenMissingBlock := false

	lines := strings.Split(combinedContent, "\n")
	i := 0
	for i < len(lines) {
		line := lines[i]
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			pendingComments = nil
			i++
			continue
		}
		if strings.HasPrefix(trimmed, "#") {
			if trimmed == missingHeader {
				if seenMissingBlock {
					i++
					for i < len(lines) {
						if strings.TrimSpace(lines[i]) == "" || !strings.HasPrefix(strings.TrimSpace(lines[i]), "#") {
							break
						}
						i++
					}
					pendingComments = nil
					continue
				} else {
					seenMissingBlock = true
				}
			}
			pendingComments = append(pendingComments, line)
			i++
			continue
		}
		// Handle the database block as a single unit.
		if strings.HasPrefix(trimmed, "database:") {
			blk := ConfigBlock{
				Key:      "database",
				Comments: pendingComments,
				Lines:    []string{line},
			}
			pendingComments = nil
			i++
			// Capture following indented lines as part of the database block.
			for i < len(lines) {
				nextLine := lines[i]
				if nextLine != strings.TrimLeft(nextLine, " ") && strings.TrimSpace(nextLine) != "" {
					blk.Lines = append(blk.Lines, nextLine)
					i++
				} else {
					break
				}
			}
			blocks["database"] = blk
			continue
		}
		// Process flat keys.
		if line == strings.TrimLeft(line, " ") && strings.Contains(trimmed, ":") {
			parts := strings.SplitN(trimmed, ":", 2)
			key := strings.TrimSpace(parts[0])
			if knownKeys[key] {
				dedupedComments := deduplicate(pendingComments)
				blk := ConfigBlock{
					Key:      key,
					Comments: dedupedComments,
					Lines:    []string{line},
				}
				pendingComments = nil
				// Capture subsequent indented lines.
				for i++; i < len(lines); i++ {
					nextLine := lines[i]
					if nextLine != strings.TrimLeft(nextLine, " ") && strings.TrimSpace(nextLine) != "" {
						blk.Lines = append(blk.Lines, nextLine)
					} else {
						break
					}
				}
				if _, exists := blocks[key]; !exists {
					blocks[key] = blk
				}
				continue
			}
		}
		pendingComments = nil
		i++
	}

	// Define sections.
	type Section struct {
		Header string
		Keys   []string
	}
	sections := []Section{
		{Header: "", Keys: []string{"port"}},
		{Header: "# Logging Configuration", Keys: []string{"logDir", "logFile"}},
		{Header: "# Database Configuration", Keys: []string{"database"}},
		{Header: "# iPXE Settings", Keys: []string{"ipxe_folder", "boot_customization_folder"}},
		{Header: "# Cloud-Init Settings", Keys: []string{"cloud_init_folder"}},
	}

	var outputLines []string
	for _, sec := range sections {
		if sec.Header != "" {
			outputLines = append(outputLines, sec.Header)
		}
		for _, key := range sec.Keys {
			if blk, ok := blocks[key]; ok {
				for _, comm := range deduplicate(blk.Comments) {
					if len(outputLines) == 0 || strings.TrimSpace(outputLines[len(outputLines)-1]) != strings.TrimSpace(comm) {
						outputLines = append(outputLines, comm)
					}
				}
				// Replace loop with a single append operation as suggested by S1011.
				outputLines = append(outputLines, blk.Lines...)
			}
		}
		outputLines = append(outputLines, "")
	}

	// Remove any consecutive duplicate lines.
	finalLines := deduplicateConsecutive(outputLines)
	organizedContent := strings.Join(finalLines, "\n")

	// Write the organized config file using candidate locations.
	if err := writeFileToCandidateLocations(filepath.Base(configPath), []byte(organizedContent)); err != nil {
		return err
	}

	return nil
}

// writeFileToCandidateLocations attempts to write a file by trying multiple candidate locations.
func writeFileToCandidateLocations(fileToCheck string, data []byte) error {
	var candidates []string

	// Candidate 1: directory of the config file (if available).
	if configPath := viper.ConfigFileUsed(); configPath != "" {
		configDir := filepath.Dir(configPath)
		candidates = append(candidates, filepath.Join(configDir, fileToCheck))
	}
	// Candidate 2: current directory under "ubuntu-autoinstall-webhook".
	candidates = append(candidates, filepath.Join(".", "ubuntu-autoinstall-webhook", fileToCheck))
	// Candidate 3: user's home directory under ubuntu-autoinstall-webhook.
	if homeDir, err := os.UserHomeDir(); err == nil {
		candidates = append(candidates, filepath.Join(homeDir, "ubuntu-autoinstall-webhook", fileToCheck))
	}
	// Candidate 4: temporary directory under ubuntu-autoinstall-webhook.
	candidates = append(candidates, filepath.Join(os.TempDir(), "ubuntu-autoinstall-webhook", fileToCheck))

	var lastErr error
	for _, candidatePath := range candidates {
		if err := os.MkdirAll(filepath.Dir(candidatePath), 0755); err != nil {
			lastErr = err
			continue
		}
		if err := os.WriteFile(candidatePath, data, 0644); err != nil {
			lastErr = err
			logger.Warningf("Failed to write file to %s: %v", candidatePath, err)
			continue
		}
		logger.Infof("Successfully wrote file to %s", candidatePath)
		return nil
	}
	return fmt.Errorf("failed to write file to all candidate locations: last error: %v", lastErr)
}

// ensureFolderExists checks if a folder exists and is writable by creating a temporary file.
func ensureFolderExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(path, 0755); err != nil {
			return false
		}
	}
	// Try writing a temporary file.
	testFile := filepath.Join(path, ".tmp")
	if err := os.WriteFile(testFile, []byte{}, 0644); err != nil {
		return false
	}
	os.Remove(testFile)
	return true
}

// getAvailableFolder attempts to use the provided folder and a set of alternatives.
func getAvailableFolder(defaultPath string, alternatives ...string) (string, error) {
	paths := append([]string{defaultPath}, alternatives...)
	for _, path := range paths {
		if !ensureFolderExists(path) {
			logger.Warningf("Cannot use folder %s trying next candidate.", path)
			continue
		}
		return path, nil
	}
	return "", fmt.Errorf("no valid folder found")
}

// validatePaths uses fallback logic for critical directories.
// It also checks for invalid characters or sequences in the original path.
func validatePaths() error {
	keys := []string{"logDir", "ipxe_folder", "boot_customization_folder", "cloud_init_folder"}
	for _, key := range keys {
		originalPath := viper.GetString(key)
		// Check for illegal sequences in the original path.
		if strings.Contains(originalPath, "..") || strings.Contains(originalPath, "~") || strings.Contains(originalPath, "//") {
			return fmt.Errorf("invalid path %s for key %s: contains illegal characters or sequences", originalPath, key)
		}

		base := filepath.Base(originalPath)
		// Candidate 2: current working directory.
		cand2 := filepath.Join(".", "ubuntu-autoinstall-webhook", base)
		// Candidate 3: user's home directory.
		home, err := os.UserHomeDir()
		if err != nil {
			home = "."
		}
		cand3 := filepath.Join(home, "ubuntu-autoinstall-webhook", base)
		// Candidate 4: temporary directory.
		cand4 := filepath.Join(os.TempDir(), "ubuntu-autoinstall-webhook", base)

		available, err := getAvailableFolder(originalPath, cand2, cand3, cand4)
		if err != nil {
			return fmt.Errorf("failed to validate path for %s: %v", key, err)
		}
		viper.Set(key, available)
	}
	return nil
}
