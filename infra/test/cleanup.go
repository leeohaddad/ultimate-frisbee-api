package test

import (
	"context"
	"fmt"
	"testing"

	adapterViper "github.com/leeohaddad/ultimate-frisbee-api/infra/adapter/config/viper"
	loggerZap "github.com/leeohaddad/ultimate-frisbee-api/infra/adapter/logger/zap"
	databasePostgres "github.com/leeohaddad/ultimate-frisbee-api/infra/database/postgres"
)

// CleanupTeams removes any teams with the given names from the configured local database.
// This is intended to be used by tests to make the environment deterministic.
func CleanupTeams(names []string) error {
	// Build a logger (non-testing helper)
	zapLogger, err := loggerZap.NewLogger()
	if err != nil {
		return fmt.Errorf("failed to build logger for cleanup: %w", err)
	}

	client := databasePostgres.NewClient(zapLogger)

	// Load config from file using viper adapter (config/local.yaml expected)
	cfg, err := adapterViper.NewConfig("config/local.yaml")
	if err != nil {
		return fmt.Errorf("failed to load config for cleanup: %w", err)
	}

	appConfigs, err := cfg.GetConfigs()
	if err != nil {
		return fmt.Errorf("failed to map configs for cleanup: %w", err)
	}

	connectionString := appConfigs.Database.ConnectionString

	if err := client.Connect(connectionString); err != nil {
		return fmt.Errorf("failed to connect to database for cleanup: %w", err)
	}
	defer client.Close()

	ctx := context.Background()

	for _, name := range names {
		// Use the same placeholder style as repository code ("?")
		query := "DELETE FROM teams WHERE name = ?"
		if _, err := client.ExecuteCommand(ctx, query, name); err != nil {
			return fmt.Errorf("failed to delete test team '%s': %w", name, err)
		}
	}

	return nil
}

// CleanupTeamsT is a test-friendly wrapper that fails the test on error.
// Keep this in the package for convenience in tests that have access to *testing.T.
func CleanupTeamsT(t *testing.T, names []string) {
	t.Helper()
	if err := CleanupTeams(names); err != nil {
		t.Fatalf("cleanup failed: %v", err)
	}
}
