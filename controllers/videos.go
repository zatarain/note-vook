package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zatarain/note-vook/models"
)

type VideosController struct {
	Database models.DataAccessInterface
}

type AddVideoContract struct {
	Title       string           `json:"title" binding:"required"`
	Description string           `json:"description"`
	Link        string           `json:"link" binding:"required,url"`
	Duration    models.TimeStamp `json:"duration" binding:"required"`
}

type EditVideoContract struct {
	Title       string           `json:"title"`
	Description string           `json:"description"`
	Link        string           `json:"link" binding:"omitempty,url"`
	Duration    models.TimeStamp `json:"duration"`
}

func CurrentUser(context *gin.Context) *models.User {
	value, _ := context.Get("user")
	return value.(*models.User)
}

func (videos *VideosController) Index(context *gin.Context) {
	user := CurrentUser(context)
	var recordset []models.Video
	videos.Database.Find(&recordset, "user_id = ?", user.ID)
	context.JSON(http.StatusOK, recordset)
}

func (videos *VideosController) Add(context *gin.Context) {
	var input AddVideoContract

	// Trying to bind input from JSON
	if binding := context.ShouldBindJSON(&input); binding != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error":  "Failed to read input",
			"reason": binding.Error(),
		})
		return
	}

	user := CurrentUser(context)
	video := models.Video{
		UserID:      user.ID,
		Title:       input.Title,
		Description: input.Description,
		Link:        input.Link,
		Duration:    input.Duration,
	}
	inserting := videos.Database.Create(&video).Error
	if inserting != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error":  "Failed to save the video",
			"reason": inserting.Error(),
		})
		return
	}

	context.JSON(http.StatusCreated, &video)
}

func (videos *VideosController) search(video *models.Video, context *gin.Context) bool {
	id := context.Param("id")
	user := CurrentUser(context)
	searching := videos.Database.Model(video).Preload("Annotations").
		First(video, "id = ? AND user_id = ?", id, user.ID).Error
	if searching != nil {
		context.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error":  "Video not found",
			"reason": searching.Error(),
		})
		return false
	}
	return true
}

func (videos *VideosController) View(context *gin.Context) {
	var video models.Video
	if !videos.search(&video, context) {
		return
	}
	context.JSON(http.StatusOK, &video)
}

func (videos *VideosController) Edit(context *gin.Context) {
	var video models.Video
	if !videos.search(&video, context) {
		return
	}

	// Trying to bind input from JSON
	var input EditVideoContract
	if binding := context.ShouldBindJSON(&input); binding != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error":  "Failed to read input",
			"reason": binding.Error(),
		})
		return
	}

	video.UpdatedAt = time.Now()
	saving := videos.Database.Model(&video).Updates(input).Error
	if saving != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error":  "Failed to save the video",
			"reason": saving.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, &video)
}

func (videos *VideosController) Delete(context *gin.Context) {
	var video models.Video
	if !videos.search(&video, context) {
		return
	}

	deleting := videos.Database.Delete(&video).Error
	if deleting != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error":  "Failed to delete the video",
			"reason": deleting.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "Video successfully deleted",
	})
}
