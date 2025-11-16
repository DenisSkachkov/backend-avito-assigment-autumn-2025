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
		return err
	}

	if len(t.Members) > 0 {
		stmtUser, err := tx.PrepareContext(ctx, `
			INSERT INTO users (user_id, username, is_active)
			VALUES ($1, $2, $3)
			ON CONFLICT (user_id) DO NOTHING
		`)
		if err != nil {
			return err
		}
		defer stmtUser.Close()

		stmtMember, err := tx.PrepareContext(ctx, `
			INSERT INTO team_members (team_name, user_id)
			VALUES ($1, $2)
		`)
		if err != nil {
			return err
		}
		defer stmtMember.Close()

		for _, u := range t.Members {
			u.Team = t.Name
			if _, err = stmtUser.ExecContext(ctx, u.Id, u.Name, u.IsActive); err != nil {
				return err
			}
			if _, err = stmtMember.ExecContext(ctx, t.Name, u.Id); err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}


func (r *teamRepo) GetTeamByName(ctx context.Context, name string) (*models.Team, error) {

    exists, err := r.TeamExists(ctx, name)
    if err != nil {
        return nil, err
    }
    if !exists {
        return nil, sql.ErrNoRows
    }

    rows, err := r.db.QueryxContext(ctx, `
        SELECT user_id, username, is_active
        FROM users
        WHERE user_id IN (
            SELECT user_id FROM team_members WHERE team_name = $1
        )
        ORDER BY user_id
    `, name)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    members := []models.User{}
    for rows.Next() {
        var u models.User
        if err := rows.StructScan(&u); err != nil {
            return nil, err
        }
        u.Team = name
        members = append(members, u)
    }

    return &models.Team{
        Name:    name,
        Members: members, 
    }, nil
}



func (r *teamRepo) TeamExists(ctx context.Context, name string) (bool, error) {
    var exists bool
    err := r.db.GetContext(ctx, &exists, `SELECT EXISTS(SELECT 1 FROM teams WHERE team_name=$1)`, name)
    return exists, err
}



func (r *teamRepo) GetTeamByUserId(ctx context.Context, userId string) (*models.Team, error) {
	var teamName string
	err := r.db.GetContext(ctx, &teamName, `SELECT team_name FROM team_members WHERE user_id = $1 LIMIT 1`, userId)
	if err != nil {
		return nil, err
	}
	return r.GetTeamByName(ctx, teamName)
}
