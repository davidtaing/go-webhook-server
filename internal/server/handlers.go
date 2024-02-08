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

		// check if the webhook is already in DB
		// ✅ - if yes return 200 OK.

		// need to return 200 OK, or the request will be treated as a failure
		// and will be retried/resent by the sender

		// create if webhook doesn't exist

		// return creation result as status
		fmt.Println("Webhook received!")
	}
}
