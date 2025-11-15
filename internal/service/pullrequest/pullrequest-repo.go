package pullrequest

import (
	"context"

	"github.com/DenisSkachkov/backend-avito-assigment-autumn-2025/internal/models"
)

type PullRequestRepository interface {
	GetPullRequestsByReviewerId(ctx context.Context, reviewerId string) ([]*models.PullRequest, error)
	GetByID(ctx context.Context, id string) (*models.PullRequest, error)
	ExistsByName(ctx context.Context, authorId, name string) (bool, error)
	Create(ctx context.Context,pr *models.PullRequest) error
	Update(ctx context.Context, pr *models.PullRequest) error
}