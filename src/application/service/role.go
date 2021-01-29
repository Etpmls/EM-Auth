package service

import (
	"context"
	"errors"
	"github.com/Etpmls/EM-Auth/src/application"
	"github.com/Etpmls/EM-Auth/src/application/model"
	"github.com/Etpmls/EM-Auth/src/application/protobuf"
	em "github.com/Etpmls/Etpmls-Micro/v2"
	"github.com/Etpmls/Etpmls-Micro/v2/define"
	"github.com/Etpmls/Etpmls-Micro/v2/library"
	em_protobuf "github.com/Etpmls/Etpmls-Micro/v2/protobuf"
	"strings"

	"google.golang.org/grpc/codes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)


type ServiceRole struct {
	protobuf.UnimplementedRoleServer
}

// Get all roles
// 获取全部角色
func (this *ServiceRole) GetAll(ctx context.Context, request *em_protobuf.Pagination) (*em_protobuf.Response, error) {
	type Role model.RoleGetOne
	var data []Role

	// 获取分页和标题
	var orm em_library.Gorm
	limit, offset := orm.GeneratePaginationLimit(int(request.Number), int(request.Size))
	var count int64
	// Get the title of the search, if not get all the data
	// 获取搜索的标题，如果没有获取全部数据
	search := request.Search

	em.DB.Model(&Role{}).Preload("Permissions").Where("name " +em.FUZZY_SEARCH+ " ?", "%"+ search +"%").Count(&count).Limit(limit).Offset(offset).Find(&data)

	return em.SuccessRpc(em.SUCCESS_Code, em.I18n.TranslateFromRequest(ctx, "SUCCESS_Get"), map[string]interface{}{application.FieldData: data, application.FieldCount: count})
}

// Create Role
// 创建角色
type validate_RoleCreate struct {
	Name string                     `json:"name" validate:"required,max=30"`
	Remark string                   `json:"remark" validate:"max=255"`
	Permissions []model.Permission `json:"permissions"`
}
func (this *ServiceRole) Create(ctx context.Context, request *protobuf.RoleCreate) (*em_protobuf.Response, error) {
	// Validate
	{
		err := em.Validator.Validate(request, &validate_RoleCreate{})
		if err != nil {
			em.LogWarn.Output(em.MessageWithLineNum(err.Error()))
			return em.ErrorRpc(codes.InvalidArgument, em.ERROR_Code, em.I18n.TranslateFromRequest(ctx, "ERROR_Validate"), nil, err)
		}
	}

	// Request -> Role
	var role model.Role
	r, err := role.InterfaceToRole(request)
	if err != nil {
		return em.ErrorRpc(codes.InvalidArgument, em.ERROR_Code, em.I18n.TranslateFromRequest(ctx, "ERROR_Create"), nil, err)
	}

	// Create
	err = em.DB.Transaction(func(tx *gorm.DB) error {
		// Insert Data
		result := tx.Create(&r)
		if result.Error != nil {
			em.LogError.Output(em.MessageWithLineNum(result.Error.Error()))
			return result.Error
		}

		return nil
	})
	if err != nil {
		return em.ErrorRpc(codes.InvalidArgument, em.ERROR_Code, em.I18n.TranslateFromRequest(ctx, "ERROR_Create"), nil, err)
	}

	e, _ := em.Kv.ReadKey(define.KvCacheEnable)
	if strings.ToLower(e) == "true" {
		em.Cache.DeleteString(application.Cache_RoleGetAll)
	}


	return em.SuccessRpc(em.SUCCESS_Code, em.I18n.TranslateFromRequest(ctx, "SUCCESS_Create"), nil)
}

// Edit Role
// 编辑角色
type validate_RoleEdit struct {
	ID        uint `json:"id" validate:"required,min=1"`
	validate_RoleCreate
}
func (this *ServiceRole) Edit(ctx context.Context, request *protobuf.RoleEdit) (*em_protobuf.Response, error) {
	// Validate
	{
		err := em.Validator.Validate(request, &validate_RoleEdit{})
		if err != nil {
			em.LogWarn.Output(em.MessageWithLineNum(err.Error()))
			return em.ErrorRpc(codes.InvalidArgument, em.ERROR_Code, em.I18n.TranslateFromRequest(ctx, "ERROR_Validate"), nil, err)
		}
	}

	// Request -> Role
	var role model.Role
	r, err := role.InterfaceToRole(request)
	if err != nil {
		return em.ErrorRpc(codes.InvalidArgument, em.ERROR_Code, em.I18n.TranslateFromRequest(ctx, "ERROR_Edit"), nil, err)
	}

	err = em.DB.Transaction(func(tx *gorm.DB) error {
		// Edit Role
		// 修改角色
		var old_r model.Role
		result := tx.First(&old_r, request.Id)
		if result.RowsAffected == 0 {
			return errors.New("No record.")
		}

		r.CreatedAt = old_r.CreatedAt
		tx.Omit(clause.Associations).Save(&r)

		// Update Relationship
		// 更新关联
		err := em.DB.Model(&r).Association("Permissions").Replace(r.Permissions)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return em.ErrorRpc(codes.InvalidArgument, em.ERROR_Code, em.I18n.TranslateFromRequest(ctx, "ERROR_Edit"), nil, err)
	}

	e, _ := em.Kv.ReadKey(define.KvCacheEnable)
	if strings.ToLower(e) == "true" {
		em.Cache.DeleteString(application.Cache_RoleGetAll)
	}


	return em.SuccessRpc(em.SUCCESS_Code, em.I18n.TranslateFromRequest(ctx, "SUCCESS_Edit"), nil)
}

// Delete Role
// 删除角色
type validate_RoleDelete struct {
	Roles []model.Role `json:"roles" validate:"required"`
}
func (this *ServiceRole) Delete(ctx context.Context, request *protobuf.RoleDelete) (*em_protobuf.Response, error) {
	// Validate
	{
		err := em.Validator.Validate(request, &validate_RoleDelete{})
		if err != nil {
			em.LogWarn.Output(em.MessageWithLineNum(err.Error()))
			return em.ErrorRpc(codes.InvalidArgument, em.ERROR_Code, em.I18n.TranslateFromRequest(ctx, "ERROR_Validate"), nil, err)
		}
	}


	var ids []int
	for _, v := range request.Roles {
		ids = append(ids, int(v.Id))
	}

	err := em.DB.Transaction(func(tx *gorm.DB) error {
		var r []model.Role
		tx.Where("id IN ?", ids).Find(&r)

		// Delete Role
		// 删除角色
		result := tx.Where("id IN ?", ids).Delete(&model.Role{})
		if result.Error != nil {
			em.LogError.Output(em.MessageWithLineNum(result.Error.Error()))
			return result.Error
		}

		// Delete Association
		// 删除关联
		err := tx.Model(&r).Association("Users").Clear()
		if err != nil {
			em.LogError.Output(em.MessageWithLineNum(err.Error()))
			return err
		}
		err = tx.Model(&r).Association("Permissions").Clear()
		if err != nil {
			em.LogError.Output(em.MessageWithLineNum(err.Error()))
			return err
		}

		return nil
	})
	if err != nil {
		return em.ErrorRpc(codes.InvalidArgument, em.ERROR_Code, em.I18n.TranslateFromRequest(ctx, "ERROR_Delete"), nil, err)
	}

	e, _ := em.Kv.ReadKey(define.KvCacheEnable)
	if strings.ToLower(e) == "true" {
		em.Cache.DeleteString(application.Cache_RoleGetAll)
	}


	return em.SuccessRpc(em.SUCCESS_Code, em.I18n.TranslateFromRequest(ctx, "SUCCESS_Delete"), nil)
}