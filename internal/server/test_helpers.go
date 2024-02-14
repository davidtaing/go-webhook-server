package server

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"

	"github.com/davidtaing/go-webhook-server/internal/database"
	l "github.com/davidtaing/go-webhook-server/internal/logger"
	"github.com/davidtaing/go-webhook-server/internal/migration"
)

const WEBHOOK_HANDLER_ENDPOINT = "/webhook"
const TEST_FIXTURES_TEMP_DIR = "./test_fixtures"
const TEST_MIGRATIONS_DIR = "../../db/migrations"

var logger = l.New()

// setupTestFixture creates a new server and runs the migrations for the given test fixture.
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

type requestOpts struct {
	handler http.HandlerFunc
	body    []byte
	method  string // this will be treated as POST if not provided
}

func sendRequest(opts requestOpts) (*httptest.ResponseRecorder, error) {
	if opts.method == "" {
		opts.method = http.MethodPost
	}

	bodyBuf := bytes.NewBuffer(opts.body)

	req, err := http.NewRequest(opts.method, WEBHOOK_HANDLER_ENDPOINT, bodyBuf)
	if err != nil {
		return nil, err
	}

	rr := httptest.NewRecorder()
	opts.handler.ServeHTTP(rr, req)

	return rr, nil
}
