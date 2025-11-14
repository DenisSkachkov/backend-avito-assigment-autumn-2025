package models

type User struct {
	Id string `json:"user_id"`
	Name string `json:"username"`
	IsActive bool `json:"is_active"`
}