package entity

import (
	"fmt"
	"strings"
	"time"
)

// Team represents a team that is part of the Ultimate Frisbee community.
type Team struct {
	Slug          string
	Name          string
	Description   string
	OriginCountry string

	CreatedAt time.Time
	CreatedBy string
	UpdatedAt time.Time
	UpdatedBy string
}

/****************/
/*  ATTRIBUTES  */
/****************/

type TeamAttribute string

type teamAttributeList struct {
	Slug          TeamAttribute
	Name          TeamAttribute
	Description   TeamAttribute
	OriginCountry TeamAttribute

	CreatedAt TeamAttribute
	CreatedBy TeamAttribute
	UpdatedAt TeamAttribute
	UpdatedBy TeamAttribute
}

// TeamAttributes represents the names of the attributes that a Team entity can have.
var TeamAttributes = &teamAttributeList{
	Slug:          "Slug",
	Name:          "Name",
	Description:   "Description",
	OriginCountry: "OriginCountry",

	CreatedAt: "CreatedAt",
	CreatedBy: "CreatedBy",
	UpdatedAt: "UpdatedAt",
	UpdatedBy: "UpdatedBy",
}

/***************/
/*    DEBUG    */
/***************/

func (team *Team) String() string {
	return team.StringWithIndentation(0)
}

func (team *Team) StringWithIndentation(indentationLevel int) string {
	if team == nil {
		return "[Team]=nil"
	}
	indentation := strings.Repeat(" ", indentationLevel)
	builder := strings.Builder{}
	builder.WriteString("[Team]\n")
	builder.WriteString(fmt.Sprintf("%sSlug: %s\n", indentation, team.Slug))
	builder.WriteString(fmt.Sprintf("%sName: %s\n", indentation, team.Name))
	builder.WriteString(fmt.Sprintf("%sDescription: %s\n", indentation, team.Description))
	builder.WriteString(fmt.Sprintf("%sOriginCountry: %s\n", indentation, team.OriginCountry))

	builder.WriteString(fmt.Sprintf("%sCreatedAt: %s\n", indentation, team.CreatedAt.String()))
	builder.WriteString(fmt.Sprintf("%sCreatedBy: %s\n", indentation, team.CreatedBy))
	builder.WriteString(fmt.Sprintf("%sUpdatedAt: %s\n", indentation, team.UpdatedAt.String()))
	builder.WriteString(fmt.Sprintf("%sUpdatedBy: %s\n", indentation, team.UpdatedBy))

	return builder.String()
}

/***************/
/*   TESTING   */
/***************/

func (team *Team) Clone() *Team {
	if team == nil {
		return nil
	}
	newTeam := &Team{
		Slug:          team.Slug,
		Name:          team.Name,
		Description:   team.Description,
		OriginCountry: team.OriginCountry,

		CreatedAt: team.CreatedAt,
		CreatedBy: team.CreatedBy,
		UpdatedAt: team.UpdatedAt,
		UpdatedBy: team.UpdatedBy,
	}

	return newTeam
}

func (team *Team) WithSlug(newSlug string) *Team {
	newTeam := team.Clone()
	newTeam.Slug = newSlug

	return newTeam
}

func (team *Team) WithName(newName string) *Team {
	newTeam := team.Clone()
	newTeam.Name = newName

	return newTeam
}

func (team *Team) WithDescription(newDescription string) *Team {
	newTeam := team.Clone()
	newTeam.Description = newDescription

	return newTeam
}

func (team *Team) WithOriginCountry(newOriginCountry string) *Team {
	newTeam := team.Clone()
	newTeam.Description = newOriginCountry

	return newTeam
}

func (team *Team) WithCreatedAt(newCreatedAt time.Time) *Team {
	newTeam := team.Clone()
	newTeam.CreatedAt = newCreatedAt

	return newTeam
}

func (team *Team) WithCreatedBy(newCreatedBy string) *Team {
	newTeam := team.Clone()
	newTeam.CreatedBy = newCreatedBy

	return newTeam
}

func (team *Team) WithUpdatedAt(newUpdatedAt time.Time) *Team {
	newTeam := team.Clone()
	newTeam.UpdatedAt = newUpdatedAt

	return newTeam
}

func (team *Team) WithUpdatedBy(newUpdatedBy string) *Team {
	newTeam := team.Clone()
	newTeam.UpdatedBy = newUpdatedBy

	return newTeam
}
