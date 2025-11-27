package dto

import (
	"time"

	"github.com/google/uuid"
)

// Request DTOs
type CreateCourseRequest struct {
	ProgramID *uuid.UUID `json:"program_id" validate:"omitempty,uuid"`
	Code      string     `json:"code" validate:"required,min=2,max=20"`
	Title     string     `json:"title" validate:"required,min=3,max=200"`
	Credits   *int       `json:"credits" validate:"omitempty,min=1,max=10"`
}

type UpdateCourseRequest struct {
	ProgramID *uuid.UUID `json:"program_id" validate:"omitempty,uuid"`
	Code      *string    `json:"code" validate:"omitempty,min=2,max=20"`
	Title     *string    `json:"title" validate:"omitempty,min=3,max=200"`
	Credits   *int       `json:"credits" validate:"omitempty,min=1,max=10"`
}

// Response DTOs
type CourseResponse struct {
	ID        uuid.UUID        `json:"id"`
	ProgramID *uuid.UUID       `json:"program_id,omitempty"`
	Code      string           `json:"code"`
	Title     string           `json:"title"`
	Credits   *int             `json:"credits,omitempty"`
	CreatedAt time.Time        `json:"created_at"`
	Program   *ProgramResponse `json:"program,omitempty"`
}
