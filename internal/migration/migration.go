package migration

import (
	"fmt"

	"github.com/davidtaing/go-webhook-server/internal/logger"
	"github.com/golang-migrate/migrate/v4"
	"github.com/spf13/cobra"
)

const DB_SCHEME = "sqlite3://"
const MIGRATIONS_URL = "file://./db/migrations"

type MigrateCmdContext struct {
	Logger   *logger.Logger
	Steps    int
	Database string
}

// Runs up migrations for the given database path
func RunUpMigrations(path string, steps int, logger *logger.Logger) {
	url := DB_SCHEME + path

	logger.Infow("Running up migrations",
		"path", path,
		"migrations", MIGRATIONS_URL,
		"steps", steps,
	)

	m, err := migrate.New(MIGRATIONS_URL, url)
	if err != nil {
		logger.Fatal("failed to create migration instance:\n", err)
	}

	defer m.Close()

	if steps == 0 {
		err = m.Up()
	} else if steps > 0 {
		err = m.Steps(steps)
	} else {
		logger.Fatal("steps must be greater than or equal to 0")
	}

	if err != nil && err.Error() == "no change" {
		fmt.Println("no change")
		return
	}

	if err != nil {
		logger.Fatal("Error applying migrations:\n", err)
		return
	}
}

func RunDownMigrations(path string, steps int, logger *logger.Logger) {
	url := DB_SCHEME + path

	logger.Infow("Running down migrations",
		"path", path,
		"migrations", MIGRATIONS_URL,
		"steps", steps,
	)

	m, err := migrate.New(MIGRATIONS_URL, url)
	if err != nil {
		logger.Fatal("failed to create migration instace: \n", err)
	}

	defer m.Close()

	if steps == 0 {
		err = m.Down()
	} else if steps > 0 {
		// negative steps denote down migrations in the go-migrate API
		err = m.Steps(-steps)
	} else {
		logger.Fatal("steps must be greater than or equal to 0")
	}

	if err != nil && err.Error() == "no change" {
		fmt.Println("no change")
		return
	}

	if err != nil {
		logger.Fatal("Error applying migrations: \n", err)
		return
	}
}

func SetupUpCmd(ctx *MigrateCmdContext) func(cmd *cobra.Command, args []string) {
	logger := ctx.Logger

	return func(cmd *cobra.Command, args []string) {
		RunUpMigrations(ctx.Database, ctx.Steps, logger)
	}
}

func SetupDownCmd(ctx *MigrateCmdContext) func(cmd *cobra.Command, args []string) {
	logger := ctx.Logger

	return func(cmd *cobra.Command, args []string) {
		RunDownMigrations(ctx.Database, ctx.Steps, logger)
	}
}
