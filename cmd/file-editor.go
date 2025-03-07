// cmd/file-editor.go
package cmd

import (
    "fmt"

    "github.com/spf13/cobra"
    "github.com/spf13/viper"
    "github.com/jdfalk/ubuntu-autoinstall-webhook/internal/fileeditor"
)

var fileEditorCmd = &cobra.Command{
    Use:   "file-editor",
    Short: "Starts the file-editor microservice",
    RunE: func(cmd *cobra.Command, args []string) error {
        fmt.Println("Starting file-editor microservice...")
        // Example: Retrieve a configuration value
        baseDir := viper.GetString("file_editor.base_dir")
        fmt.Println("File editor base directory:", baseDir)

        // Instantiate the file-editor service.
        feService := fileeditor.NewService()
        // For example, call a method (stubbed for now).
        if err := feService.WriteIpxeFile("AA:BB:CC:DD:EE:FF", []byte("#!ipxe\necho Booting...")); err != nil {
            return err
        }
        fmt.Println("File-editor microservice started successfully.")
        return nil
    },
}

func init() {
    rootCmd.AddCommand(fileEditorCmd)
}
