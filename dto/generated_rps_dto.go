package dto

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// Request DTOs
type CreateGeneratedRPSRequest struct {
	TemplateVersionID *uuid.UUID `json:"template_version_id" validate:"required,uuid"`
	CourseID          *uuid.UUID `json:"course_id" validate:"required,uuid"`
	GeneratedBy       *uuid.UUID `json:"generated_by" validate:"omitempty,uuid"`
}

type UpdateGeneratedRPSRequest struct {
	Status          *string        `json:"status" validate:"omitempty,oneof=queued processing done failed"`
	Result          datatypes.JSON `json:"result" validate:"omitempty"`
	ExportedFileURL *string        `json:"exported_file_url" validate:"omitempty,url"`
	AIMetadata      datatypes.JSON `json:"ai_metadata" validate:"omitempty"`
}

type UpdateGeneratedRPSStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=queued processing done failed"`
}

// CompleteGenerationRequest - used by internal worker to submit generation result
type CompleteGenerationRequest struct {
	JobID           uuid.UUID      `json:"job_id" validate:"required,uuid"`
	Status          string         `json:"status" validate:"required,oneof=done failed"`
	Result          datatypes.JSON `json:"result" validate:"omitempty"`
	ExportedFileURL *string        `json:"exported_file_url" validate:"omitempty"`
	AIMetadata      datatypes.JSON `json:"ai_metadata" validate:"omitempty"`
	ErrorMessage    *string        `json:"error_message" validate:"omitempty"`
}

// GenerateRPSRequest - request body for POST /generate
type GenerateRPSRequest struct {
	TemplateVersionID uuid.UUID        `json:"template_version_id" validate:"required,uuid"`
	CourseID          uuid.UUID        `json:"course_id" validate:"required,uuid"`
	Options           *GenerateOptions `json:"options" validate:"omitempty"`
}

// GenerateOptions - options for RPS generation
type GenerateOptions struct {
	Language  string                 `json:"language" validate:"omitempty,oneof=id en"`
	Tone      string                 `json:"tone" validate:"omitempty,oneof=formal informal"`
	Overrides map[string]interface{} `json:"overrides" validate:"omitempty"`
}

// GenerateRPSResponse - response for POST /generate
type GenerateRPSResponse struct {
	JobID  uuid.UUID `json:"job_id"`
	Status string    `json:"status"`
}

// Response DTOs
type GeneratedRPSResponse struct {
	ID                uuid.UUID                `json:"id"`
	TemplateVersionID *uuid.UUID               `json:"template_version_id,omitempty"`
	CourseID          *uuid.UUID               `json:"course_id,omitempty"`
	GeneratedBy       *uuid.UUID               `json:"generated_by,omitempty"`
	Status            string                   `json:"status"`
	Result            datatypes.JSON           `json:"result,omitempty"`
	ExportedFileURL   *string                  `json:"exported_file_url,omitempty"`
	AIMetadata        datatypes.JSON           `json:"ai_metadata,omitempty"`
	CreatedAt         time.Time                `json:"created_at"`
	UpdatedAt         time.Time                `json:"updated_at"`
	TemplateVersion   *TemplateVersionResponse `json:"template_version,omitempty"`
	Course            *CourseResponse          `json:"course,omitempty"`
	Generator         *UserResponse            `json:"generator,omitempty"`
}
