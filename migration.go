package main

import (
	"github.com/leeohaddad/ultimate-frisbee-api/domain/port/config"
	"github.com/leeohaddad/ultimate-frisbee-api/infra/database/postgres"
)

const migrationPath = "./infra/database/migrations"

func runMigration(applicationConfig *config.Application) {
	applicationLogger := getLogger()

	err := postgres.RunMigrations(migrationPath, applicationConfig.Database.ConnectionString)
	if err != nil {
		applicationLogger.Errorf("Error executing migrations: %w", err)
		panic(err)
	}

	applicationLogger.Info("Migrations executed successfully!")
}
