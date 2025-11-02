package result

import (
	"github.com/leeohaddad/ultimate-frisbee-api/domain/entity"
)

type GetAllTeams struct {
	Teams []*entity.Team
}

type GetTeamByName struct {
	Team *entity.Team
}

type CreateTeam struct {
	Team *entity.Team
}

type UpdateTeam struct {
	Team *entity.Team
}
