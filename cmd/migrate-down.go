package cmd

import (
	l "github.com/davidtaing/go-webhook-server/internal/logger"
	"github.com/davidtaing/go-webhook-server/internal/migration"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/cobra"
)

var downCmd = &cobra.Command{
	Use:   "migrate-down",
	Short: "Applies down migrations",
	Long: `Applies down migrations to the database.
	If the --steps flag is not provided, all pending down migrations will be applied.
	Due to the way the sqlite3 driver works, down migrations may create a new database file if the database does not exist.
	`,
}

func init() {
	ctx := &migration.MigrateCmdContext{
		Logger: l.New(),
	}

	downCmd.Flags().StringVarP(&ctx.Database, "database", "d", "", "path to the sqlite3 database file")
	downCmd.MarkFlagRequired("database")
	downCmd.Flags().IntVarP(&ctx.Steps, "steps", "s", 0, "Number of steps to migrate. If no steps are provided, all migrations will be applied.")

	downCmd.Run = migration.SetupDownCmd(ctx)

	rootCmd.AddCommand(downCmd)
}
