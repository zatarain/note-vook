package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zatarain/note-vook/models"
)

type AnnotationsController struct {
	Database models.DataAccessInterface
}

type AddAnnotationContract struct {
	VideoID uint             `json:"video_id" binding:"required"`
	Type    uint             `json:"type"`
	Title   string           `json:"title" binding:"required"`
	Notes   string           `json:"notes"`
	Start   models.TimeStamp `json:"start" binding:"required,ltefield=End"`
	End     models.TimeStamp `json:"end" binding:"required"`
}

type EditAnnotationContract struct {
	Type  uint             `json:"type"`
	Title string           `json:"title"`
	Notes string           `json:"notes"`
	Start models.TimeStamp `json:"start" binding:"ltefield=End"`
	End   models.TimeStamp `json:"end"`
}

func (annotations *AnnotationsController) findVideo(context *gin.Context, video *models.Video, id uint) bool {
	user := CurrentUser(context)
	searching := annotations.Database.First(video, "id = ? AND user_id = ?", id, user.ID).Error
	if searching != nil {
		context.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error":  "Video not found",
			"reason": searching.Error(),
		})
		return false
	}
	return true
}

func (annotations *AnnotationsController) search(context *gin.Context, annotation *models.Annotation) bool {
	user := CurrentUser(context)
	id := context.Param("id")

	searching := annotations.Database.
		Joins("Video").First(annotation, "annotations.id = ? AND user_id = ?", id, user.ID).Error
	if searching != nil {
		context.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error":  "Annotation not found",
			"reason": searching.Error(),
		})
		return false
	}
	return true
}

func isInRange(timestamp models.TimeStamp, duration models.TimeStamp) bool {
	return timestamp >= 0 && timestamp <= duration
}

func (annotations *AnnotationsController) CheckInterval(
	context *gin.Context,
	start models.TimeStamp,
	end models.TimeStamp,
	duration models.TimeStamp,
) bool {
	if !isInRange(start, duration) || !isInRange(end, duration) {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":  "Invalid time interval",
			"reason": "start and end must be positive and less or equal than video duration",
		})
		return false
	}
	return true
}

func (annotations *AnnotationsController) Add(context *gin.Context) {
	// Try to bind the input from JSON
	var input AddAnnotationContract

	if binding := context.ShouldBindJSON(&input); binding != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error":  "Failed to read input",
			"reason": binding.Error(),
		})
		return
	}

	// Check if the video exists for the current user
	video := models.Video{}
	if !annotations.findVideo(context, &video, input.VideoID) {
		return
	}

	// Check the Start and End are valid within the video Duration
	if !annotations.CheckInterval(context, input.Start, input.End, video.Duration) {
		return
	}

	// Insert the record to the Database
	annotation := models.Annotation{
		VideoID: input.VideoID,
		Type:    input.Type,
		Title:   input.Title,
		Notes:   input.Notes,
		Start:   input.Start,
		End:     input.End,
	}
	inserting := annotations.Database.Create(&annotation).Error
	if inserting != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error":  "Failed to save the annotation",
			"reason": inserting.Error(),
		})
		return
	}

	// Send status created with the new annotation
	context.JSON(http.StatusCreated, &annotation)
}

func (annotations *AnnotationsController) Edit(context *gin.Context) {
	// Look for the annotation we want to edit
	var annotation models.Annotation
	if !annotations.search(context, &annotation) {
		return
	}

	// Try to bind the input from JSON
	var input EditAnnotationContract

	if binding := context.ShouldBindJSON(&input); binding != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error":  "Failed to read input",
			"reason": binding.Error(),
		})
		return
	}

	// Check the Start and End are valid within the video Duration
	if !annotations.CheckInterval(context, input.Start, input.End, annotation.Video.Duration) {
		return
	}

	annotation.UpdatedAt = time.Now()
	saving := annotations.Database.Model(&annotation).Updates(input).Error
	if saving != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error":  "Failed to save the annotation",
			"reason": saving.Error(),
		})
		return
	}

	// Send success status with new values
	context.JSON(http.StatusOK, &annotation)
}

func (annotations *AnnotationsController) Delete(context *gin.Context) {
	// Look for the annotation we want to delete
	var annotation models.Annotation
	if !annotations.search(context, &annotation) {
		return
	}

	// Try to delete the annotation from database
	deleting := annotations.Database.Delete(&annotation).Error
	if deleting != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error":  "Failed to delete the annotation",
			"reason": deleting.Error(),
		})
		return
	}

	// Send success message
	context.JSON(http.StatusOK, gin.H{
		"message": "Annotation successfully deleted",
	})
}
