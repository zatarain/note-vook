package models

import "time"

type Annotation struct {
	ID        int       `json:"id" gorm:"primary_key"`
	VideoID   int       `json:"user_id" gorm:"index:idx_video"`
	Type      int       `json:"type" gorm:"index:idx_type"`
	Title     string    `json:"title"`
	Notes     string    `json:"notes"`
	Start     TimeStamp `json:"start"`
	End       TimeStamp `json:"end"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
