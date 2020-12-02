package service

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Etpmls/EM-Auth/src/application"
	"github.com/Etpmls/EM-Auth/src/application/client"
	"github.com/Etpmls/EM-Auth/src/application/model"
	"github.com/Etpmls/EM-Auth/src/application/protobuf"
	em "github.com/Etpmls/Etpmls-Micro"
	"github.com/Etpmls/Etpmls-Micro/library"
	em_protobuf "github.com/Etpmls/Etpmls-Micro/protobuf"
	"github.com/Etpmls/Etpmls-Micro/utils"
	"github.com/go-redis/redis/v8"
	"google.golang.org/grpc/codes"
	"gorm.io/gorm"
	"strconv"
	"strings"
)

type ServiceUser struct {
	protobuf.UnimplementedUserServer
}

// User Register
// 用户注册
func (this *ServiceUser) Register(ctx context.Context, request *protobuf.UserRegister) (*em_protobuf.Response, error) {
	return &em_protobuf.Response{
		Code:    em.SUCCESS_Code,
		Status:  em.SUCCESS_Status,
		Message: em_library.I18n.TranslateFromRequest(ctx, "ERROR_MESSAGE_RegistrationClosed"),
	}, nil
}


// User Login
// 用户登录
type validate_UserLogin struct {
	Username string `json:"username" validate:"required,max=255"`
	Password string `json:"password" validate:"required,max=255"`
}
func (this *ServiceUser) Login(ctx context.Context, request *protobuf.UserLogin) (*em_protobuf.Response, error)  {
	// Validate
	{
		var vd validate_UserLogin
		err := em_library.Validator.Validate(request, &vd)
		if err != nil {
			em.LogWarn.Output(em_utils.MessageWithLineNum(err.Error()))
			return em.ErrorRpc(codes.InvalidArgument, em.ERROR_Code, em_library.I18n.TranslateFromRequest(ctx, "ERROR_Validate"), nil, err)
		}
	}

	// Verify Username & Password
	var us model.User
	usr, err := us.Verify(request.Username, request.Password)
	if err != nil {
		em.LogInfo.Output(em_utils.MessageWithLineNum("Verify user failed!"))
		return em.ErrorRpc(codes.InvalidArgument, em.ERROR_Code, em_library.I18n.TranslateFromRequest(ctx, "ERROR_Login"), nil, err)
	}

	//JWT
	token, err := us.UserGetToken(usr.ID, usr.Username)
	if err != nil {
		em.LogError.Output(em_utils.MessageWithLineNum("Get Token failed! Error:" + err.Error()))
		return em.ErrorRpc(codes.InvalidArgument, em.ERROR_Code, em_library.I18n.TranslateFromRequest(ctx, "ERROR_Login"), nil, err)
	}

	//Return Token
	resData := make(map[string]string)
	resData["token"] = token

	return em.SuccessRpc(em.SUCCESS_Code, em_library.I18n.TranslateFromRequest(ctx, "SUCCESS_Login"), resData)
}


// User Logout
// 用户登出
func (this *ServiceUser) Logout(ctx context.Context, request *em_protobuf.Empty) (*em_protobuf.Response, error) {
	return em.SuccessRpc(em.SUCCESS_Code, em_library.I18n.TranslateFromRequest(ctx, "SUCCESS_Logout"), nil)
}

