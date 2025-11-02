package result

import (
	"github.com/leeohaddad/ultimate-frisbee-api/domain/entity"
)

type GetTeamGameCaptains struct {
	GameCaptains []*entity.Person
}
