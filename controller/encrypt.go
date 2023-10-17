package controller

import (
	"gin-gorm-clean-template/common"
	"gin-gorm-clean-template/dto"
	"gin-gorm-clean-template/entity"
	"gin-gorm-clean-template/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type EncryptController interface {
	CreateEncrypt(ctx *gin.Context)
	GetAllEncrypt(ctx *gin.Context)
}

type encryptController struct {
	jwtService     service.JWTService
	encryptService service.EncryptService
}

func NewEncryptController(us service.EncryptService, jwts service.JWTService) EncryptController {
	return &encryptController{
		encryptService: us,
		jwtService:     jwts,
	}
}

func (uc *encryptController) CreateEncrypt(ctx *gin.Context) {
	var encrypt dto.EncryptCreateDto
	err := ctx.ShouldBind(&encrypt)

	if err != nil {
		res := common.BuildErrorResponse("Gagal Memproses Request", err.Error(), common.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	token := ctx.MustGet("token").(string)
	userID, err := uc.jwtService.GetUserIDByToken(token)

	if err != nil {
		res := common.BuildErrorResponse("Gagal Memproses Request", err.Error(), common.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	// AES

	encrypt.IDCard.Filename = userID.String() + "-aes-" + encrypt.IDCard.Filename
	encrypt.CV.Filename = userID.String() + "-aes-" + encrypt.CV.Filename
	encrypt.Video.Filename = userID.String() + "-aes-" + encrypt.Video.Filename

	ctx.SaveUploadedFile(encrypt.IDCard, "uploads/id-card/"+encrypt.IDCard.Filename)
	ctx.SaveUploadedFile(encrypt.CV, "uploads/cv/"+encrypt.CV.Filename)
	ctx.SaveUploadedFile(encrypt.Video, "uploads/video/"+encrypt.Video.Filename)

	resultAES, err := uc.encryptService.CreateEncrypt(ctx.Request.Context(), encrypt, userID, "AES", "1s")
	if err != nil {
		res := common.BuildErrorResponse("Gagal Menambahkan Encrypt", err.Error(), common.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	// RC4

	encrypt.IDCard.Filename = userID.String() + "-rc4-" + encrypt.IDCard.Filename
	encrypt.CV.Filename = userID.String() + "-rc4-" + encrypt.CV.Filename
	encrypt.Video.Filename = userID.String() + "-rc4-" + encrypt.Video.Filename

	ctx.SaveUploadedFile(encrypt.IDCard, "uploads/id-card/"+encrypt.IDCard.Filename)
	ctx.SaveUploadedFile(encrypt.CV, "uploads/cv/"+encrypt.CV.Filename)
	ctx.SaveUploadedFile(encrypt.Video, "uploads/video/"+encrypt.Video.Filename)

	resultRC4, err := uc.encryptService.CreateEncrypt(ctx.Request.Context(), encrypt, userID, "RC4", "1s")
	if err != nil {
		res := common.BuildErrorResponse("Gagal Menambahkan Encrypt", err.Error(), common.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	// DES

	encrypt.IDCard.Filename = userID.String() + "-des-" + encrypt.IDCard.Filename
	encrypt.CV.Filename = userID.String() + "-des-" + encrypt.CV.Filename
	encrypt.Video.Filename = userID.String() + "-des-" + encrypt.Video.Filename

	ctx.SaveUploadedFile(encrypt.IDCard, "uploads/id-card/"+encrypt.IDCard.Filename)
	ctx.SaveUploadedFile(encrypt.CV, "uploads/cv/"+encrypt.CV.Filename)
	ctx.SaveUploadedFile(encrypt.Video, "uploads/video/"+encrypt.Video.Filename)

	resultDES, err := uc.encryptService.CreateEncrypt(ctx.Request.Context(), encrypt, userID, "DES", "1s")
	if err != nil {
		res := common.BuildErrorResponse("Gagal Menambahkan Encrypt", err.Error(), common.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	var result [3]entity.Encrypt
	result[0] = resultAES
	result[1] = resultRC4
	result[2] = resultDES

	res := common.BuildResponse(true, "Berhasil Menambahkan Encrypt", result)
	ctx.JSON(http.StatusOK, res)
}

func (uc *encryptController) GetAllEncrypt(ctx *gin.Context) {
	token := ctx.MustGet("token").(string)
	userID, err := uc.jwtService.GetUserIDByToken(token)

	if err != nil {
		res := common.BuildErrorResponse("Gagal Memproses Request", err.Error(), common.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	result, err := uc.encryptService.GetAllEncrypt(ctx.Request.Context(), userID)

	if err != nil {
		res := common.BuildErrorResponse("Gagal Mendapatkan List Encrypt", err.Error(), common.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := common.BuildResponse(true, "Berhasil Mendapatkan List Encrypt", result)
	ctx.JSON(http.StatusOK, res)
}
