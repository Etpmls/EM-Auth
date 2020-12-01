// +build mysql

package database

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"unique;notnull"`
	Password string `gorm:"notnull"`
	Avatar Attachment	`gorm:"polymorphic:Owner;polymorphicValue:user-avatar"`
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
	Path	string
	Auth int
	Remark string
	Roles []Role `gorm:"many2many:role_permissions;"`
}
