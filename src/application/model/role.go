package model

import (
	"encoding/json"
	"github.com/Etpmls/EM-Auth/src/application"
	em "github.com/Etpmls/Etpmls-Micro"
	"github.com/Etpmls/Etpmls-Micro/library"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"time"
)

type Role struct {
	ID        uint           `json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
	Name string              `json:"name"`
	Remark string            `json:"remark"`
	Users []User             `gorm:"many2many:role_users" json:"users"`
	Permissions []Permission `gorm:"many2many:role_permissions;" json:"permissions"`
}

type RoleGetOne struct {
	ID        uint           `json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
	Name      string         `json:"name"`
	Remark    string         `json:"remark"`
	// Users []User `gorm:"many2many:role_users" json:"users"`
	Permissions []Permission `gorm:"many2many:role_permissions;" json:"permissions"`
}

// interface conversion Role
// interface转换Role
func (this *Role) InterfaceToRole(i interface{}) (Role, error) {
	var r Role
	us, err := json.Marshal(i)
	if err != nil {
		em.LogError.Output(em.MessageWithLineNum("Object to JSON failed!" + err.Error()))
		return Role{}, err
	}
	err = json.Unmarshal(us, &r)
	if err != nil {
		em.LogError.Output(em.MessageWithLineNum("JSON conversion object failed!" + err.Error()))
		return Role{}, err
	}
	return r, nil
}

// Get all Roles
// 获取全部角色
func (this *Role) GetAll() ([]Role, error) {
	if em_library.Config.App.EnableCache {
		return this.getAll_Cache()
	} else {
		return this.getAll_NoCache()
	}
}
func (this *Role) getAll_NoCache() ([]Role, error) {
	var data []Role

	em.DB.Find(&data)

	if em_library.Config.App.EnableCache {
		b, err := json.Marshal(data)
		if err != nil {
			em.LogError.Output(em.MessageWithLineNum(err.Error()))
			return nil, err
		}
		em.Cache.SetString(application.Cache_RoleGetAll, string(b), 0)
	}

	return data, nil
}
func (this *Role) getAll_Cache() ([]Role, error) {
	j, err := em.Cache.GetString(application.Cache_RoleGetAll)
	if err != nil {
		if err == redis.Nil {
			return this.getAll_NoCache()
		}
		return nil, err
	}

	var roles []Role
	err = json.Unmarshal([]byte(j), &roles)
	if err != nil {
		em.LogError.Output(em.MessageWithLineNum(err.Error()))
		em.Cache.DeleteString(application.Cache_RoleGetAll)
		return nil, err
	}

	return roles, nil
}