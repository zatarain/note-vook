package models

import (
	"time"
)

type Video struct {
	ID          uint      `json:"id" gorm:"primary_key"`
	UserID      uint      `json:"user_id" gorm:"index:idx_user;index:unq_user_video,unique"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Link        string    `json:"link" gorm:"index:unq_user_video,unique"`
	Duration    TimeStamp `json:"duration"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime:milli"`

	// Associations
	Annotations []Annotation `json:"annotations" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
