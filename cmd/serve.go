package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/jdfalk/ubuntu-autoinstall-webhook/internal/db"
	"github.com/jdfalk/ubuntu-autoinstall-webhook/internal/logger"
	"github.com/jdfalk/ubuntu-autoinstall-webhook/internal/server"
)

/*
InitDBFunc is the exported function variable for initializing the database.
It is set to the real db.InitDB by default but can be overridden for testing.
*/
var InitDBFunc = db.InitDB

// serveCmd represents the "serve" command, which starts the webhook server.
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the webhook server",
	Run: func(cmd *cobra.Command, args []string) {
		// Use default port "8080" if not set.
		port := viper.GetString("port")
		if port == "" {
			port = "8080"
		}

		// Initialize the database connection using the exported InitDBFunc.
		err := InitDBFunc()
		if err != nil {
			logger.Errorf("Error initializing database: %v", err.Error())
			return
		}
		// Ensure the database is closed when the server stops using the shared SafeClose helper.
		defer db.SafeClose()

		// Log the startup message.
		logger.Infof("Webhook server running on port %s", port)
		// Start the server on the specified port.
		server.StartServer(port)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
