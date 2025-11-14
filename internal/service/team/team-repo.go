package team

import (
	"context"

	"github.com/DenisSkachkov/backend-avito-assigment-autumn-2025/internal/models"
)

type TeamRepository interface {
	GetTeamByUserId(ctx context.Context, userId string) (*models.Team, error)
}