package repository

import (
	"context"
	"gin-gorm-clean-template/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EncryptRepository interface {
	CreateEncrypt(ctx context.Context, Encrypt entity.Encrypt) (entity.Encrypt, error)
	GetAllEncrypt(ctx context.Context, userID uuid.UUID) ([]entity.Encrypt, error)
	FindEncryptByID(ctx context.Context, EncryptID uuid.UUID) (entity.Encrypt, error)
}

type EncryptConnection struct {
	connection *gorm.DB
}

func NewEncryptRepository(db *gorm.DB) EncryptRepository {
	return &EncryptConnection{
		connection: db,
	}
}

func (db *EncryptConnection) CreateEncrypt(ctx context.Context, Encrypt entity.Encrypt) (entity.Encrypt, error) {
	Encrypt.ID = uuid.New()
	uc := db.connection.Create(&Encrypt)
	if uc.Error != nil {
		return entity.Encrypt{}, uc.Error
	}
	return Encrypt, nil
}

func (db *EncryptConnection) GetAllEncrypt(ctx context.Context, userID uuid.UUID) ([]entity.Encrypt, error) {
	var listEncrypt []entity.Encrypt
	tx := db.connection.Where("user_id = ?", userID).Take(&listEncrypt)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return listEncrypt, nil
}

func (db *EncryptConnection) FindEncryptByID(ctx context.Context, EncryptID uuid.UUID) (entity.Encrypt, error) {
	var Encrypt entity.Encrypt
	ux := db.connection.Where("id = ?", EncryptID).Take(&Encrypt)
	if ux.Error != nil {
		return Encrypt, ux.Error
	}
	return Encrypt, nil
}
