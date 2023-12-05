package main

import (
	"gin-gorm-clean-template/common"
	"gin-gorm-clean-template/config"
	"gin-gorm-clean-template/controller"
	"gin-gorm-clean-template/middleware"
	"gin-gorm-clean-template/repository"
	"gin-gorm-clean-template/routes"
	"gin-gorm-clean-template/service"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		res := common.BuildErrorResponse("Gagal Terhubung ke Server", err.Error(), common.EmptyObj{})
		(*gin.Context).JSON((&gin.Context{}), http.StatusBadGateway, res)
		return
	}

	var (
		db *gorm.DB = config.SetupDatabaseConnection()

		jwtService service.JWTService = service.NewJWTService()

		userRepository    repository.UserRepository    = repository.NewUserRepository(db)
		encryptRepository repository.EncryptRepository = repository.NewEncryptRepository(db)

		userService    service.UserService    = service.NewUserService(userRepository, encryptRepository)
		encryptService service.EncryptService = service.NewEncryptService(encryptRepository, userRepository)

		userController    controller.UserController    = controller.NewUserController(userService, jwtService)
		encryptController controller.EncryptController = controller.NewEncryptController(encryptService, jwtService)
	)

	server := gin.Default()
	server.Use(middleware.CORSMiddleware())

	routes.UserRoutes(server, userController, jwtService)
	routes.EncryptRoutes(server, encryptController, jwtService)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	server.Run("127.0.0.1:" + port)
}
