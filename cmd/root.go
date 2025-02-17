package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

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

// Init config (loads environment variables and config file)
func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	viper.SetConfigName("config") // Config file name (without extension)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/webhook/") // Look in system-wide config dir

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
