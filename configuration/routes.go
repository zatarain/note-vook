package configuration

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/zatarain/note-vook/controllers"
)

func Setup(server gin.IRouter) {
	users := &controllers.UsersController{
		Database:       Database,
		SecretTokenKey: os.Getenv("SECRET_TOKEN_KEY"),
	}
	server.HEAD("/health", controllers.HealthCheck)
	server.POST("/signup", users.Signup)
	server.POST("/login", users.Login)
	server.GET("/videos", users.Authorise, controllers.GetVideos)
}
