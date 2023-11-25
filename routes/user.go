package routes

import (
	"gin-gorm-clean-template/controller"
	"gin-gorm-clean-template/middleware"
	"gin-gorm-clean-template/service"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine, UserController controller.UserController, jwtService service.JWTService) {
	userRoutes := router.Group("/api/user")
	{
		userRoutes.POST("", UserController.RegisterUser)
		userRoutes.POST("/login", UserController.LoginUser)
		userRoutes.GET("", middleware.Authenticate(jwtService), UserController.GetAllUser)
		userRoutes.DELETE("/", middleware.Authenticate(jwtService), UserController.DeleteUser)
		userRoutes.PUT("/", middleware.Authenticate(jwtService), UserController.UpdateUser)
		userRoutes.GET("/me", middleware.Authenticate(jwtService), UserController.MeUser)

		userRoutes.POST("/request_data", middleware.Authenticate(jwtService), UserController.SendEmailEncrypt)
		userRoutes.GET("/response_email", UserController.SendEmailResponse)
		userRoutes.POST("/asymmetric_encypt", middleware.Authenticate(jwtService), UserController.AsymmetricEncrypt)
		userRoutes.POST("/asymmetric_decrypt", middleware.Authenticate(jwtService), UserController.AsymmetricDecrypt)
		userRoutes.GET("/get_list_requesting_user", middleware.Authenticate(jwtService), UserController.GetListRequestingUser)
	}
}
