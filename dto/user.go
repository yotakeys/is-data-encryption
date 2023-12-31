package dto

import (
	"github.com/google/uuid"
)

type UserCreateDto struct {
	ID       uuid.UUID `gorm:"primary_key" json:"id" form:"id"`
	Email    string    `json:"email" form:"email" binding:"required"`
	Password string    `json:"password" form:"password" binding:"required"`
}

type UserUpdateDto struct {
	ID       uuid.UUID `gorm:"primary_key" json:"id" form:"id"`
	Email    string    `json:"email" form:"email"`
	Password string    `json:"password" form:"password"`
}

type UserLoginDTO struct {
	Email    string `json:"email" form:"email" binding:"email"`
	Password string `json:"password" form:"password" binding:"required"`
}

type UserEmail struct {
	Email string `json:"email" form:"email" binding:"email"`
}

type AsymmetricUserReqeustDTO struct {
	RequestingUserEmail string `json:"requesting_user_email" form:"requesting_user_email" binding:"requesting_user_email"`
	RequestedUserEmail  string `json:"requested_user_email" form:"requested_user_email" binding:"requested_user_email"`
}
