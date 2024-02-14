package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/davidtaing/go-webhook-server/internal/database"
	l "github.com/davidtaing/go-webhook-server/internal/logger"
	"github.com/davidtaing/go-webhook-server/internal/migration"
)

const WEBHOOK_HANDLER_ENDPOINT = "/webhook"
const TEST_FIXTURES_TEMP_DIR = "./test_fixtures"
const TEST_MIGRATIONS_DIR = "../../db/migrations"

var logger = l.New()

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

	sendRequest(h, eventJSON)
	rr, err := sendRequest(h, eventJSON)

	if err != nil {
		t.Fatal(err)
	}

	// expect a 200 OK status code, as non OK statuses will cause the duplicate event to be retried by the sender
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func setupTestFixture(dbName string) (s *server, cleanup func()) {
	dbPath := filepath.Join(TEST_FIXTURES_TEMP_DIR, dbName)

	// Create a temp directory for the test fixtures if it doesn't exist
	if _, err := os.Stat(TEST_FIXTURES_TEMP_DIR); os.IsNotExist(err) {
		err := os.Mkdir(TEST_FIXTURES_TEMP_DIR, 0755)
		if err != nil {
			logger.Fatal(err)
		}
	}

	migrationOpts := migration.MigrationOpts{
		DatabasePath:   dbPath,
		MigrationsPath: TEST_MIGRATIONS_DIR,
	}

	err := migration.RunUpMigrations(migrationOpts, logger)

	if err != nil {
		logger.Fatal(err)
	}

	s = &server{
		db:     database.Open("./"+dbPath, logger),
		logger: logger,
	}

	cleanup = func() {
		s.db.Close()
		err := os.RemoveAll(dbPath)
		if err != nil {
			logger.Fatal(err)
		}
	}

	return s, cleanup
}

func sendRequest(handler http.HandlerFunc, buf []byte) (*httptest.ResponseRecorder, error) {
	bodyBuf := bytes.NewBuffer(buf)

	req, err := http.NewRequest(http.MethodPost, WEBHOOK_HANDLER_ENDPOINT, bodyBuf)
	if err != nil {
		return nil, err
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	return rr, nil
}
