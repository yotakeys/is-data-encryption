package service

import (
	"context"
	"gin-gorm-clean-template/dto"
	"gin-gorm-clean-template/encrypt"
	"gin-gorm-clean-template/entity"
	"gin-gorm-clean-template/repository"
	"math/rand"

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

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytesRmndr(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

func (us *encryptService) CreateEncrypt(ctx *gin.Context, encryptDTO dto.EncryptCreateDto, userID uuid.UUID, encryptMethod string, encryptTime string) (entity.Encrypt, error) {

	encryptDTO.IDCard.Filename = "uploads/id-card/" + RandStringBytesRmndr(5) + encryptDTO.IDCard.Filename
	encryptDTO.CV.Filename = "uploads/cv/" + RandStringBytesRmndr(5) + encryptDTO.CV.Filename
	encryptDTO.Video.Filename = "uploads/video/" + RandStringBytesRmndr(5) + encryptDTO.Video.Filename

	ctx.SaveUploadedFile(encryptDTO.IDCard, encryptDTO.IDCard.Filename)
	ctx.SaveUploadedFile(encryptDTO.CV, encryptDTO.CV.Filename)
	ctx.SaveUploadedFile(encryptDTO.Video, encryptDTO.Video.Filename)

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

		encrypted_path_idcard, data, err := encrypt.AESEncrypt(encryptDTO.IDCard.Filename)
		encryptDTO.IDCard.Filename = encrypted_path_idcard
		if err != nil || data == nil {
			return entity.Encrypt{}, err
		}

		encrypted_path_cv, data, err := encrypt.AESEncrypt(encryptDTO.CV.Filename)
		encryptDTO.CV.Filename = encrypted_path_cv
		if err != nil || data == nil {
			return entity.Encrypt{}, err
		}

		encrypted_path_video, data, err := encrypt.AESEncrypt(encryptDTO.Video.Filename)
		encryptDTO.Video.Filename = encrypted_path_video
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

		encrypted_path_idcard, data, err := encrypt.RC4Encrypt(encryptDTO.IDCard.Filename)
		encryptDTO.IDCard.Filename = encrypted_path_idcard
		if err != nil || data == nil {
			return entity.Encrypt{}, err
		}

		encrypted_path_cv, data, err := encrypt.RC4Encrypt(encryptDTO.CV.Filename)
		encryptDTO.CV.Filename = encrypted_path_cv
		if err != nil || data == nil {
			return entity.Encrypt{}, err
		}

		encrypted_path_video, data, err := encrypt.RC4Encrypt(encryptDTO.Video.Filename)
		encryptDTO.Video.Filename = encrypted_path_video
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

		encrypted_path_idcard, data, err := encrypt.DESEncrypt(encryptDTO.IDCard.Filename)
		encryptDTO.IDCard.Filename = encrypted_path_idcard
		if err != nil || data == nil {
			return entity.Encrypt{}, err
		}

		encrypted_path_cv, data, err := encrypt.DESEncrypt(encryptDTO.CV.Filename)
		encryptDTO.CV.Filename = encrypted_path_cv
		if err != nil || data == nil {
			return entity.Encrypt{}, err
		}

		encrypted_path_video, data, err := encrypt.DESEncrypt(encryptDTO.Video.Filename)
		encryptDTO.Video.Filename = encrypted_path_video
		if err != nil || data == nil {
			return entity.Encrypt{}, err
		}

	}

	encrypt := entity.Encrypt{
		Name:          encryptDTO.Name,
		PhoneNumber:   encryptDTO.PhoneNumber,
		CVUrl:         encryptDTO.CV.Filename,
		IDCardUrl:     encryptDTO.IDCard.Filename,
		VideoUrl:      encryptDTO.Video.Filename,
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
			decrypt_idcard, err := encrypt.AESDecrypt(data.IDCardUrl)
			if err != nil {
				return nil, err
			}
			decrypt_cv, err := encrypt.AESDecrypt(data.CVUrl)
			if err != nil {
				return nil, err
			}
			decrypt_video, err := encrypt.AESDecrypt(data.VideoUrl)
			if err != nil {
				return nil, err
			}
			datas[i].Name = decrypt_name
			datas[i].PhoneNumber = decrypt_phone
			datas[i].IDCardUrl = decrypt_idcard
			datas[i].CVUrl = decrypt_cv
			datas[i].VideoUrl = decrypt_video
		} else if data.EncryptMethod == "RC4" {
			decrypt_name, err := encrypt.RC4Decrypt(data.Name)
			if err != nil {
				return nil, err
			}
			decrypt_phone, err := encrypt.RC4Decrypt(data.PhoneNumber)
			if err != nil {
				return nil, err
			}
			decrypt_idcard, err := encrypt.RC4Decrypt(data.IDCardUrl)
			if err != nil {
				return nil, err
			}
			decrypt_cv, err := encrypt.RC4Decrypt(data.CVUrl)
			if err != nil {
				return nil, err
			}
			decrypt_video, err := encrypt.RC4Decrypt(data.VideoUrl)
			if err != nil {
				return nil, err
			}
			datas[i].Name = decrypt_name
			datas[i].PhoneNumber = decrypt_phone
			datas[i].IDCardUrl = decrypt_idcard
			datas[i].CVUrl = decrypt_cv
			datas[i].VideoUrl = decrypt_video

		} else if data.EncryptMethod == "DES" {
			decrypt_name, err := encrypt.DESDecrypt(data.Name)
			if err != nil {
				return nil, err
			}
			decrypt_phone, err := encrypt.DESDecrypt(data.PhoneNumber)
			if err != nil {
				return nil, err
			}
			decrypt_idcard, err := encrypt.DESDecrypt(data.IDCardUrl)
			if err != nil {
				return nil, err
			}
			decrypt_cv, err := encrypt.DESDecrypt(data.CVUrl)
			if err != nil {
				return nil, err
			}
			decrypt_video, err := encrypt.DESDecrypt(data.VideoUrl)
			if err != nil {
				return nil, err
			}
			datas[i].Name = decrypt_name
			datas[i].PhoneNumber = decrypt_phone
			datas[i].IDCardUrl = decrypt_idcard
			datas[i].CVUrl = decrypt_cv
			datas[i].VideoUrl = decrypt_video
		}
	}

	return datas, nil
}
