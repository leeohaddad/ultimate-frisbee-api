package result

import (
	"github.com/leeohaddad/ultimate-frisbee-api/domain/entity"
)

type GetTeamMembershipsByRole struct {
	Memberships []entity.Membership
}
