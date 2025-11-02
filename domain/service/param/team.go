package param

import (
	"github.com/leeohaddad/ultimate-frisbee-api/domain/entity"
	"github.com/leeohaddad/ultimate-frisbee-api/domain/port/repository"
)

// TODO(lhaddad): turn entities into pointers

type GetAllTeams struct {
	Repository repository.Team
}

type GetTeamByName struct {
	Name string

	Repository repository.Team
}

type CreateTeam struct {
	Team *entity.Team

	Repository repository.Team
}

type UpdateTeam struct {
	Team              *entity.Team
	UpdatedAttributes []entity.TeamAttribute

	Repository repository.Team
}
