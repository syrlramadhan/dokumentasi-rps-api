package models

import "github.com/google/uuid"

type Program struct {
	ID   uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Code string    `json:"code" gorm:"type:text;unique;not null"`
	Name string    `json:"name" gorm:"type:text;not null"`
}
