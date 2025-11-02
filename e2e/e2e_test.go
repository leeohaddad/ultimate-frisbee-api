package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	baseURL    = "http://127.0.0.1:42007"
	apiVersion = "v1"
)

// Team represents the API team payload structure
type Team struct {
	Slug          string  `json:"slug"`
	Name          string  `json:"name"`
	Description   *string `json:"description,omitempty"`
	OriginCountry *string `json:"originCountry,omitempty"`
	CreatedBy     *string `json:"createdBy,omitempty"`
	CreatedAt     *string `json:"createdAt,omitempty"`
	UpdatedBy     *string `json:"updatedBy,omitempty"`
	UpdatedAt     *string `json:"updatedAt,omitempty"`
}

// CreateTeamRequest represents the payload for creating a team
type CreateTeamRequest struct {
	Slug          string `json:"slug"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	OriginCountry string `json:"originCountry"`
	CreatedBy     string `json:"createdBy"`
}

// UpdateTeamRequest represents the payload for updating a team
type UpdateTeamRequest struct {
	Description   string `json:"description"`
	OriginCountry string `json:"originCountry"`
	UpdatedBy     string `json:"updatedBy"`
}

// Test helper functions
func makeHTTPRequest(method, url string, body interface{}) (*http.Response, error) {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	return client.Do(req)
}

func parseJSONResponse(resp *http.Response, target interface{}) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.Unmarshal(body, target)
}

func getStringResponse(resp *http.Response) (string, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	return string(body), nil
}

// End-to-End Test Suite
func TestUltimateFrisbeeAPI_E2E(t *testing.T) {
	t.Run("Health Check", testHealthCheck)
	t.Run("Team CRUD Operations", testTeamCRUDOperations)
	t.Run("Error Handling", testErrorHandling)
	t.Run("Data Validation", testDataValidation)
}

func testHealthCheck(t *testing.T) {
	t.Log("Testing Health Check endpoint...")

	resp, err := makeHTTPRequest("GET", fmt.Sprintf("%s/%s/health/", baseURL, apiVersion), nil)
	require.NoError(t, err, "Health check request should not fail")
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode, "Health check should return 200 OK")

	body, err := getStringResponse(resp)
	require.NoError(t, err, "Should be able to read health check response")

	t.Logf("Health check response: %s", body)
}

func testTeamCRUDOperations(t *testing.T) {
	testTeamName := "test-e2e-team"

	t.Run("Create Team", func(t *testing.T) {
		t.Log("Testing Create Team endpoint...")

		createRequest := CreateTeamRequest{
			Slug:          "test-e2e-team",
			Name:          testTeamName,
			Description:   "End-to-end test team for API testing",
			OriginCountry: "Canada",
			CreatedBy:     "e2e-test",
		}

		resp, err := makeHTTPRequest("POST", fmt.Sprintf("%s/%s/teams/", baseURL, apiVersion), createRequest)
		require.NoError(t, err, "Create team request should not fail")
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode, "Create team should return 201 Created")

		var createdTeam Team
		err = parseJSONResponse(resp, &createdTeam)
		require.NoError(t, err, "Should be able to parse created team response")

		assert.Equal(t, "test-e2e-team", createdTeam.Slug, "Created team should have correct slug")
		assert.Equal(t, testTeamName, createdTeam.Name, "Created team should have correct name")
		assert.NotNil(t, createdTeam.Description, "Created team should have description")
		assert.Equal(t, createRequest.Description, *createdTeam.Description, "Created team should have correct description")
		assert.NotNil(t, createdTeam.OriginCountry, "Created team should have origin country")
		assert.Equal(t, createRequest.OriginCountry, *createdTeam.OriginCountry, "Created team should have correct origin country")

		t.Logf("Successfully created team: %+v", createdTeam)
	})

	t.Run("Get All Teams", func(t *testing.T) {
		t.Log("Testing Get All Teams endpoint...")

		resp, err := makeHTTPRequest("GET", fmt.Sprintf("%s/%s/teams/", baseURL, apiVersion), nil)
		require.NoError(t, err, "Get all teams request should not fail")
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode, "Get all teams should return 200 OK")

		var teams []Team
		err = parseJSONResponse(resp, &teams)
		require.NoError(t, err, "Should be able to parse teams response")

		assert.GreaterOrEqual(t, len(teams), 4, "Should have at least 4 teams (3 seeded + 1 test team)")

		// Check if our test team is in the list
		var foundTestTeam bool
		for _, team := range teams {
			if team.Name == testTeamName {
				foundTestTeam = true
				assert.Equal(t, "test-e2e-team", team.Slug, "Test team should have correct slug")
				break
			}
		}
		assert.True(t, foundTestTeam, "Test team should be in the list of all teams")

		t.Logf("Successfully retrieved %d teams", len(teams))
	})

	t.Run("Get Team by Name", func(t *testing.T) {
		t.Log("Testing Get Team by Name endpoint...")

		resp, err := makeHTTPRequest("GET", fmt.Sprintf("%s/%s/teams/%s/", baseURL, apiVersion, testTeamName), nil)
		require.NoError(t, err, "Get team by name request should not fail")
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode, "Get team by name should return 200 OK")

		var team Team
		err = parseJSONResponse(resp, &team)
		require.NoError(t, err, "Should be able to parse team response")

		assert.Equal(t, "test-e2e-team", team.Slug, "Retrieved team should have correct slug")
		assert.Equal(t, testTeamName, team.Name, "Retrieved team should have correct name")
		assert.NotNil(t, team.Description, "Retrieved team should have description")
		assert.NotNil(t, team.OriginCountry, "Retrieved team should have origin country")

		t.Logf("Successfully retrieved team: %+v", team)
	})

	t.Run("Update Team", func(t *testing.T) {
		t.Log("Testing Update Team endpoint...")

		updateRequest := UpdateTeamRequest{
			Description:   "Updated description for end-to-end test team",
			OriginCountry: "United States",
			UpdatedBy:     "e2e-test-updater",
		}

		resp, err := makeHTTPRequest("PUT", fmt.Sprintf("%s/%s/teams/%s/", baseURL, apiVersion, testTeamName), updateRequest)
		require.NoError(t, err, "Update team request should not fail")
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode, "Update team should return 200 OK")

		var updatedTeam Team
		err = parseJSONResponse(resp, &updatedTeam)
		require.NoError(t, err, "Should be able to parse updated team response")

		assert.Equal(t, "test-e2e-team", updatedTeam.Slug, "Updated team should have correct slug")
		assert.Equal(t, testTeamName, updatedTeam.Name, "Updated team should have correct name")
		assert.NotNil(t, updatedTeam.Description, "Updated team should have description")
		assert.Equal(t, updateRequest.Description, *updatedTeam.Description, "Updated team should have new description")
		assert.NotNil(t, updatedTeam.OriginCountry, "Updated team should have origin country")
		assert.Equal(t, updateRequest.OriginCountry, *updatedTeam.OriginCountry, "Updated team should have new origin country")

		t.Logf("Successfully updated team: %+v", updatedTeam)
	})

	t.Run("Get Seeded Teams", func(t *testing.T) {
		t.Log("Testing that seeded teams are available...")

		seededTeams := []string{"Ultimate Warriors", "Disc Dynamos", "Flying Circus"}

		for _, teamName := range seededTeams {
			resp, err := makeHTTPRequest("GET", fmt.Sprintf("%s/%s/teams/%s/", baseURL, apiVersion, teamName), nil)
			require.NoError(t, err, "Get seeded team request should not fail")
			defer resp.Body.Close()

			assert.Equal(t, http.StatusOK, resp.StatusCode, "Get seeded team should return 200 OK")

			var team Team
			err = parseJSONResponse(resp, &team)
			require.NoError(t, err, "Should be able to parse seeded team response")

			assert.Equal(t, teamName, team.Name, "Seeded team should have correct name")
			assert.NotEmpty(t, team.Slug, "Seeded team should have slug")

			t.Logf("Successfully verified seeded team: %s", teamName)
		}
	})
}

func testErrorHandling(t *testing.T) {
	t.Run("Get Non-Existent Team", func(t *testing.T) {
		t.Log("Testing error handling for non-existent team...")

		resp, err := makeHTTPRequest("GET", fmt.Sprintf("%s/%s/teams/non-existent-team/", baseURL, apiVersion), nil)
		require.NoError(t, err, "Request should not fail")
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode, "Non-existent team should return 404 Not Found")

		body, err := getStringResponse(resp)
		require.NoError(t, err, "Should be able to read error response")

		assert.Contains(t, body, "no team with name 'non-existent-team' was found", "Error message should be descriptive")

		t.Logf("Correctly handled non-existent team error: %s", body)
	})

	t.Run("Update Non-Existent Team", func(t *testing.T) {
		t.Log("Testing error handling for updating non-existent team...")

		updateRequest := UpdateTeamRequest{
			Description:   "This should fail",
			OriginCountry: "Canada",
			UpdatedBy:     "test",
		}

		resp, err := makeHTTPRequest("PUT", fmt.Sprintf("%s/%s/teams/non-existent-team/", baseURL, apiVersion), updateRequest)
		require.NoError(t, err, "Request should not fail")
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode, "Update non-existent team should return 404 Not Found")

		body, err := getStringResponse(resp)
		require.NoError(t, err, "Should be able to read error response")

		assert.Contains(t, body, "no team with name 'non-existent-team' was found", "Error message should be descriptive")

		t.Logf("Correctly handled non-existent team update error: %s", body)
	})
}

func testDataValidation(t *testing.T) {
	t.Run("Create Team with Invalid Data", func(t *testing.T) {
		t.Log("Testing data validation for team creation...")

		// Test with missing required fields
		invalidRequest := map[string]string{
			"name": "", // Empty name should be invalid
		}

		resp, err := makeHTTPRequest("POST", fmt.Sprintf("%s/%s/teams/", baseURL, apiVersion), invalidRequest)
		require.NoError(t, err, "Request should not fail")
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "Invalid team data should return 400 Bad Request")

		body, err := getStringResponse(resp)
		require.NoError(t, err, "Should be able to read error response")

		t.Logf("Correctly handled validation error: %s", body)
	})

	t.Run("Create Team with Malformed JSON", func(t *testing.T) {
		t.Log("Testing malformed JSON handling...")

		// Send malformed JSON
		req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s/teams/", baseURL, apiVersion), bytes.NewBufferString(`{"invalid": json}`))
		require.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		require.NoError(t, err, "Request should not fail")
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "Malformed JSON should return 400 Bad Request")

		body, err := getStringResponse(resp)
		require.NoError(t, err, "Should be able to read error response")

		t.Logf("Correctly handled malformed JSON error: %s", body)
	})
}
