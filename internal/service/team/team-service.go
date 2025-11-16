package team

import (
	"context"
	"database/sql"
	"errors"

	"github.com/DenisSkachkov/backend-avito-assigment-autumn-2025/internal/models"
	"github.com/DenisSkachkov/backend-avito-assigment-autumn-2025/internal/service"
	"github.com/DenisSkachkov/backend-avito-assigment-autumn-2025/internal/service/user"
)

type TeamService struct {
	teamrepo TeamRepository
	userrepo user.UserRepository
}

func NewTeamService(teamRepo TeamRepository, userRepo user.UserRepository) *TeamService {
	return &TeamService{teamrepo: teamRepo, userrepo: userRepo}
}


 
func (s *TeamService) CreateTeam(ctx context.Context, t *models.Team) (*models.Team, error) {
    exists, err := s.teamrepo.TeamExists(ctx, t.Name)
    if err != nil {
        return nil, err
    }
    if exists {
        return nil, service.ErrTeamExists
    }

    if err := s.teamrepo.CreateTeam(ctx, t); err != nil {
        return nil, err
    }
    return t, nil
}

func (s *TeamService) GetTeam(ctx context.Context, name string) (*models.Team, error) {
	tm, err := s.teamrepo.GetTeamByName(ctx, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, service.ErrNotFound
		}
		return nil, err
	}
	return tm, nil
}