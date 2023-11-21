package repository

import (
	"context"
	"gin-gorm-clean-template/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	RegisterUser(ctx context.Context, user entity.User) (entity.User, error)
	GetAllUser(ctx context.Context) ([]entity.User, error)
	FindUserByEmail(ctx context.Context, email string) (entity.User, error)
	FindUserByID(ctx context.Context, userID uuid.UUID) (entity.User, error)
	DeleteUser(ctx context.Context, userID uuid.UUID) error
	UpdateUser(ctx context.Context, user entity.User) error
	CreateAsymmetric(ctx context.Context, asymmetric entity.Asymmetric) (entity.Asymmetric, error)
}

type userConnection struct {
	connection *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userConnection{
		connection: db,
	}
}

func (db *userConnection) RegisterUser(ctx context.Context, user entity.User) (entity.User, error) {
	user.ID = uuid.New()
	uc := db.connection.Create(&user)
	if uc.Error != nil {
		return entity.User{}, uc.Error
	}
	return user, nil
}

func (db *userConnection) GetAllUser(ctx context.Context) ([]entity.User, error) {
	var listUser []entity.User
	tx := db.connection.Find(&listUser)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return listUser, nil
}

func (db *userConnection) FindUserByEmail(ctx context.Context, email string) (entity.User, error) {
	var user entity.User
	ux := db.connection.Where("email = ?", email).Take(&user)
	if ux.Error != nil {
		return user, ux.Error
	}
	return user, nil
}

func (db *userConnection) FindUserByID(ctx context.Context, userID uuid.UUID) (entity.User, error) {
	var user entity.User
	ux := db.connection.Where("id = ?", userID).Take(&user)
	if ux.Error != nil {
		return user, ux.Error
	}
	return user, nil
}

func (db *userConnection) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	uc := db.connection.Delete(&entity.User{}, &userID)
	if uc.Error != nil {
		return uc.Error
	}
	return nil
}

func (db *userConnection) UpdateUser(ctx context.Context, user entity.User) error {
	uc := db.connection.Updates(&user)
	if uc.Error != nil {
		return uc.Error
	}
	return nil
}

func (db *userConnection) CreateAsymmetric(ctx context.Context, asymmetric entity.Asymmetric) (entity.Asymmetric, error) {
	asymmetric.ID = uuid.New()
	uc := db.connection.Create(&asymmetric)
	if uc.Error != nil {
		return entity.Asymmetric{}, uc.Error
	}
	return asymmetric, nil
}

func (db *userConnection) FindAsymmetric(ctx context.Context, userID uuid.UUID) ([]entity.Asymmetric, error) {
	var listAsymmetric []entity.Asymmetric
	tx := db.connection.Where("user_id = ?", userID).Find(&listAsymmetric)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return listAsymmetric, nil
}
