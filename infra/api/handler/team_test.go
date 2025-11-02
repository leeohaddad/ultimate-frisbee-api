//go:build integration
// +build integration

package handler_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/leeohaddad/ultimate-frisbee-api/domain/entity"
	repositoryPostgres "github.com/leeohaddad/ultimate-frisbee-api/infra/adapter/repository/postgres"
	"github.com/leeohaddad/ultimate-frisbee-api/infra/api/handler"
	handlerParam "github.com/leeohaddad/ultimate-frisbee-api/infra/api/handler/param"
	handlerResult "github.com/leeohaddad/ultimate-frisbee-api/infra/api/handler/result"
	"github.com/leeohaddad/ultimate-frisbee-api/infra/api/payload"
	databasePostgres "github.com/leeohaddad/ultimate-frisbee-api/infra/database/postgres"
	"github.com/leeohaddad/ultimate-frisbee-api/infra/test"
	"github.com/leeohaddad/ultimate-frisbee-api/infra/test/fixture"
)

// TODO(lhaddad): remove the fixture functions from here also and use the ones from package test.fixture instead.

func GetDefaultFixtureTeam(t *testing.T) *entity.Team {
	t.Helper()

	return fixture.GetFakeTeam().
		WithSlug("bra-sp-some-test-team").
		WithName("Some Test Team").
		WithDescription("A test team.").
		WithOriginCountry("BRA").
		WithCreatedBy("Someone who created the test team").
		WithUpdatedBy("Someone who updated the test team")
}

func GetAnotherFixtureTeam(t *testing.T) *entity.Team {
	t.Helper()

	return fixture.GetFakeTeam().
		WithSlug("another-test-team").
		WithName("Another Test Team").
		WithDescription("Another test team.").
		WithOriginCountry("BRA").
		WithCreatedBy("Someone who created the other test team").
		WithUpdatedBy("Someone who updated the other test team")
}

