package models

import (
	"time"

	"github.com/google/uuid"
)

type Course struct {
	ID        uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ProgramID *uuid.UUID `json:"program_id" gorm:"type:uuid"`
	Code      string     `json:"code" gorm:"type:text;not null"`
	Title     string     `json:"title" gorm:"type:text;not null"`
	Credits   *int       `json:"credits" gorm:"type:int"`
	CreatedAt time.Time  `json:"created_at" gorm:"default:now()"`

	// Relations
	Program *Program `json:"program,omitempty" gorm:"foreignKey:ProgramID"`
}
