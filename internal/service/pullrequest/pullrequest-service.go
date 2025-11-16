package pullrequest

import (
	"context"
	"math/rand"
	"time"

	"github.com/DenisSkachkov/backend-avito-assigment-autumn-2025/internal/models"
	"github.com/DenisSkachkov/backend-avito-assigment-autumn-2025/internal/service"
	"github.com/DenisSkachkov/backend-avito-assigment-autumn-2025/internal/service/team"
	"github.com/DenisSkachkov/backend-avito-assigment-autumn-2025/internal/service/user"
)

type PullRequestService struct {
	user user.UserRepository
	pr PullRequestRepository
	team team.TeamRepository
}

func NewPullRequestService(user user.UserRepository, pr PullRequestRepository, team team.TeamRepository) *PullRequestService {
	return &PullRequestService{user: user,pr: pr,team: team}
}

func (p *PullRequestService) GetReview(ctx context.Context, id string) ([]*models.PullRequest, error) {
	if _,err := p.user.GetUserById(ctx, id); err != nil {
		return nil, service.ErrNotFound
	}
	return p.pr.GetPullRequestsByReviewerId(ctx, id)
}

func selectTwoReviewers(members []models.User, authorId string) []string {
	reviewers := make([]string, 0)
	for _, v := range members {
		if v.Id == authorId{
			continue
		}
		reviewers = append(reviewers, v.Id)
	}

	if len(reviewers) <= 2 {
		return reviewers
	}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(reviewers), func(i, j int) {
		reviewers[i], reviewers[j] = reviewers[j], reviewers[i]
	})

	return reviewers[:2]
} 

func (p *PullRequestService) Create(ctx context.Context, prId string, authorId string, prName string) (*models.PullRequest, error) {
	user, err := p.user.GetUserById(ctx, authorId)
	if err != nil {
		return nil, service.ErrNotFound
	}

	exists, err := p.pr.ExistsByName(ctx, authorId, prName)
	if err != nil {
		return nil, err
	} 
	if exists {
		return nil, service.ErrPullRequestExists
	}

	team, err := p.team.GetTeamByUserId(ctx, authorId)
	if err != nil {
		return nil, service.ErrNotFound
	}

	reviewers := selectTwoReviewers(team.Members, authorId)

	pr := &models.PullRequest{
		PullRequestId: prId,
		PullRequestName: prName,
		AuthorId: user.Id,
		Status: "OPEN",
		AssignedReviewers: reviewers,
	}

	if err := p.pr.Create(ctx, pr); err != nil {
		return nil, err
	}

	return pr, nil
}

func (p *PullRequestService) Merge(ctx context.Context, prId string) (*models.PullRequest, error) {
	pr, err := p.pr.GetByID(ctx, prId)
	if err != nil {
		return nil, service.ErrNotFound
	}

	if pr.Status == "MERGED" {
		return pr, nil
	}

	pr.Status = "MERGED"
	mergedTime := time.Now().UTC()
	pr.MergedAt = &mergedTime

	if err := p.pr.Update(ctx, pr); err != nil {
		return nil, err
	}

	return pr,nil


}

func (p *PullRequestService) ReassignReviewer(ctx context.Context, prId, oldReviewerId string) (*models.PullRequest, string, error) {
	pr, err := p.pr.GetByID(ctx, prId) 
	if err != nil {
		return nil, "", service.ErrNotFound
	}
	
	if pr.Status == "MERGED" {
		return nil, "", service.ErrPullRequestMerged
	}

	idx := -1
	for k, v := range pr.AssignedReviewers {
		if v == oldReviewerId {
			idx = k
			break
		}
	}

	if idx == -1{
		return nil, "", service.ErrNotAssigned
	}

	if _, err := p.user.GetUserById(ctx, oldReviewerId);err != nil {
		return nil, "", service.ErrNotFound
	}

	team, err := p.team.GetTeamByUserId(ctx, oldReviewerId)
	if err != nil {
		return nil, "", service.ErrNotFound
	}

	candidateID := ""
    for _, u := range team.Members {
        if !u.IsActive {
            continue
        }
        if u.Id == pr.AuthorId {
            continue
        }
        if u.Id == oldReviewerId {
            continue
        }

        alreadyAssigned := false
        for _, r := range pr.AssignedReviewers {
            if r == u.Id {
                alreadyAssigned = true
                break
            }
        }
        if alreadyAssigned {
            continue
        }

        candidateID = u.Id
        break
    }

    if candidateID == "" {
        return nil, "", service.ErrNoCandidate
    }

    pr.AssignedReviewers[idx] = candidateID

    if err := p.pr.Update(ctx, pr); err != nil {
        return nil, "", err
    }

    return pr, candidateID, nil

}