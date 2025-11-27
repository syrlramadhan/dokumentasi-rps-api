package models

import (
	"time"

	"github.com/google/uuid"
)

type Template struct {
	ID          uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ProgramID   *uuid.UUID `json:"program_id" gorm:"type:uuid"`
	Name        string     `json:"name" gorm:"type:text;not null"`
	Description *string    `json:"description" gorm:"type:text"`
	CreatedBy   *uuid.UUID `json:"created_by" gorm:"type:uuid"`
	CreatedAt   time.Time  `json:"created_at" gorm:"default:now()"`
	IsActive    bool       `json:"is_active" gorm:"default:true"`

	// Relations
	Program *Program `json:"program,omitempty" gorm:"foreignKey:ProgramID"`
	Creator *User    `json:"creator,omitempty" gorm:"foreignKey:CreatedBy"`
}
