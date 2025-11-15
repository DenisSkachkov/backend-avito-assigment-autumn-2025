package postgres

import (
	"context"
	"database/sql"

	"github.com/DenisSkachkov/backend-avito-assigment-autumn-2025/internal/models"
	"github.com/DenisSkachkov/backend-avito-assigment-autumn-2025/internal/service/pullrequest"
	"github.com/jmoiron/sqlx"
)

type prRepo struct {
	db *sqlx.DB
}


func NewPRRepo(db *sqlx.DB) pullrequest.PullRequestRepository {
	return &prRepo{db: db}
}

func (r *prRepo) GetPullRequestsByReviewerId(ctx context.Context, reviewerID string) ([]*models.PullRequest, error) {
    rows, err := r.db.QueryxContext(ctx, `
        SELECT pr.pull_request_id, pr.pull_request_name, pr.author_id, pr.status, pr.merged_at
        FROM pull_requests pr
        JOIN pr_reviewers rr ON rr.pull_request_id = pr.pull_request_id
        WHERE rr.reviewer_id = $1
        ORDER BY pr.pull_request_id
    `, reviewerID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var res []*models.PullRequest
    for rows.Next() {
        var pr models.PullRequest
        var mergedAt sql.NullTime
        if err := rows.Scan(&pr.PullRequestId, &pr.PullRequestName, &pr.AuthorId, &pr.Status, &mergedAt); err != nil {
            return nil, err
        }
        if mergedAt.Valid {
            t := mergedAt.Time
            pr.MergedAt = &t
        }
        // получить reviewers
        revs, err := r.GetReviewersForPR(ctx, pr.PullRequestId)
        if err != nil {
            return nil, err
        }
        pr.AssignedReviewers = revs
        res = append(res, &pr)
    }
    return res, nil
}

func (r *prRepo) GetByID(ctx context.Context, id string) (*models.PullRequest, error) {
    var pr models.PullRequest
    var mergedAt sql.NullTime
    err := r.db.QueryRowxContext(ctx, `
        SELECT pull_request_id, pull_request_name, author_id, status, merged_at
        FROM pull_requests
        WHERE pull_request_id = $1
    `, id).Scan(&pr.PullRequestId, &pr.PullRequestName, &pr.AuthorId, &pr.Status, &mergedAt)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, sql.ErrNoRows
        }
        return nil, err
    }
    if mergedAt.Valid {
        t := mergedAt.Time
        pr.MergedAt = &t
    }
    revs, err := r.GetReviewersForPR(ctx, pr.PullRequestId)
    if err != nil {
        return nil, err
    }
    pr.AssignedReviewers = revs
    return &pr, nil
}

func (r *prRepo) ExistsByName(ctx context.Context, authorId, name string) (bool, error) {
    var exists bool
    err := r.db.GetContext(ctx, &exists, `
        SELECT EXISTS (SELECT 1 FROM pull_requests WHERE author_id=$1 AND pull_request_name=$2)
    `, authorId, name)
    return exists, err
}

func (r *prRepo) Create(ctx context.Context, pr *models.PullRequest) error {
    tx, err := r.db.BeginTxx(ctx, nil)
    if err != nil {
        return err
    }
    defer func() {
        if err != nil {
            tx.Rollback()
        }
    }()

    _, err = tx.ExecContext(ctx, `
        INSERT INTO pull_requests (pull_request_id, pull_request_name, author_id, status, merged_at)
        VALUES ($1, $2, $3, $4, $5)
    `, pr.PullRequestId, pr.PullRequestName, pr.AuthorId, pr.Status, pr.MergedAt)
    if err != nil {
        if isUniqueViolation(err) {
            tx.Rollback()
            return err
        }
        tx.Rollback()
        return err
    }

    // assign reviewers
    if len(pr.AssignedReviewers) > 0 {
        stmt, err := tx.PrepareContext(ctx, `INSERT INTO pr_reviewers (pull_request_id, reviewer_id) VALUES ($1, $2)`)
        if err != nil {
            tx.Rollback()
            return err
        }
        defer stmt.Close()

        for _, rID := range pr.AssignedReviewers {
            if _, err = stmt.ExecContext(ctx, pr.PullRequestId, rID); err != nil {
                tx.Rollback()
                return err
            }
        }
    }

    if err = tx.Commit(); err != nil {
        return err
    }
    return nil
}

func (r *prRepo) Update(ctx context.Context, pr *models.PullRequest) error {
    tx, err := r.db.BeginTxx(ctx, nil)
    if err != nil {
        return err
    }
    defer func() {
        if err != nil {
            tx.Rollback()
        }
    }()

    _, err = tx.ExecContext(ctx, `
        UPDATE pull_requests SET pull_request_name=$1, author_id=$2, status=$3, merged_at=$4
        WHERE pull_request_id=$5
    `, pr.PullRequestName, pr.AuthorId, pr.Status, pr.MergedAt, pr.PullRequestId)
    if err != nil {
        tx.Rollback()
        return err
    }

    // Обновляем список ревьюверов: проще всего — удалить старые и вставить новые
    if _, err = tx.ExecContext(ctx, `DELETE FROM pr_reviewers WHERE pull_request_id = $1`, pr.PullRequestId); err != nil {
        tx.Rollback()
        return err
    }
    if len(pr.AssignedReviewers) > 0 {
        stmt, err := tx.PrepareContext(ctx, `INSERT INTO pr_reviewers (pull_request_id, reviewer_id) VALUES ($1, $2)`)
        if err != nil {
            tx.Rollback()
            return err
        }
        defer stmt.Close()

        for _, rID := range pr.AssignedReviewers {
            if _, err = stmt.ExecContext(ctx, pr.PullRequestId, rID); err != nil {
                tx.Rollback()
                return err
            }
        }
    }

    if err = tx.Commit(); err != nil {
        return err
    }
    return nil
}

func (r *prRepo) GetReviewersForPR(ctx context.Context, prID string) ([]string, error) {
    rows, err := r.db.QueryxContext(ctx, `SELECT reviewer_id FROM pr_reviewers WHERE pull_request_id = $1 ORDER BY reviewer_id`, prID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var res []string
    for rows.Next() {
        var id string
        if err := rows.Scan(&id); err != nil {
            return nil, err
        }
        res = append(res, id)
    }
    return res, nil
}
