package result

import (
	"github.com/leeohaddad/ultimate-frisbee-api/domain/entity"
)

type GetAllPeople struct {
	People []*entity.Person
}

type GetPersonByUserName struct {
	Person *entity.Person
}
