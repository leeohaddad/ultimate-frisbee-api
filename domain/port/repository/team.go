package repository

import (
	"context"

	"errors"

	"github.com/leeohaddad/ultimate-frisbee-api/domain/entity"
)

// ErrAlreadyExists is returned by repository implementations when an entity
// cannot be created because a unique constraint (eg. name) already exists.
var ErrAlreadyExists = errors.New("repository: already exists")

type Team interface {
	GetAllTeams(context context.Context) ([]*entity.Team, error)
	GetTeamByName(context context.Context, name string) (*entity.Team, error)
	CreateTeam(context context.Context, team *entity.Team) (*entity.Team, error)
	UpdateTeam(context context.Context, team *entity.Team, updatedAttributes []entity.TeamAttribute) (*entity.Team, error)
}
