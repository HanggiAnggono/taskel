package model

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        uint   `gorm:"primaryKey"`
	Username  string `gorm:"unique"`
	Name      string
	Email     *string
	Password  string `gorm:"type:varchar(255)"`
	Watches   []*Task `gorm:"many2many:task_users"`
	Role      *Role `gorm:"OnDelete:SET NULL"`
	RoleID    *uint
	CreatedAt time.Time
	UpdatedAt time.Time
}

func UserHashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	return string(hashed), err
}

func UserComparePassword(password, hashed string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))

	return err == nil
}
