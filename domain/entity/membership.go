package entity

import (
	"fmt"
	"strings"
	"time"
)

// Membership represents a relationship between a person and a team.
type Membership struct {
	Team   *Team
	Person *Person
	Role   string

	StartDate time.Time
	EndDate   time.Time

	CreatedAt time.Time
	CreatedBy string
	UpdatedAt time.Time
	UpdatedBy string
}

/****************/
/*  ATTRIBUTES  */
/****************/

type MembershipAttribute string

type membershipAttributeList struct {
	Team   MembershipAttribute
	Person MembershipAttribute
	Role   MembershipAttribute

	StartDate MembershipAttribute
	EndDate   MembershipAttribute

	CreatedAt MembershipAttribute
	CreatedBy MembershipAttribute
	UpdatedAt MembershipAttribute
	UpdatedBy MembershipAttribute
}

// MembershipAttributes represents the names of the attributes that a Membership entity can have.
var MembershipAttributes = &membershipAttributeList{
	Team:   "Team",
	Person: "Person",
	Role:   "Role",

	StartDate: "StartDate",
	EndDate:   "EndDate",

	CreatedAt: "CreatedAt",
	CreatedBy: "CreatedBy",
	UpdatedAt: "UpdatedAt",
	UpdatedBy: "UpdatedBy",
}

/***************/
/*    DEBUG    */
/***************/

func (membership *Membership) String() string {
	return membership.StringWithIndentation(0)
}

func (membership *Membership) StringWithIndentation(indentationLevel int) string {
	if membership == nil {
		return "[Membership]=nil"
	}
	indentation := strings.Repeat(" ", indentationLevel)
	builder := strings.Builder{}
	builder.WriteString("[Membership]\n")
	team := membership.Team.StringWithIndentation(indentationLevel + 2)
	builder.WriteString(fmt.Sprintf("%s Team: %s\n", indentation, team))
	person := membership.Team.StringWithIndentation(indentationLevel + 2)
	builder.WriteString(fmt.Sprintf("%s Person: %s\n", indentation, person))
	builder.WriteString(fmt.Sprintf("%s Role: '%s'\n", indentation, membership.Role))

	builder.WriteString(fmt.Sprintf("%s StartDate: '%s'\n", indentation, membership.StartDate))
	builder.WriteString(fmt.Sprintf("%s EndDate: '%s'\n", indentation, membership.EndDate))

	builder.WriteString(fmt.Sprintf("%s CreatedAt: %s\n", indentation, membership.CreatedAt.String()))
	builder.WriteString(fmt.Sprintf("%s CreatedBy: %s\n", indentation, membership.CreatedBy))
	builder.WriteString(fmt.Sprintf("%s UpdatedAt: %s\n", indentation, membership.UpdatedAt.String()))
	builder.WriteString(fmt.Sprintf("%s UpdatedBy: %s\n", indentation, membership.UpdatedBy))

	return builder.String()
}

/***************/
/*   TESTING   */
/***************/

func (membership *Membership) Clone() *Membership {
	if membership == nil {
		return nil
	}
	newMembership := &Membership{
		Team:   membership.Team.Clone(),
		Person: membership.Person.Clone(),
		Role:   membership.Role,

		StartDate: membership.StartDate,
		EndDate:   membership.EndDate,

		CreatedAt: membership.CreatedAt,
		CreatedBy: membership.CreatedBy,
		UpdatedAt: membership.UpdatedAt,
		UpdatedBy: membership.UpdatedBy,
	}

	return newMembership
}

func (membership *Membership) WithTeam(newTeam *Team) *Membership {
	newMembership := membership.Clone()
	newMembership.Team = newTeam

	return newMembership
}

func (membership *Membership) WithPerson(newPerson *Person) *Membership {
	newMembership := membership.Clone()
	newMembership.Person = newPerson

	return newMembership
}

func (membership *Membership) WithRole(newRole string) *Membership {
	newMembership := membership.Clone()
	newMembership.Role = newRole

	return newMembership
}

func (membership *Membership) WithStartDate(newStartDate time.Time) *Membership {
	newMembership := membership.Clone()
	newMembership.StartDate = newStartDate

	return newMembership
}

func (membership *Membership) WithEndDate(newEndDate time.Time) *Membership {
	newMembership := membership.Clone()
	newMembership.EndDate = newEndDate

	return newMembership
}

func (membership *Membership) WithCreatedAt(newCreatedAt time.Time) *Membership {
	newMembership := membership.Clone()
	newMembership.CreatedAt = newCreatedAt

	return newMembership
}

func (membership *Membership) WithCreatedBy(newCreatedBy string) *Membership {
	newMembership := membership.Clone()
	newMembership.CreatedBy = newCreatedBy

	return newMembership
}

func (membership *Membership) WithUpdatedAt(newUpdatedAt time.Time) *Membership {
	newMembership := membership.Clone()
	newMembership.UpdatedAt = newUpdatedAt

	return newMembership
}

func (membership *Membership) WithUpdatedBy(newUpdatedBy string) *Membership {
	newMembership := membership.Clone()
	newMembership.UpdatedBy = newUpdatedBy

	return newMembership
}
