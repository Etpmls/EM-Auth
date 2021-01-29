package model

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Etpmls/EM-Auth/src/application"
	em "github.com/Etpmls/Etpmls-Micro/v2"
	"github.com/Etpmls/Etpmls-Micro/v2/define"
	"github.com/go-redis/redis/v8"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"strings"
	"time"
)

type User struct {
	ID        uint           `json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
	Username  string         `json:"username"`
	Password  string         `json:"password"`
	Avatar    Attachment     `gorm:"-" json:"avatar"`
	Roles     []Role         `gorm:"many2many:role_users" json:"roles"`
}

type UserGetOne struct {
	ID        uint                  `json:"id"`
	CreatedAt time.Time             `json:"created_at"`
	UpdatedAt time.Time             `json:"updated_at"`
	DeletedAt gorm.DeletedAt        `json:"deleted_at"`
	Username  string                `json:"username"`
	Password  string                `json:"-"`
	Avatar    Attachment            `gorm:"-" json:"avatar"`
	Roles     []Role `gorm:"many2many:role_users" json:"roles"`
}

// Get token by ID&username
// 通过ID&用户名获取Token
func (this *User) UserGetToken(userId uint, username string) (string, error) {
	return em.Micro.Auth.CreateGeneralToken(int(userId), username)
}

// Get user by token
// 根据token获取用户
func (this *User) GetUserByToken(token string) (u User, err error) {
	// 从Token获取ID
	id, err := em.Micro.Auth.GetIdByToken(token)
	if err != nil {
		em.LogError.Output(em.MessageWithLineNum_OneRecord(err.Error()))
		return u, err
	}
	// 从Token获取username
	username, err  := em.Micro.Auth.GetIssuerByToken(token)
	if err != nil {
		em.LogError.Output(em.MessageWithLineNum(err.Error()))
		return u, err
	}
	// 获取用户
	var data User
	result := em.DB.Where("id = ? AND username = ?", id, username).First(&data)
	if !(result.RowsAffected > 0) {
		return u, em.LogError.OutputAndReturnError(em.MessageWithLineNum("The current user was not found in the database!"))
	}

	return data, nil
}

// Obtain user id based on token
// 根据token获取用户id
func (this *User) GetUserIdByToken(token string) (id int, err error) {
	// 从Token获取ID
	id, err = em.Micro.Auth.GetIdByToken(token)
	if err != nil {
		em.LogError.Output(em.MessageWithLineNum(err.Error()))
		return 0, err
	}

	return id, nil
}

// Obtain user id based on Request
// 根据Request获取用户id
func (this *User) GetUserIdByRequest(ctx context.Context) (id int, err error) {
	token, err := em.Micro.Auth.GetTokenFromCtx(ctx)
	if err != nil {
		return 0, err
	}

	// Get ID from Token
	// 从Token获取ID
	id, err = em.Micro.Auth.GetIdByToken(token)
	if err != nil {
		return 0, em.LogError.OutputAndReturnError(em.MessageWithLineNum(err.Error()))
	}

	return id, nil
}

// Verify user logic
// 验证用户逻辑
func (this *User) Verify(username string, password string) (u User, err error) {
	//Search User
	var user User
	em.DB.Where("username = ?", username).First(&user)
	if !(user.ID > 0) {
		em.LogWarn.Output(em.MessageWithLineNum_OneRecord("The username does not exist! Username:" + username))
		return u, errors.New("The username does not exist!")
	}

	//Password is wrong
	b, err := this.VerifyPassword(password, user.Password)
	if err != nil || !b {
		em.LogInfo.Output(em.MessageWithLineNum("Verification failed or wrong password!"))
		return u, errors.New("Verification failed or wrong password!")
	}

	return user, err
}

// Verify user password
// 验证用户密码
func (this *User) VerifyPassword(password, hash string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		em.LogInfo.Output(em.MessageWithLineNum(err.Error()))
		return false, err
	}
	return true, err
}

// Bcrypt Password
// 加密密码
func (this *User) BcryptPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// interface conversion User
// interface转换User
func (this *User) InterfaceToUser(i interface{}) (User, error) {
	var u User
	us, err := json.Marshal(i)
	if err != nil {
		em.LogError.Output(em.MessageWithLineNum("Object to JSON failed! err:" + err.Error()))
		return User{}, err
	}
	err = json.Unmarshal(us, &u)
	if err != nil {
		em.LogError.Output(em.MessageWithLineNum("JSON conversion object failed! err:" + err.Error()))
		return User{}, err
	}
	return u, nil
}

// interface conversion UserGetOne
// interface转换UserGetOne
func (this *User) InterfaceToUserGetOne(i interface{}) (UserGetOne, error) {
	var u UserGetOne
	us, err := json.Marshal(i)
	if err != nil {
		em.LogError.Output(em.MessageWithLineNum("Object to JSON failed! err:" + err.Error()))
		return UserGetOne{}, err
	}
	err = json.Unmarshal(us, &u)
	if err != nil {
		em.LogError.Output(em.MessageWithLineNum("JSON conversion object failed! err:" + err.Error()))
		return UserGetOne{}, err
	}
	return u, nil
}

// Get all Users
// 获取全部用户
func (this *User) GetAll() ([]User, error) {
	e, _ := em.Kv.ReadKey(define.KvCacheEnable)
	if strings.ToLower(e) != "true" {
		return this.getAll_NoCache(strings.ToLower(e) == "true")
	}
	return this.getAll_Cache(strings.ToLower(e) == "true")
}
func (this *User) getAll_NoCache(enableCache bool) ([]User, error) {
	var data []User

	em.DB.Find(&data)

	if enableCache {
		b, err := json.Marshal(data)
		if err != nil {
			em.LogError.Output(em.MessageWithLineNum(err.Error()))
			return nil, err
		}
		em.Cache.SetString(application.Cache_UserGetAll, string(b), 0)
	}

	return data, nil
}
func (this *User) getAll_Cache(enableCache bool) ([]User, error) {
	j, err := em.Cache.GetString(application.Cache_UserGetAll)
	if err != nil {
		if err == redis.Nil {
			return this.getAll_NoCache(enableCache)
		}
		return nil, err
	}

	var users []User
	err = json.Unmarshal([]byte(j), &users)
	if err != nil {
		em.LogError.Output(em.MessageWithLineNum(err.Error()))
		em.Cache.DeleteString(application.Cache_UserGetAll)
		return nil, err
	}

	return users, nil
}