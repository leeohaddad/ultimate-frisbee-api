package service

import (
	"context"
	"fmt"

	domainServiceParam "github.com/leeohaddad/ultimate-frisbee-api/domain/service/param"
	domainServiceResult "github.com/leeohaddad/ultimate-frisbee-api/domain/service/result"
)

func GetAllTeams(
	context context.Context,
	param domainServiceParam.GetAllTeams,
) (domainServiceResult.GetAllTeams, error) {
	teams, err := param.Repository.GetAllTeams(context)
	if err != nil {
		return domainServiceResult.GetAllTeams{
			Teams: teams,
		}, fmt.Errorf("failed to fetch all teams from repository: %w", err)
	}

	return domainServiceResult.GetAllTeams{
		Teams: teams,
	}, nil
}

func GetTeamByName(
	context context.Context,
	param domainServiceParam.GetTeamByName,
) (domainServiceResult.GetTeamByName, error) {
	team, err := param.Repository.GetTeamByName(context, param.Name)
	if err != nil {
		return domainServiceResult.GetTeamByName{
			Team: team,
		}, fmt.Errorf("failed to fetch team by name '%s' from repository: %w", param.Name, err)
	}

	return domainServiceResult.GetTeamByName{
		Team: team,
	}, nil
}

func CreateTeam(
	context context.Context,
	param domainServiceParam.CreateTeam,
) (domainServiceResult.CreateTeam, error) {
	team, err := param.Repository.CreateTeam(context, param.Team)
	if err != nil {
		return domainServiceResult.CreateTeam{
			Team: team,
		}, fmt.Errorf("failed to create team with name '%s' in repository: %w", param.Team.Name, err)
	}

	return domainServiceResult.CreateTeam{
		Team: team,
	}, nil
}

func UpdateTeam(
	context context.Context,
	param domainServiceParam.UpdateTeam,
) (domainServiceResult.UpdateTeam, error) {
	team, err := param.Repository.UpdateTeam(context, param.Team, param.UpdatedAttributes)
	if err != nil {
		return domainServiceResult.UpdateTeam{
			Team: team,
		}, fmt.Errorf("failed to update team with name '%s' in repository: %w", param.Team.Name, err)
	}

	return domainServiceResult.UpdateTeam{
		Team: team,
	}, nil
}
