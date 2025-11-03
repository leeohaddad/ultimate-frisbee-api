package entity

import (
	"fmt"
	"strings"
	"time"
)

// Person represents a person that interacts with the Ultimate Frisbee community.
type Person struct {
	UserName      string
	Name          string
	Email         string
	PhoneNumber   string
	WFDFNumber    string
	OriginCountry string

	CreatedAt time.Time
	CreatedBy string
	UpdatedAt time.Time
	UpdatedBy string
}

/****************/
/*  ATTRIBUTES  */
/****************/

type PersonAttribute string

type personAttributeList struct {
	Email         PersonAttribute
	PhoneNumber   PersonAttribute
	WFDFNumber    PersonAttribute
	OriginCountry PersonAttribute

	Name PersonAttribute

	CreatedAt PersonAttribute
	CreatedBy PersonAttribute
	UpdatedAt PersonAttribute
	UpdatedBy PersonAttribute
}

// PersonAttributes represents the names of the attributes that a Person entity can have.
var PersonAttributes = &personAttributeList{
	Email:         "Email",
	PhoneNumber:   "PhoneNumber",
	WFDFNumber:    "WFDFNumber",
	OriginCountry: "OriginCountry",

	Name: "Name",

	CreatedAt: "CreatedAt",
	CreatedBy: "CreatedBy",
	UpdatedAt: "UpdatedAt",
	UpdatedBy: "UpdatedBy",
}

/***************/
/*    DEBUG    */
/***************/

func (person *Person) String() string {
	return person.StringWithIndentation(0)
}

func (person *Person) StringWithIndentation(indentationLevel int) string {
	if person == nil {
		return "[Person]=nil"
	}
	indentation := strings.Repeat(" ", indentationLevel)
	builder := strings.Builder{}
	builder.WriteString("[Person]\n")
	builder.WriteString(fmt.Sprintf("%s User Name: %s\n", indentation, person.UserName))
	builder.WriteString(fmt.Sprintf("%s Name: %s\n", indentation, person.Name))
	builder.WriteString(fmt.Sprintf("%s Email: %s\n", indentation, person.Email))
	builder.WriteString(fmt.Sprintf("%s Phone Number: %s\n", indentation, person.PhoneNumber))
	builder.WriteString(fmt.Sprintf("%s WFDF Number: %s\n", indentation, person.WFDFNumber))
	builder.WriteString(fmt.Sprintf("%s Origin Country: %s\n", indentation, person.OriginCountry))

	builder.WriteString(fmt.Sprintf("%s CreatedAt: %s\n", indentation, person.CreatedAt.String()))
	builder.WriteString(fmt.Sprintf("%s CreatedBy: %s\n", indentation, person.CreatedBy))
	builder.WriteString(fmt.Sprintf("%s UpdatedAt: %s\n", indentation, person.UpdatedAt.String()))
	builder.WriteString(fmt.Sprintf("%s UpdatedBy: %s\n", indentation, person.UpdatedBy))

	return builder.String()
}

/***************/
/*   TESTING   */
/***************/

func (person *Person) Clone() *Person {
	if person == nil {
		return nil
	}
	newPerson := &Person{
		UserName:      person.UserName,
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

	return newPerson
}

func (person *Person) WithUserName(newUserName string) *Person {
	newPerson := person.Clone()
	newPerson.UserName = newUserName

	return newPerson
}

func (person *Person) WithName(newName string) *Person {
	newPerson := person.Clone()
	newPerson.Name = newName

	return newPerson
}

func (person *Person) WithEmail(newEmail string) *Person {
	newPerson := person.Clone()
	newPerson.Email = newEmail

	return newPerson
}

func (person *Person) WithPhoneNumber(newPhoneNumber string) *Person {
	newPerson := person.Clone()
	newPerson.PhoneNumber = newPhoneNumber

	return newPerson
}

func (person *Person) WithWFDFNumber(newWFDFNumber string) *Person {
	newPerson := person.Clone()
	newPerson.WFDFNumber = newWFDFNumber

	return newPerson
}

func (person *Person) WithOriginCountry(newOriginCountry string) *Person {
	newPerson := person.Clone()
	newPerson.OriginCountry = newOriginCountry

	return newPerson
}

func (person *Person) WithCreatedAt(newCreatedAt time.Time) *Person {
	newPerson := person.Clone()
	newPerson.CreatedAt = newCreatedAt

	return newPerson
}

func (person *Person) WithCreatedBy(newCreatedBy string) *Person {
	newPerson := person.Clone()
	newPerson.CreatedBy = newCreatedBy

	return newPerson
}

func (person *Person) WithUpdatedAt(newUpdatedAt time.Time) *Person {
	newPerson := person.Clone()
	newPerson.UpdatedAt = newUpdatedAt

	return newPerson
}

func (person *Person) WithUpdatedBy(newUpdatedBy string) *Person {
	newPerson := person.Clone()
	newPerson.UpdatedBy = newUpdatedBy

	return newPerson
}
