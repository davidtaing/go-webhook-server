package models

type Webhook struct {
	ID      string `json:"id"`
	Event   string `json:"event"`
	Payload string `json:"payload"`
}
