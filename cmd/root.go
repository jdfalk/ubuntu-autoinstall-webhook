package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
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
	Key      string
	Comments []string
	Lines    []string // Includes the key line and any indented child lines.
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

// initConfig loads configuration from file/environment and then processes it.
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
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.user", "user")
	viper.SetDefault("database.password", "password")
	viper.SetDefault("database.dbname", "autoinstall")
	viper.SetDefault("database.sslmode", "disable")
	viper.SetDefault("database.max_open_conns", 100)
	viper.SetDefault("database.max_idle_conns", 10)
	viper.SetDefault("database.conn_max_lifetime", 3600)
	viper.SetDefault("ipxe_folder", "/var/www/html/ipxe")
	viper.SetDefault("boot_customization_folder", "/var/www/html/ipxe/boot")
	viper.SetDefault("cloud_init_folder", "/var/www/html/cloud-init/")

	// Read config file if available.
	if err := viper.ReadInConfig(); err == nil {
		logger.Infof("%s", "Using config file: "+viper.ConfigFileUsed())
	}

	// Process (ensure + organize) the configuration file.
	if err := processConfigFile(); err != nil {
		logger.Errorf("%s", fmt.Sprintf("Failed to process config file: %v", err))
	}

	// Validate paths with fallback logic.
	if err := validatePaths(fs); err != nil {
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

// processConfigFile merges ensuring missing config entries exist and organizing the file.
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
	}

	// Build a set of keys present in the config.
	configSet := make(map[string]bool)
	for _, line := range strings.Split(existingContent, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "#") {
			trimmed = strings.TrimPrefix(trimmed, "#")
			trimmed = strings.TrimSpace(trimmed)
		}
		if strings.Contains(trimmed, ":") {
			parts := strings.SplitN(trimmed, ":", 2)
			key := strings.TrimSpace(parts[0])
			if !strings.Contains(key, " ") {
				configSet[key] = true
			}
		}
	}

	exists := func(key string) bool {
		_, ok := configSet[key]
		return ok
	}

	var missingEntries []string
	header := "# Missing configuration options (added automatically)"
	if !exists("port") {
		missingEntries = append(missingEntries, "# port: 25000")
	}
	if !exists("logDir") {
		missingEntries = append(missingEntries, "# logDir: \"/opt/custom-logs\"")
	}
	if !exists("logFile") {
		missingEntries = append(missingEntries, "# logFile: \"autoinstall_report.log\"")
	}
	if !exists("database") {
		missingEntries = append(missingEntries, "# database:")
		missingEntries = append(missingEntries, "#   host: \"cockroachdb\"")
		missingEntries = append(missingEntries, "#   port: 26257")
		missingEntries = append(missingEntries, "#   user: \"admin\"")
		missingEntries = append(missingEntries, "#   password: \"securepassword\"")
		missingEntries = append(missingEntries, "#   dbname: \"autoinstall\"")
		missingEntries = append(missingEntries, "#   sslmode: \"disable\"")
		missingEntries = append(missingEntries, "#   max_open_conns: 100")
		missingEntries = append(missingEntries, "#   max_idle_conns: 10")
		missingEntries = append(missingEntries, "#   conn_max_lifetime: 3600")
	}
	if !exists("ipxe_folder") {
		missingEntries = append(missingEntries, "# ipxe_folder: \"/var/www/html/ipxe\"")
	}
	if !exists("boot_customization_folder") {
		missingEntries = append(missingEntries, "# boot_customization_folder: \"/var/www/html/ipxe/boot\"")
	}
	if !exists("cloud_init_folder") {
		missingEntries = append(missingEntries, "# cloud_init_folder: \"/var/www/html/cloud-init/\"")
	}

	if len(missingEntries) > 0 {
		missingEntries = append([]string{header}, missingEntries...)
	}
	combinedContent := existingContent + "\n" + strings.Join(missingEntries, "\n") + "\n"

	// Organize the config file using fixed sections.
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
				i++
				if key == "database" {
					for i < len(lines) {
						nextLine := lines[i]
						if nextLine != strings.TrimLeft(nextLine, " ") && strings.TrimSpace(nextLine) != "" {
							blk.Lines = append(blk.Lines, nextLine)
							i++
						} else {
							break
						}
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

	// Define sections. For each section, sort keys alphabetically.
	type Section struct {
		Header string
		Keys   []string
	}
	sections := []Section{
		{Header: "", Keys: []string{"port"}},
		{Header: "# Logging Configuration", Keys: []string{"logDir", "logFile"}},
		{Header: "# Database Configuration", Keys: []string{"database"}},
		{Header: "# iPXE Settings", Keys: []string{"boot_customization_folder", "ipxe_folder"}},
		{Header: "# Cloud-Init Settings", Keys: []string{"cloud_init_folder"}},
	}

	var outputLines []string
	for _, sec := range sections {
		// Sort keys alphabetically.
		sortedKeys := make([]string, len(sec.Keys))
		copy(sortedKeys, sec.Keys)
		sort.Strings(sortedKeys)

		// Check if the header is already present among the blocks for this section.
		headerAlreadyPresent := false
		for _, key := range sortedKeys {
			if blk, ok := blocks[key]; ok {
				if len(blk.Comments) > 0 && strings.TrimSpace(blk.Comments[0]) == sec.Header {
					headerAlreadyPresent = true
					break
				}
			}
		}
		// Add the header only if it isn't already present.
		if sec.Header != "" && !headerAlreadyPresent {
			outputLines = append(outputLines, sec.Header)
		}
		for _, key := range sortedKeys {
			if blk, ok := blocks[key]; ok {
				// Deduplicate block comments and avoid printing consecutive duplicates.
				deduped := deduplicate(blk.Comments)
				for _, comm := range deduped {
					if len(outputLines) == 0 || strings.TrimSpace(outputLines[len(outputLines)-1]) != strings.TrimSpace(comm) {
						outputLines = append(outputLines, comm)
					}
				}
				outputLines = append(outputLines, blk.Lines...)
			}
		}
		outputLines = append(outputLines, "")
	}

	// Remove any consecutive duplicate lines from the final output.
	finalLines := deduplicateConsecutive(outputLines)
	organizedContent := strings.Join(finalLines, "\n")

	// Write the organized config file using candidate locations.
	if err := writeFileToCandidateLocations(filepath.Base(configPath), []byte(organizedContent)); err != nil {
		return err
	}

	return nil
}

// writeFileToCandidateLocations attempts to write a file by trying four candidate locations.
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
			logger.Warningf("%s", "Failed to write file to "+candidatePath+": "+err.Error())
			continue
		}
		logger.Infof("%s", "Successfully wrote file to "+candidatePath)
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
			logger.Warningf("%s", "Cannot use folder "+path+" trying next candidate.")
			continue
		}
		return path, nil
	}
	return "", fmt.Errorf("no valid folder found")
}

// validatePaths uses fallback logic for critical directories.
// It also checks for invalid characters or sequences in the original path.
func validatePaths(fs FileSystem) error {
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
