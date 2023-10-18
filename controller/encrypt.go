package controller

import (
	"gin-gorm-clean-template/common"
	"gin-gorm-clean-template/dto"
	"gin-gorm-clean-template/entity"
	"gin-gorm-clean-template/service"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type EncryptController interface {
	CreateEncrypt(ctx *gin.Context)
	GetAllEncrypt(ctx *gin.Context)
	GetFile(ctx *gin.Context)
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

	resultAES, err := uc.encryptService.CreateEncrypt(ctx, encrypt, userID, "AES")
	if err != nil {
		res := common.BuildErrorResponse("Gagal Menambahkan Encrypt", err.Error(), common.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	// RC4

	resultRC4, err := uc.encryptService.CreateEncrypt(ctx, encrypt, userID, "RC4")
	if err != nil {
		res := common.BuildErrorResponse("Gagal Menambahkan Encrypt", err.Error(), common.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	// DES

	resultDES, err := uc.encryptService.CreateEncrypt(ctx, encrypt, userID, "DES")
	if err != nil {
		res := common.BuildErrorResponse("Gagal Menambahkan Encrypt", err.Error(), common.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	var result [3]entity.Encrypt
	result[0] = resultAES
	result[1] = resultRC4
	result[2] = resultDES

	res := common.BuildResponse(true, "Berhasil Menambahkan Encrypt", resultRC4)
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

func (uc *encryptController) GetFile(ctx *gin.Context) {
	filePath := ctx.Query("file_path")

	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		res := common.BuildErrorResponse("Gagal Mendapatkan File", err.Error(), common.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	ctx.File(filePath)
}
