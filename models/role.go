package model

import "gorm.io/gorm"

type Role struct {
	gorm.Model
	Name string
	Permissions []Permission `gorm:"many2many:role_permissions"`
}