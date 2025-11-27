package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type AuditLog struct {
	ID         int64          `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID     *uuid.UUID     `json:"user_id" gorm:"type:uuid"`
	Action     string         `json:"action" gorm:"type:text;not null"`
	TargetType *string        `json:"target_type" gorm:"type:text"`
	TargetID   *uuid.UUID     `json:"target_id" gorm:"type:uuid"`
	Payload    datatypes.JSON `json:"payload" gorm:"type:jsonb"`
	CreatedAt  time.Time      `json:"created_at" gorm:"default:now()"`

	// Relations
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}
