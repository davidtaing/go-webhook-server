package models

type Webhook struct {
	ID      string `json:"id"`
	Source  string `json:"source"`
	Event   string `json:"event"`
	Created string `json:"created"`
	Updated string `json:"updated"`
	Payload string `json:"payload"`
}
