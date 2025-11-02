package param

import (
	"github.com/leeohaddad/ultimate-frisbee-api/domain/port/repository"
)

type GetTeamMembershipsByRole struct {
	TeamSlug string
	Role     string

	Repository repository.Membership
}
