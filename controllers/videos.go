package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zatarain/note-vook/models"
)

type VideosController struct {
	Database models.DataAccessInterface
}

func (videos *VideosController) Index(context *gin.Context) {
	var recordset []models.Video
	value, _ := context.Get("user")
	user := value.(*models.User)
	videos.Database.Find(&recordset, "user_id = ?", user.ID)
	context.JSON(http.StatusOK, recordset)
}

func (videos *VideosController) Add(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{"message": "hello"})
}
