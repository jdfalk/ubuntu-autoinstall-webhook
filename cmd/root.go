package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Global variables for flags
var configFile string
var logDir string

// Root command
var rootCmd = &cobra.Command{
	Use:   "webhook",
	Short: "Ubuntu Autoinstall Webhook CLI",
	Long:  "A webhook service for capturing Ubuntu Autoinstall reports",
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

// Init function to set up global flags
func init() {
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "Path to the config file")
	rootCmd.PersistentFlags().StringVar(&logDir, "logDir", "", "Directory for log storage")

	// Ensure config is loaded before executing any command
	cobra.OnInitialize(func() { initConfig(OsFs{}) })
}

// Load configuration
func initConfig(fs FileSystem) {
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		viper.AddConfigPath("/etc/webhook/")
	}

	// Read logDir from CLI flag if set
	if logDir != "" {
		viper.Set("logDir", logDir)
	}

	// Set default values
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

	// Read config file if available
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	// Ensure the config file has the correct options
	ensureConfigFile()

	// Validate paths
	if err := validatePaths(fs); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Enable ENV variables (WEBHOOK_PORT, WEBHOOK_LOGDIR, etc.)
	viper.AutomaticEnv()
}

// ensureConfigFile ensures the config file has the correct options by appending missing configuration lines.
// It attempts to write to multiple candidate locations if needed.
func ensureConfigFile() {
	configPath := viper.ConfigFileUsed()
	if configPath == "" {
		configPath = "config.yaml"
	}

	// Read existing config (if any)
	existingConfig, err := os.ReadFile(configPath)
	if err != nil {
		// if the file doesn't exist or cannot be read, assume empty content
		existingConfig = []byte("")
	}

	// Build a set of keys that are present—either defined or commented out.
	configSet := make(map[string]bool)
	for _, line := range strings.Split(string(existingConfig), "\n") {
		trimmed := strings.TrimSpace(line)
		// Remove a leading '#' if present.
		if strings.HasPrefix(trimmed, "#") {
			trimmed = strings.TrimPrefix(trimmed, "#")
			trimmed = strings.TrimSpace(trimmed)
		}
		if strings.Contains(trimmed, ":") {
			parts := strings.SplitN(trimmed, ":", 2)
			key := strings.TrimSpace(parts[0])
			// For nested keys (child keys), we only store the top-level if appropriate.
			// For example, a line "database:" or "database: something" should register "database".
			if !strings.Contains(key, " ") {
				configSet[key] = true
			}
		}
	}

	// Helper function to check if a key exists in the config.
	exists := func(key string) bool {
		_, ok := configSet[key]
		return ok
	}

	// Build the list of missing entries.
	var newEntries []string

	// Add header only if at least one missing entry will be appended.
	header := "# Missing configuration options (added automatically)"

	// Check top-level options
	if !exists("port") {
		newEntries = append(newEntries, "# port: 25000")
	}
	if !exists("logDir") {
		newEntries = append(newEntries, "# logDir: \"/opt/custom-logs\"")
	}
	if !exists("logFile") {
		newEntries = append(newEntries, "# logFile: \"autoinstall_report.log\"")
	}

	// For database settings, if the parent "database" is missing, then add the whole block.
	if !exists("database") {
		newEntries = append(newEntries, "# database:")
		newEntries = append(newEntries, "#   host: \"cockroachdb\"")
		newEntries = append(newEntries, "#   port: 26257")
		newEntries = append(newEntries, "#   user: \"admin\"")
		newEntries = append(newEntries, "#   password: \"securepassword\"")
		newEntries = append(newEntries, "#   dbname: \"autoinstall\"")
		newEntries = append(newEntries, "#   sslmode: \"disable\"")
		newEntries = append(newEntries, "#   max_open_conns: 100")
		newEntries = append(newEntries, "#   max_idle_conns: 10")
		newEntries = append(newEntries, "#   conn_max_lifetime: 3600")
	}

	// Check remaining options
	if !exists("ipxe_folder") {
		newEntries = append(newEntries, "# ipxe_folder: \"/var/www/html/ipxe\"")
	}
	if !exists("boot_customization_folder") {
		newEntries = append(newEntries, "# boot_customization_folder: \"/var/www/html/ipxe/boot\"")
	}
	if !exists("cloud_init_folder") {
		newEntries = append(newEntries, "# cloud_init_folder: \"/var/www/html/cloud-init/\"")
	}

	if len(newEntries) == 0 {
		// Nothing to add.
		return
	}

	// Prepend header
	newEntries = append([]string{header}, newEntries...)

	// Define candidate candidate paths for writing the updated config file.
	candidates := []string{
		configPath,
		"./config.yaml",
		"/tmp/custom-logs/config.yaml",
	}

	var writeErr error = fmt.Errorf("no candidate succeeded")
	// Attempt to write the new entries using one of the candidate paths.
	for _, path := range candidates {
		if err := tryWriteConfig(path, newEntries); err == nil {
			writeErr = nil
			break
		} else {
			fmt.Printf("Failed to write config to %s: %v\n", path, err)
		}
	}
	if writeErr != nil {
		fmt.Println("Failed to update config file in all candidate locations:", writeErr)
	}

	// Clean up duplicate blocks in the config file
	if err := cleanUpConfigFile(configPath); err != nil {
		fmt.Println("Failed to clean up config file:", err)
	}
}

