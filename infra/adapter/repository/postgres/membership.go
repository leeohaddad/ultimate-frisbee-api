package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/leeohaddad/ultimate-frisbee-api/infra/helper"

	"github.com/leeohaddad/ultimate-frisbee-api/domain/entity"
	repositoryPort "github.com/leeohaddad/ultimate-frisbee-api/domain/port/repository"
	postgresDatabase "github.com/leeohaddad/ultimate-frisbee-api/infra/database/postgres"
)

// Enforce that MembershipRepository implements the repositoryPort.Membership interface.
var _ repositoryPort.Membership = (*MembershipRepository)(nil)

type MembershipRepository struct {
	client postgresDatabase.Client
}

// membership is a representation on how the membership is retrieved from the database.
type membership struct {
	TeamSlug       string `pg:"team_slug"`
	PersonUserName string `pg:"person_user_name"`
	Role           string `pg:"role"`
	StartDate      string `pg:"start_date"`
	EndDate        string `pg:"end_date"`

	CreatedAt string `pg:"created_at"`
	CreatedBy string `pg:"created_by"`
	UpdatedAt string `pg:"updated_at"`
	UpdatedBy string `pg:"updated_by"`
}

// NewMembershipRepository instantiates a new membership repository for postgres.
func NewMembershipRepository(client postgresDatabase.Client) *MembershipRepository {
	return &MembershipRepository{
		client: client,
	}
}

func (repository *MembershipRepository) GetMembershipsByTeamSlug(context context.Context, teamSlug string) ([]entity.Membership, error) {
	query := `select
              team_slug,
              person_username,
              role
            from
              memberships
            where
              team_slug = ?`

	// Execute query in DB
	var fetchedMemberships []membership
	queryResult, err := repository.client.ExecuteQuery(context, &fetchedMemberships, query, teamSlug)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve game captains from team %s: %w", teamSlug, err)
	}

	// Query executed successfully but no entity found for this ID
	if queryResult.RowsReturned == 0 {
		return []entity.Membership{}, nil
	}

	return membershipsToMembershipEntities(fetchedMemberships)
}

func membershipsToMembershipEntities(memberships []membership) ([]entity.Membership, error) {
	var membershipEntities []entity.Membership

	for _, membership := range memberships {
		membershipEntity, err := membershipToMembershipEntity(membership)
		if err != nil {
			return []entity.Membership{}, fmt.Errorf("failed to convert membership to membership entity: %w", err)
		}
		membershipEntities = append(membershipEntities, *membershipEntity)
	}

	return membershipEntities, nil
}

func membershipToMembershipEntity(membership membership) (*entity.Membership, error) {
	startDate, err := time.Parse(helper.DefaultTimeLayout, membership.StartDate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse start date: %w", err)
	}
	endDate, err := time.Parse(helper.DefaultTimeLayout, membership.EndDate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse end date: %w", err)
	}

	createdAt, err := time.Parse(helper.DefaultTimeLayout, membership.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to parse created at: %w", err)
	}
	updatedAt, err := time.Parse(helper.DefaultTimeLayout, membership.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to parse updated at: %w", err)
	}

	return &entity.Membership{
		Team:   &entity.Team{Slug: membership.TeamSlug},
		Person: &entity.Person{UserName: membership.PersonUserName},
		Role:   membership.Role,

		StartDate: startDate,
		EndDate:   endDate,

		CreatedAt: createdAt,
		CreatedBy: membership.CreatedBy,
		UpdatedAt: updatedAt,
		UpdatedBy: membership.UpdatedBy,
	}, nil
}
