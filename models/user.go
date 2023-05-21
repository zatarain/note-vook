package models

import (
	"fmt"
	"time"
)

type User struct {
	ID        int       `json:"id" gorm:"primary_key"`
	Nickname  string    `json:"nickname" gorm:"unique"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (user *User) String() string {
	return fmt.Sprintf(
		"ID = %d, Nickname = '%s', Created At = '%s', Updated At = '%s'",
		user.ID,
		user.Nickname,
		user.CreatedAt.Format(time.RFC1123),
		user.UpdatedAt.Format(time.RFC1123),
	)
}
