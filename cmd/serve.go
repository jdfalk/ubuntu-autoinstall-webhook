package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/jdfalk/ubuntu-autoinstall-webhook/internal/db"
	"github.com/jdfalk/ubuntu-autoinstall-webhook/internal/logger"
	"github.com/jdfalk/ubuntu-autoinstall-webhook/internal/server"
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
		// Initialize the database connection.
		// This function should be defined in the db package.
		// It should return an error if the connection fails.
		err := db.InitDB()
		if err != nil {
			logger.Errorf("Error initializing database: %v", err.Error())
			return
		}
		// Ensure the database is closed when the server stops.
		defer db.CloseDB()
		// Log the startup message.
		logger.Infof("Webhook server running on port %s", port)
		// Start the server using the proper function.
		server.StartServer(port)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
