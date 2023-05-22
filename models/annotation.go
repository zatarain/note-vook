package models

import (
	"time"
)

type Annotation struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	VideoID   uint      `json:"video_id" gorm:"index:idx_video"`
	Type      uint      `json:"type" gorm:"index:idx_type"`
	Title     string    `json:"title"`
	Notes     string    `json:"notes"`
	Start     TimeStamp `json:"start"`
	End       TimeStamp `json:"end"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Associations
	Video *Video `json:"video" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
