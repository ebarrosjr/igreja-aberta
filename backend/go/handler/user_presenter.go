package handler

import (
	"time"

	"github.com/ebarrosjr/igreja-aberta/backend/go/database"
)

type userResponse struct {
	ID             int        `json:"id"`
	CongregationID int        `json:"congregation_id"`
	Name           string     `json:"name"`
	Email          string     `json:"email"`
	Status         string     `json:"status"`
	LastLoginAt    *time.Time `json:"last_login_at"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

func newUserResponse(user *database.User) userResponse {
	return userResponse{
		ID:             user.Id,
		CongregationID: user.CongregationId,
		Name:           user.Name,
		Email:          user.Email,
		Status:         user.Status,
		LastLoginAt:    user.LastLoginAt,
		CreatedAt:      user.CreatedAt,
		UpdatedAt:      user.UpdatedAt,
	}
}
