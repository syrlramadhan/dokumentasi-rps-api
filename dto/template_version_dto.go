package dto

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// Request DTOs
type CreateTemplateVersionRequest struct {
	TemplateID *uuid.UUID     `json:"template_id" validate:"required,uuid"`
	Version    int            `json:"version" validate:"required,min=1"`
	Definition datatypes.JSON `json:"definition" validate:"required"`
	CreatedBy  *uuid.UUID     `json:"created_by" validate:"omitempty,uuid"`
}

type UpdateTemplateVersionRequest struct {
	Version    *int           `json:"version" validate:"omitempty,min=1"`
	Definition datatypes.JSON `json:"definition" validate:"omitempty"`
}

// Response DTOs
type TemplateVersionResponse struct {
	ID         uuid.UUID         `json:"id"`
	TemplateID *uuid.UUID        `json:"template_id,omitempty"`
	Version    int               `json:"version"`
	Definition datatypes.JSON    `json:"definition"`
	CreatedBy  *uuid.UUID        `json:"created_by,omitempty"`
	CreatedAt  time.Time         `json:"created_at"`
	Template   *TemplateResponse `json:"template,omitempty"`
	Creator    *UserResponse     `json:"creator,omitempty"`
}
