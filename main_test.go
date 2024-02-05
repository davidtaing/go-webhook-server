package main

import (
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
		req, err := http.NewRequest(tt.method, "/webhook", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(WebhookHandler)

		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusMethodNotAllowed {
			t.Errorf("handler returned wrong status code for method %s: got %v want %v", tt.method, status, http.StatusMethodNotAllowed)
		}
	}
}
