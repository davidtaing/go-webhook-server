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
const MIGRATIONS_URL = "file://./db/migrations"

type MigrateCmdContext struct {
	Logger   *logger.Logger
	Steps    int
	Database string
}

// Run up migrations on the SQLite database. If the database does not exist, a new database will be created.
func RunUpMigrations(path string, steps int, logger *logger.Logger) error {
	url := DB_SCHEME + path

	logger.Infow("Running up migrations",
		"path", path,
		"migrations", MIGRATIONS_URL,
		"steps", steps,
	)

	m, err := migrate.New(MIGRATIONS_URL, url)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}

	defer m.Close()

	if steps < 0 {
		return fmt.Errorf("steps must be greater than or equal to 0")
	}

	if steps == 0 {
		err = m.Up()
	} else {
		err = m.Steps(steps)
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
func RunDownMigrations(path string, steps int, logger *logger.Logger) error {
	url := DB_SCHEME + path

	logger.Infow("Running down migrations",
		"path", path,
		"migrations", MIGRATIONS_URL,
		"steps", steps,
	)

	m, err := migrate.New(MIGRATIONS_URL, url)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}

	defer m.Close()

	if steps < 0 {
		return fmt.Errorf("steps must be greater than or equal to 0")
	}

	if steps == 0 {
		err = m.Down()
	} else {
		// negative steps denote down migrations in the go-migrate API
		err = m.Steps(-steps)
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

func SetupUpCmd(ctx *MigrateCmdContext) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		err := RunUpMigrations(ctx.Database, ctx.Steps, ctx.Logger)

		if err != nil {
			ctx.Logger.Fatal(err)
		}
	}
}

func SetupDownCmd(ctx *MigrateCmdContext) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		err := RunDownMigrations(ctx.Database, ctx.Steps, ctx.Logger)

		if err != nil {
			ctx.Logger.Fatal(err)
		}
	}
}
