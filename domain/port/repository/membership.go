package repository

import (
	"context"

	"github.com/leeohaddad/ultimate-frisbee-api/domain/entity"
)

type Membership interface {
	GetMembershipsByTeamSlug(context context.Context, teamSlug string) ([]entity.Membership, error)
}
