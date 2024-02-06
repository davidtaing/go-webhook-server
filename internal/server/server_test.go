package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, "/webhook", nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			env := Env{}
			handler := http.HandlerFunc(env.WebhookHandler)

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != http.StatusMethodNotAllowed {
				t.Errorf("handler returned wrong status code for method %s: got %v want %v", tt.method, status, http.StatusMethodNotAllowed)
			}
		})
	}
}

func TestWebhookHandler_HandleDuplicateEvent(t *testing.T) {
	event := WebhookEvent{
		ID:    1,
		Event: "test",
	}

	eventJSON, err := json.Marshal(event)
	if err != nil {
		t.Fatal(err)
	}

	eventBuffer := bytes.NewBuffer(eventJSON)

	_, err = http.NewRequest(http.MethodPost, "/webhook", eventBuffer)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest(http.MethodPost, "/webhook", eventBuffer)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	env := Env{}
	handler := http.HandlerFunc(env.WebhookHandler)

	handler.ServeHTTP(rr, req)

	// expect a 200 OK status code, as non OK statuses will cause the duplicate event to be retried by the sender
	t.Run("returns 200 OK upon duplicated event", func(t *testing.T) {
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}
	})
}
