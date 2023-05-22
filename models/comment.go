package model

import (
	"time"

	"gorm.io/gorm"
)

type Comment struct {
	ID          uint `gorm:"primaryKey"`
	CommentedID uint
	CommenterID uint
	Commenter   *User
	Comment     string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt
	Upvotes     uint
	Downvotes   uint
}
