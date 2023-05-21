package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zatarain/note-vook/models"
)

type AnnotationsController struct {
	Database models.DataAccessInterface
}

type AddAnnotationContract struct {
	VideoID int              `json:"video_id" binding:"required"`
	Type    int              `json:"type"`
	Title   string           `json:"title" binding:"required"`
	Notes   string           `json:"notes"`
	Start   models.TimeStamp `json:"start" binding:"required,ltefield=End"`
	End     models.TimeStamp `json:"end" binding:"required"`
}

func (annotations *AnnotationsController) findVideo(context *gin.Context, video *models.Video, id int) bool {
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
			"error":  "Invalid interval",
			"reason": "`start` and `end` must be positive and less or equal than video `duration`",
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

	context.JSON(http.StatusCreated, &annotation)
}