// Get current user
// 获取当前用户
func (this *ServiceUser) GetCurrent(ctx context.Context, request *protobuf.UserGetCurrent) (*em_protobuf.Response, error) {
	if em_library.Config.App.Cache {
		return this.getCurrent_Cache(ctx, request)
	} else {
		return this.getCurrent_NoCache(ctx, request)
	}
}
func (this *ServiceUser) getCurrent_NoCache(ctx context.Context, request *protobuf.UserGetCurrent) (*em_protobuf.Response, error) {
	// Get User By request
	var user model.User
	u, err := user.GetUserByToken(request.Token)
	if err != nil {
		return em.ErrorRpc(codes.InvalidArgument, em.ERROR_Code, em_library.I18n.TranslateFromRequest(ctx, "ERROR_GetUser"), nil, err)
	}

	// Filter some field
	filter_user, err := user.InterfaceToUserGetOne(u)
	if err != nil {
		return em.ErrorRpc(codes.InvalidArgument, em.ERROR_Code, em_library.I18n.TranslateFromRequest(ctx, "ERROR_GetUser"), nil, em.LogError.OutputAndReturnError(em_utils.MessageWithLineNum(err.Error())))
	}

	// Ignore the avatar tag in the User structure
	type tmp struct {
		model.UserGetOne
		Avatar string     `json:"avatar"`
		Roles []string `json:"roles"`
	}
	var userApi = tmp{ UserGetOne: filter_user }

	// Avatar
	// 1.Get token By Request
	var path string
	token, err := em.NewAuth().GetTokenFromCtx(ctx)
	if err == nil {
		path = client.NewClient().User_GetAvatar(token, uint32(u.ID), application.Relationship_User_Avatar)
	}

	userApi.Avatar = path
	// Roles
	var r []model.Role
	_ = em.DB.Model(&u).Association("Roles").Find(&r)
	for _, v := range r {
		userApi.Roles = append(userApi.Roles, v.Name)
	}

	if em_library.Config.App.Cache {
		b, err := json.Marshal(userApi)
		if err != nil {
			em.LogError.Output(err)
		} else {
			var m = make(map[string]string)
			m[strconv.Itoa(int(u.ID))] = string(b)
			em_library.Cache.SetHash(application.Cache_UserGetCurrent, m)
		}
	}

	return em.SuccessRpc(em.SUCCESS_Code, em_library.I18n.TranslateFromRequest(ctx, "SUCCESS_GetUser"), userApi)
}
func (this *ServiceUser) getCurrent_Cache(ctx context.Context, request *protobuf.UserGetCurrent) (*em_protobuf.Response, error) {
	var user model.User
	id, err := user.GetUserIdByToken(request.Token)
	if err != nil {
		return em.ErrorRpc(codes.InvalidArgument, em.ERROR_Code, em_library.I18n.TranslateFromRequest(ctx, "ERROR_GetUser"), nil, em.LogError.OutputAndReturnError(em_utils.MessageWithLineNum(err.Error())))
	}

	str, err := em_library.Cache.GetHash(application.Cache_UserGetCurrent, strconv.Itoa(int(id)))
	if err != nil {
		if err == redis.Nil {
			return this.getCurrent_NoCache(ctx, request)
		}
		return em.ErrorRpc(codes.InvalidArgument, em.ERROR_Code, em_library.I18n.TranslateFromRequest(ctx, "ERROR_GetUser"), nil, em.LogError.OutputAndReturnError(em_utils.MessageWithLineNum(err.Error())))
	}

	type tmp struct {
		model.UserGetOne
		Avatar string     `json:"avatar"`
		Roles []string `json:"roles"`
	}
	var userApi tmp
	err = json.Unmarshal([]byte(str), &userApi)
	if err != nil {
		em.LogError.Output(em_utils.MessageWithLineNum(err.Error()))
		em_library.Cache.DeleteHash(application.Cache_UserGetCurrent, strconv.Itoa(int(id)))
	}

	return em.SuccessRpc(em.SUCCESS_Code, em_library.I18n.TranslateFromRequest(ctx, "SUCCESS_GetUser"), userApi)
}

// Get all user
// 获取全部用户
func (this *ServiceUser) GetAll(ctx context.Context, request *em_protobuf.Pagination) (*em_protobuf.Response, error) {
	// 重写ApiUserGetAllV2的Roles字段，防止泄露隐私字段信息
	type Role model.RoleGetOne
	type User struct {
		model.UserGetOne
		Roles []Role `gorm:"many2many:role_users" json:"roles"`
	}
	var data []User

	// 获取分页和标题
	var orm em_library.Gorm
	limit, offset := orm.GeneratePaginationLimit(int(request.Number), int(request.Size))
	var count int64
	// Get the title of the search, if not get all the data
	// 获取搜索的标题，如果没有获取全部数据
	search := request.Search

	em.DB.Model(&User{}).Preload("Roles").Where("username " +em.FUZZY_SEARCH+ " ?", "%"+ search +"%").Count(&count).Limit(limit).Offset(offset).Find(&data)

	m := map[string]interface{}{"data": data, em_library.Config.Field.Pagination.Count: count}

	return em.SuccessRpc(em.SUCCESS_Code, em_library.I18n.TranslateFromRequest(ctx, "SUCCESS_Get"), m)
}

