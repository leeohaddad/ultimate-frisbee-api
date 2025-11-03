package main

import (
	"context"
	"fmt"

	"github.com/leeohaddad/ultimate-frisbee-api/domain/port/config"
	"github.com/leeohaddad/ultimate-frisbee-api/infra/adapter/logger/zap"
	"github.com/leeohaddad/ultimate-frisbee-api/infra/database/seeds"
)

func runSeed(applicationConfig *config.Application) {
	applicationLogger, err := zap.NewLogger()
	if err != nil {
		fmt.Printf("error while trying to set up logger: %s\n", err.Error())
		return
	}

	applicationLogger.Info("starting database seeding...")

	databaseClient := getDatabase(applicationConfig, applicationLogger)
	if databaseClient == nil {
		applicationLogger.Error("failed to connect to database")
		return
	}

	repositories := getRepositories(applicationConfig, databaseClient)
	ctx := context.Background()

	// Seed teams
	if err := seeds.SeedTeams(ctx, repositories.Team, applicationLogger); err != nil {
		applicationLogger.WithError(err).Error("failed to seed teams")
		return
	}

	// Seed people
	if err := seeds.SeedPeople(ctx, repositories.Person, applicationLogger); err != nil {
		applicationLogger.WithError(err).Error("failed to seed people")
		return
	}

	applicationLogger.Info("database seeding completed!")
}
