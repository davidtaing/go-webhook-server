package cmd

import (
	l "github.com/davidtaing/go-webhook-server/internal/logger"
	"github.com/davidtaing/go-webhook-server/internal/migration"

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
	logger := l.New()
	opts := &migration.MigrationOpts{}

	downCmd.Flags().StringVarP(&opts.DatabasePath, "database", "d", "", "path to the sqlite3 database file")
	downCmd.MarkFlagRequired("database")
	downCmd.Flags().IntVarP(&opts.Steps, "steps", "s", 0, "Number of steps to migrate. If no steps are provided, all migrations will be applied.")

	downCmd.Run = migration.SetupDownCmd(*opts, logger)

	rootCmd.AddCommand(downCmd)
}
