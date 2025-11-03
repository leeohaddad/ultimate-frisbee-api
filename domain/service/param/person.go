package param

import (
	"github.com/leeohaddad/ultimate-frisbee-api/domain/port/repository"
)

// TODO(lhaddad): turn entities into pointers

type GetAllPeople struct {
	Repository repository.Person
}

type GetPersonByUserName struct {
	UserName string

	Repository repository.Person
}
