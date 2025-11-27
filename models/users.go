package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Username    string    `json:"username" gorm:"type:text;unique;not null"`
	Email       *string   `json:"email" gorm:"type:text;unique"`
	DisplayName *string   `json:"display_name" gorm:"type:text"`
	Role        string    `json:"role" gorm:"type:text;not null"` // 'admin'|'kaprodi'|'dosen'|'viewer'
	CreatedAt   time.Time `json:"created_at" gorm:"default:now()"`
}
