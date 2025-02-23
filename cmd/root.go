package cmd

import (
	"fmt"
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
		fmt.Println(err)
		os.Exit(1)
	}
}

// Init function to set up global flags
func init() {
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "Path to the config file")
	rootCmd.PersistentFlags().StringVar(&logDir, "logDir", "", "Directory for log storage")

	// Ensure config is loaded before executing any command
	cobra.OnInitialize(initConfig)
}

// Load configuration
func initConfig() {
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
	if err := validatePaths(OsFs{}); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Enable ENV variables (WEBHOOK_PORT, WEBHOOK_LOGDIR, etc.)
	viper.AutomaticEnv()
}

// ensureConfigFile ensures the config file has the correct options
func ensureConfigFile() {
	configPath := viper.ConfigFileUsed()
	if configPath == "" {
		configPath = "config.yaml"
	}

	existingConfig, err := os.ReadFile(configPath)
	if err != nil {
		return
	}

	existingLines := strings.Split(string(existingConfig), "\n")
	configSet := make(map[string]bool)
	for _, line := range existingLines {
		if strings.Contains(line, ":") {
			key := strings.Split(line, ":")[0]
			configSet[strings.TrimSpace(key)] = true
		}
	}

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

	newEntries := []string{}
	for _, line := range missingConfig {
		key := strings.Split(line, ":")[0]
		key = strings.TrimPrefix(key, "# ")
		if !configSet[key] {
			newEntries = append(newEntries, line)
		}
	}

	if len(newEntries) > 1 {
		f, err := os.OpenFile(configPath, os.O_APPEND|os.O_WRONLY, 0644)
		if err == nil {
			defer f.Close()
			f.WriteString("\n" + strings.Join(newEntries, "\n") + "\n")
		}
	}
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
			return fmt.Errorf("Invalid path %s: contains illegal characters or sequences", p)
		}

		absPath, err := filepath.Abs(p)
		if err != nil {
			return fmt.Errorf("Invalid path %s: %v", p, err)
		}

		info, err := fs.Stat(absPath)
		if os.IsNotExist(err) {
			err = fs.MkdirAll(absPath, 0755)
			if err != nil {
				return fmt.Errorf("Failed to create directory %s: %v", absPath, err)
			}
		} else if err != nil {
			return fmt.Errorf("Error accessing path %s: %v", absPath, err)
		} else if !info.IsDir() {
			return fmt.Errorf("Path %s is not a directory", absPath)
		}
	}
	return nil
}
