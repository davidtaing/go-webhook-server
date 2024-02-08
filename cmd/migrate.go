package cmd

import (
	l "github.com/davidtaing/go-webhook-server/internal/logger"
	"github.com/spf13/cobra"

	"github.com/golang-migrate/migrate/v4"
)

var steps int
var path string
var logger = l.New()

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Applies schema migrations to the specified sqlite database",
}

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Applies up migrations to the database.",
	Long: `Applies up migrations to the database.
	If the --steps flag is not provided, all pending up migrations will be applied.`,
	Run: func(cmd *cobra.Command, args []string) {
		logger.Infow("Running up migrations",
			"path", path,
			"steps", steps,
		)

		m, err := migrate.New("file://db/migrations", "sqlite3://"+path)
		if err != nil {
			logger.Fatal("Error creating migration instance: %v\n", err)
			return
		}

		defer m.Close()

		if steps > 0 {
			if err := m.Steps(steps); err != nil {
				logger.Fatal("Error applying migration steps: %v\n", err)
				return
			}
		}

		if err := m.Up(); err != nil {
			logger.Fatal("Error applying migrations: %v\n", err)
			return
		}
	},
}

var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Applies down migrations",
	Run: func(cmd *cobra.Command, args []string) {
		logger.Infow("Running down migrations",
			"path", path,
			"steps", steps,
		)

		m, err := migrate.New("file://db/migrations", "sqlite3://"+path)
		if err != nil {
			logger.Fatal("Error creating migration instance: %v\n", err)
			return
		}

		defer m.Close()

		if steps > 0 {
			if err := m.Steps(-steps); err != nil {
				logger.Fatal("Error applying migration steps: %v\n", err)
				return
			}
		}

		if err := m.Down(); err != nil {
			logger.Fatal("Error applying migrations: %v\n", err)
			return
		}
	},
}

func init() {
	migrateCmd.AddCommand(upCmd)
	migrateCmd.AddCommand(downCmd)
	rootCmd.AddCommand(migrateCmd)

	upCmd.Flags().StringVar(&path, "path", "", "Path to the migration files")
	upCmd.Flags().IntVar(&steps, "steps", 0, "Number of steps to migrate. If no steps are provided, all migrations will be applied.")

	downCmd.Flags().StringVar(&path, "path", "", "Path to the migration files")
	downCmd.Flags().IntVar(&steps, "steps", 0, "Number of steps to migrate. If no steps are provided, all migrations will be applied.")
}
