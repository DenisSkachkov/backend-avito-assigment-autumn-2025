package user

import (
	"context"
	"github.com/DenisSkachkov/backend-avito-assigment-autumn-2025/internal/models"
)
type UserRepository interface {
	GetUserById(ctx context.Context, id string) (*models.User, error)
	Update(ctx context.Context, u *models.User) error
}