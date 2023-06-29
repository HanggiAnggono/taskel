package model

import (
	"fmt"
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

func (t *Task) BeforeCreate(tx *gorm.DB) (err error) {
	var lastTask Task
	if result := tx.Last(&lastTask); result.Error == nil {
		t.ID = lastTask.ID + 1
		t.Key = fmt.Sprintf("TASK-%d", t.ID)
	} else if result.Error == gorm.ErrRecordNotFound {
		t.ID = 1
	} else {
		err = result.Error
	}
	return
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
