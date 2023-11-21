package dto

import (
	"mime/multipart"

	"github.com/google/uuid"
)

type EncryptCreateDto struct {
	ID             uuid.UUID            `gorm:"primary_key" json:"id" form:"id"`
	Name           string               `json:"name" form:"name" binding:"required"`
	PhoneNumber    string               `json:"phone_number" form:"phone_number" binding:"required"`
	IDCard         multipart.FileHeader `json:"id_card" form:"id_card" binding:"required"`
	CV             multipart.FileHeader `json:"cv" form:"cv" binding:"required"`
	Video          multipart.FileHeader `json:"video" form:"video" binding:"required"`
	UserID         uuid.UUID            `json:"user_id" form:"user_id"`
	IDCardFileName string               `json:"id_card_filename" form:"id_card_filename" binding:"required"`
	CVFileName     string               `json:"cv_filename" form:"cv_filename" binding:"required"`
	VideoFileName  string               `json:"video_filename" form:"video_filename" binding:"required"`
}

type DecryptRSAResponseDTO struct {
	Name        string `json:"name" form:"name" binding:"required"`
	PhoneNumber string `json:"phone_number" form:"phone_number" binding:"required"`
	IDCard      string `json:"id_card" form:"id_card" binding:"required"`
	CV          string `json:"cv" form:"cv" binding:"required"`
	Video       string `json:"video" form:"video" binding:"required"`
}
