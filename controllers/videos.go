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
	Title       string           `json:"title" binding:"required"`
	Description string           `json:"description" binding:"required"`
	Link        string           `json:"link" binding:"required"`
	Duration    models.TimeStamp `json:"duration" binding:"required"`
}

type EditVideoContract struct {
	Title       string           `json:"title"`
	Description string           `json:"description"`
	Link        string           `json:"link"`
	Duration    models.TimeStamp `json:"duration"`
}

func (videos *VideosController) Add(context *gin.Context) {
	var input AddVideoContract

	// Trying to bind input from JSON
	if binding := context.ShouldBindJSON(&input); binding != nil {
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
		Duration:    int64(input.Duration.Duration),
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

func (videos *VideosController) View(context *gin.Context) {
	id := context.Param("id")
	user := CurrentUser(context)
	var video models.Video
	searching := videos.Database.Where("id = ? AND user_id = ?", id, user.ID).First(&video).Error
	if searching != nil {
		context.JSON(http.StatusNotFound, gin.H{
			"summary": "Video not found",
			"details": searching.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, video)
}

func (videos *VideosController) Edit(context *gin.Context) {
	id := context.Param("id")
	user := CurrentUser(context)
	var video models.Video
	searching := videos.Database.Where("id = ? AND user_id = ?", id, user.ID).First(&video).Error
	if searching != nil {
		context.JSON(http.StatusNotFound, gin.H{
			"summary": "Video not found",
			"details": searching.Error(),
		})
		return
	}

	// Trying to bind input from JSON
	var input EditVideoContract
	if binding := context.ShouldBindJSON(&input); binding != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"summary": "Failed to read input",
			"details": binding.Error(),
		})
		return
	}
	videos.Database.Model(&video).Updates(input)

	context.JSON(http.StatusOK, gin.H{
		"summary": "Video successfully updated",
		"data":    video,
	})
}

func (videos *VideosController) Delete(context *gin.Context) {
	id := context.Param("id")
	user := CurrentUser(context)
	var video models.Video
	searching := videos.Database.Where("id = ? AND user_id = ?", id, user.ID).First(&video).Error
	if searching != nil {
		context.JSON(http.StatusNotFound, gin.H{
			"summary": "Video not found",
			"details": searching.Error(),
		})
		return
	}

	videos.Database.Delete(&video)

	context.JSON(http.StatusOK, gin.H{
		"summary": "Video successfully deleted",
	})
}
