// cmd/file-editor.go
package cmd

import (
	"context"
	"fmt"

	"github.com/jdfalk/ubuntu-autoinstall-webhook/internal/fileeditor"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var fileEditorCmd = &cobra.Command{
	Use:   "file-editor",
	Short: "Starts the file-editor microservice",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Starting file-editor microservice...")

		// Example: Retrieve a configuration value
		baseDir := viper.GetString("file_editor.base_dir")
		fmt.Println("File editor base directory:", baseDir)

		// Create context with cancellation for graceful shutdown
		ctx, cancel := context.WithCancel(cmd.Context())
		defer cancel()

		// Instantiate the file-editor service
		feService := fileeditor.NewService()

		// Start the service (including background tasks)
		if err := feService.Start(ctx); err != nil {
			return fmt.Errorf("failed to start file editor service: %w", err)
		}

		// For example, call a method (stubbed for now)
		if err := feService.WriteIpxeFile(ctx, "AA:BB:CC:DD:EE:FF", []byte("#!ipxe\necho Booting...")); err != nil {
			return err
		}

		fmt.Println("File-editor microservice started successfully.")

		// Block until context is canceled (e.g., by SIGINT)
		<-ctx.Done()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(fileEditorCmd)
}
