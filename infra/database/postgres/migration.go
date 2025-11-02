package postgres

import (
	"errors"
	"fmt"

	migrate "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // import the postgres driver used by go-migrate
	_ "github.com/golang-migrate/migrate/v4/source/file"       // import the file driver used by go-migrate
)

func RunMigrations(migrationRelativePath, connectionString string) error {
	// use "file://" instead of "file:///" to point the file on the relative path to the binary
	// ref: https://github.com/golang-migrate/migrate/tree/master/source/file
	migrationDir := fmt.Sprintf("file://%s", migrationRelativePath)

	migrationTool, err := migrate.New(migrationDir, connectionString)
	if err != nil {
		return fmt.Errorf("failed to create migration engine: %w", err)
	}
	defer migrationTool.Close()

	err = migrationTool.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to run migration engine: %w", err)
	}

	return nil
}

// RunMigrationsDown runs all down migrations (rollback)
func RunMigrationsDown(migrationRelativePath, connectionString string) error {
	migrationDir := fmt.Sprintf("file://%s", migrationRelativePath)

	migrationTool, err := migrate.New(migrationDir, connectionString)
	if err != nil {
		return fmt.Errorf("failed to create migration engine: %w", err)
	}
	defer migrationTool.Close()

	err = migrationTool.Down()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to run down migrations: %w", err)
	}

	return nil
}

// GetMigrationVersion returns the current migration version
func GetMigrationVersion(migrationRelativePath, connectionString string) (uint, bool, error) {
	migrationDir := fmt.Sprintf("file://%s", migrationRelativePath)

	migrationTool, err := migrate.New(migrationDir, connectionString)
	if err != nil {
		return 0, false, fmt.Errorf("failed to create migration engine: %w", err)
	}
	defer migrationTool.Close()

	version, dirty, err := migrationTool.Version()
	if err != nil {
		return 0, false, fmt.Errorf("failed to get migration version: %w", err)
	}

	return version, dirty, nil
}
