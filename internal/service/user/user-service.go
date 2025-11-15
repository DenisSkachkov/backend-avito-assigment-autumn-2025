package user

import (
	"context"

	"github.com/DenisSkachkov/backend-avito-assigment-autumn-2025/internal/models"
	"github.com/DenisSkachkov/backend-avito-assigment-autumn-2025/internal/service"
)

type UserService struct {
	repo UserRepository
}

func (s *UserService) SetActive(ctx context.Context, id string, active bool) (*models.User, error) {
	user, err := s.repo.GetUserById(ctx, id)
	if err != nil {
		return nil, service.ErrNotFound
	}
	user.IsActive = active
	return user, s.repo.Update(ctx, user)
}
