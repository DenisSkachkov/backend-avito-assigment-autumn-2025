package postgres

import (
	"github.com/DenisSkachkov/backend-avito-assigment-autumn-2025/internal/service/pullrequest"
	"github.com/DenisSkachkov/backend-avito-assigment-autumn-2025/internal/service/team"
	"github.com/DenisSkachkov/backend-avito-assigment-autumn-2025/internal/service/user"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	PRs pullrequest.PullRequestRepository
	Users user.UserRepository
	Teams team.TeamRepository
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Teams: NewTeamRepo(db),
		Users: NewUserRepo(db),
		PRs: NewPRRepo(db),
	}
}