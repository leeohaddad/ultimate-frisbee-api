package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/leeohaddad/ultimate-frisbee-api/domain/entity"
	repositoryPort "github.com/leeohaddad/ultimate-frisbee-api/domain/port/repository"
	postgresDatabase "github.com/leeohaddad/ultimate-frisbee-api/infra/database/postgres"
)

// Enforce that PersonRepository implements the repositoryPort.Person interface.
var _ repositoryPort.Person = (*PersonRepository)(nil)

type PersonRepository struct {
	client postgresDatabase.Client
}

// person is a representation on how the person is retrieved from the database.
type person struct {
	ID            string `pg:"id"`
	Name          string `pg:"name"`
	Email         string `pg:"email"`
	PhoneNumber   string `pg:"phone_number"`
	WFDFNumber    string `pg:"wfdf_number"`
	OriginCountry string `pg:"origin_country"`

	CreatedAt time.Time `pg:"created_at"`
	CreatedBy string    `pg:"created_by"`
	UpdatedAt time.Time `pg:"updated_at"`
	UpdatedBy string    `pg:"updated_by"`
}

// NewPersonRepository instantiates a new person repository for postgres.
func NewPersonRepository(client postgresDatabase.Client) *PersonRepository {
	return &PersonRepository{
		client: client,
	}
}

func (repository *PersonRepository) GetAllPeople(context context.Context) ([]*entity.Person, error) {
	query := `select
              user_name,
			  name,
              email,
			  phone_number,
			  wfdf_number,
			  origin_country,

			  created_by,
              created_at,
              updated_at,
              updated_by
            from
              people`

	// Execute query in DB
	var fetchedPeople []person
	queryResult, err := repository.client.ExecuteQuery(context, &fetchedPeople, query)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve all people: %w", err)
	}

	// Query executed successfully but no entity found for this ID
	if queryResult.RowsReturned == 0 {
		return nil, nil
	}

	return peopleToPersonEntities(fetchedPeople), nil
}

func peopleToPersonEntities(people []person) []*entity.Person {
	personEntities := make([]*entity.Person, 0)

	for _, person := range people {
		personEntities = append(personEntities, personToPersonEntity(person))
	}

	return personEntities
}

func personToPersonEntity(person person) *entity.Person {
	// Rows are scanned directly into Go types by the DB client. createdAt/updatedAt
	// are already time.Time so we can use them as-is.
	return &entity.Person{
		UserName:  person.ID,
		Name:      person.Name,
		Email:     person.Email,
		CreatedAt: person.CreatedAt,
		UpdatedAt: person.UpdatedAt,
		UpdatedBy: person.UpdatedBy,
	}
}
