package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/leeohaddad/ultimate-frisbee-api/domain/entity"
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

	// Create sample teams
	teams := []struct {
		slug          string
		name          string
		description   string
		originCountry string
		createdBy     string
	}{
		{
			slug:          "ultimate-warriors",
			name:          "Ultimate Warriors",
			description:   "A competitive ultimate frisbee team from California",
			originCountry: "USA",
			createdBy:     "admin",
		},
		{
			slug:          "disc-dynamos",
			name:          "Disc Dynamos",
			description:   "Professional ultimate frisbee team from New York",
			originCountry: "USA",
			createdBy:     "admin",
		},
		{
			slug:          "flying-circus",
			name:          "Flying Circus",
			description:   "European championship ultimate frisbee team",
			originCountry: "Germany",
			createdBy:     "admin",
		},
	}

	ctx := context.Background()
	for _, team := range teams {
		// Check if team already exists
		existing, err := repositories.Team.GetTeamByName(ctx, team.name)
		if err != nil {
			applicationLogger.WithError(err).Errorf("error checking if team %s exists", team.name)
			continue
		}

		if existing != nil {
			applicationLogger.Infof("team %s already exists, skipping", team.name)
			continue
		}

		// Create new team
		teamEntity := &entity.Team{
			Slug:          team.slug,
			Name:          team.name,
			Description:   team.description,
			OriginCountry: team.originCountry,
			CreatedBy:     team.createdBy,
			CreatedAt:     time.Now(),
			UpdatedBy:     team.createdBy,
			UpdatedAt:     time.Now(),
		}

		createdTeam, err := repositories.Team.CreateTeam(ctx, teamEntity)
		if err != nil {
			applicationLogger.WithError(err).Errorf("failed to create team %s", team.name)
			continue
		}

		applicationLogger.Infof("successfully created team: %s", createdTeam.Name)
	}

	// Create sample people
	people := []struct {
		userName      string
		name          string
		email         string
		phoneNumber   string
		wfdfNumber    string
		originCountry string
		createdBy     string
	}{
		{
			userName:      "notdougz",
			name:          "Douglas Olvieira",
			email:         "doug@gmail.com",
			phoneNumber:   "(11) 98765-4321",
			wfdfNumber:    "12",
			originCountry: "Brazil",
			createdBy:     "admin",
		},
		{
			userName:      "allanbm100",
			name:          "Allan Moreira",
			email:         "allan@gmail.com",
			phoneNumber:   "(11) 98765-4321",
			wfdfNumber:    "34",
			originCountry: "Brazil",
			createdBy:     "admin",
		},
		{
			userName:      "Iolivieri",
			name:          "Isabella Olivieri",
			email:         "bella@gmail.com",
			phoneNumber:   "(11) 98765-4321",
			wfdfNumber:    "56",
			originCountry: "Brazil",
			createdBy:     "admin",
		},
	}
	for _, person := range people {
		// Check if person already exists
		existing, err := repositories.Person.GetPersonByUserName(ctx, person.userName)
		if err != nil {
			applicationLogger.WithError(err).Errorf("error checking if person %s exists", person.userName)
			continue
		}

		if existing != nil {
			applicationLogger.Infof("person %s already exists, skipping", person.userName)
			continue
		}

		// Create new person
		personEntity := &entity.Person{
			UserName:      person.userName,
			Name:          person.name,
			Email:         person.email,
			PhoneNumber:   person.phoneNumber,
			WFDFNumber:    person.wfdfNumber,
			OriginCountry: person.originCountry,
			CreatedBy:     person.createdBy,
			CreatedAt:     time.Now(),
			UpdatedBy:     person.createdBy,
			UpdatedAt:     time.Now(),
		}

		createdPerson, err := repositories.Person.CreatePerson(ctx, personEntity)
		if err != nil {
			applicationLogger.WithError(err).Errorf("failed to create person %s", person.name)
			continue
		}

		applicationLogger.Infof("successfully created person: %s", createdPerson.Name)
	}

	applicationLogger.Info("database seeding completed!")
}
