/**
 * Sends a test webhook event to the specified URL.
 * This isn't meant to be a production-ready webhook sender, but rather a simple util fire a webhook event.
 */
package sender

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type WebhookPayload struct {
	ID    string `json:"id"`
	Event string `json:"event"`
}

func SendEvent(event WebhookPayload, URL string) {
	eventJSON, err := json.Marshal(event)
	if err != nil {
		// Handle error
		return
	}

	resp, err := http.Post(URL, "application/json", bytes.NewBuffer(eventJSON))
	if err != nil {
		return
	}
	defer resp.Body.Close()
}
