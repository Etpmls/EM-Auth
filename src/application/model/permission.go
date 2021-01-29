package model

import (
	"encoding/json"
	"github.com/Etpmls/EM-Auth/src/application"
	em "github.com/Etpmls/Etpmls-Micro/v2"
	"github.com/Etpmls/Etpmls-Micro/v2/define"
	"strings"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"time"
)

type Permission struct {
	ID        uint           `json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
	Name      string         `json:"name"`
	Auth      int            `json:"auth"`
	Method    string		 `json:"method"`
	Path      string         `json:"path"`
	Remark    string         `json:"remark"`
	Roles     []Role         `gorm:"many2many:role_permissions" json:"roles"`
}

type PermissionGetOne struct {
	ID        uint           `json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
	Name      string         `json:"name"`
	Auth      int            `json:"auth"`
	Method    string		 `json:"method"`
	Path      string         `json:"path"`
	Remark    string         `json:"remark"`
	// Roles []Role             `gorm:"many2many:role_permissions" json:"roles"`
}

func (this *Permission) InterfaceToPermission(i interface{}) (Permission, error) {
	var p Permission
	us, err := json.Marshal(i)
	if err != nil {
		em.LogError.Output(em.MessageWithLineNum("Object to JSON failed!" + err.Error()))
		return Permission{}, err
	}
	err = json.Unmarshal(us, &p)
	if err != nil {
		em.LogError.Output(em.MessageWithLineNum("JSON conversion object failed!" + err.Error()))
		return Permission{}, err
	}
	return p, nil
}

// Get all permissions
// 获取全部权限
func (this *Permission) GetAll() ([]Permission, error) {
	e, _ := em.Kv.ReadKey(define.KvCacheEnable)
	if strings.ToLower(e) != "true" {
		return this.getAll_NoCache(strings.ToLower(e) == "true")
	}

	return this.getAll_Cache(strings.ToLower(e) == "true")
}
func (this *Permission) getAll_NoCache(enableCache bool) ([]Permission, error) {
	var data []Permission


	em.DB.Find(&data)

	if enableCache {
		b, err := json.Marshal(data)
		if err != nil {
			em.LogError.Output(em.MessageWithLineNum(err.Error()))
			return nil, err
		}
		em.Cache.SetString(application.Cache_PermissionGetAll, string(b), 0)
	}


	return data, nil
}
func (this *Permission) getAll_Cache(enableCache bool) ([]Permission, error) {
	j, err := em.Cache.GetString(application.Cache_PermissionGetAll)
	if err != nil {
		if err == redis.Nil {
			return this.getAll_NoCache(enableCache)
		}
		return nil, err
	}

	var permissions []Permission
	err = json.Unmarshal([]byte(j), &permissions)
	if err != nil {
		em.LogError.Output(em.MessageWithLineNum(err.Error()))
		em.Cache.DeleteString(application.Cache_PermissionGetAll)
		return nil, err
	}

	return permissions, nil
}