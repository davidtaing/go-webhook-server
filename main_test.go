package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWebhookHandler_MethodNotAllowed(t *testing.T) {
	methods := []string{
		http.MethodGet,
		http.MethodPut,
		http.MethodDelete,
		http.MethodPatch,
		http.MethodOptions,
		http.MethodHead,
		http.MethodConnect,
		http.MethodTrace,
	}

	for _, method := range methods {
		req, err := http.NewRequest(method, "/webhook", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(WebhookHandler)

		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusMethodNotAllowed {
			t.Errorf("handler returned wrong status code for method %s: got %v want %v", method, status, http.StatusMethodNotAllowed)
		}
	}
}
