package test

import (
	"context"
	"strings"
	"testing"

	"github.com/leeohaddad/ultimate-frisbee-api/infra/test/fixture"

	"github.com/leeohaddad/ultimate-frisbee-api/infra/database/postgres"
	"github.com/stretchr/testify/require"
)

type FixtureScenario struct {
	Description    string
	FixtureQueries []fixture.Query
	InputData      map[string]interface{}
	OutputData     map[string]interface{}
}

func RunFixtureScenarios(t *testing.T, scenarios []FixtureScenario, scenarioValidator func(t *testing.T, testContext context.Context, client postgres.Client, scenario FixtureScenario)) {
	t.Helper()

	for _, scenario := range scenarios {
		t.Run(scenario.Description, func(t *testing.T) {
			client := GetDatabaseClient(t)
			defer client.Close()

			var err error
			testContext := context.Background()

			for _, fixtureQuery := range scenario.FixtureQueries {
				_, err = client.ExecuteQuery(testContext, nil, fixtureQuery.GetQuery(), fixtureQuery.GetParameters()...)
				if err != nil {
					if !strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
						require.NoError(t, err)
					}
				}
			}

			scenarioValidator(t, testContext, client, scenario)
		})
	}
}
