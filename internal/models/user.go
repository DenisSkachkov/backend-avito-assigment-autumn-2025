package models

type User struct {
	Id string `json:"user_id"`
	Name string `json:"username"`
	Team string `json:"team_name"`
	IsActive bool `json:"is_active"`
}