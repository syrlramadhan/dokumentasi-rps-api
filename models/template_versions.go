package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type TemplateVersion struct {
	ID         uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	TemplateID *uuid.UUID     `json:"template_id" gorm:"type:uuid"`
	Version    int            `json:"version" gorm:"not null"`
	Definition datatypes.JSON `json:"definition" gorm:"type:jsonb;not null"` // struktur RPS semi-terstruktur
	CreatedBy  *uuid.UUID     `json:"created_by" gorm:"type:uuid"`
	CreatedAt  time.Time      `json:"created_at" gorm:"default:now()"`

	// Relations
	Template *Template `json:"template,omitempty" gorm:"foreignKey:TemplateID"`
	Creator  *User     `json:"creator,omitempty" gorm:"foreignKey:CreatedBy"`
}
