package application

import (
	"context"
	"fmt"

	serviceParam "github.com/leeohaddad/ultimate-frisbee-api/application/param"
	serviceResult "github.com/leeohaddad/ultimate-frisbee-api/application/result"

	"github.com/leeohaddad/ultimate-frisbee-api/domain/entity"
	domainService "github.com/leeohaddad/ultimate-frisbee-api/domain/service"
	domainServiceParam "github.com/leeohaddad/ultimate-frisbee-api/domain/service/param"
)

func GetTeamGameCaptains(context context.Context, param serviceParam.GetTeamGameCaptains) (serviceResult.GetTeamGameCaptains, error) {
	var gameCaptains []*entity.Person
	// TODO(lhaddad): register roles programmatically instead of relying on magic strings
	role := "Game Captain"

	// Retrieve all team memberships with the specified role
	result, err := domainService.GetTeamMembershipsByRole(context, domainServiceParam.GetTeamMembershipsByRole{
		Role: role,

		Repository: param.MembershipRepository,
	})
	if err != nil {
		return serviceResult.GetTeamGameCaptains{
			GameCaptains: []*entity.Person{},
		}, fmt.Errorf("failed to list game captains through domain service: %w", err)
	}

	// Extract game captains Person objects from obtained memberships
	for _, membership := range result.Memberships {
		gameCaptains = append(gameCaptains, membership.Person)
	}

	// Hydrate game captains Person objects data
	for i, gameCaptain := range gameCaptains {
		if gameCaptain.Name == "" {
			personResult, err := domainService.GetPersonByUserName(context, domainServiceParam.GetPersonByUserName{
				UserName: gameCaptain.UserName,

				Repository: param.PersonRepository,
			})
			if err != nil {
				return serviceResult.GetTeamGameCaptains{
					GameCaptains: []*entity.Person{},
				}, fmt.Errorf("failed to retrieve person '%s' through domain service: %w", gameCaptain.UserName, err)
			}

			gameCaptains[i] = personResult.Person
		}
	}

	return serviceResult.GetTeamGameCaptains{
		GameCaptains: gameCaptains,
	}, nil
}
