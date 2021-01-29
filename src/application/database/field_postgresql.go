// +build postgresql

package database

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
	Avatar string	`gorm:"-"`
	Roles []Role `gorm:"many2many:role_users;"`
}

type Role struct {
	gorm.Model
	Name string
	Remark string
	Users []User `gorm:"many2many:role_users;"`
	Permissions []Permission `gorm:"many2many:role_permissions;"`
}

type Permission struct {
	gorm.Model
	Name string
	Auth int
	Method string
	Path	string
	Remark string
	Roles []Role `gorm:"many2many:role_permissions;"`
}
