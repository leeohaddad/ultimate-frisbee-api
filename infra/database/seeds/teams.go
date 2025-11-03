package seeds

import (
	"context"
	"time"

	"github.com/leeohaddad/ultimate-frisbee-api/domain/entity"
	"github.com/leeohaddad/ultimate-frisbee-api/domain/port/logger"
	"github.com/leeohaddad/ultimate-frisbee-api/domain/port/repository"
)

// SeedTeams creates sample teams in the database
func SeedTeams(ctx context.Context, teamRepo repository.Team, logger logger.Logger) error {
	// Sample teams data
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

	logger.Info("starting team seeding...")

	for _, team := range teams {
		// Check if team already exists
		existing, err := teamRepo.GetTeamByName(ctx, team.name)
		if err != nil {
			logger.WithError(err).Errorf("error checking if team %s exists", team.name)
			continue
		}

		if existing != nil {
			logger.Infof("team %s already exists, skipping", team.name)
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

		createdTeam, err := teamRepo.CreateTeam(ctx, teamEntity)
		if err != nil {
			logger.WithError(err).Errorf("failed to create team %s", team.name)
			continue
		}

		logger.Infof("successfully created team: %s", createdTeam.Name)
	}

	logger.Info("team seeding completed!")
	return nil
}
