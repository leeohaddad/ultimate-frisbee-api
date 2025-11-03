package seeds

import (
	"context"
	"time"

	"github.com/leeohaddad/ultimate-frisbee-api/domain/entity"
	"github.com/leeohaddad/ultimate-frisbee-api/domain/port/logger"
	"github.com/leeohaddad/ultimate-frisbee-api/domain/port/repository"
)

// SeedPeople creates sample people in the database
func SeedPeople(ctx context.Context, personRepo repository.Person, logger logger.Logger) error {
	// Sample people data
	people := []struct {
		userName      string
		name          string
		email         string
		phoneNumber   string
		wfdfNumber    string
		originCountry string
		createdBy     string
	}{
		{
			userName:      "notdougz",
			name:          "Douglas Olvieira",
			email:         "doug@gmail.com",
			phoneNumber:   "(11) 98765-4321",
			wfdfNumber:    "12",
			originCountry: "Brazil",
			createdBy:     "admin",
		},
		{
			userName:      "allanbm100",
			name:          "Allan Moreira",
			email:         "allan@gmail.com",
			phoneNumber:   "(11) 98765-4321",
			wfdfNumber:    "34",
			originCountry: "Brazil",
			createdBy:     "admin",
		},
		{
			userName:      "Iolivieri",
			name:          "Isabella Olivieri",
			email:         "bella@gmail.com",
			phoneNumber:   "(11) 98765-4321",
			wfdfNumber:    "56",
			originCountry: "Brazil",
			createdBy:     "admin",
		},
	}

	logger.Info("starting people seeding...")

	for _, person := range people {
		// Check if person already exists
		existing, err := personRepo.GetPersonByUserName(ctx, person.userName)
		if err != nil {
			logger.WithError(err).Errorf("error checking if person %s exists", person.userName)
			continue
		}

		if existing != nil {
			logger.Infof("person %s already exists, skipping", person.userName)
			continue
		}

		// Create new person
		personEntity := &entity.Person{
			UserName:      person.userName,
			Name:          person.name,
			Email:         person.email,
			PhoneNumber:   person.phoneNumber,
			WFDFNumber:    person.wfdfNumber,
			OriginCountry: person.originCountry,
			CreatedBy:     person.createdBy,
			CreatedAt:     time.Now(),
			UpdatedBy:     person.createdBy,
			UpdatedAt:     time.Now(),
		}

		createdPerson, err := personRepo.CreatePerson(ctx, personEntity)
		if err != nil {
			logger.WithError(err).Errorf("failed to create person %s", person.name)
			continue
		}

		logger.Infof("successfully created person: %s", createdPerson.Name)
	}

	logger.Info("people seeding completed!")
	return nil
}
