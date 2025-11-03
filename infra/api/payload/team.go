package payload

import (
	"time"

	"github.com/leeohaddad/ultimate-frisbee-api/domain/entity"
	"github.com/leeohaddad/ultimate-frisbee-api/infra/helper"
)

type Team struct {
	Slug          string  `json:"slug"`
	Name          string  `json:"name"`
	Description   *string `json:"description"`
	OriginCountry *string `json:"originCountry"`

	CreatedBy *string `json:"createdBy"`
	CreatedAt *string `json:"createdAt"`
	UpdatedBy *string `json:"updatedBy"`
	UpdatedAt *string `json:"updatedAt"`
}

func ValidateCreateTeamInput(team *Team) (bool, string) {
	currentEntity := "Team"

	if helper.IsNilOrEmpty(&team.Name) {
		return false, helper.ErrorMessageInField(currentEntity, "Name")
	}

	if helper.IsNilOrEmpty(team.Description) {
		return false, helper.ErrorMessageInField(currentEntity, "Description")
	}

	if helper.IsNilOrEmpty(team.OriginCountry) {
		return false, helper.ErrorMessageInField(currentEntity, "Origin Country")
	}

	if helper.IsNilOrEmpty(team.CreatedBy) {
		return false, helper.ErrorMessageInField(currentEntity, "Created By")
	}

	return true, ""
}

func ValidateUpdateTeamInput(team *Team, name string) (bool, string) {
	if name == "" {
		return false, "team name defined in the path variable is empty"
	}

	if team.Name != "" && team.Name != name {
		return false, "updating the team name is not allowed"
	}

	if (team.Description == nil || *team.Description == "") &&
		(team.OriginCountry == nil || *team.OriginCountry == "") &&
		(team.UpdatedBy == nil || *team.UpdatedBy == "") {
		return false, "at least one of the following fields should not be empty: [Description, OriginCountry, UpdatedBy]"
	}

	return true, ""
}

func GetFilledTeamAttributesForUpdate(team *Team) []entity.TeamAttribute {
	var attributes []entity.TeamAttribute

	if team.Description != nil {
		attributes = append(attributes, entity.TeamAttributes.Description)
	}

	if team.OriginCountry != nil {
		attributes = append(attributes, entity.TeamAttributes.OriginCountry)
	}

	if team.UpdatedBy != nil {
		attributes = append(attributes, entity.TeamAttributes.UpdatedBy)
	}

	return attributes
}

func TeamToTeamEntity(team Team) *entity.Team {
	var description string
	if team.Description != nil {
		description = *team.Description
	}

	var originCountry string
	if team.OriginCountry != nil {
		originCountry = *team.OriginCountry
	}

	var createdBy string
	if team.CreatedBy != nil {
		createdBy = *team.CreatedBy
	}

	var createdAt time.Time
	if team.CreatedAt != nil {
		var err error
		createdAt, err = time.Parse(helper.DefaultTimeLayout, *team.CreatedAt)
		if err != nil {
			createdAt = time.Time{}
		}
	}

	var updatedBy string
	if team.UpdatedBy != nil {
		updatedBy = *team.UpdatedBy
	}

	var updatedAt time.Time
	if team.UpdatedAt != nil {
		var err error
		updatedAt, err = time.Parse(helper.DefaultTimeLayout, *team.UpdatedAt)
		if err != nil {
			updatedAt = time.Time{}
		}
	}

	return &entity.Team{
		Slug:          team.Slug,
		Name:          team.Name,
		Description:   description,
		OriginCountry: originCountry,
		CreatedBy:     createdBy,
		CreatedAt:     createdAt,
		UpdatedBy:     updatedBy,
		UpdatedAt:     updatedAt,
	}
}

func TeamEntityToTeam(teamEntity *entity.Team) Team {
	createdAt := teamEntity.CreatedAt.Format(helper.DefaultTimeLayout)
	updatedAt := teamEntity.UpdatedAt.Format(helper.DefaultTimeLayout)
	return Team{
		Slug:          teamEntity.Slug,
		Name:          teamEntity.Name,
		Description:   &teamEntity.Description,
		OriginCountry: &teamEntity.OriginCountry,
		CreatedBy:     &teamEntity.CreatedBy,
		CreatedAt:     &createdAt,
		UpdatedBy:     &teamEntity.UpdatedBy,
		UpdatedAt:     &updatedAt,
	}
}

func TeamEntitiesToTeams(teamEntities []*entity.Team) []Team {
	teams := make([]Team, 0)

	for _, teamEntity := range teamEntities {
		teams = append(teams, TeamEntityToTeam(teamEntity))
	}

	if teams == nil {
		return []Team{}
	}

	return teams
}
