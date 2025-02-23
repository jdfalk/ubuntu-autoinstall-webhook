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

	// Determine which keys are already set
	configSet := make(map[string]bool)
	for _, line := range strings.Split(string(existingConfig), "\n") {
		if strings.Contains(line, ":") {
			key := strings.Split(line, ":")[0]
			configSet[strings.TrimSpace(key)] = true
		}
	}

	// Define the missing configuration options to be appended.
	missingConfig := []string{
		"# Missing configuration options (added automatically)",
		"# port: 25000",
		"# logDir: \"/opt/custom-logs\"",
		"# logFile: \"autoinstall_report.log\"",
		"# database:",
		"#   host: \"cockroachdb\"",
		"#   port: 26257",
		"#   user: \"admin\"",
		"#   password: \"securepassword\"",
		"#   dbname: \"autoinstall\"",
		"#   sslmode: \"disable\"",
		"#   max_open_conns: 100",
		"#   max_idle_conns: 10",
		"#   conn_max_lifetime: 3600",
		"# ipxe_folder: \"/var/www/html/ipxe\"",
		"# boot_customization_folder: \"/var/www/html/ipxe/boot\"",
		"# cloud_init_folder: \"/var/www/html/cloud-init/\"",
	}

	// Build the list of new entries that need to be appended.
	newEntries := []string{}
	for _, line := range missingConfig {
		key := strings.Split(line, ":")[0]
		key = strings.TrimPrefix(key, "# ")
		if !configSet[key] {
			newEntries = append(newEntries, line)
		}
	}

	// If there are missing options, attempt to write them in one of several candidate locations.
	if len(newEntries) > 1 {
		candidates := []string{
			configPath,
			"./config.yaml",
			"/tmp/custom-logs/config.yaml",
		}

		writeErr := fmt.Errorf("no candidate succeeded")
		for _, path := range candidates {
			if err := tryWriteConfig(path, newEntries); err == nil {
				// Successfully wrote to one candidate, stop trying.
				writeErr = nil
				break
			} else {
				// If writing failed, log the error and try the next candidate.
				fmt.Printf("Failed to write config to %s: %v\n", path, err)
			}
		}
		if writeErr != nil {
			fmt.Println("Failed to update config file in all candidate locations:", writeErr)
		}
	}
}

// tryWriteConfig attempts to append newEntries to the file at candidate path.
// It creates the file if it doesn't exist.
func tryWriteConfig(path string, entries []string) error {
	// If the candidate directory does not exist, try to create it.
	dir := strings.TrimSuffix(path, "/config.yaml")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

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
