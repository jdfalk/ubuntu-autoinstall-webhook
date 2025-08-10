// cmd/database.go
package cmd

import (
	"context"
	"fmt"

	"github.com/jdfalk/ubuntu-autoinstall-webhook/internal/database"
	"github.com/spf13/cobra"
)

var databaseCmd = &cobra.Command{
	Use:   "database",
	Short: "Starts the database microservice",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Starting database microservice...")
		dbService := database.NewService()
		// Pass a context to Connect
		if err := dbService.Connect(context.Background()); err != nil {
			return err
		}
		if err := dbService.MigrateSchema(context.Background()); err != nil {
			return err
		}
		fmt.Println("Database microservice started successfully.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(databaseCmd)
}
