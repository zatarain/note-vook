package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zatarain/note-vook/models"
)

type VideosController struct {
	Database models.DataAccessInterface
}

func CurrentUser(context *gin.Context) *models.User {
	value, _ := context.Get("user")
	return value.(*models.User)
}

func (videos *VideosController) Index(context *gin.Context) {
	var recordset []models.Video
	user := CurrentUser(context)
	videos.Database.Find(&recordset, "user_id = ?", user.ID)
	context.JSON(http.StatusOK, recordset)
}

type AddVideoContract struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Link        string `json:"link"`
	Duration    int    `json:"duration"`
}

func (videos *VideosController) Add(context *gin.Context) {
	var input AddVideoContract

	// Trying to bind input from JSON
	if binding := context.BindJSON(&input); binding != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"summary": "Failed to read input",
			"details": binding.Error(),
		})
		return
	}

	user := CurrentUser(context)
	video := &models.Video{
		UserID:      user.ID,
		Title:       input.Title,
		Description: input.Description,
		Link:        input.Link,
		Duration:    input.Duration,
	}
	inserting := videos.Database.Create(&video).Error
	if inserting != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"summary": "Failed to save the video",
			"details": inserting.Error(),
		})
		return
	}

	context.JSON(http.StatusCreated, gin.H{
		"message": "Video successfully added",
		"data":    video,
	})
}
