package fixture

import (
	"github.com/leeohaddad/ultimate-frisbee-api/domain/entity"
)

const (
	// FakeTeamDefaultSlug is the default slug for a fake team.
	FakeTeamDefaultSlug = "bra-sp-my-team-slug"
	// FakeTeamDefaultName is the default name for a fake team.
	FakeTeamDefaultName = "My Team Name"
	// FakeTeamDefaultDescription is the default description for a fake team.
	FakeTeamDefaultDescription = "This is my team description."
	// FakeTeamDefaultOriginCountry is the default origin country for a fake team.
	FakeTeamDefaultOriginCountry = "BRA"

	// FakeTeamAnotherSlug is another name for a fake slug.
	FakeTeamAnotherSlug = "bra-sp-another-team-slug"
	// FakeTeamAnotherName is another name for a fake team.
	FakeTeamAnotherName = "Another Team Name"
)

func GetFakeTeam() *entity.Team {
	return &entity.Team{
		Slug:          FakeTeamDefaultSlug,
		Name:          FakeTeamDefaultName,
		Description:   FakeTeamDefaultDescription,
		OriginCountry: FakeTeamDefaultOriginCountry,
	}
}

func GenerateTeamQueries(teams ...*entity.Team) []Query {
	queries := make([]Query, 0)

	for _, team := range teams {
		if team == nil {
			continue
		}
		queries = append(queries, GenerateCustomQuery(
			"insert into teams(slug, name, description, origin_country, created_by, updated_by) values (?, ?, ?, ?, ?, ?)",
			team.Slug, team.Name, team.Description, team.OriginCountry, team.CreatedBy, team.UpdatedBy,
		))
	}

	return queries
}

func GetDefaultFixtureTeam() *entity.Team {
	team := GetFakeTeam()

	team.Slug = FakeTeamDefaultSlug
	team.Name = FakeTeamDefaultName
	team.Description = "My Awesome Team"
	team.OriginCountry = "Brazil"

	return team
}

func GetAnotherFixtureTeam() *entity.Team {
	team := GetFakeTeam()

	team.Name = FakeTeamAnotherName
	team.Description = "Another Not So Awesome Team"
	team.OriginCountry = "Brazil"

	return team
}

func GetDefaultFixtureTeamUpdated() *entity.Team {
	team := GetDefaultFixtureTeam()

	team.Name = "My Updated Team"
	team.Description = "My Awesome Updated Team"

	return team
}
