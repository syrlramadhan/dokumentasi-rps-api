package dto

import "github.com/google/uuid"

// Request DTOs
type CreateProgramRequest struct {
	Code string `json:"code" validate:"required,min=2,max=20"`
	Name string `json:"name" validate:"required,min=3,max=100"`
}

type UpdateProgramRequest struct {
	Code *string `json:"code" validate:"omitempty,min=2,max=20"`
	Name *string `json:"name" validate:"omitempty,min=3,max=100"`
}

// Response DTOs
type ProgramResponse struct {
	ID   uuid.UUID `json:"id"`
	Code string    `json:"code"`
	Name string    `json:"name"`
}
