package models

import (
	"time"

	"gorm.io/gorm"
)

type Video struct {
	gorm.Model
	ID          int       `json:"id" gorm:"primary_key"`
	UserID      int       `json:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Link        string    `json:"link"`
	Duration    int       `json:"duration"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
