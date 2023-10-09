package entity

import (
	"time"

	"gorm.io/gorm"
)

type Timestamp struct {
	CreatedAt 		time.Time 	`json:"created_at" default:"CURRENT_TIMESTAMP"`
	UpdatedAt 		time.Time 	`json:"updated_at"`
	DeletedAt 		gorm.DeletedAt
}

type Authorization struct {
	Token string `json:"token"`
	Role  string `json:"role"`
}