package controller

import (
	"gin-gorm-clean-template/common"
	"gin-gorm-clean-template/dto"
	"gin-gorm-clean-template/entity"
	"gin-gorm-clean-template/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mashingan/smapping"
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
	var encrypt2 dto.EncryptCreateDto
	var encrypt3 dto.EncryptCreateDto
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

	err = smapping.FillStruct(&encrypt2, smapping.MapFields(&encrypt))

	if err != nil {
		res := common.BuildErrorResponse("Gagal Memproses Request", err.Error(), common.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	err = smapping.FillStruct(&encrypt3, smapping.MapFields(&encrypt))

	if err != nil {
		res := common.BuildErrorResponse("Gagal Memproses Request", err.Error(), common.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	// AES

	resultAES, err := uc.encryptService.CreateEncrypt(ctx, encrypt, userID, "AES", "1s")
	if err != nil {
		res := common.BuildErrorResponse("Gagal Menambahkan Encrypt", err.Error(), common.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	// RC4

	resultRC4, err := uc.encryptService.CreateEncrypt(ctx, encrypt2, userID, "RC4", "1s")
	if err != nil {
		res := common.BuildErrorResponse("Gagal Menambahkan Encrypt", err.Error(), common.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	// DES

	resultDES, err := uc.encryptService.CreateEncrypt(ctx, encrypt3, userID, "DES", "1s")
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