// Create user
// 创建用户
type validate_UserCreate struct {
	validate_UserLogin
	Roles []model.Role `json:"roles" validate:"required"`
}
func (this *ServiceUser) Create(ctx context.Context, request *protobuf.UserCreate) (*em_protobuf.Response, error) {
	// Validate
	{
		var vd validate_UserCreate
		err := em_library.Validator.Validate(request, &vd)
		if err != nil {
			em.LogWarn.Output(em_utils.MessageWithLineNum(err.Error()))
			return em.ErrorRpc(codes.InvalidArgument, em.ERROR_Code, em_library.I18n.TranslateFromRequest(ctx, "ERROR_Validate"), nil, err)
		}
	}

	// Request -> User
	var user model.User
	u, err := user.InterfaceToUser(request)
	if err != nil {
		return em.ErrorRpc(codes.InvalidArgument, em.ERROR_Code, em_library.I18n.TranslateFromRequest(ctx, "ERROR_Create"), nil, err)
	}

	// Check if Username exists
	// 检查Username是否存在
	var count_username int64
	em.DB.Model(&model.User{}).Where("username = ?", u.Username).Count(&count_username)
	if count_username != 0 {
		em.LogInfo.Output(em_utils.MessageWithLineNum("Username already exists"))
		return em.ErrorRpc(codes.InvalidArgument, em.ERROR_Code, em_library.I18n.TranslateFromRequest(ctx, "ERROR_MESSAGE_DuplicateUserName"), nil, errors.New("Username already exists"))
	}

	// Check if the role exists
	// 检查role是否存在
	var role_ids []uint
	for _, v := range u.Roles {
		role_ids = append(role_ids, v.ID)
	}
	var count int64
	em.DB.Model(&model.Role{}).Where("id IN ?", role_ids).Count(&count)
	if int(count) != len(role_ids) {
		em.LogError.Output(em_utils.MessageWithLineNum("Role does not exist"))
		return em.ErrorRpc(codes.InvalidArgument, em.ERROR_Code, em_library.I18n.TranslateFromRequest(ctx, "ERROR_Create"), nil, errors.New("Role does not exist"))
	}

	// Create User
	// 创建用户
	err = em.DB.Transaction(func(tx *gorm.DB) error {
		// Bcrypt Password
		u.Password, err = user.BcryptPassword(u.Password)
		if err != nil {
			return em.LogError.OutputAndReturnError(em_utils.MessageWithLineNum("Password encryption failed" + err.Error()))
		}

		// Create User
		result := tx.Create(&u)
		if result.Error != nil {
			return em.LogError.OutputAndReturnError(em_utils.MessageWithLineNum("Create user failed" + result.Error.Error()))
		}

		return nil
	})
	if err != nil {
		return em.ErrorRpc(codes.InvalidArgument, em.ERROR_Code, em_library.I18n.TranslateFromRequest(ctx, "ERROR_Create"), nil, err)
	}

	// Delete Cache
	// 删除缓存
	if em_library.Config.App.Cache {
		em_library.Cache.DeleteString(application.Cache_UserGetAll)
	}

	data, err := user.InterfaceToUserGetOne(u)
	if err != nil {
		// No need to return
		em.LogError.Output(em_utils.MessageWithLineNum(err.Error()))
	}

	return em.SuccessRpc(em.SUCCESS_Code, em_library.I18n.TranslateFromRequest(ctx, "SUCCESS_Create"), data)
}

