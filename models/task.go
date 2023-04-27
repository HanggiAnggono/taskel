package model

import (
	"time"
)

type Task struct {
	ID          uint /* `gorm:"primaryKey"` */
	Title       string
	Description *string
	Status      string
	UserID      *uint
	User        *User
	CreatedAt   time.Time
	UpdatedAt   time.Time
	// DeletedAt 	gorm.DeletedAt
}
