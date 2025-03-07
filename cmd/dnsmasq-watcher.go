// cmd/dnsmasq-watcher.go
package cmd

import (
    "fmt"

    "github.com/spf13/cobra"
    "github.com/jdfalk/ubuntu-autoinstall-webhook/internal/dnsmasqwatcher"
)

var dnsmasqWatcherCmd = &cobra.Command{
    Use:   "dnsmasq-watcher",
    Short: "Starts the dnsmasq-watcher microservice",
    RunE: func(cmd *cobra.Command, args []string) error {
        fmt.Println("Starting dnsmasq-watcher microservice...")
        watcher := dnsmasqwatcher.NewService()
        if err := watcher.Start(); err != nil {
            return err
        }
        fmt.Println("Dnsmasq-watcher microservice started successfully.")
        return nil
    },
}

func init() {
    rootCmd.AddCommand(dnsmasqWatcherCmd)
}