func TestTeamHandler_GetAllTeams(t *testing.T) {
	t.Parallel()

	scenarios := []test.FixtureScenario{
		{
			Description:    "should return no team when database is clean",
			FixtureQueries: []fixture.Query{},
			InputData:      map[string]interface{}{},
			OutputData: map[string]interface{}{
				"expectedStatusCode":     http.StatusOK,
				"expectedResponseType":   handlerResult.ResponseBodyTypes.JSON,
				"expectedJSONResponse":   []payload.Team{},
				"expectedStringResponse": "",
			},
		},
		{
			Description:    "should return one team when database is filled with one team",
			FixtureQueries: fixture.GenerateTeamQueries(GetDefaultFixtureTeam(t)),
			InputData:      map[string]interface{}{},
			OutputData: map[string]interface{}{
				"expectedStatusCode":   http.StatusOK,
				"expectedResponseType": handlerResult.ResponseBodyTypes.JSON,
				"expectedJSONResponse": []payload.Team{
					payload.TeamEntityToTeam(GetDefaultFixtureTeam(t)),
				},
				"expectedStringResponse": "",
			},
		},
		{
			Description: "should return two teams when database is filled with two teams",
			FixtureQueries: fixture.GenerateTeamQueries(
				GetDefaultFixtureTeam(t),
				GetAnotherFixtureTeam(t),
			),
			InputData: map[string]interface{}{},
			OutputData: map[string]interface{}{
				"expectedStatusCode":   http.StatusOK,
				"expectedResponseType": handlerResult.ResponseBodyTypes.JSON,
				"expectedJSONResponse": []payload.Team{
					payload.TeamEntityToTeam(GetDefaultFixtureTeam(t)),
					payload.TeamEntityToTeam(GetAnotherFixtureTeam(t)),
				},
				"expectedStringResponse": "",
			},
		},
	}

	test.RunFixtureScenarios(
		t, scenarios,
		func(t *testing.T, testContext context.Context, client databasePostgres.Client, scenario test.FixtureScenario) {
			t.Helper()

			expectedStatusCode, ok := scenario.OutputData["expectedStatusCode"].(int)
			require.True(t, ok)
			expectedResponseType, ok := scenario.OutputData["expectedResponseType"].(handlerResult.ResponseBodyType)
			require.True(t, ok)
			expectedTeams, ok := scenario.OutputData["expectedJSONResponse"].([]payload.Team)
			require.True(t, ok)
			expectedMessage, ok := scenario.OutputData["expectedStringResponse"].(string)
			require.True(t, ok)

			teamRepository := repositoryPostgres.NewTeamRepository(client)

			result := handler.GetAllTeamsHandlerV1(testContext, handlerParam.GetAllTeamsHandlerV1{
				Repository: teamRepository,
			})

			switch result.ResponseType {
			case handlerResult.ResponseBodyTypes.JSON:
				obtainedTeams, ok := result.JSONResponse.([]payload.Team)
				require.True(t, ok)
				require.Len(t, obtainedTeams, len(expectedTeams))

				indexedObtainedTeams := map[string]payload.Team{}
				for _, obtainedTeam := range obtainedTeams {
					indexedObtainedTeams[obtainedTeam.Slug] = obtainedTeam
				}
				for _, expectedTeam := range expectedTeams {
					obtainedTeam, hasKey := indexedObtainedTeams[expectedTeam.Slug]
					require.True(t, hasKey)
					require.Equal(t, expectedTeam.Name, obtainedTeam.Name)
					require.Equal(t, expectedTeam.Description, obtainedTeam.Description)
					require.Equal(t, expectedTeam.OriginCountry, obtainedTeam.OriginCountry)
					// require.Equal(t, expectedTeam.OriginCountryID, obtainedTeam.OriginCountryID)
					require.Equal(t, *expectedTeam.CreatedBy, *obtainedTeam.CreatedBy)
					require.NotNil(t, obtainedTeam.CreatedAt)
					require.Equal(t, expectedTeam.UpdatedBy, obtainedTeam.UpdatedBy)
					require.NotNil(t, obtainedTeam.UpdatedAt)
				}
			case handlerResult.ResponseBodyTypes.String:
				require.Contains(t, result.StringResponse, expectedMessage)
			}
			require.Equal(t, expectedResponseType, result.ResponseType)
			require.Equal(t, expectedStatusCode, result.StatusCode)
		},
	)
}

