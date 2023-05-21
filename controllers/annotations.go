package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/zatarain/note-vook/models"
)

type Annotations struct {
	Database *models.DataAccessInterface
}

type AddAnnotationContract struct {
	VideoID int              `json:"video_id" binding:"required"`
	Type    int              `json:"type"`
	Title   string           `json:"title" binding:"required"`
	Notes   string           `json:"notes"`
	Start   models.TimeStamp `json:"start" binding:"required"`
	End     models.TimeStamp `json:"end" binding:"required"`
}

func (annotations *Annotations) Add(context *gin.Context) {
	// Try to bind the input

	// Check if the video exists for the current user

	// Check the Start and End are valid within the video Duration

	// Insert the record to the Database
}
