package repository

import (
	"context"

	"github.com/leeohaddad/ultimate-frisbee-api/domain/entity"
)

type Person interface {
	GetAllPeople(context context.Context) ([]*entity.Person, error)
	// GetPersonByUserName(context context.Context, ID string) (*entity.Person, error)
	// CreatePerson(context context.Context, person *entity.Person) (*entity.Person, error)
	// UpdatePerson(context context.Context, person *entity.Person, updatedAttributes []entity.PersonAttribute) (*entity.Person, error)
}
