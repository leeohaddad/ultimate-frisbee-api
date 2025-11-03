package param

import (
	"github.com/leeohaddad/ultimate-frisbee-api/domain/port/repository"
)

type GetAllPeopleHandlerV1 struct {
	Repository repository.Person
}
