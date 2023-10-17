package controller

import (
	"gin-gorm-clean-template/common"
	"gin-gorm-clean-template/dto"
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

	ctx.SaveUploadedFile(encrypt.IDCard, "uploads/id-card/"+encrypt.IDCard.Filename)
	ctx.SaveUploadedFile(encrypt.CV, "uploads/cv/"+encrypt.CV.Filename)
	ctx.SaveUploadedFile(encrypt.Video, "uploads/video/"+encrypt.Video.Filename)

	result, err := uc.encryptService.CreateEncrypt(ctx.Request.Context(), encrypt)
	if err != nil {
		res := common.BuildErrorResponse("Gagal Menambahkan Encrypt", err.Error(), common.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := common.BuildResponse(true, "Berhasil Menambahkan Encrypt", result)
	ctx.JSON(http.StatusOK, res)
}

func (uc *encryptController) GetAllEncrypt(ctx *gin.Context) {
	result, err := uc.encryptService.GetAllEncrypt(ctx.Request.Context())
	if err != nil {
		res := common.BuildErrorResponse("Gagal Mendapatkan List Encrypt", err.Error(), common.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := common.BuildResponse(true, "Berhasil Mendapatkan List Encrypt", result)
	ctx.JSON(http.StatusOK, res)
}
