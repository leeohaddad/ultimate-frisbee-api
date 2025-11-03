package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/leeohaddad/ultimate-frisbee-api/domain/port/repository"
	postgresRepositories "github.com/leeohaddad/ultimate-frisbee-api/infra/adapter/repository/postgres"
	postgresDatabase "github.com/leeohaddad/ultimate-frisbee-api/infra/database/postgres"

	"github.com/leeohaddad/ultimate-frisbee-api/domain/port/logger"
	"github.com/leeohaddad/ultimate-frisbee-api/infra/adapter/logger/zap"

	"github.com/leeohaddad/ultimate-frisbee-api/domain/port/config"
	"github.com/leeohaddad/ultimate-frisbee-api/infra/adapter/config/viper"
)

const (
	// EntrypointAPI represents the entrypoint for the API.
	EntrypointAPI = "api"
	// EntrypointMigration represents the entrypoint for the migration tool.
	EntrypointMigration = "migration"
	// EntrypointSeed represents the entrypoint for the database seeding tool.
	EntrypointSeed = "seed"
)

func main() {
	var startEntrypoint string
	var configFile string

	flag.StringVar(&startEntrypoint, "e", EntrypointAPI, "Define which entrypoint will be called on Ultimate Frisbee API. Options: [api, migration, seed] Default: [api]")
	flag.StringVar(&configFile, "config", "./config/local.yaml", "Path to the configuration file to be used by the service.")
	flag.Parse()

	applicationConfig := getConfig(configFile)

	switch startEntrypoint {
	case EntrypointAPI:
		runAPI(applicationConfig)
	case EntrypointMigration:
		runMigration(applicationConfig)
	case EntrypointSeed:
		runSeed(applicationConfig)
	default:
		panic(fmt.Errorf("unknown entrypoint: %s", startEntrypoint))
	}
}

func getConfig(configFile string) *config.Application {
	config, err := viper.NewConfig(configFile)
	if err != nil {
		fmt.Printf("error while trying to set up configs: %s\n", err.Error())

		return nil
	}

	applicationConfig, err := config.GetConfigs()
	if err != nil {
		fmt.Printf("error while trying to get configs: %s\n", err.Error())

		return nil
	}

	return applicationConfig
}

func getLogger() logger.Logger {
	applicationLogger, err := zap.NewLogger()

	if err != nil {
		fmt.Printf("error while trying to set up API zap: %s\n", err.Error())

		return nil
	}

	return applicationLogger
}

func getDatabase(applicationConfig *config.Application, logger logger.Logger) postgresDatabase.Client {
	databaseClient := postgresDatabase.NewClient(logger)
	err := databaseClient.Connect(applicationConfig.Database.ConnectionString)

	if err != nil {
		logger.WithError(err).Error("error while trying to set up database connection")

		return nil
	}

	return databaseClient
}

func getRepositories(applicationConfig *config.Application, databaseClient postgresDatabase.Client) repository.Collection {
	return repository.Collection{
		Team:   postgresRepositories.NewTeamRepository(databaseClient),
		Person: postgresRepositories.NewPersonRepository(databaseClient),
		// Tournament: postgresRepositories.NewRepository(databaseClient),
	}
}

func waitShutdown() {
	channel := make(chan os.Signal, 1)
	signal.Notify(channel, os.Interrupt, syscall.SIGTERM)
	sig := <-channel
	fmt.Printf("captured signal %v\n", sig)
}
