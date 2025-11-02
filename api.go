package main

import (
	"sync"

	"github.com/leeohaddad/ultimate-frisbee-api/domain/port/config"

	"github.com/leeohaddad/ultimate-frisbee-api/infra/api"
)

const AppName string = "Ultimate Frisbee API"

func runAPI(applicationConfig *config.Application) {
	applicationLogger := getLogger()

	// Configure API
	applicationLogger.Infof("configuring the %s...", AppName)
	databaseClient := getDatabase(applicationConfig, applicationLogger)
	repositories := getRepositories(applicationConfig, databaseClient)
	apiApp := api.NewApp(applicationLogger, applicationConfig, repositories)
	applicationLogger.Infof("the %s was configured successfully", AppName)

	// Start API
	applicationLogger.Infof("starting the %s...", AppName)
	var wg sync.WaitGroup
	wg.Add(1)
	if err := apiApp.Start(&wg); err != nil {
		applicationLogger.WithError(err).Errorf("error while starting the %s", AppName)
		return
	}
	applicationLogger.Infof("the %s was started successfully", AppName)

	// Stop API (when shutdown signal is received)
	waitShutdown()
	apiApp.Stop()
	applicationLogger.Infof("gracefully shutting down the %s", AppName)
}
