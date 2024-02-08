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

func SetupUpCmd(ctx *MigrateCmdContext) func(cmd *cobra.Command, args []string) {
	logger := ctx.Logger

	return func(cmd *cobra.Command, args []string) {
		url := DB_SCHEME + ctx.Database

		logger.Infow("Running up migrations",
			"path", ctx.Database,
			"migrations", MIGRATIONS_URL,
			"steps", ctx.Steps,
		)

		m, err := migrate.New(MIGRATIONS_URL, url)
		if err != nil {
			logger.Fatal("failed to create migration instance:\n", err)
		}

		defer m.Close()

		if ctx.Steps > 0 {
			err = m.Steps(ctx.Steps)
		} else {
			err = m.Up()
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
}

func SetupDownCmd(ctx *MigrateCmdContext) func(cmd *cobra.Command, args []string) {
	logger := ctx.Logger

	return func(cmd *cobra.Command, args []string) {
		url := DB_SCHEME + ctx.Database

		logger.Infow("Running down migrations",
			"path", ctx.Database,
			"migrations", MIGRATIONS_URL,
			"steps", ctx.Steps,
		)

		m, err := migrate.New(MIGRATIONS_URL, url)
		if err != nil {
			logger.Fatal("failed to create migration instace: \n", err)
		}

		defer m.Close()

		if ctx.Steps > 0 {
			err = m.Steps(-ctx.Steps)
		} else {
			err = m.Down()
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
}
