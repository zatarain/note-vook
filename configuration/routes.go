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

	videos := &controllers.VideosController{
		Database: Database,
	}

	annotations := &controllers.AnnotationsController{
		Database: Database,
	}

	server.HEAD("/health", controllers.HealthCheck)
	server.POST("/signup", users.Signup)
	server.POST("/login", users.Login)

	// Authorised end-points
	server.GET("/videos", users.Authorise, videos.Index)
	server.POST("/videos", users.Authorise, videos.Add)
	server.GET("/videos/:id", users.Authorise, videos.View)
	server.PATCH("/videos/:id", users.Authorise, videos.Edit)
	server.DELETE("/videos/:id", users.Authorise, videos.Delete)

	server.POST("/annotations", users.Authorise, annotations.Add)
	server.DELETE("/annotations/:id", users.Authorise, annotations.Delete)
}
