package models

import (
	"time"
)

type Video struct {
	ID          int       `json:"id" gorm:"primary_key"`
	UserID      int       `json:"user_id" gorm:"index:idx_user;index:unq_user_video,unique"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Link        string    `json:"link" gorm:"index:unq_user_video,unique"`
	Duration    TimeStamp `json:"duration"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
