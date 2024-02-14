package migration

import (
	"fmt"

	"github.com/davidtaing/go-webhook-server/internal/logger"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/cobra"
)

const DB_SCHEME = "sqlite3://"
const MIGRATIONS_SCHEME = "file://"
const MIGRATIONS_DIR = "./db/migrations"

type MigrationOpts struct {
	Steps          int
	DatabasePath   string
	MigrationsPath string // defaults to ./db/migrations if this is a nil value. Override for testing as working directory is set to the file location.
}

// Run up migrations on the SQLite database. If the database does not exist, a new database will be created.
func RunUpMigrations(opts MigrationOpts, logger *logger.Logger) error {
	url := DB_SCHEME + opts.DatabasePath
	migrationURL := getMigrationURL(opts.MigrationsPath)

	logger.Debugw("Running up migrations",
		"dbName", opts.DatabasePath,
		"migrations", migrationURL,
		"steps", opts.Steps,
	)

	m, err := migrate.New(migrationURL, url)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}

	defer m.Close()

	if opts.Steps < 0 {
		return fmt.Errorf("steps must be greater than or equal to 0")
	}

	if opts.Steps == 0 {
		err = m.Up()
	} else {
		err = m.Steps(opts.Steps)
	}

	if err != nil && err.Error() == "no change" {
		fmt.Println("no change")
		return nil
	}

	if err != nil {
		return fmt.Errorf("error applying migrations: %w", err)
	}

	return nil
}

// Run down migrations on the SQLite database. If the database does not exist, a new database will be created.
func RunDownMigrations(opts MigrationOpts, logger *logger.Logger) error {
	url := DB_SCHEME + opts.DatabasePath
	migrationURL := getMigrationURL(opts.MigrationsPath)

	logger.Debugw("Running down migrations",
		"dbName", opts.DatabasePath,
		"migrations", migrationURL,
		"steps", opts.Steps,
	)

	m, err := migrate.New(migrationURL, url)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}

	defer m.Close()

	if opts.Steps < 0 {
		return fmt.Errorf("steps must be greater than or equal to 0")
	}

	if opts.Steps == 0 {
		err = m.Down()
	} else {
		// negative steps denote down migrations in the go-migrate API
		err = m.Steps(-opts.Steps)
	}

	if err != nil && err.Error() == "no change" {
		fmt.Println("no change")
		return nil
	}

	if err != nil {
		return fmt.Errorf("error applying migrations: %w", err)
	}

	return nil

}

func SetupUpCmd(opts MigrationOpts, logger *logger.Logger) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		err := RunUpMigrations(opts, logger)

		if err != nil {
			logger.Fatal(err)
		}
	}
}

func SetupDownCmd(opts MigrationOpts, logger *logger.Logger) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		err := RunDownMigrations(opts, logger)

		if err != nil {
			logger.Fatal(err)
		}
	}
}

// Returns the URL for the migration directory. If an empty path is provided, the default path is used.
//
// Use the path parameter to override the default path for testing.
func getMigrationURL(path string) string {
	migrationDir := MIGRATIONS_DIR

	if path != "" {
		migrationDir = path
	}

	return MIGRATIONS_SCHEME + migrationDir
}
