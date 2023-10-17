package service

import (
	"context"
	"gin-gorm-clean-template/dto"
	"gin-gorm-clean-template/encrypt"
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

	if encryptMethod == "AES" {
		encrypten_name, data, err := encrypt.AESEncrypt(encryptDTO.Name)
		encryptDTO.Name = encrypten_name
		if err != nil || data == nil {
			return entity.Encrypt{}, err
		}
		encrypten_phone, data, err := encrypt.AESEncrypt(encryptDTO.PhoneNumber)
		encryptDTO.PhoneNumber = encrypten_phone
		if err != nil || data == nil {
			return entity.Encrypt{}, err
		}
	} else if encryptMethod == "RC4" {
		encrypten_name, data, err := encrypt.RC4Encrypt(encryptDTO.Name)
		encryptDTO.Name = encrypten_name
		if err != nil || data == nil {
			return entity.Encrypt{}, err
		}
		encrypten_phone, data, err := encrypt.RC4Encrypt(encryptDTO.PhoneNumber)
		encryptDTO.PhoneNumber = encrypten_phone
		if err != nil || data == nil {
			return entity.Encrypt{}, err
		}
	} else if encryptMethod == "DES" {
		encrypten_name, data, err := encrypt.DESEncrypt(encryptDTO.Name)
		encryptDTO.Name = encrypten_name
		if err != nil || data == nil {
			return entity.Encrypt{}, err
		}
		encrypten_phone, data, err := encrypt.DESEncrypt(encryptDTO.PhoneNumber)
		encryptDTO.PhoneNumber = encrypten_phone
		if err != nil || data == nil {
			return entity.Encrypt{}, err
		}
	}

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

	datas, err := us.encryptRepository.GetAllEncrypt(ctx, userID)

	if err != nil {
		return nil, err
	}

	for i, data := range datas {
		if data.EncryptMethod == "AES" {
			decrypt_name, err := encrypt.AESDecrypt(data.Name)
			if err != nil {
				return nil, err
			}
			decrypt_phone, err := encrypt.AESDecrypt(data.PhoneNumber)
			if err != nil {
				return nil, err
			}
			datas[i].Name = decrypt_name
			datas[i].PhoneNumber = decrypt_phone
		} else if data.EncryptMethod == "RC4" {
			decrypt_name, err := encrypt.RC4Decrypt(data.Name)
			if err != nil {
				return nil, err
			}
			decrypt_phone, err := encrypt.RC4Decrypt(data.PhoneNumber)
			if err != nil {
				return nil, err
			}
			datas[i].Name = decrypt_name
			datas[i].PhoneNumber = decrypt_phone

		} else if data.EncryptMethod == "DES" {
			decrypt_name, err := encrypt.DESDecrypt(data.Name)
			if err != nil {
				return nil, err
			}
			decrypt_phone, err := encrypt.DESDecrypt(data.PhoneNumber)
			if err != nil {
				return nil, err
			}
			datas[i].Name = decrypt_name
			datas[i].PhoneNumber = decrypt_phone
		}
	}

	return datas, nil
}
