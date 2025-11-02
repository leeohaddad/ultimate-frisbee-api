//go:build integration
// +build integration

package postgres_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/leeohaddad/ultimate-frisbee-api/domain/entity"
	repositoryPostgres "github.com/leeohaddad/ultimate-frisbee-api/infra/adapter/repository/postgres"
	databasePostgres "github.com/leeohaddad/ultimate-frisbee-api/infra/database/postgres"
	"github.com/leeohaddad/ultimate-frisbee-api/infra/test"
	"github.com/leeohaddad/ultimate-frisbee-api/infra/test/fixture"
)

func TestTeamRepository_GetTeamByName(t *testing.T) {
	t.Parallel()

	scenarios := []test.FixtureScenario{
		{
			Description:    "should return no team when the database has no team with the specified name",
			FixtureQueries: fixture.GenerateTeamQueries(fixture.GetAnotherFixtureTeam()),
			InputData: map[string]interface{}{
				"name": fixture.GetDefaultFixtureTeam().Name,
			},
			OutputData: map[string]interface{}{
				"expectedTeam": (*entity.Team)(nil),
			},
		},
		{
			Description:    "should return the desired team when the database has a team with this ID",
			FixtureQueries: fixture.GenerateTeamQueries(fixture.GetDefaultFixtureTeam()),
			InputData: map[string]interface{}{
				"name": fixture.GetDefaultFixtureTeam().Name,
			},
			OutputData: map[string]interface{}{
				"expectedTeam": fixture.GetDefaultFixtureTeam(),
			},
		},
	}

	test.RunFixtureScenarios(
		t, scenarios,
		func(t *testing.T, testContext context.Context, client databasePostgres.Client, scenario test.FixtureScenario) {
			t.Helper()
			teamRepository := repositoryPostgres.NewTeamRepository(client)

			// Prepare dependencies (arrange)
			id, ok := scenario.InputData["name"].(string)
			require.True(t, ok)
			expectedTeam, ok := scenario.OutputData["expectedTeam"].(*entity.Team)
			require.True(t, ok)

			// Execute method to fetch the entity
			obtainedTeam, err := teamRepository.GetTeamByName(testContext, id)
			require.NoError(t, err)

			// Check if fetched entity is filled correctly (assert)
			if expectedTeam == nil {
				require.Nil(t, obtainedTeam)

				return
			}
			require.Equal(t, expectedTeam.Name, obtainedTeam.Name)
			require.Equal(t, expectedTeam.Description, obtainedTeam.Description)
			require.Equal(t, expectedTeam.OriginCountry, obtainedTeam.OriginCountry)
		},
	)
}
