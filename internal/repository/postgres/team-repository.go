package postgres

import (
	"context"
	"database/sql"
	"github.com/DenisSkachkov/backend-avito-assigment-autumn-2025/internal/models"
	"github.com/DenisSkachkov/backend-avito-assigment-autumn-2025/internal/service/team"
	"github.com/jmoiron/sqlx"
)

type teamRepo struct {
	db *sqlx.DB
}

func NewTeamRepo(db *sqlx.DB) team.TeamRepository {
	return &teamRepo{db: db}
}

func (r *teamRepo) CreateTeam(ctx context.Context, t *models.Team) error {
    tx, err := r.db.BeginTxx(ctx, nil)
    if err != nil {
        return err
    }

    defer func() {
        if err != nil {
            tx.Rollback()
        }
    }()

    if _, err = tx.ExecContext(ctx, `INSERT INTO teams (team_name) VALUES ($1)`, t.Name); err != nil {
        tx.Rollback()
        return err
    }

    if len(t.Members) > 0 {
        stmt, err := tx.PrepareContext(ctx, `INSERT INTO team_members (team_name, user_id) VALUES ($1, $2)`)
        if err != nil {
            tx.Rollback()
            return err
        }
        defer stmt.Close()

        for _, u := range t.Members {
            if _, err = stmt.ExecContext(ctx, t.Name, u.Id); err != nil {
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

func (r *teamRepo) GetTeamByName(ctx context.Context, name string) (*models.Team, error) {
    // Проверим, есть ли команда
    var exists bool
    err := r.db.GetContext(ctx, &exists, `SELECT EXISTS(SELECT 1 FROM teams WHERE team_name=$1)`, name)
    if err != nil {
        return nil, err
    }
    if !exists {
        return nil, sql.ErrNoRows
    }

    // Соберём участников
    rows, err := r.db.QueryxContext(ctx, `
        SELECT u.user_id AS id, u.username, u.is_active
        FROM team_members tm
        JOIN users u ON u.user_id = tm.user_id
        WHERE tm.team_name = $1
        ORDER BY u.user_id
    `, name)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    members := make([]models.User, 0)
    for rows.Next() {
        var u models.User
        if err := rows.StructScan(&u); err != nil {
            return nil, err
        }
        members = append(members, u)
    }

    team := &models.Team{
        Name:    name,
        Members: members,
    }
    return team, nil
}

func (r *teamRepo) GetTeamByUserId(ctx context.Context, userId string) (*models.Team, error) {
    var teamName string
    err := r.db.GetContext(ctx, &teamName, `SELECT team_name FROM team_members WHERE user_id = $1 LIMIT 1`, userId)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, sql.ErrNoRows
        }
        return nil, err
    }
    return r.GetTeamByName(ctx, teamName)
}