func TestTeamHandler_GetTeamByName(t *testing.T) {
	t.Parallel()

	scenarios := []test.FixtureScenario{
		{
			Description:    "should return no team when the database has no team with the specified name",
			FixtureQueries: []fixture.Query{},
			InputData: map[string]interface{}{
				"Name": "unknown-team",
			},
			OutputData: map[string]interface{}{
				"expectedStatusCode":     http.StatusNotFound,
				"expectedResponseType":   handlerResult.ResponseBodyTypes.String,
				"expectedJSONResponse":   payload.Team{},
				"expectedStringResponse": "no team with name 'unknown-team' was found in the repository",
			},
		},
		{
			Description: "should return no team when the database has teams but only with other name",
			FixtureQueries: fixture.GenerateTeamQueries(
				GetDefaultFixtureTeam(t),
			),
			InputData: map[string]interface{}{
				"Name": "unknown-team",
			},
			OutputData: map[string]interface{}{
				"expectedStatusCode":     http.StatusNotFound,
				"expectedResponseType":   handlerResult.ResponseBodyTypes.String,
				"expectedJSONResponse":   payload.Team{},
				"expectedStringResponse": "no team with name 'unknown-team' was found in the repository",
			},
		},
		{
			Description: "should return the desired team when the database has a team with the specified name",
			FixtureQueries: fixture.GenerateTeamQueries(
				GetDefaultFixtureTeam(t),
			),
			InputData: map[string]interface{}{
				"Name": GetDefaultFixtureTeam(t).Name,
			},
			OutputData: map[string]interface{}{
				"expectedStatusCode":     http.StatusOK,
				"expectedResponseType":   handlerResult.ResponseBodyTypes.JSON,
				"expectedJSONResponse":   payload.TeamEntityToTeam(GetDefaultFixtureTeam(t)),
				"expectedStringResponse": "",
			},
		},
	}

	test.RunFixtureScenarios(
		t, scenarios,
		func(t *testing.T, testContext context.Context, client databasePostgres.Client, scenario test.FixtureScenario) {
			t.Helper()

			name, ok := scenario.InputData["Name"].(string)
			require.True(t, ok)
			expectedStatusCode, ok := scenario.OutputData["expectedStatusCode"].(int)
			require.True(t, ok)
			expectedResponseType, ok := scenario.OutputData["expectedResponseType"].(handlerResult.ResponseBodyType)
			require.True(t, ok)
			expectedTeam, ok := scenario.OutputData["expectedJSONResponse"].(payload.Team)
			require.True(t, ok)
			expectedMessage, ok := scenario.OutputData["expectedStringResponse"].(string)
			require.True(t, ok)

			teamRepository := repositoryPostgres.NewTeamRepository(client)

			result := handler.GetTeamByNameHandlerV1(testContext, handlerParam.GetTeamByNameHandlerV1{
				Name:       name,
				Repository: teamRepository,
			})

			switch result.ResponseType {
			case handlerResult.ResponseBodyTypes.JSON:
				obtainedTeam, ok := result.JSONResponse.(payload.Team)
				require.True(t, ok)
				require.Equal(t, expectedTeam.Name, obtainedTeam.Name)
				require.Equal(t, expectedTeam.Description, obtainedTeam.Description)
				require.Equal(t, expectedTeam.OriginCountry, obtainedTeam.OriginCountry)
				// require.Equal(t, expectedTeam.OriginCountryID, obtainedTeam.OriginCountryID)
				require.Equal(t, expectedTeam.CreatedBy, obtainedTeam.CreatedBy)
				require.NotNil(t, obtainedTeam.CreatedAt)
				require.Equal(t, expectedTeam.UpdatedBy, obtainedTeam.UpdatedBy)
				require.NotNil(t, obtainedTeam.UpdatedAt)
			case handlerResult.ResponseBodyTypes.String:
				require.Contains(t, result.StringResponse, expectedMessage)
			}
			require.Equal(t, expectedResponseType, result.ResponseType)
			require.Equal(t, expectedStatusCode, result.StatusCode)
		},
	)
}

