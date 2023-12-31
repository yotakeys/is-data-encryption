package service

import (
	"bytes"
	"context"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"gin-gorm-clean-template/dto"
	"gin-gorm-clean-template/encrypt"
	"gin-gorm-clean-template/entity"
	"gin-gorm-clean-template/helpers"
	"gin-gorm-clean-template/repository"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type EncryptService interface {
	CreateEncrypt(ctx *gin.Context, encryptDTO dto.EncryptCreateDto, userID uuid.UUID, encryptMethod string) (entity.Encrypt, error)
	GetAllEncrypt(ctx context.Context, userID uuid.UUID) ([]dto.DecryptResponseDto, error)
	VerifyDigitalSignature(ctx *gin.Context, req dto.VerifyDigitalSignatureRequest) (dto.VerifyDigitalSignatureResponse, error)
}

type encryptService struct {
	encryptRepository repository.EncryptRepository
	userRepository    repository.UserRepository
}

func NewEncryptService(ur repository.EncryptRepository, us repository.UserRepository) EncryptService {
	return &encryptService{
		encryptRepository: ur,
		userRepository:    us,
	}
}

type (
	DigitalSignature struct {
		Name     string
		Email    string
		DateTime string
	}

	DigitalSignatureContent struct {
		Data      []byte
		Signature []byte
		PublicKey []byte
	}
)

const (
	DataCommentKey      = "Hash"
	SignatureCommentKey = "Signature"
	PublicKeyCommentKey = "PublicKey"
)

func (us *encryptService) CreateEncrypt(ctx *gin.Context, encryptDTO dto.EncryptCreateDto, userID uuid.UUID, encryptMethod string) (entity.Encrypt, error) {
	user, err := us.userRepository.FindUserByID(ctx, userID)
	if err != nil {
		return entity.Encrypt{}, err
	}
	IDCardPath := "uploads/id-card/" + helpers.RandStringBytesRmndr(12) + encryptDTO.IDCardFileName
	CVPath := "uploads/cv/" + helpers.RandStringBytesRmndr(12) + encryptDTO.CVFileName
	VideoPath := "uploads/video/" + helpers.RandStringBytesRmndr(12) + encryptDTO.VideoFileName

	ctx.SaveUploadedFile(&encryptDTO.IDCard, IDCardPath)
	ctx.SaveUploadedFile(&encryptDTO.CV, CVPath)
	ctx.SaveUploadedFile(&encryptDTO.Video, VideoPath)

	CvFile, err := os.ReadFile(CVPath)
	if err != nil {
		return entity.Encrypt{}, err
	}

	messageAdd, err := helpers.WriteDigitalSignature(CvFile, user, encryptDTO.Name)
	if err != nil {
		return entity.Encrypt{}, err
	}
	os.WriteFile(CVPath, messageAdd, 0644)

	var encryptTime float64

	if encryptMethod == "AES" {
		encrypten_name, data, err := encrypt.AESEncrypt(encryptDTO.Name, user.SymmetricKeyAes)
		if err != nil || data == nil {
			return entity.Encrypt{}, err
		}
		floatNumber, err := strconv.ParseFloat(data["elapsed"].(string)[:len(data["elapsed"].(string))-1], 64)
		if err != nil || data == nil {
			return entity.Encrypt{}, err
		}
		encryptTime = encryptTime + floatNumber

		encryptDTO.Name = encrypten_name

		encrypten_phone, data, err := encrypt.AESEncrypt(encryptDTO.PhoneNumber, user.SymmetricKeyAes)

		if err != nil || data == nil {
			return entity.Encrypt{}, err
		}
		floatNumber, err = strconv.ParseFloat(data["elapsed"].(string)[:len(data["elapsed"].(string))-1], 64)
		if err != nil || data == nil {
			return entity.Encrypt{}, err
		}
		encryptTime = encryptTime + floatNumber

		encryptDTO.PhoneNumber = encrypten_phone
		if err != nil || data == nil {
			return entity.Encrypt{}, err
		}

		encrypted_path_idcard, data, err := encrypt.AESEncrypt(IDCardPath, user.SymmetricKeyAes)
		if err != nil || data == nil {
			return entity.Encrypt{}, err
		}
		floatNumber, err = strconv.ParseFloat(data["elapsed"].(string)[:len(data["elapsed"].(string))-1], 64)
		if err != nil || data == nil {
			return entity.Encrypt{}, err
		}
		encryptTime = encryptTime + floatNumber
		IDCardPath = encrypted_path_idcard
		if err != nil || data == nil {
			return entity.Encrypt{}, err
		}

		encrypted_path_cv, data, err := encrypt.AESEncrypt(CVPath, user.SymmetricKeyAes)
		if err != nil || data == nil {
			return entity.Encrypt{}, err
		}
		floatNumber, err = strconv.ParseFloat(data["elapsed"].(string)[:len(data["elapsed"].(string))-1], 64)
		if err != nil || data == nil {
			return entity.Encrypt{}, err
		}
		encryptTime = encryptTime + floatNumber
		CVPath = encrypted_path_cv
		if err != nil || data == nil {
			return entity.Encrypt{}, err
		}

		encrypted_path_video, data, err := encrypt.AESEncrypt(VideoPath, user.SymmetricKeyAes)
		if err != nil || data == nil {
			return entity.Encrypt{}, err
		}
		floatNumber, err = strconv.ParseFloat(data["elapsed"].(string)[:len(data["elapsed"].(string))-1], 64)
		if err != nil || data == nil {
			return entity.Encrypt{}, err
		}
		encryptTime = encryptTime + floatNumber
		VideoPath = encrypted_path_video
		if err != nil || data == nil {
			return entity.Encrypt{}, err
		}

	} else if encryptMethod == "RC4" {
		encrypten_name, data, err := encrypt.RC4Encrypt(encryptDTO.Name, user.SymmetricKeyRc4)
		if err != nil || data == nil {
			return entity.Encrypt{}, err
		}
		floatNumber, err := strconv.ParseFloat(data["elapsed"].(string)[:len(data["elapsed"].(string))-1], 64)
		if err != nil || data == nil {
			return entity.Encrypt{}, err
		}
		encryptTime = encryptTime + floatNumber
		encryptDTO.Name = encrypten_name
		if err != nil || data == nil {
			return entity.Encrypt{}, err
		}
		encrypten_phone, data, err := encrypt.RC4Encrypt(encryptDTO.PhoneNumber, user.SymmetricKeyRc4)
		if err != nil || data == nil {
			return entity.Encrypt{}, err
		}
		floatNumber, err = strconv.ParseFloat(data["elapsed"].(string)[:len(data["elapsed"].(string))-1], 64)
		if err != nil || data == nil {
			return entity.Encrypt{}, err
		}
		encryptTime = encryptTime + floatNumber
		encryptDTO.PhoneNumber = encrypten_phone
		if err != nil || data == nil {
			return entity.Encrypt{}, err
		}

		encrypted_path_idcard, data, err := encrypt.RC4Encrypt(IDCardPath, user.SymmetricKeyRc4)
		if err != nil || data == nil {
			return entity.Encrypt{}, err
		}
		floatNumber, err = strconv.ParseFloat(data["elapsed"].(string)[:len(data["elapsed"].(string))-1], 64)
		if err != nil || data == nil {
			return entity.Encrypt{}, err
		}
		encryptTime = encryptTime + floatNumber
		IDCardPath = encrypted_path_idcard
		if err != nil || data == nil {
			return entity.Encrypt{}, err
		}

		encrypted_path_cv, data, err := encrypt.RC4Encrypt(CVPath, user.SymmetricKeyRc4)
		if err != nil || data == nil {
			return entity.Encrypt{}, err
		}
		floatNumber, err = strconv.ParseFloat(data["elapsed"].(string)[:len(data["elapsed"].(string))-1], 64)
		if err != nil || data == nil {
			return entity.Encrypt{}, err
		}
		encryptTime = encryptTime + floatNumber
		CVPath = encrypted_path_cv
		if err != nil || data == nil {
			return entity.Encrypt{}, err
		}

		encrypted_path_video, data, err := encrypt.RC4Encrypt(VideoPath, user.SymmetricKeyRc4)
		if err != nil || data == nil {
			return entity.Encrypt{}, err
		}
		floatNumber, err = strconv.ParseFloat(data["elapsed"].(string)[:len(data["elapsed"].(string))-1], 64)
		if err != nil || data == nil {
			return entity.Encrypt{}, err
		}
		encryptTime = encryptTime + floatNumber
		VideoPath = encrypted_path_video
		if err != nil || data == nil {
			return entity.Encrypt{}, err
		}

	} else if encryptMethod == "DES" {
		encrypten_name, data, err := encrypt.DESEncrypt(encryptDTO.Name, user.SymmetricKeyDes)
		if err != nil || data == nil {
			return entity.Encrypt{}, err
		}
		floatNumber, err := strconv.ParseFloat(data["elapsed"].(string)[:len(data["elapsed"].(string))-1], 64)
		if err != nil || data == nil {
			return entity.Encrypt{}, err
		}
		encryptTime = encryptTime + floatNumber
		encryptDTO.Name = encrypten_name
		if err != nil || data == nil {
			return entity.Encrypt{}, err
		}
		encrypten_phone, data, err := encrypt.DESEncrypt(encryptDTO.PhoneNumber, user.SymmetricKeyDes)
		if err != nil || data == nil {
			return entity.Encrypt{}, err
		}
		floatNumber, err = strconv.ParseFloat(data["elapsed"].(string)[:len(data["elapsed"].(string))-1], 64)
		if err != nil || data == nil {
			return entity.Encrypt{}, err
		}
		encryptTime = encryptTime + floatNumber
		encryptDTO.PhoneNumber = encrypten_phone
		if err != nil || data == nil {
			return entity.Encrypt{}, err
		}

		encrypted_path_idcard, data, err := encrypt.DESEncrypt(IDCardPath, user.SymmetricKeyDes)
		if err != nil || data == nil {
			return entity.Encrypt{}, err
		}
		floatNumber, err = strconv.ParseFloat(data["elapsed"].(string)[:len(data["elapsed"].(string))-1], 64)
		if err != nil || data == nil {
			return entity.Encrypt{}, err
		}
		encryptTime = encryptTime + floatNumber
		IDCardPath = encrypted_path_idcard
		if err != nil || data == nil {
			return entity.Encrypt{}, err
		}

		encrypted_path_cv, data, err := encrypt.DESEncrypt(CVPath, user.SymmetricKeyDes)
		if err != nil || data == nil {
			return entity.Encrypt{}, err
		}
		floatNumber, err = strconv.ParseFloat(data["elapsed"].(string)[:len(data["elapsed"].(string))-1], 64)
		if err != nil || data == nil {
			return entity.Encrypt{}, err
		}
		encryptTime = encryptTime + floatNumber
		CVPath = encrypted_path_cv
		if err != nil || data == nil {
			return entity.Encrypt{}, err
		}

		encrypted_path_video, data, err := encrypt.DESEncrypt(VideoPath, user.SymmetricKeyDes)
		if err != nil || data == nil {
			return entity.Encrypt{}, err
		}
		floatNumber, err = strconv.ParseFloat(data["elapsed"].(string)[:len(data["elapsed"].(string))-1], 64)
		if err != nil || data == nil {
			return entity.Encrypt{}, err
		}
		encryptTime = encryptTime + floatNumber
		VideoPath = encrypted_path_video
		if err != nil || data == nil {
			return entity.Encrypt{}, err
		}
	}

	encrypt := entity.Encrypt{
		Name:          encryptDTO.Name,
		PhoneNumber:   encryptDTO.PhoneNumber,
		CVUrl:         CVPath,
		IDCardUrl:     IDCardPath,
		VideoUrl:      VideoPath,
		EncryptMethod: encryptMethod,
		EncryptTime:   strconv.FormatFloat(encryptTime, 'f', -1, 64),
		UserID:        userID,
	}

	return us.encryptRepository.CreateEncrypt(ctx.Request.Context(), encrypt)
}

func (us *encryptService) GetAllEncrypt(ctx context.Context, userID uuid.UUID) ([]dto.DecryptResponseDto, error) {

	user, err := us.userRepository.FindUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	datas, err := us.encryptRepository.GetAllEncrypt(ctx, userID)

	if err != nil {
		return nil, err
	}

	for i, data := range datas {
		if data.EncryptMethod == "AES" {
			decrypt_name, err := encrypt.AESDecrypt(data.Name, user.SymmetricKeyAes)
			if err != nil {
				return nil, err
			}
			decrypt_phone, err := encrypt.AESDecrypt(data.PhoneNumber, user.SymmetricKeyAes)
			if err != nil {
				return nil, err
			}
			decrypt_idcard, err := encrypt.AESDecrypt(data.IDCardUrl, user.SymmetricKeyAes)
			if err != nil {
				return nil, err
			}
			decrypt_cv, err := encrypt.AESDecrypt(data.CVUrl, user.SymmetricKeyAes)
			if err != nil {
				return nil, err
			}
			decrypt_video, err := encrypt.AESDecrypt(data.VideoUrl, user.SymmetricKeyAes)
			if err != nil {
				return nil, err
			}
			datas[i].Name = decrypt_name
			datas[i].PhoneNumber = decrypt_phone
			datas[i].IDCardUrl = decrypt_idcard
			datas[i].CVUrl = decrypt_cv
			datas[i].VideoUrl = decrypt_video
		} else if data.EncryptMethod == "RC4" {
			decrypt_name, err := encrypt.RC4Decrypt(data.Name, user.SymmetricKeyRc4)
			if err != nil {
				return nil, err
			}
			decrypt_phone, err := encrypt.RC4Decrypt(data.PhoneNumber, user.SymmetricKeyRc4)
			if err != nil {
				return nil, err
			}
			decrypt_idcard, err := encrypt.RC4Decrypt(data.IDCardUrl, user.SymmetricKeyRc4)
			if err != nil {
				return nil, err
			}
			decrypt_cv, err := encrypt.RC4Decrypt(data.CVUrl, user.SymmetricKeyRc4)
			if err != nil {
				return nil, err
			}
			decrypt_video, err := encrypt.RC4Decrypt(data.VideoUrl, user.SymmetricKeyRc4)
			if err != nil {
				return nil, err
			}
			datas[i].Name = decrypt_name
			datas[i].PhoneNumber = decrypt_phone
			datas[i].IDCardUrl = decrypt_idcard
			datas[i].CVUrl = decrypt_cv
			datas[i].VideoUrl = decrypt_video

		} else if data.EncryptMethod == "DES" {
			decrypt_name, err := encrypt.DESDecrypt(data.Name, user.SymmetricKeyDes)
			if err != nil {
				return nil, err
			}
			decrypt_phone, err := encrypt.DESDecrypt(data.PhoneNumber, user.SymmetricKeyDes)
			if err != nil {
				return nil, err
			}
			decrypt_idcard, err := encrypt.DESDecrypt(data.IDCardUrl, user.SymmetricKeyDes)
			if err != nil {
				return nil, err
			}
			decrypt_cv, err := encrypt.DESDecrypt(data.CVUrl, user.SymmetricKeyDes)
			if err != nil {
				return nil, err
			}
			decrypt_video, err := encrypt.DESDecrypt(data.VideoUrl, user.SymmetricKeyDes)
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

	response := make([]dto.DecryptResponseDto, len(datas))
	for i, data := range datas {
		response[i].ID = data.ID
		response[i].Name = data.Name
		response[i].PhoneNumber = data.PhoneNumber
		response[i].CVUrl = data.CVUrl
		response[i].IDCardUrl = data.IDCardUrl
		response[i].VideoUrl = data.VideoUrl
		response[i].EncryptMethod = data.EncryptMethod
		response[i].EncryptTime = data.EncryptTime
		response[i].UserID = data.UserID
	}

	return response, nil
}

func (us *encryptService) VerifyDigitalSignature(ctx *gin.Context, req dto.VerifyDigitalSignatureRequest) (dto.VerifyDigitalSignatureResponse, error) {
	var data DigitalSignature
	parts := strings.Split(req.Files.Filename, ".")
	ext := parts[len(parts)-1]
	CVPath := "uploads/cv/" + helpers.RandStringBytesRmndr(12) + req.Files.Filename + "." + ext
	ctx.SaveUploadedFile(&req.Files, CVPath)
	file, err := os.ReadFile(CVPath)
	if err != nil {
		return dto.VerifyDigitalSignatureResponse{}, err
	}

	readContent, err := helpers.ReadDigitalSignature(file)
	if err != nil {
		return dto.VerifyDigitalSignatureResponse{}, err
	}

	decryptPubKey := readContent.PublicKey
	if len(decryptPubKey) == 0 {
		return dto.VerifyDigitalSignatureResponse{}, errors.New("digital signature is invalid")
	}

	pubKey, err := encrypt.ParsePublicKeyFromPEM(string(decryptPubKey))
	if err != nil {
		return dto.VerifyDigitalSignatureResponse{}, err
	}

	datas := bytes.Split(file, []byte("\n%"+DataCommentKey))
	if len(datas) < 2 {
		return dto.VerifyDigitalSignatureResponse{}, errors.New("digital signature is missing")
	}

	content := datas[0]
	signature := readContent.Data
	hash := sha256.Sum256(content)

	err = json.Unmarshal(readContent.Signature, &data)
	if err != nil {
		return dto.VerifyDigitalSignatureResponse{}, err
	}

	err = rsa.VerifyPKCS1v15(pubKey, crypto.SHA256, hash[:], signature)
	if err != nil {
		return dto.VerifyDigitalSignatureResponse{}, errors.New("digital signature is invalid")
	}

	return dto.VerifyDigitalSignatureResponse{
		DigitalSignatureRes: dto.DigitalSignatureRes{
			Name:  data.Name,
			Email: data.Email,
			Date:  data.DateTime,
		},
	}, nil
}