// Edit user
// 编辑用户
type validate_UserEdit struct {
	ID uint             `json:"id" validate:"required"`
	Username string     `json:"username" validate:"required,max=255"`
	Roles []model.Role `json:"roles" validate:"required"`
}
func (this *ServiceUser) Edit(ctx context.Context, request *protobuf.UserEdit) (*em_protobuf.Response, error) {
	// Validate
	var vd validate_UserEdit
	err := em_utils.ChangeType(request, &vd)
	if err != nil {
		em.LogError.Output(em_utils.MessageWithLineNum(err.Error()))
		return em.ErrorRpc(codes.InvalidArgument, em.ERROR_Code, em_library.I18n.TranslateFromRequest(ctx, "ERROR_Edit"), nil, err)
	}
	err = em_library.Validator.ValidateStruct(vd)
	if err != nil {
		return em.ErrorRpc(codes.InvalidArgument, em.ERROR_Code, err.Error(), nil, err)
	}

	// Request -> User
	var user model.User
	u, err := user.InterfaceToUser(request)
	if err != nil {
		return em.ErrorRpc(codes.InvalidArgument, em.ERROR_Code, em_library.I18n.TranslateFromRequest(ctx, "ERROR_Edit"), nil, err)
	}

	// Find if the user exists
	// 查找该用户是否存在
	var form model.User
	result := em.DB.First(&form, request.Id)
	if result.RowsAffected == 0 {
		em.LogWarn.Output(em_utils.MessageWithLineNum("No user record"))
		return em.ErrorRpc(codes.InvalidArgument, em.ERROR_Code, em_library.I18n.TranslateFromRequest(ctx, "ERROR_Edit"), nil, errors.New("No user record"))
	}

	// Check if Username exists
	// 检查Username是否存在
	var count_username int64
	em.DB.Model(&model.User{}).Where("username = ?", u.Username).Not(request.Id).Count(&count_username)
	if count_username != 0 {
		em.LogDebug.Output(em_utils.MessageWithLineNum("The user name already exists!"))
		return em.ErrorRpc(codes.InvalidArgument, em.ERROR_Code, em_library.I18n.TranslateFromRequest(ctx, "ERROR_MESSAGE_DuplicateUserName"), nil, errors.New("Username already exists"))
	}

	// Check if the role exists
	// 检查role是否存在
	var role_ids []uint
	for _, v := range u.Roles {
		role_ids = append(role_ids, v.ID)
	}
	var count int64
	em.DB.Model(&model.Role{}).Where("id IN ?", role_ids).Count(&count)
	if int(count) != len(role_ids) {
		em.LogWarn.Output(em_utils.MessageWithLineNum("Role does not exist"))
		return em.ErrorRpc(codes.InvalidArgument, em.ERROR_Code, em_library.I18n.TranslateFromRequest(ctx, "ERROR_Edit"), nil, errors.New("Role does not exist"))
	}

	// If user set new password
	if len(u.Password) > 0 {
		var user model.User
		u.Password, err = user.BcryptPassword(u.Password)
		if err != nil {
			em.LogError.Output(em_utils.MessageWithLineNum(err.Error()))
			return em.ErrorRpc(codes.InvalidArgument, em.ERROR_Code, em_library.I18n.TranslateFromRequest(ctx, "ERROR_Edit"), nil, err)
		}
	}

	err = em.DB.Transaction(func(tx *gorm.DB) error {
		// Replace association
		// 替换关联
		var roleslist []model.Role
		for _, v := range u.Roles {
			roleslist = append(roleslist, model.Role{
				ID:          v.ID,
			})
		}
		err = tx.Model(&model.User{ID: u.ID}).Association("Roles").Replace(roleslist)
		if err != nil {
			return em.LogError.OutputAndReturnError(em_utils.MessageWithLineNum(err.Error()))
		}

		// Update operation, the updates method will not affect the association
		// 更新操作，updates方法不会影响关联
		result := tx.Model(&model.User{}).Where(u.ID).Updates(u)
		if result.Error != nil {
			return em.LogError.OutputAndReturnError(em_utils.MessageWithLineNum(result.Error.Error()))
		}

		return nil
	})
	if err != nil {
		return em.ErrorRpc(codes.InvalidArgument, em.ERROR_Code, em_library.I18n.TranslateFromRequest(ctx, "ERROR_Edit"), nil, err)
	}

	// Delete Cache
	// 删除缓存
	if em_library.Config.App.Cache {
		em_library.Cache.DeleteString(application.Cache_UserGetAll)
		em_library.Cache.DeleteHash(application.Cache_UserGetCurrent, strconv.Itoa(int(request.Id)))
	}

	data, err := user.InterfaceToUserGetOne(u)
	if err != nil {
		// No need to return
		em.LogError.Output(em_utils.MessageWithLineNum(err.Error()))
	}

	return em.SuccessRpc(em.SUCCESS_Code, em_library.I18n.TranslateFromRequest(ctx, "SUCCESS_Edit"), data)
}

