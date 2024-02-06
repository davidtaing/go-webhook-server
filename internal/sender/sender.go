package sender

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/davidtaing/go-webhook-server/internal/server"
)

func SendEvent(event server.WebhookEvent, URL string) {
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