// tryWriteConfig attempts to append entries to the file at candidate path.
// It creates the file and its directory if they do not exist.
func tryWriteConfig(path string, entries []string) error {
	// Determine the directory (assume file name "config.yaml")
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	// Append entries with proper newlines.
	_, err = f.WriteString("\n" + strings.Join(entries, "\n") + "\n")
	return err
}

// validatePaths validates the paths for logDir and other directories
func validatePaths(fs FileSystem) error {
	paths := []string{
		viper.GetString("logDir"),
		viper.GetString("ipxe_folder"),
		viper.GetString("boot_customization_folder"),
		viper.GetString("cloud_init_folder"),
	}

	for _, p := range paths {
		if p == "" {
			continue
		}

		// Sanitize the path to prevent path injection
		if strings.Contains(p, "..") || strings.Contains(p, "~") || strings.Contains(p, "//") {
			return fmt.Errorf("invalid path %s: contains illegal characters or sequences", p)
		}

		absPath, err := filepath.Abs(p)
		if err != nil {
			return fmt.Errorf("invalid path %s: %v", p, err)
		}

		info, err := fs.Stat(absPath)
		if os.IsNotExist(err) {
			err = fs.MkdirAll(absPath, 0755)
			if err != nil {
				return fmt.Errorf("failed to create directory %s: %v", absPath, err)
			}
		} else if err != nil {
			return fmt.Errorf("error accessing path %s: %v", absPath, err)
		} else if !info.IsDir() {
			return fmt.Errorf("path %s is not a directory", absPath)
		}
	}
	return nil
}

// cleanUpConfigFile reads the config file at the given path and removes duplicate
// "Missing configuration options (added automatically)" blocks. A block is defined as that header
// plus any following commented lines. Only the first occurrence is kept.
func cleanUpConfigFile(path string) error {
	// Read the current content.
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")
	var cleaned []string
	// Map to track blocks already seen.
	blockMap := make(map[string]bool)

	// Define the header string marker.
	headerMarker := "# Missing configuration options (added automatically)"

	// Iterate over all lines.
	for i := 0; i < len(lines); {
		line := strings.TrimSpace(lines[i])
		// If we see a header, assume it starts a block.
		if line == headerMarker {
			// Gather the entire block: header + subsequent comment lines (until a non-comment or blank line is reached).
			var blockLines []string
			blockLines = append(blockLines, lines[i])
			i++
			for i < len(lines) {
				trimmedNext := strings.TrimSpace(lines[i])
				// Stop block if we hit a blank line or a new section header.
				if trimmedNext == "" || (!strings.HasPrefix(trimmedNext, "#") && trimmedNext != headerMarker) {
					break
				}
				// If the very next header appears (i.e. same headerMarker), break as well.
				if trimmedNext == headerMarker {
					break
				}
				blockLines = append(blockLines, lines[i])
				i++
			}
			// Combine block lines.
			blockText := strings.Join(blockLines, "\n")
			// Only keep the block if it hasn't been seen before.
			if !blockMap[blockText] {
				blockMap[blockText] = true
				cleaned = append(cleaned, blockLines...)
			} else {
				// Skip this duplicate block.
				// (Do nothing – block not appended.)
			}
			// Optionally append a blank line after the block.
			cleaned = append(cleaned, "")
		} else {
			// For non-header lines, always keep them.
			cleaned = append(cleaned, lines[i])
			i++
		}
	}

	newContent := strings.Join(cleaned, "\n")
	// Only write if there are changes.
	if newContent != string(content) {
		return os.WriteFile(path, []byte(newContent), 0644)
	}
	return nil
}
