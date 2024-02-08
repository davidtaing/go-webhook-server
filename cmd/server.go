package cmd

import (
	"github.com/davidtaing/go-webhook-server/internal/server"
	"github.com/spf13/cobra"
)

// Starts the webhook server on port 8080
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Starts the webhook server on port 8080",
	Run: func(cmd *cobra.Command, args []string) {
		server.Run()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
