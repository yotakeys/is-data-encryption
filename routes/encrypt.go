package routes

import (
	"gin-gorm-clean-template/controller"
	"gin-gorm-clean-template/middleware"
	"gin-gorm-clean-template/service"

	"github.com/gin-gonic/gin"
)

func EncryptRoutes(router *gin.Engine, EncryptController controller.EncryptController, jwtService service.JWTService) {
	encryptRoutes := router.Group("/api/encrypt")
	{
		encryptRoutes.POST("", middleware.Authenticate(jwtService), EncryptController.CreateEncrypt)
		encryptRoutes.GET("", middleware.Authenticate(jwtService), EncryptController.GetAllEncrypt)
		encryptRoutes.GET("/file", EncryptController.GetFile)
	}
}
