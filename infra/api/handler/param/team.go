package param

import (
	"github.com/leeohaddad/ultimate-frisbee-api/domain/port/repository"
	"github.com/leeohaddad/ultimate-frisbee-api/infra/api/payload"
)

type GetAllTeamsHandlerV1 struct {
	Repository repository.Team
}

type GetTeamByNameHandlerV1 struct {
	Name string

	Repository repository.Team
}

type CreateTeamHandlerV1 struct {
	Payload payload.Team

	Repository repository.Team
}

type UpdateTeamHandlerV1 struct {
	Name    string
	Payload payload.Team

	Repository repository.Team
}
