package cmd

import (
	"fmt"
	"os"

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

	// Read config file if available
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	// Enable ENV variables (WEBHOOK_PORT, WEBHOOK_LOGDIR, etc.)
	viper.AutomaticEnv()
}
