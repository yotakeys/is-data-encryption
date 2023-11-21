package entity

import (
	"github.com/google/uuid"
)

type Asymmetric struct {
	ID uuid.UUID `gorm:"primary_key;not_null" json:"id"`

	RequestedUserID  uuid.UUID `json:"requested_user_id"`
	RequestingUserID uuid.UUID `json:"requesting_user_id"`

	RequestedUser  User `gorm:"foreignKey:RequestedUserID" json:"requested_user"`
	RequestingUser User `gorm:"foreignKey:RequestingUserID" json:"requesting_user"`

	Name        string `json:"name" binding:"name"`
	PhoneNumber string `json:"phone_number" binding:"phone_number"`
	IDCardUrl   string `json:"id_card_url" binding:"id_card_url"`
	CVUrl       string `json:"cv_url" binding:"cv_url"`
	VideoUrl    string `json:"video_url" binding:"video_url"`

	Timestamp
}
