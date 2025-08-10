// cmd/configuration.go
package cmd

import (
	"fmt"

	"github.com/jdfalk/ubuntu-autoinstall-webhook/internal/configuration"
	"github.com/spf13/cobra"
)

var configurationCmd = &cobra.Command{
	Use:   "configuration",
	Short: "Starts the configuration microservice",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Starting configuration microservice...")
		configService := configuration.NewService()
		cfg, err := configService.LoadConfig()
		if err != nil {
			return err
		}
		fmt.Printf("Loaded configuration: %+v\n", cfg)
		fmt.Println("Configuration microservice started successfully.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(configurationCmd)
}
