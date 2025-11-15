package team

import (
	"context"

	"github.com/DenisSkachkov/backend-avito-assigment-autumn-2025/internal/models"
)

type TeamRepository interface {
	GetTeamByUserId(ctx context.Context, userId string) (*models.Team, error)
	GetTeamByName(ctx context.Context, name string) (*models.Team, error)
	CreateTeam(ctx context.Context, t *models.Team) error
}