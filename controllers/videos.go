package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zatarain/note-vook/models"
)

type VideosController struct {
	Database models.DataAccessInterface
}

type AddVideoContract struct {
	Title       string           `json:"title" binding:"required"`
	Description string           `json:"description"`
	Link        string           `json:"link" binding:"required"`
	Duration    models.TimeStamp `json:"duration" binding:"required"`
}

type EditVideoContract struct {
	Title       string           `json:"title"`
	Description string           `json:"description"`
	Link        string           `json:"link"`
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
	log.Println("Actual result within the handler: ", recordset)
	context.JSON(http.StatusOK, gin.H{"data": recordset})
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

func (videos *VideosController) View(context *gin.Context) {
	id := context.Param("id")
	user := CurrentUser(context)
	var video models.Video
	searching := videos.Database.First(&video, "id = ? AND user_id = ?", id, user.ID).Error
	if searching != nil {
		context.JSON(http.StatusNotFound, gin.H{
			"summary": "Video not found",
			"details": searching.Error(),
		})
		return
	}
	recordset := []models.Video{video}
	context.JSON(http.StatusOK, gin.H{"data": recordset})
	context.JSON(http.StatusOK, gin.H{"data": video})
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
