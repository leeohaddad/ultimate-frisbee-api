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
	Username      string `pg:"username"`
	ID            string `pg:"id"` // TODO: remove me
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
              username,
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

func (repository *PersonRepository) GetPersonByUserName(context context.Context, userName string) (*entity.Person, error) {
	query := `select
			  username,
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
			  people
			where
			  username = ? limit 1`

	var fetchedPerson person
	queryResult, err := repository.client.ExecuteQuery(context, &fetchedPerson, query, userName)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve person by username: %w", err)
	}

	// Query executed successfully but no entity found for this username
	if queryResult.RowsReturned == 0 {
		return nil, nil
	}

	return personToPersonEntity(fetchedPerson), nil
}

func (repository *PersonRepository) CreatePerson(context context.Context, personEntity *entity.Person) (*entity.Person, error) {
	query := `insert into people (
			  username,
			  name,
			  email,
			  phone_number,
			  wfdf_number,
			  origin_country,

			  created_by,
			  created_at,
			  updated_at,
			  updated_by
			) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := repository.client.ExecuteCommand(
		context,
		query,

		personEntity.UserName,
		personEntity.Name,
		personEntity.Email,
		personEntity.PhoneNumber,
		personEntity.WFDFNumber,
		personEntity.OriginCountry,

		personEntity.CreatedBy,
		personEntity.CreatedAt,
		personEntity.UpdatedAt,
		personEntity.UpdatedBy,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create person: %w", err)
	}

	return personEntity, nil
}

func personToPersonEntity(person person) *entity.Person {
	// Rows are scanned directly into Go types by the DB client. createdAt/updatedAt
	// are already time.Time so we can use them as-is.
	return &entity.Person{
		UserName:      person.Username,
		Name:          person.Name,
		Email:         person.Email,
		PhoneNumber:   person.PhoneNumber,
		WFDFNumber:    person.WFDFNumber,
		OriginCountry: person.OriginCountry,

		CreatedAt: person.CreatedAt,
		CreatedBy: person.CreatedBy,
		UpdatedAt: person.UpdatedAt,
		UpdatedBy: person.UpdatedBy,
	}
}

func peopleToPersonEntities(people []person) []*entity.Person {
	personEntities := make([]*entity.Person, 0)

	for _, person := range people {
		personEntities = append(personEntities, personToPersonEntity(person))
	}

	return personEntities
}
