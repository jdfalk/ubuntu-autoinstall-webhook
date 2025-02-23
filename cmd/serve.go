package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/jdfalk/ubuntu-autoinstall-webhook/internal/server"
	"github.com/jdfalk/ubuntu-autoinstall-webhook/internal/server/logger"
)

// Serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the webhook server",
	Run: func(cmd *cobra.Command, args []string) {
		// Use default port "8080" if not set.
		port := viper.GetString("port")
		if port == "" {
			port = "8080"
		}
		// Log the startup message.
		logger.Infof("Webhook server running on port %s", port)
		// Start the server using the proper function.
		server.StartServer(port)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
