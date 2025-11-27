package dto

import (
	"time"

	"github.com/google/uuid"
)

// Request DTOs
type CreateTemplateRequest struct {
	ProgramID   *uuid.UUID `json:"program_id" validate:"omitempty,uuid"`
	Name        string     `json:"name" validate:"required,min=3,max=100"`
	Description *string    `json:"description" validate:"omitempty,max=500"`
	CreatedBy   *uuid.UUID `json:"created_by" validate:"omitempty,uuid"`
	IsActive    *bool      `json:"is_active"`
}

type UpdateTemplateRequest struct {
	ProgramID   *uuid.UUID `json:"program_id" validate:"omitempty,uuid"`
	Name        *string    `json:"name" validate:"omitempty,min=3,max=100"`
	Description *string    `json:"description" validate:"omitempty,max=500"`
	IsActive    *bool      `json:"is_active"`
}

// Response DTOs
type TemplateResponse struct {
	ID          uuid.UUID        `json:"id"`
	ProgramID   *uuid.UUID       `json:"program_id,omitempty"`
	Name        string           `json:"name"`
	Description *string          `json:"description,omitempty"`
	CreatedBy   *uuid.UUID       `json:"created_by,omitempty"`
	CreatedAt   time.Time        `json:"created_at"`
	IsActive    bool             `json:"is_active"`
	Program     *ProgramResponse `json:"program,omitempty"`
	Creator     *UserResponse    `json:"creator,omitempty"`
}
