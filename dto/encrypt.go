package dto

import (
	"mime/multipart"

	"github.com/google/uuid"
)

type EncryptCreateDto struct {
	ID          uuid.UUID             `gorm:"primary_key" json:"id" form:"id"`
	Name        string                `json:"name" form:"name" binding:"required"`
	PhoneNumber string                `json:"phone_number" form:"phone_number" binding:"required"`
	IDCard      *multipart.FileHeader `json:"id_card" form:"id_card" binding:"required"`
	CV          *multipart.FileHeader `json:"cv" form:"cv" binding:"required"`
	Video       *multipart.FileHeader `json:"video" form:"video" binding:"required"`
	UserID      uuid.UUID             `json:"user_id" form:"user_id"`
}