// Delete user
// 删除用户
type validate_UserDelete struct {
	Users []model.User `json:"users" validate:"required"`
}
func (this *ServiceUser) Delete(ctx context.Context, request *protobuf.UserDelete) (*em_protobuf.Response, error) {
	// Validate
	var vd validate_UserDelete
	err := em_utils.ChangeType(request, &vd)
	if err != nil {
		em.LogError.Output(em_utils.MessageWithLineNum(err.Error()))
		return em.ErrorRpc(codes.InvalidArgument, em.ERROR_Code, em_library.I18n.TranslateFromRequest(ctx, "ERROR_Delete"), nil, err)
	}
	err = em_library.Validator.ValidateStruct(vd)
	if err != nil {
		return em.ErrorRpc(codes.InvalidArgument, em.ERROR_Code, err.Error(), nil, err)
	}

	var ids []int
	for _, v := range request.Users {
		ids = append(ids, int(v.Id))
	}

	// Find if admin is included in ids
	// 查找ids中是否包含admin
	b := em_utils.CheckIfSliceContainsInt(1, ids)
	if b {
		em.LogWarn.Output(em_utils.MessageWithLineNum("Cannot delete administrator"))
		return em.ErrorRpc(codes.InvalidArgument, em.ERROR_Code, em_library.I18n.TranslateFromRequest(ctx, "ERROR_MESSAGE_ProhibitOperationOfAdministratorUsers"), nil, errors.New("Cannot delete administrator"))
	}

	err = em.DB.Transaction(func(tx *gorm.DB) error {
		var u []model.User
		tx.Where("id IN ?", ids).Find(&u)

		// 删除用户
		result := tx.Delete(&u)
		if result.Error != nil {
			em.LogError.Output(em_utils.MessageWithLineNum(result.Error.Error()))
			return result.Error
		}

		// 删除关联
		err = tx.Model(&u).Association("Roles").Clear()
		if err != nil {
			em.LogError.Output(em_utils.MessageWithLineNum(err.Error()))
			return err
		}

		return nil
	})
	if err != nil {
		return em.ErrorRpc(codes.InvalidArgument, em.ERROR_Code, em_library.I18n.TranslateFromRequest(ctx, "ERROR_Delete"), nil, err)
	}

	// Delete Cache
	// 删除缓存
	if em_library.Config.App.Cache {
		em_library.Cache.DeleteString(application.Cache_UserGetAll)
		var tmp []string
		for _, v := range ids {
			tmp = append(tmp, strconv.Itoa(int(v)))
		}
		em_library.Cache.DeleteHash(application.Cache_UserGetCurrent, strings.Join(tmp, " "))
	}

	return em.SuccessRpc(em.SUCCESS_Code, em_library.I18n.TranslateFromRequest(ctx, "SUCCESS_Delete"), nil)
}

// Update user information
// 更新用户信息
type validate_UserUpdateInformation struct {
	Password string           `json:"password" validate:"omitempty,min=6,max=50"`
	Avatar   model.Attachment `json:"avatar"`
}
func (this *ServiceUser) UpdateInformation(ctx context.Context, request *protobuf.UserUpdateInformation) (*em_protobuf.Response, error) {
	// Validate
	{
		err := em_library.Validator.Validate(request, &validate_UserUpdateInformation{})
		if err != nil {
			em.LogWarn.Output(em_utils.MessageWithLineNum(err.Error()))
			return em.ErrorRpc(codes.InvalidArgument, em.ERROR_Code, em_library.I18n.TranslateFromRequest(ctx, "ERROR_Validate"), nil, err)
		}
	}

	// Get User id
	var user model.User
	id, err := user.GetUserIdByRequest(ctx)
	if err != nil {
		return em.ErrorRpc(codes.InvalidArgument, em.ERROR_Code, em_library.I18n.TranslateFromRequest(ctx, "ERROR_Update"), nil, err)
	}

	// Request -> User
	u, err := user.InterfaceToUser(request)
	if err != nil {
		return em.ErrorRpc(codes.InvalidArgument, em.ERROR_Code, em_library.I18n.TranslateFromRequest(ctx, "ERROR_Update"), nil, err)
	}

	// Update
	err = em.DB.Transaction(func(tx *gorm.DB) error {

		// Create avatar attachment
		err := client.NewClient().User_CreateAvatar(ctx, request.GetAvatar().GetPath(), uint32(id), application.Relationship_User_Avatar)
		if err != nil {
			return err
		}

		// Update password if exists
		if len(u.Password) > 0 {
			u.Password, err = user.BcryptPassword(u.Password)
		}

		result := tx.Model(&model.User{ID: id}).Updates(&u)
		if result.Error != nil {
			em.LogError.Output(em_utils.MessageWithLineNum(result.Error.Error()))
			return result.Error
		}

		return nil
	})
	if err != nil {
		return em.ErrorRpc(codes.InvalidArgument, em.ERROR_Code, em_library.I18n.TranslateFromRequest(ctx, "ERROR_Update"), nil, err)
	}

	if em_library.Config.App.Cache {
		em_library.Cache.DeleteString(application.Cache_UserGetAll)
		em_library.Cache.DeleteHash(application.Cache_UserGetCurrent, strconv.Itoa(int(id)))
	}

	return em.SuccessRpc(em.SUCCESS_Code, em_library.I18n.TranslateFromRequest(ctx, "SUCCESS_Update"), nil)
}
