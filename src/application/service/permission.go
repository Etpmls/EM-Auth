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
)


type ServicePermission struct {
	protobuf.UnimplementedPermissionServer
}

// Get all permissions
// 获取全部权限
func (this *ServicePermission) GetAll(ctx context.Context, request *em_protobuf.Pagination) (*em_protobuf.Response, error) {
	type Permission model.PermissionGetOne
	var data []Permission

	// 获取分页和标题
	var orm em_library.Gorm
	limit, offset := orm.GeneratePaginationLimit(int(request.Number), int(request.Size))
	var count int64
	// Get the title of the search, if not get all the data
	// 获取搜索的标题，如果没有获取全部数据
	search := request.Search

	em.DB.Model(&Permission{}).Where("name " + em.FUZZY_SEARCH + " ?", "%"+ search +"%").Count(&count).Limit(limit).Offset(offset).Find(&data)

	return em.SuccessRpc(em.SUCCESS_Code, em.I18n.TranslateFromRequest(ctx, "SUCCESS_Get"), map[string]interface{}{application.FieldData: data, application.FieldCount: count})
}

// Create Permission
// 创建权限
type validate_PermissionCreate struct {
	Name string `json:"name" validate:"required,max=255"`
	Auth int	`json:"auth" validate:"min=0,max=10"`
	Method string `json:"method" validate:"required,min=1"`
	Path	string	`json:"path" validate:"required,max=255"`
	Remark string `json:"remark" validate:"max=255"`
}
func (this *ServicePermission) Create(ctx context.Context, request *protobuf.PermissionCreate) (*em_protobuf.Response, error) {
	// Validate
	{
		err := em.Validator.Validate(request, &validate_PermissionCreate{})
		if err != nil {
			em.LogWarn.Output(em.MessageWithLineNum(err.Error()))
			return em.ErrorRpc(codes.InvalidArgument, em.ERROR_Code, em.I18n.TranslateFromRequest(ctx, "ERROR_Validate"), nil, err)
		}
	}

	// Request -> Permission
	var permission model.Permission
	p, err := permission.InterfaceToPermission(request)
	if err != nil {
		return em.ErrorRpc(codes.InvalidArgument, em.ERROR_Code, em.I18n.TranslateFromRequest(ctx, "ERROR_Create"), nil, err)
	}
	
	// Create
	err = em.DB.Transaction(func(tx *gorm.DB) error {
		// Insert Data
		result := tx.Create(&p)
		if result.Error != nil {
			em.LogError.Output(em.MessageWithLineNum(result.Error.Error()))
			return result.Error
		}

		return nil
	})
	if err != nil {
		em.LogError.Output(em.MessageWithLineNum(err.Error()))
		return em.ErrorRpc(codes.InvalidArgument, em.ERROR_Code, em.I18n.TranslateFromRequest(ctx, "ERROR_Create"), nil, err)
	}

	e, _ := em.Kv.ReadKey(define.KvCacheEnable)
	if strings.ToLower(e) == "true" {
		em.Cache.DeleteString(application.Cache_PermissionGetAll)
	}


	return em.SuccessRpc(em.SUCCESS_Code, em.I18n.TranslateFromRequest(ctx, "SUCCESS_Create"), nil)
}

// Edit Permission
// 修改权限
type validate_PermissionEdit struct {
	ID        uint `json:"id" validate:"required,min=1"`
	validate_PermissionCreate
}
func (this *ServicePermission) Edit(ctx context.Context, request *protobuf.PermissionEdit) (*em_protobuf.Response, error) {
	// Validate
	{
		err := em.Validator.Validate(request, &validate_PermissionEdit{})
		if err != nil {
			em.LogWarn.Output(em.MessageWithLineNum(err.Error()))
			return em.ErrorRpc(codes.InvalidArgument, em.ERROR_Code, em.I18n.TranslateFromRequest(ctx, "ERROR_Validate"), nil, err)
		}
	}

	// Request -> Permission
	var permission model.Permission
	p, err := permission.InterfaceToPermission(request)
	if err != nil {
		return em.ErrorRpc(codes.InvalidArgument, em.ERROR_Code, em.I18n.TranslateFromRequest(ctx, "ERROR_Edit"), nil, err)
	}

	err = em.DB.Transaction(func(tx *gorm.DB) error {
		var old_p model.Permission
		result := tx.First(&old_p, request.Id)
		if result.RowsAffected == 0 {
			return errors.New("No record.")
		}

		p.CreatedAt = old_p.CreatedAt
		tx.Save(&p)

		return nil
	})
	if err != nil {
		em.LogError.Output(em.MessageWithLineNum(err.Error()))
		return em.ErrorRpc(codes.InvalidArgument, em.ERROR_Code, em.I18n.TranslateFromRequest(ctx, "ERROR_Edit"), nil, err)
	}

	e, _ := em.Kv.ReadKey(define.KvCacheEnable)
	if strings.ToLower(e) == "true" {
		em.Cache.DeleteString(application.Cache_PermissionGetAll)
	}


	return em.SuccessRpc(em.SUCCESS_Code, em.I18n.TranslateFromRequest(ctx, "SUCCESS_Edit"), nil)
}

// Delete Permission
// 删除权限
type validate_PermissionDelete struct {
	Permissions []model.Permission `json:"permissions" validate:"required"`
}
func (this *ServicePermission) Delete(ctx context.Context, request *protobuf.PermissionDelete) (*em_protobuf.Response, error) {
	// Validate
	{
		err := em.Validator.Validate(request, &validate_PermissionDelete{})
		if err != nil {
			em.LogWarn.Output(em.MessageWithLineNum(err.Error()))
			return em.ErrorRpc(codes.InvalidArgument, em.ERROR_Code, em.I18n.TranslateFromRequest(ctx, "ERROR_Validate"), nil, err)
		}
	}

	var ids []int
	for _, v := range request.Permissions {
		ids = append(ids, int(v.Id))
	}

	err := em.DB.Transaction(func(tx *gorm.DB) error {
		var p []model.Permission
		tx.Where("id IN ?", ids).Find(&p)

		// Delete Permission
		// 删除权限
		result := tx.Delete(&p)
		if result.Error != nil {
			return result.Error
		}

		// Delete Association
		// 删除关联
		err := tx.Model(&p).Association("Roles").Clear()
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		em.LogError.Output(em.MessageWithLineNum(err.Error()))
		return em.ErrorRpc(codes.InvalidArgument, em.ERROR_Code, em.I18n.TranslateFromRequest(ctx, "ERROR_Delete"), nil, err)
	}

	e, _ := em.Kv.ReadKey(define.KvCacheEnable)
	if strings.ToLower(e) == "true" {
		em.Cache.DeleteString(application.Cache_PermissionGetAll)
	}


	return em.SuccessRpc(em.SUCCESS_Code, em.I18n.TranslateFromRequest(ctx, "SUCCESS_Delete"), nil)
}

// Get Advanced verify permissions
// 获取高级认证权限
func (this *ServicePermission) GetAdvancedVerify(ctx context.Context, request *em_protobuf.Empty) (*em_protobuf.Response, error) {
	type Permission model.PermissionGetOne
	var data []Permission
	em.DB.Where("auth = ?", application.Auth_AdvancedVerify).Find(&data)
	return em.SuccessRpc(em.SUCCESS_Code, em.I18n.TranslateFromRequest(ctx, "SUCCESS_Get"), data)
}
