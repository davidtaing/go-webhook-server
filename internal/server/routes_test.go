package server

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/davidtaing/go-webhook-server/internal/repository"
)

func TestWebhookHandler_MethodNotAllowed(t *testing.T) {
	tests := []struct {
		name   string
		method string
	}{
		{"GET", http.MethodGet},
		{"PUT", http.MethodPut},
		{"DELETE", http.MethodDelete},
		{"PATCH", http.MethodPatch},
	}

	srv, cleanup := setupTestFixture("TestWebhookHandler_MethodNotAllowed.db")
	defer cleanup()

	h := srv.handleWebhook()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr, err := sendRequest(requestOpts{
				handler: h,
				method:  tt.method,
				body:    []byte{},
			})

			if err != nil {
				t.Fatal(err)
			}

			if status := rr.Code; status != http.StatusMethodNotAllowed {
				t.Errorf("handler returned wrong status code for method %s: got %v want %v", tt.method, status, http.StatusMethodNotAllowed)
			}
		})
	}
}

func TestWebhookHandler_HandleDuplicateEvent(t *testing.T) {
	srv, cleanup := setupTestFixture("./TestWebhookHandler_HandleDuplicateEvent.db")
	defer cleanup()

	h := srv.handleWebhook()

	event := struct {
		ID    string `json:"id"`
		Event string `json:"event"`
	}{
		ID:    "1",
		Event: "test",
	}

	eventJSON, err := json.Marshal(event)
	if err != nil {
		t.Fatal(err)
	}

	requestOpts := requestOpts{
		handler: h,
		body:    eventJSON,
	}

	sendRequest(requestOpts)
	rr, err := sendRequest(requestOpts)

	if err != nil {
		t.Fatal(err)
	}

	t.Run("recieve 200 status", func(t *testing.T) {
		// expect a 200 OK status code, as non OK statuses will cause the duplicate event to be retried by the sender
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}
	})

	t.Run("only one webhook in database", func(t *testing.T) {
		repo := &repository.WebhookRepository{DB: srv.db}
		filters := map[string]interface{}{
			"id": event.ID,
		}

		webhooks, err := repo.Get(&filters)
		if err != nil {
			t.Fatal(err)
		}

		if len(webhooks) != 1 {
			t.Errorf("expected 1 webhook in database, got %d", len(webhooks))
		}
	})
}
