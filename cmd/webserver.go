// cmd/webserver.go
package cmd

import (
	"fmt"

	"github.com/jdfalk/ubuntu-autoinstall-webhook/internal/webserver"
	"github.com/spf13/cobra"
)

var webserverCmd = &cobra.Command{
	Use:   "webserver",
	Short: "Starts the webserver microservice",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Starting webserver microservice...")
		ws := webserver.NewService()
		if err := ws.Start(); err != nil {
			return err
		}
		fmt.Println("Webserver microservice started successfully.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(webserverCmd)
}
