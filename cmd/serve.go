package cmd

import (
	"log"

	"github.com/jdfalk/ubuntu-autoinstall-webhook/internal/server"
	"github.com/jdfalk/ubuntu-autoinstall-webhook/internal/server/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the webhook server",
	Run: func(cmd *cobra.Command, args []string) {
		port := viper.GetString("port")
		// Log both to the file and the SQL database.
		logger.Info("Webhook server running on port %s", port)
		log.Fatal(server.StartServer(port))
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
