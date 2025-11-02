package entity

import (
	"fmt"
	"strings"
	"time"
)

// Person represents a person that interacts with the Ultimate Frisbee community.
type Person struct {
	UserName string
	ID       string
	IDType   string

	Name string

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
	ID     PersonAttribute
	IDType PersonAttribute

	Name PersonAttribute

	CreatedAt PersonAttribute
	CreatedBy PersonAttribute
	UpdatedAt PersonAttribute
	UpdatedBy PersonAttribute
}

// PersonAttributes represents the names of the attributes that a Person entity can have.
var PersonAttributes = &personAttributeList{
	ID:     "ID",
	IDType: "IDType",

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
	builder.WriteString(fmt.Sprintf("%s ID: %s\n", indentation, person.ID))
	builder.WriteString(fmt.Sprintf("%s IDType: %s\n", indentation, person.IDType))

	builder.WriteString(fmt.Sprintf("%s Name: %s\n", indentation, person.Name))

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
		UserName: person.UserName,
		ID:       person.ID,
		IDType:   person.IDType,

		Name: person.Name,

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

func (person *Person) WithID(newID string) *Person {
	newPerson := person.Clone()
	newPerson.ID = newID

	return newPerson
}

func (person *Person) WithIDType(newIDType string) *Person {
	newPerson := person.Clone()
	newPerson.IDType = newIDType

	return newPerson
}

func (person *Person) WithName(newName string) *Person {
	newPerson := person.Clone()
	newPerson.Name = newName

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
