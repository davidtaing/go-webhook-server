package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/davidtaing/go-webhook-server/internal/database"
	"github.com/davidtaing/go-webhook-server/internal/logger"
	"github.com/davidtaing/go-webhook-server/internal/models"
)

const WEBHOOK_HANDLER_ENDPOINT = "/webhook"

func setupTestServer(path string) *server {
	l := logger.New()

	s := &server{
		db:     database.Open(path, l),
		logger: l,
	}

	return s
}

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

	srv := setupTestServer("")
	defer srv.db.Close()

	h := srv.handleWebhook()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, WEBHOOK_HANDLER_ENDPOINT, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			h.ServeHTTP(rr, req)

			if status := rr.Code; status != http.StatusMethodNotAllowed {
				t.Errorf("handler returned wrong status code for method %s: got %v want %v", tt.method, status, http.StatusMethodNotAllowed)
			}
		})
	}
}

func TestWebhookHandler_HandleDuplicateEvent(t *testing.T) {
	srv := setupTestServer("./TestWebhookHandler_HandleDuplicateEvent.db")
	defer srv.db.Close()

	h := srv.handleWebhook()

	event := models.Webhook{
		ID:    "1",
		Event: "test",
	}

	eventJSON, err := json.Marshal(event)
	if err != nil {
		t.Fatal(err)
	}

	eventBuffer := bytes.NewBuffer(eventJSON)

	req, err := http.NewRequest(http.MethodPost, WEBHOOK_HANDLER_ENDPOINT, eventBuffer)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)
	h.ServeHTTP(rr, req)

	// expect a 200 OK status code, as non OK statuses will cause the duplicate event to be retried by the sender
	t.Run("returns 200 OK upon duplicated event", func(t *testing.T) {
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}
	})
}
