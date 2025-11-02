package postgres

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/leeohaddad/ultimate-frisbee-api/domain/entity"
	repositoryPort "github.com/leeohaddad/ultimate-frisbee-api/domain/port/repository"
	postgresDatabase "github.com/leeohaddad/ultimate-frisbee-api/infra/database/postgres"
)

// Enforce that TeamRepository implements the repositoryPort.Team interface.
var _ repositoryPort.Team = (*TeamRepository)(nil)

type TeamRepository struct {
	client postgresDatabase.Client
}

// team is a representation on how the team is retrieved from the database.
type team struct {
	Slug          string    `pg:"slug"`
	Name          string    `pg:"name"`
	Description   string    `pg:"description"`
	OriginCountry string    `pg:"origin_country"`
	CreatedBy     string    `pg:"created_by"`
	CreatedAt     time.Time `pg:"created_at"`
	UpdatedBy     string    `pg:"updated_by"`
	UpdatedAt     time.Time `pg:"updated_at"`
}

// NewTeamRepository instantiates a new team repository for postgres.
func NewTeamRepository(client postgresDatabase.Client) *TeamRepository {
	return &TeamRepository{
		client: client,
	}
}

func (repository *TeamRepository) GetAllTeams(context context.Context) ([]*entity.Team, error) {
	query := `select
              slug,
              name,
              description,
              origin_country,
              created_at,
              created_by,
              updated_at,
              updated_by
            from
              teams`

	// Execute query in DB
	var fetchedTeams []team
	queryResult, err := repository.client.ExecuteQuery(context, &fetchedTeams, query)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve all teams: %w", err)
	}

	// Query executed successfully but no entity found for this ID
	if queryResult.RowsReturned == 0 {
		return nil, nil
	}

	return teamsToTeamEntities(fetchedTeams), nil
}

func (repository *TeamRepository) GetTeamByName(context context.Context, name string) (*entity.Team, error) {
	query := `select
              slug,
              name,
              description,
              origin_country,
              created_at,
              created_by,
              updated_at,
              updated_by
            from
              teams
            where
              name = ? limit 1`

	// Execute query in DB
	var fetchedTeam team
	queryResult, err := repository.client.ExecuteQuery(context, &fetchedTeam, query, name)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve team %s: %w", name, err)
	}

	// Query executed successfully but no entity found for this ID
	if queryResult.RowsReturned == 0 {
		return nil, nil
	}

	return teamToTeamEntity(fetchedTeam), nil
}

func (repository *TeamRepository) CreateTeam(
	context context.Context,
	teamEntity *entity.Team,
) (*entity.Team, error) {
	// Convert entity to database representation
	teamToInsert := team{
		Slug:          teamEntity.Slug,
		Name:          teamEntity.Name,
		Description:   teamEntity.Description,
		OriginCountry: teamEntity.OriginCountry,
		CreatedBy:     teamEntity.CreatedBy,
		// CreatedAt/UpdatedAt are handled by the DB defaults on insert
		UpdatedBy: teamEntity.UpdatedBy,
	}
	// Insert and RETURNING to fetch the inserted row in one statement. This
	// avoids a separate SELECT and guarantees we obtain DB-defaulted columns.
	query := `insert into teams (
	 slug,
	 name,
	 description,
	 origin_country,
	 created_by,
	 updated_by
   ) values (?, ?, ?, ?, ?, ?) returning
	 slug,
	 name,
	 description,
	 origin_country,
	 created_at,
	 created_by,
	 updated_at,
	 updated_by`

	var inserted team
	queryResult, err := repository.client.ExecuteQuery(
		context,
		&inserted,
		query,
		teamToInsert.Slug,
		teamToInsert.Name,
		teamToInsert.Description,
		teamToInsert.OriginCountry,
		teamToInsert.CreatedBy,
		teamToInsert.UpdatedBy,
	)
	if err != nil {
		// If the error is a duplicate key on the team name, return a sentinel
		// error that callers can detect and convert to a 409 Conflict.
		if strings.Contains(err.Error(), "teams_name_key") || strings.Contains(err.Error(), "duplicate key value") {
			return nil, fmt.Errorf("%w: %v", repositoryPort.ErrAlreadyExists, err)
		}

		return nil, fmt.Errorf("failed to create team: %w", err)
	}
	if queryResult == nil || queryResult.RowsReturned == 0 {
		return nil, fmt.Errorf("no rows were returned after inserting team '%s'", teamEntity.Name)
	}

	return teamToTeamEntity(inserted), nil
}