func TestTeamHandler_CreateTeam(t *testing.T) {
	t.Parallel()

	scenarios := []test.FixtureScenario{
		{
			Description:    "should create a team and return it correctly",
			FixtureQueries: []fixture.Query{},
			InputData: map[string]interface{}{
				"team": payload.TeamEntityToTeam(GetDefaultFixtureTeam(t)),
			},
			OutputData: map[string]interface{}{
				"expectedStatusCode":     http.StatusCreated,
				"expectedResponseType":   handlerResult.ResponseBodyTypes.JSON,
				"expectedJSONResponse":   payload.TeamEntityToTeam(GetDefaultFixtureTeam(t)),
				"expectedStringResponse": "",
			},
		},
		// TODO(lhaddad): change this behavior to return HTTP 200 instead of HTTP 500 if it already exists.
		{
			Description:    "should fail to create a team if it already exists",
			FixtureQueries: fixture.GenerateTeamQueries(GetDefaultFixtureTeam(t)),
			InputData: map[string]interface{}{
				"team": payload.TeamEntityToTeam(GetDefaultFixtureTeam(t)),
			},
			OutputData: map[string]interface{}{
				"expectedStatusCode":     http.StatusInternalServerError,
				"expectedResponseType":   handlerResult.ResponseBodyTypes.String,
				"expectedJSONResponse":   payload.TeamEntityToTeam(GetDefaultFixtureTeam(t)),
				"expectedStringResponse": fmt.Sprintf("failed to create team with name '%s' in application service", GetDefaultFixtureTeam(t).Name),
			},
		},
		{
			Description:    "should fail to create a team if another team with the same slug already exists",
			FixtureQueries: fixture.GenerateTeamQueries(GetAnotherFixtureTeam(t).WithSlug(GetDefaultFixtureTeam(t).Slug)),
			InputData: map[string]interface{}{
				"team": payload.TeamEntityToTeam(GetDefaultFixtureTeam(t)),
			},
			OutputData: map[string]interface{}{
				"expectedStatusCode":     http.StatusInternalServerError,
				"expectedResponseType":   handlerResult.ResponseBodyTypes.String,
				"expectedJSONResponse":   payload.Team{},
				"expectedStringResponse": fmt.Sprintf("failed to create team with name '%s' in application service", GetDefaultFixtureTeam(t).Name),
			},
		},
		{
			Description:    "should fail to create a team if the payload team name is empty",
			FixtureQueries: []fixture.Query{},
			InputData: map[string]interface{}{
				"team": payload.TeamEntityToTeam(GetDefaultFixtureTeam(t).WithName("")),
			},
			OutputData: map[string]interface{}{
				"expectedStatusCode":     http.StatusBadRequest,
				"expectedResponseType":   handlerResult.ResponseBodyTypes.String,
				"expectedJSONResponse":   payload.Team{},
				"expectedStringResponse": "the team's 'name' should not be empty",
			},
		},
		{
			Description:    "should fail to create a team if the payload team description is empty",
			FixtureQueries: []fixture.Query{},
			InputData: map[string]interface{}{
				"team": payload.TeamEntityToTeam(GetDefaultFixtureTeam(t).WithDescription("")),
			},
			OutputData: map[string]interface{}{
				"expectedStatusCode":     http.StatusBadRequest,
				"expectedResponseType":   handlerResult.ResponseBodyTypes.String,
				"expectedJSONResponse":   payload.Team{},
				"expectedStringResponse": "the team 'Description' should not be empty",
			},
		},
	}

	test.RunFixtureScenarios(
		t, scenarios,
		func(t *testing.T, testContext context.Context, client databasePostgres.Client, scenario test.FixtureScenario) {
			t.Helper()

			team, ok := scenario.InputData["team"].(payload.Team)
			require.True(t, ok)
			expectedStatusCode, ok := scenario.OutputData["expectedStatusCode"].(int)
			require.True(t, ok)
			expectedResponseType, ok := scenario.OutputData["expectedResponseType"].(handlerResult.ResponseBodyType)
			require.True(t, ok)
			expectedTeam, ok := scenario.OutputData["expectedJSONResponse"].(payload.Team)
			require.True(t, ok)
			expectedMessage, ok := scenario.OutputData["expectedStringResponse"].(string)
			require.True(t, ok)

			teamRepository := repositoryPostgres.NewTeamRepository(client)

			result := handler.CreateTeamHandlerV1(testContext, handlerParam.CreateTeamHandlerV1{
				Payload:    team,
				Repository: teamRepository,
			})

			switch result.ResponseType {
			case handlerResult.ResponseBodyTypes.JSON:
				obtainedTeam, ok := result.JSONResponse.(payload.Team)
				require.True(t, ok)
				require.Equal(t, expectedTeam.Name, obtainedTeam.Name)
				require.Equal(t, expectedTeam.Description, obtainedTeam.Description)
				require.Equal(t, expectedTeam.OriginCountry, obtainedTeam.OriginCountry)
				// require.Equal(t, expectedTeam.OriginCountryID, obtainedTeam.OriginCountryID)
				require.Equal(t, expectedTeam.CreatedBy, obtainedTeam.CreatedBy)
				require.NotNil(t, obtainedTeam.CreatedAt)
				require.Equal(t, expectedTeam.UpdatedBy, obtainedTeam.UpdatedBy)
				require.NotNil(t, obtainedTeam.UpdatedAt)
			case handlerResult.ResponseBodyTypes.String:
				require.Contains(t, result.StringResponse, expectedMessage)
			}
			require.Equal(t, expectedResponseType, result.ResponseType)
			require.Equal(t, expectedStatusCode, result.StatusCode)
		},
	)
}

