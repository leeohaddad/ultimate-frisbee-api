package service

import (
	"context"
	"fmt"

	"github.com/leeohaddad/ultimate-frisbee-api/domain/entity"
	domainServiceParam "github.com/leeohaddad/ultimate-frisbee-api/domain/service/param"
	domainServiceResult "github.com/leeohaddad/ultimate-frisbee-api/domain/service/result"
)

func GetTeamMembershipsByRole(
	context context.Context,
	param domainServiceParam.GetTeamMembershipsByRole,
) (domainServiceResult.GetTeamMembershipsByRole, error) {
	memberships, err := param.Repository.GetMembershipsByTeamSlug(context, param.TeamSlug)
	if err != nil {
		return domainServiceResult.GetTeamMembershipsByRole{
			Memberships: []entity.Membership{},
		}, fmt.Errorf("failed to fetch all memberships of team '%s' from repository: %w", param.TeamSlug, err)
	}

	membershipsWithRole := []entity.Membership{}
	for _, membership := range memberships {
		if membership.Role == param.Role {
			membershipsWithRole = append(membershipsWithRole, membership)
		}
	}

	return domainServiceResult.GetTeamMembershipsByRole{
		Memberships: membershipsWithRole,
	}, nil
}
