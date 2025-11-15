package team

import (
	"context"
	"github.com/DenisSkachkov/backend-avito-assigment-autumn-2025/internal/models"
	"github.com/DenisSkachkov/backend-avito-assigment-autumn-2025/internal/service"
	"github.com/DenisSkachkov/backend-avito-assigment-autumn-2025/internal/service/user"
)

type TeamService struct {
	teamrepo TeamRepository
	userrepo user.UserRepository
}

func (t *TeamService) CreateTeam(ctx context.Context, team *models.Team) (*models.Team, error) {
	if _,err := t.teamrepo.GetTeamByName(ctx, team.Name); err != nil {
		return nil, service.ErrTeamExists
	}

	if len(team.Members) > 0 {
		if err := t.userrepo.CreateUsers(ctx, team.Members); err != nil {
			return nil, err
		}
	}

	if err := t.teamrepo.CreateTeam(ctx, team); err != nil {
		return nil, err
	}

	return team, nil
}

func (t *TeamService) GetTeam(ctx context.Context, name string) (*models.Team, error) {
    team, err := t.teamrepo.GetTeamByName(ctx, name)
    if err != nil {
        return nil, err
    }
    if team == nil {
        return nil, service.ErrNotFound
    }

    return team, nil
}