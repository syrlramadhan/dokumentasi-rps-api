package dto

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// Request DTOs
type CreateAuditLogRequest struct {
	UserID     *uuid.UUID     `json:"user_id" validate:"omitempty,uuid"`
	Action     string         `json:"action" validate:"required,min=1,max=100"`
	TargetType *string        `json:"target_type" validate:"omitempty,max=50"`
	TargetID   *uuid.UUID     `json:"target_id" validate:"omitempty,uuid"`
	Payload    datatypes.JSON `json:"payload" validate:"omitempty"`
}

type AuditLogFilterRequest struct {
	UserID     *uuid.UUID `json:"user_id" validate:"omitempty,uuid"`
	Action     *string    `json:"action" validate:"omitempty"`
	TargetType *string    `json:"target_type" validate:"omitempty"`
	TargetID   *uuid.UUID `json:"target_id" validate:"omitempty,uuid"`
	StartDate  *time.Time `json:"start_date" validate:"omitempty"`
	EndDate    *time.Time `json:"end_date" validate:"omitempty"`
}

// Response DTOs
type AuditLogResponse struct {
	ID         int64          `json:"id"`
	UserID     *uuid.UUID     `json:"user_id,omitempty"`
	Action     string         `json:"action"`
	TargetType *string        `json:"target_type,omitempty"`
	TargetID   *uuid.UUID     `json:"target_id,omitempty"`
	Payload    datatypes.JSON `json:"payload,omitempty"`
	CreatedAt  time.Time      `json:"created_at"`
	User       *UserResponse  `json:"user,omitempty"`
}
