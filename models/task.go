package model

import (
	"time"

	"gorm.io/gorm"
)

type Task struct {
	ID          uint   /* `gorm:"primaryKey"` */
	Key         string `gorm:"uniqueIndex"`
	Title       string
	Description *string
	Status      string
	UserID      *uint
	User        *User
	Watchers    []*User    `gorm:"many2many:task_users"`
	Comments    []*Comment `gorm:"polymorphic:Commentable"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt
}

func (t *Task) StatusName() string {
	switch t.Status {
	case "inprogress":
		return "In Progress"
	case "done":
		return "Done"
	default:
		return "To Do"
	}
}

func (t *Task) StatusColor() string {
	switch t.Status {
	case "inprogress":
		return "bg-warning"
	case "done":
		return "bg-success"
	default:
		return "bg-secondary"
	}
}