func TestTeamHandler_UpdateTeam(t *testing.T) {
	t.Parallel()

	scenarios := []test.FixtureScenario{
		{
			Description:    "should update all attributes correctly and return the entity",
			FixtureQueries: fixture.GenerateTeamQueries(GetDefaultFixtureTeam(t)),
			InputData: map[string]interface{}{
				"Name": GetDefaultFixtureTeam(t).Name,
				"team": payload.TeamEntityToTeam(GetDefaultFixtureTeam(t).
					WithDescription(GetAnotherFixtureTeam(t).Description),
				),
			},
			OutputData: map[string]interface{}{
				"expectedStatusCode":   http.StatusOK,
				"expectedResponseType": handlerResult.ResponseBodyTypes.JSON,
				"expectedJSONResponse": payload.TeamEntityToTeam(GetDefaultFixtureTeam(t).
					WithDescription(GetAnotherFixtureTeam(t).Description),
				),
				"expectedStringResponse": "",
			},
		},
		{
			Description:    "should update nothing because ID does not exist on database",
			FixtureQueries: fixture.GenerateTeamQueries(GetDefaultFixtureTeam(t)),
			InputData: map[string]interface{}{
				"Name": "unknown-team",
				"team": payload.TeamEntityToTeam(GetDefaultFixtureTeam(t).
					WithName("unknown-team").
					WithDescription(GetAnotherFixtureTeam(t).Description),
				),
			},
			OutputData: map[string]interface{}{
				"expectedStatusCode":     http.StatusNotFound,
				"expectedResponseType":   handlerResult.ResponseBodyTypes.String,
				"expectedJSONResponse":   payload.Team{},
				"expectedStringResponse": "no team with name 'unknown-team' was found in the repository",
			},
		},
		{
			Description:    "should fail to update team if trying to update the team name",
			FixtureQueries: fixture.GenerateTeamQueries(GetDefaultFixtureTeam(t)),
			InputData: map[string]interface{}{
				"Name": GetDefaultFixtureTeam(t).Name,
				"team": payload.TeamEntityToTeam(GetDefaultFixtureTeam(t).
					WithDescription(GetAnotherFixtureTeam(t).Description),
				),
			},
			OutputData: map[string]interface{}{
				"expectedStatusCode":     http.StatusBadRequest,
				"expectedResponseType":   handlerResult.ResponseBodyTypes.String,
				"expectedJSONResponse":   payload.Team{},
				"expectedStringResponse": "updating the team name is not allowed",
			},
		},
		{
			Description:    "should fail to update team if the query string team name is empty",
			FixtureQueries: fixture.GenerateTeamQueries(GetDefaultFixtureTeam(t)),
			InputData: map[string]interface{}{
				"Name": "",
				"team": payload.TeamEntityToTeam(GetDefaultFixtureTeam(t).
					WithDescription(GetAnotherFixtureTeam(t).Description),
				),
			},
			OutputData: map[string]interface{}{
				"expectedStatusCode":     http.StatusBadRequest,
				"expectedResponseType":   handlerResult.ResponseBodyTypes.String,
				"expectedJSONResponse":   payload.Team{},
				"expectedStringResponse": "team name defined in the path variable is not valid",
			},
		},
		{
			Description:    "should fail to update team if all the updatable fields are empty",
			FixtureQueries: fixture.GenerateTeamQueries(GetDefaultFixtureTeam(t)),
			InputData: map[string]interface{}{
				"Name": GetDefaultFixtureTeam(t).Name,
				"team": payload.TeamEntityToTeam(GetDefaultFixtureTeam(t).
					WithDescription("").
					WithOriginCountry("").
					// WithOriginCountry(&entity.Country{ID: ""}).
					WithUpdatedBy(""),
				),
			},
			OutputData: map[string]interface{}{
				"expectedStatusCode":     http.StatusBadRequest,
				"expectedResponseType":   handlerResult.ResponseBodyTypes.String,
				"expectedJSONResponse":   payload.Team{},
				"expectedStringResponse": "at least one of the following fields should not be empty: [Description, OriginCountry, UpdatedBy]",
			},
		},
	}

	type entryWithUpdatedAt struct {
		UpdatedAt time.Time
	}

	test.RunFixtureScenarios(
		t, scenarios,
		func(t *testing.T, testContext context.Context, client databasePostgres.Client, scenario test.FixtureScenario) {
			t.Helper()

			name, ok := scenario.InputData["Name"].(string)
			require.True(t, ok)
			team, ok := scenario.InputData["team"].(payload.Team)
			require.True(t, ok)
			expectedStatusCode, ok := scenario.OutputData["expectedStatusCode"].(int)
			require.True(t, ok)
			expectedResponseType, ok := scenario.OutputData["expectedResponseType"].(handlerResult.ResponseBodyType)
			require.True(t, ok)
			expectedTeam, ok := scenario.OutputData["expectedJSONResponse"].(payload.Team)
			require.True(t, ok)
			expectedMessage, ok := scenario.OutputData["expectedStringResponse"].(string)
			require.True(t, ok)

			teamRepository := repositoryPostgres.NewTeamRepository(client)

			var updateAtBeforeSave entryWithUpdatedAt
			_, err := client.ExecuteQuery(testContext, &updateAtBeforeSave, "select updated_at from teams where name = ?", team.Name)
			require.NoError(t, err)
			time.Sleep(50 * time.Millisecond)

			result := handler.UpdateTeamHandlerV1(testContext, handlerParam.UpdateTeamHandlerV1{
				Repository: teamRepository,
				Name:       name,
				Payload:    team,
			})

			var updateAtAfterSave entryWithUpdatedAt
			_, err = client.ExecuteQuery(testContext, &updateAtAfterSave, "select updated_at from teams where name = ?", team.Name)
			require.NoError(t, err)

			switch result.ResponseType {
			case handlerResult.ResponseBodyTypes.JSON:
				obtainedTeam, ok := result.JSONResponse.(payload.Team)
				require.True(t, ok)
				require.Equal(t, expectedTeam.Name, obtainedTeam.Name)
				require.Equal(t, expectedTeam.Description, obtainedTeam.Description)
				require.Equal(t, expectedTeam.OriginCountry, obtainedTeam.OriginCountry)
				// require.Equal(t, expectedTeam.OriginCountryID, obtainedTeam.OriginCountryID)
				require.Equal(t, expectedTeam.CreatedBy, obtainedTeam.CreatedBy)
				require.NotNil(t, obtainedTeam.CreatedAt)
				require.Equal(t, expectedTeam.UpdatedBy, obtainedTeam.UpdatedBy)
				require.NotNil(t, obtainedTeam.UpdatedAt)
				require.True(t, updateAtAfterSave.UpdatedAt.After(updateAtBeforeSave.UpdatedAt),
					"%s should be greater than %s", updateAtAfterSave.UpdatedAt, updateAtBeforeSave.UpdatedAt)
			case handlerResult.ResponseBodyTypes.String:
				require.Contains(t, result.StringResponse, expectedMessage)
			}
			require.Equal(t, expectedResponseType, result.ResponseType)
			require.Equal(t, expectedStatusCode, result.StatusCode)
		},
	)
}
