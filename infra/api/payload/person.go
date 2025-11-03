package payload

import (
	"time"

	"github.com/leeohaddad/ultimate-frisbee-api/domain/entity"
	"github.com/leeohaddad/ultimate-frisbee-api/infra/helper"
)

type Person struct {
	UserName      string  `json:"userName"`
	Name          string  `json:"name"`
	Email         *string `json:"email"`
	PhoneNumber   *string `json:"phoneNumber"`
	WFDFNumber    *string `json:"wfdfNumber"`
	OriginCountry *string `json:"originCountry"`

	CreatedBy *string `json:"createdBy"`
	CreatedAt *string `json:"createdAt"`
	UpdatedBy *string `json:"updatedBy"`
	UpdatedAt *string `json:"updatedAt"`
}

func ValidateCreatePersonInput(person *Person) (bool, string) {
	currentEntity := "Person"

	if helper.IsNilOrEmpty(&person.UserName) {
		return false, helper.ErrorMessageInField(currentEntity, "UserName")
	}

	if helper.IsNilOrEmpty(&person.Name) {
		return false, helper.ErrorMessageInField(currentEntity, "Name")
	}

	if helper.IsNilOrEmpty(person.Email) {
		return false, helper.ErrorMessageInField(currentEntity, "Email")
	}

	if helper.IsNilOrEmpty(person.PhoneNumber) {
		return false, helper.ErrorMessageInField(currentEntity, "PhoneNumber")
	}

	if helper.IsNilOrEmpty(person.WFDFNumber) {
		return false, helper.ErrorMessageInField(currentEntity, "WFDFNumber")
	}

	if helper.IsNilOrEmpty(person.OriginCountry) {
		return false, helper.ErrorMessageInField(currentEntity, "OriginCountry")
	}

	if helper.IsNilOrEmpty(person.CreatedBy) {
		return false, helper.ErrorMessageInField(currentEntity, "CreatedBy")
	}

	return true, ""
}

func ValidateUpdatePersonInput(person *Person, userName string) (bool, string) {
	if userName == "" {
		return false, "user name defined in the path variable is empty"
	}

	if person.UserName != "" && person.UserName != userName {
		return false, "updating the user name is not allowed"
	}

	if (person.Email == nil || *person.Email == "") &&
		(person.PhoneNumber == nil || *person.PhoneNumber == "") &&
		(person.WFDFNumber == nil || *person.WFDFNumber == "") &&
		(person.OriginCountry == nil || *person.OriginCountry == "") &&
		(person.UpdatedBy == nil || *person.UpdatedBy == "") {
		return false, "at least one of the following fields should not be empty: [Email, PhoneNumber, WFDFNumber, OriginCountry, UpdatedBy]"
	}

	return true, ""
}

func GetFilledPersonAttributesForUpdate(person *Person) []entity.PersonAttribute {
	var attributes []entity.PersonAttribute

	if person.Email != nil {
		attributes = append(attributes, entity.PersonAttributes.Email)
	}

	if person.PhoneNumber != nil {
		attributes = append(attributes, entity.PersonAttributes.PhoneNumber)
	}

	if person.WFDFNumber != nil {
		attributes = append(attributes, entity.PersonAttributes.WFDFNumber)
	}

	if person.OriginCountry != nil {
		attributes = append(attributes, entity.PersonAttributes.OriginCountry)
	}

	if person.UpdatedBy != nil {
		attributes = append(attributes, entity.PersonAttributes.UpdatedBy)
	}

	return attributes
}

func PersonToPersonEntity(person Person) *entity.Person {
	var email string
	if person.Email != nil {
		email = *person.Email
	}

	var phoneNumber string
	if person.PhoneNumber != nil {
		phoneNumber = *person.PhoneNumber
	}

	var wfdfNumber string
	if person.WFDFNumber != nil {
		wfdfNumber = *person.WFDFNumber
	}

	var originCountry string
	if person.OriginCountry != nil {
		originCountry = *person.OriginCountry
	}

	var createdBy string
	if person.CreatedBy != nil {
		createdBy = *person.CreatedBy
	}

	var createdAt time.Time
	if person.CreatedAt != nil {
		var err error
		createdAt, err = time.Parse(helper.DefaultTimeLayout, *person.CreatedAt)
		if err != nil {
			createdAt = time.Time{}
		}
	}

	var updatedBy string
	if person.UpdatedBy != nil {
		updatedBy = *person.UpdatedBy
	}

	var updatedAt time.Time
	if person.UpdatedAt != nil {
		var err error
		updatedAt, err = time.Parse(helper.DefaultTimeLayout, *person.UpdatedAt)
		if err != nil {
			updatedAt = time.Time{}
		}
	}

	return &entity.Person{
		UserName:      person.UserName,
		Name:          person.Name,
		Email:         email,
		PhoneNumber:   phoneNumber,
		WFDFNumber:    wfdfNumber,
		OriginCountry: originCountry,

		CreatedBy: createdBy,
		CreatedAt: createdAt,
		UpdatedBy: updatedBy,
		UpdatedAt: updatedAt,
	}
}

func PersonEntityToPerson(personEntity *entity.Person) Person {
	createdAt := personEntity.CreatedAt.Format(helper.DefaultTimeLayout)
	updatedAt := personEntity.UpdatedAt.Format(helper.DefaultTimeLayout)
	return Person{
		UserName:      personEntity.UserName,
		Name:          personEntity.Name,
		Email:         &personEntity.Email,
		PhoneNumber:   &personEntity.PhoneNumber,
		WFDFNumber:    &personEntity.WFDFNumber,
		OriginCountry: &personEntity.OriginCountry,

		CreatedBy: &personEntity.CreatedBy,
		CreatedAt: &createdAt,
		UpdatedBy: &personEntity.UpdatedBy,
		UpdatedAt: &updatedAt,
	}
}

func PersonEntitiesToPeople(personEntities []*entity.Person) []Person {
	People := make([]Person, 0)

	for _, personEntity := range personEntities {
		People = append(People, PersonEntityToPerson(personEntity))
	}

	if People == nil {
		return []Person{}
	}

	return People
}
