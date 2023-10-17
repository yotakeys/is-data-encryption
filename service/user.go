package service

import (
	"context"
	"gin-gorm-clean-template/dto"
	"gin-gorm-clean-template/entity"
	"gin-gorm-clean-template/repository"
)

type EncryptService interface {
	CreateEncrypt(ctx context.Context, encryptDTO dto.EncryptCreateDto) (entity.Encrypt, error)
	GetAllEncrypt(ctx context.Context) ([]entity.Encrypt, error)
}

type encryptService struct {
	encryptRepository repository.EncryptRepository
}

func NewEncryptService(ur repository.EncryptRepository) EncryptService {
	return &encryptService{
		encryptRepository: ur,
	}
}

func (us *encryptService) CreateEncrypt(ctx context.Context, encryptDTO dto.EncryptCreateDto) (entity.Encrypt, error) {
	encrypt := entity.Encrypt{
		Name:          encryptDTO.Name,
		PhoneNumber:   encryptDTO.PhoneNumber,
		CVUrl:         "uploads/cv/" + encryptDTO.CV.Filename,
		IDCardUrl:     "uploads/id-card/" + encryptDTO.IDCard.Filename,
		VideoUrl:      "uploads/video/" + encryptDTO.Video.Filename,
		EncryptMethod: "LMAO",
		EncryptTime:   "LMAO",
	}

	return us.encryptRepository.CreateEncrypt(ctx, encrypt)
}

func (us *encryptService) GetAllEncrypt(ctx context.Context) ([]entity.Encrypt, error) {
	return us.encryptRepository.GetAllEncrypt(ctx)
}
