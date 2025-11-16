package models

type User struct {
	Id string `db:"user_id" json:"user_id"`
	Name string `db:"username" json:"username"`
	Team string `db:"team_name" json:"team_name"`
	IsActive bool `db:"is_active" json:"is_active"`
}