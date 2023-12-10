package entity

import (
	// "gin-gorm-clean-template/helpers"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID              uuid.UUID `gorm:"primary_key;not_null" json:"id"`
	Email           string    `json:"email" binding:"email"`
	Password        string    `json:"password"`
	SymmetricKeyAes string    `json:"symmetric_key_aes"`
	SymmetricKeyDes string    `json:"symmetric_key_des"`
	SymmetricKeyRc4 string    `json:"symmetric_key_rc4"`
	PublicKey       string    `json:"public_key"`
	PrivateKey      string    `json:"private_key"`

	Timestamp
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	var err error
	bytes, err := bcrypt.GenerateFromPassword([]byte(u.Password), 4)
	if err != nil {
		return err
	}
	u.Password = string(bytes)
	return nil
}
