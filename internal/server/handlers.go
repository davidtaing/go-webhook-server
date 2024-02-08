package server

import (
	"fmt"
	"net/http"
)

func (s *server) handleWebhook() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}

		fmt.Println("Webhook received!")
	}
}
