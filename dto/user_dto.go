package dto

import (
	"time"

	"github.com/google/uuid"
)

// Request DTOs
type CreateUserRequest struct {
	Username    string  `json:"username" validate:"required,min=3,max=50"`
	Email       *string `json:"email" validate:"omitempty,email"`
	DisplayName *string `json:"display_name" validate:"omitempty,max=100"`
	Role        string  `json:"role" validate:"required,oneof=admin kaprodi dosen viewer"`
}

type UpdateUserRequest struct {
	Username    *string `json:"username" validate:"omitempty,min=3,max=50"`
	Email       *string `json:"email" validate:"omitempty,email"`
	DisplayName *string `json:"display_name" validate:"omitempty,max=100"`
	Role        *string `json:"role" validate:"omitempty,oneof=admin kaprodi dosen viewer"`
}

// Response DTOs
type UserResponse struct {
	ID          uuid.UUID `json:"id"`
	Username    string    `json:"username"`
	Email       *string   `json:"email,omitempty"`
	DisplayName *string   `json:"display_name,omitempty"`
	Role        string    `json:"role"`
	CreatedAt   time.Time `json:"created_at"`
}
