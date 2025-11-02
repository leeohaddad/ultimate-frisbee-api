package service

import (
	"context"
	"fmt"

	domainServiceParam "github.com/leeohaddad/ultimate-frisbee-api/domain/service/param"
	domainServiceResult "github.com/leeohaddad/ultimate-frisbee-api/domain/service/result"
)

func GetPersonByUserName(
	context context.Context,
	param domainServiceParam.GetPersonByUserName,
) (domainServiceResult.GetPersonByUserName, error) {
	person, err := param.Repository.GetPersonByUserName(context, param.UserName)
	if err != nil {
		return domainServiceResult.GetPersonByUserName{
			Person: nil,
		}, fmt.Errorf("failed to fetch person '%s' from repository: %w", param.UserName, err)
	}

	return domainServiceResult.GetPersonByUserName{
		Person: person,
	}, nil
}
