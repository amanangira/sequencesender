// Package migrations - A shim based on golang-migrate to execute in CLI.
package migrations

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"os"
	"sequencesender"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
)

var developmentFlag bool

var rootCmd = &cobra.Command{
	Use:   "migrate",
	Short: "SequenceSender database migration tool",
	Long:  "A CLI tool for managing database migrations for this application",
}

var migrateUpCmd = &cobra.Command{
	Use:   "up",
	Short: "Apply all pending migrations",
	Run: func(cmd *cobra.Command, args []string) {
		m := getMigrator()
		defer m.Close()

		if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			log.Fatal("failed to apply migrations:", err)
		}

		fmt.Println("Migrations applied successfully")
	},
}

var migrateDownCmd = &cobra.Command{
	Use:   "down",
	Short: "Rollback the last migration",
	Run: func(cmd *cobra.Command, args []string) {
		m := getMigrator()
		defer m.Close()

		if err := m.Steps(-1); err != nil && err != migrate.ErrNoChange {
			log.Fatal("failed to rollback migration:", err)
		}
		fmt.Println("Migration rolled back successfully")
	},
}

var migrateToCmd = &cobra.Command{
	Use:   "to [version]",
	Short: "Migrate to a specific version",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		version, err := strconv.ParseUint(args[0], 10, 32)
		if err != nil {
			log.Fatal("invalid version number:", err)
		}

		m := getMigrator()
		defer m.Close()

		if err := m.Migrate(uint(version)); err != nil && err != migrate.ErrNoChange {
			log.Fatal("failed to migrate to version", version, ":", err)
		}
		fmt.Printf("Migrated to version %d successfully\n", version)
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Get the current migration version",
	Run: func(cmd *cobra.Command, args []string) {
		m := getMigrator()
		defer m.Close()

		version, dirty, err := m.Version()
		if err != nil {
			log.Fatal("failed to get version:", err)
		}

		if dirty {
			fmt.Printf("Current version: %d (dirty)\n", version)
		} else {
			fmt.Printf("Current version: %d\n", version)
		}
	},
}

func getMigrator() *migrate.Migrate {
	// Load env if --development flag is set
	if developmentFlag {
		if err := godotenv.Load(); err != nil {
			slog.Error("failed to load .env file", "error", err)
		} else {
			slog.Info("environment variables loaded from .env file")
		}
	}

	dbURL := os.Getenv(sequencesender.EnvDBURLKey)
	if dbURL == "" {
		log.Fatal(sequencesender.EnvDBURLKey + " environment variable is required")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal("failed to create database driver:", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations/versions",
		"postgres",
		driver,
	)
	if err != nil {
		log.Fatal("failed to create migrator:", err)
	}

	return m
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&developmentFlag, "development", false, "Load environment variables from .env file for development")

	rootCmd.AddCommand(migrateUpCmd)
	rootCmd.AddCommand(migrateDownCmd)
	rootCmd.AddCommand(migrateToCmd)
	rootCmd.AddCommand(versionCmd)
}
