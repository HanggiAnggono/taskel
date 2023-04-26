package model

import "time"

type User struct {
	ID        uint   `gorm:"primaryKey"`
	Username  string `gorm:"unique;not null"`
	Name      string
	Email     *string
	CreatedAt time.Time
	UpdatedAt time.Time
}
