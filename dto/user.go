package dto

import (
	"github.com/google/uuid"
)

type UserCreateDto struct {
	ID        	uuid.UUID   `gorm:"primary_key" json:"id" form:"id"`
	Name 		string 		`json:"name" form:"name" binding:"required"`
	Email 		string 		`json:"email" form:"email" binding:"required"`
	NoTelp 		string 		`json:"no_telp" form:"no_telp" binding:"required"`
	Password 	string  	`json:"password" form:"password" binding:"required"`
}

type UserUpdateDto struct {
	ID        	uuid.UUID   `gorm:"primary_key" json:"id" form:"id"`
	Name 		string 		`json:"name" form:"name"`
	Email 		string 		`json:"email" form:"email"`
	NoTelp 		string 		`json:"no_telp" form:"no_telp"`
	Password 	string  	`json:"password" form:"password"`
}

type UserLoginDTO struct {
	Email 		string 		`json:"email" form:"email" binding:"email"`
	Password 	string  	`json:"password" form:"password" binding:"required"`
}