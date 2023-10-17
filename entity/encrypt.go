package entity

import (
	"github.com/google/uuid"
)

type Encrypt struct {
	ID            uuid.UUID `gorm:"primary_key;not_null" json:"id"`
	Name          string    `json:"name" binding:"name"`
	PhoneNumber   string    `json:"phone_number" binding:"phone_number"`
	IDCardUrl     string    `json:"id_card_url" binding:"id_card_url"`
	CVUrl         string    `json:"cv_url" binding:"cv_url"`
	VideoUrl      string    `json:"video_url" binding:"video_url"`
	EncryptMethod string    `json:"encrypt_method" binding:"encrypt_method"`
	EncryptTime   string    `json:"encrypt_time" binding:"encrypt_time"`

	Timestamp
}
