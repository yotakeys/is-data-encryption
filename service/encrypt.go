package service

import (
	"context"
	"gin-gorm-clean-template/dto"
	"gin-gorm-clean-template/entity"
	"gin-gorm-clean-template/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type EncryptService interface {
	CreateEncrypt(ctx *gin.Context, encryptDTO dto.EncryptCreateDto, userID uuid.UUID, encryptMethod string, encryptTime string) (entity.Encrypt, error)
	GetAllEncrypt(ctx context.Context, userID uuid.UUID) ([]entity.Encrypt, error)
}

type encryptService struct {
	encryptRepository repository.EncryptRepository
}

func NewEncryptService(ur repository.EncryptRepository) EncryptService {
	return &encryptService{
		encryptRepository: ur,
	}
}

func (us *encryptService) CreateEncrypt(ctx *gin.Context, encryptDTO dto.EncryptCreateDto, userID uuid.UUID, encryptMethod string, encryptTime string) (entity.Encrypt, error) {
	encryptDTO.IDCard.Filename = userID.String() + "-" + encryptMethod + "-" + encryptDTO.IDCard.Filename
	encryptDTO.CV.Filename = userID.String() + "-" + encryptMethod + "-" + encryptDTO.CV.Filename
	encryptDTO.Video.Filename = userID.String() + "-" + encryptMethod + "-" + encryptDTO.Video.Filename

	ctx.SaveUploadedFile(encryptDTO.IDCard, "uploads/id-card/"+encryptDTO.IDCard.Filename)
	ctx.SaveUploadedFile(encryptDTO.CV, "uploads/cv/"+encryptDTO.CV.Filename)
	ctx.SaveUploadedFile(encryptDTO.Video, "uploads/video/"+encryptDTO.Video.Filename)

	encrypt := entity.Encrypt{
		Name:          encryptDTO.Name,
		PhoneNumber:   encryptDTO.PhoneNumber,
		CVUrl:         "uploads/cv/" + encryptDTO.CV.Filename,
		IDCardUrl:     "uploads/id-card/" + encryptDTO.IDCard.Filename,
		VideoUrl:      "uploads/video/" + encryptDTO.Video.Filename,
		EncryptMethod: encryptMethod,
		EncryptTime:   encryptTime,
		UserID:        userID,
	}

	return us.encryptRepository.CreateEncrypt(ctx.Request.Context(), encrypt)
}

func (us *encryptService) GetAllEncrypt(ctx context.Context, userID uuid.UUID) ([]entity.Encrypt, error) {
	return us.encryptRepository.GetAllEncrypt(ctx, userID)
}
