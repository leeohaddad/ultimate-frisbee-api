package param

import (
	"github.com/leeohaddad/ultimate-frisbee-api/domain/port/repository"
)

type GetTeamGameCaptains struct {
	TeamSlug string

	MembershipRepository repository.Membership
	PersonRepository     repository.Person
}
