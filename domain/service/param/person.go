package param

import (
	"github.com/leeohaddad/ultimate-frisbee-api/domain/port/repository"
)

type GetPersonByUserName struct {
	UserName string

	Repository repository.Person
}
