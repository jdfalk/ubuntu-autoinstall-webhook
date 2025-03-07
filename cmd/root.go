// cmd/root.go
package cmd

import (
    "fmt"
    "os"

    "github.com/spf13/cobra"
    "github.com/spf13/viper"
)

var cfgFile string

// rootCmd is the base command when called without any subcommands.
var rootCmd = &cobra.Command{
    Use:   "ubuntu-autoinstall-webhook",
    Short: "A multi-microservice client for Ubuntu autoinstall",
    Long: `A client and server microservice application for automating Ubuntu autoinstall.
It supports multiple microservices communicating via gRPC, with configuration managed via Viper.`,
    Run: func(cmd *cobra.Command, args []string) {
        // When no subcommand is provided, show help.
        cmd.Help()
    },
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}

func init() {
    cobra.OnInitialize(initConfig)
    // Persistent flag for specifying the config file.
    rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is config.yaml)")
}

func initConfig() {
    if cfgFile != "" {
        // Use config file from the flag.
        viper.SetConfigFile(cfgFile)
    } else {
        // Search for config.yaml in the current directory.
        viper.AddConfigPath(".")
        viper.SetConfigName("config")
    }
    viper.AutomaticEnv() // Read in environment variables that match.

    if err := viper.ReadInConfig(); err == nil {
        fmt.Println("Using config file:", viper.ConfigFileUsed())
    } else {
        fmt.Println("Error reading config file:", err)
    }
}
