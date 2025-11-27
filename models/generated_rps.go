package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type GeneratedRPS struct {
	ID                uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	TemplateVersionID *uuid.UUID     `json:"template_version_id" gorm:"type:uuid"`
	CourseID          *uuid.UUID     `json:"course_id" gorm:"type:uuid"`
	GeneratedBy       *uuid.UUID     `json:"generated_by" gorm:"type:uuid"`      // user yang memicu
	Status            string         `json:"status" gorm:"type:text;not null"`   // queued|processing|done|failed
	Result            datatypes.JSON `json:"result" gorm:"type:jsonb"`           // final RPS structured
	ExportedFileURL   *string        `json:"exported_file_url" gorm:"type:text"` // S3 link jika ada
	AIMetadata        datatypes.JSON `json:"ai_metadata" gorm:"type:jsonb"`      // ringkasan: model name, temperature, tokens, prompt_id (Mongo)
	CreatedAt         time.Time      `json:"created_at" gorm:"default:now()"`
	UpdatedAt         time.Time      `json:"updated_at" gorm:"default:now()"`

	// Relations
	TemplateVersion *TemplateVersion `json:"template_version,omitempty" gorm:"foreignKey:TemplateVersionID"`
	Course          *Course          `json:"course,omitempty" gorm:"foreignKey:CourseID"`
	Generator       *User            `json:"generator,omitempty" gorm:"foreignKey:GeneratedBy"`
}
