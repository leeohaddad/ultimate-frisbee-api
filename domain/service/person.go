package service

import (
	"context"
	"fmt"

	domainServiceParam "github.com/leeohaddad/ultimate-frisbee-api/domain/service/param"
	domainServiceResult "github.com/leeohaddad/ultimate-frisbee-api/domain/service/result"
)

func GetAllPeople(
	context context.Context,
	param domainServiceParam.GetAllPeople,
) (domainServiceResult.GetAllPeople, error) {
	People, err := param.Repository.GetAllPeople(context)
	if err != nil {
		return domainServiceResult.GetAllPeople{
			People: People,
		}, fmt.Errorf("failed to fetch all People from repository: %w", err)
	}

	return domainServiceResult.GetAllPeople{
		People: People,
	}, nil
}

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
