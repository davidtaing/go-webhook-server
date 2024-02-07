package cmd

import (
	"github.com/davidtaing/go-webhook-server/internal/sender"
	"github.com/davidtaing/go-webhook-server/internal/server"
	"github.com/spf13/cobra"
)

// Sends a test event to the webhook server
var senderCmd = &cobra.Command{
	Use:   "sender",
	Short: "Sends a test event to the webhook server",
	Run: func(cmd *cobra.Command, args []string) {
		URL := "http://localhost:8080/webhook"
		event := server.WebhookEvent{
			ID:    "123",
			Event: "user_registered",
		}

		sender.SendEvent(event, URL)
	},
}

func init() {
	rootCmd.AddCommand(senderCmd)
}
