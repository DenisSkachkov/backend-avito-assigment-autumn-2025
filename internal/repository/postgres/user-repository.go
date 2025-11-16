package postgres

import (
	"context"
	"database/sql"
	"github.com/lib/pq"
	"fmt"
	"github.com/DenisSkachkov/backend-avito-assigment-autumn-2025/internal/models"
	"github.com/DenisSkachkov/backend-avito-assigment-autumn-2025/internal/service/user"
	"github.com/jmoiron/sqlx"
)

type userRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) user.UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) GetUserById(ctx context.Context, id string) (*models.User, error) {
    var u models.User
    err := r.db.GetContext(ctx, &u, `
        SELECT u.user_id, u.username, u.is_active, tm.team_name
        FROM users u
        LEFT JOIN team_members tm ON u.user_id = tm.user_id
        WHERE u.user_id = $1
    `, id)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, sql.ErrNoRows
        }
        return nil, err
    }
    return &u, nil
}

func (r *userRepo) Update(ctx context.Context, u *models.User) error {
    _, err := r.db.ExecContext(ctx, `UPDATE users SET username=$1, is_active=$2 WHERE user_id=$3`, u.Name, u.IsActive, u.Id)
    return err
}

func (r *userRepo) CreateUsers(ctx context.Context, users []models.User) error {
    if len(users) == 0 {
        return nil
    }

    tx, err := r.db.BeginTxx(ctx, nil)
    if err != nil {
        return err
    }
    defer func() {
        if err != nil {
            tx.Rollback()
        }
    }()

    stmt, err := tx.PrepareContext(ctx, `INSERT INTO users (user_id, username, is_active) VALUES ($1, $2, $3)`)
    if err != nil {
        tx.Rollback()
        return err
    }
    defer stmt.Close()

    for _, u := range users {
        if _, err = stmt.ExecContext(ctx, u.Id, u.Name, u.IsActive); err != nil {
            if isUniqueViolation(err) {
                tx.Rollback()
                return fmt.Errorf("user exists: %w", err)
            }
            tx.Rollback()
            return err
        }
    }

    if err = tx.Commit(); err != nil {
        return err
    }
    return nil
}

func isUniqueViolation(err error) bool {
    if err == nil {
        return false
    }
    if pqErr, ok := err.(*pq.Error); ok {
        return pqErr.Code == "23505"
    }
    return false
}