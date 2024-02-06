package sender

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// Abitrary typedef to represent an event
type Event struct {
	ID    int    `json:"id"`
	Event string `json:"event"`
}

func SendEvent(event Event, URL string) {
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