func (repository *TeamRepository) UpdateTeam(
	context context.Context,
	teamEntity *entity.Team,
	updatedAttributes []entity.TeamAttribute,
) (*entity.Team, error) {
	// Build update query dynamically based on updatedAttributes
	setClauses := []string{}
	params := []interface{}{}
	for _, attr := range updatedAttributes {
		switch attr {
		case entity.TeamAttributes.Description:
			setClauses = append(setClauses, "description = ?")
			params = append(params, teamEntity.Description)
		case entity.TeamAttributes.OriginCountry:
			setClauses = append(setClauses, "origin_country = ?")
			params = append(params, teamEntity.OriginCountry)
		case entity.TeamAttributes.UpdatedBy:
			setClauses = append(setClauses, "updated_by = ?")
			params = append(params, teamEntity.UpdatedBy)
		}
	}
	// Always set updated_at to now()
	setClauses = append(setClauses, "updated_at = now()")
	query := "update teams set " + stringJoin(setClauses, ", ") + " where name = ?"
	params = append(params, teamEntity.Name)
	res, err := repository.client.ExecuteCommand(context, query, params...)
	if err != nil {
		return nil, fmt.Errorf("failed to update team: %w", err)
	}
	// If nothing was updated, return nil so handler can return 404
	if res == nil || res.RowsAffected == 0 {
		return nil, nil
	}
	// Return the updated team by fetching it back
	return repository.GetTeamByName(context, teamEntity.Name)
}

// stringJoin is a helper to join []string with a separator
func stringJoin(elems []string, sep string) string {
	if len(elems) == 0 {
		return ""
	}
	result := elems[0]
	for _, s := range elems[1:] {
		result += sep + s
	}
	return result
}

func teamsToTeamEntities(teams []team) []*entity.Team {
	teamEntities := make([]*entity.Team, 0)

	for _, team := range teams {
		teamEntities = append(teamEntities, teamToTeamEntity(team))
	}

	return teamEntities
}

func teamToTeamEntity(team team) *entity.Team {
	// Rows are scanned directly into Go types by the DB client. createdAt/updatedAt
	// are already time.Time so we can use them as-is.
	return &entity.Team{
		Slug:          team.Slug,
		Name:          team.Name,
		Description:   team.Description,
		OriginCountry: team.OriginCountry,
		CreatedAt:     team.CreatedAt,
		CreatedBy:     team.CreatedBy,
		UpdatedAt:     team.UpdatedAt,
		UpdatedBy:     team.UpdatedBy,
	}
}

/* func buildSortedParamListForUpdateTeamQuery(
	teamEntity *entity.Team,
	updatedAttributes []entity.TeamAttribute,
) []interface{} {
	params := make([]interface{}, 0)
	checkAttributeMap := map[entity.TeamAttribute]bool{}

	for _, attribute := range updatedAttributes {
		checkAttributeMap[attribute] = false
	}

	// Name attribute
	if _, ok := checkAttributeMap[entity.TeamAttributes.Name]; ok {
		params = append(params, teamEntity.Name)
	} else {
		params = append(params, nil)
	}

	// Description attribute
	if _, ok := checkAttributeMap[entity.TeamAttributes.Description]; ok {
		params = append(params, teamEntity.Description)
	} else {
		params = append(params, nil)
	}

	return append(params, teamEntity.ID)
} */
