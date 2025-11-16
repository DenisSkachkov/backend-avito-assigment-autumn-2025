package dto

type UserDTO struct {
    UserID   string `json:"user_id"`
    IsActive bool   `json:"is_active"`
}