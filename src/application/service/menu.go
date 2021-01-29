package service

import (
	"context"
	"encoding/json"
	"github.com/Etpmls/EM-Auth/src/application"
	"github.com/Etpmls/EM-Auth/src/application/protobuf"
	em "github.com/Etpmls/Etpmls-Micro/v2"
	"github.com/Etpmls/Etpmls-Micro/v2/define"
	em_protobuf "github.com/Etpmls/Etpmls-Micro/v2/protobuf"
	"github.com/go-redis/redis/v8"
	"google.golang.org/grpc/codes"
	"strings"
)

type ServiceMenu struct {
	protobuf.UnimplementedMenuServer
}

// Get all menu
// 获取全部菜单
func (this *ServiceMenu) GetAll(ctx context.Context, request *em_protobuf.Empty) (*em_protobuf.Response, error) {
	e, _ := em.Kv.ReadKey(define.KvCacheEnable)
	if strings.ToLower(e) != "true" {
		return this.getAll_NoCache(ctx, request, strings.ToLower(e) == "true")
	}

	return this.getAll_Cache(ctx, request, strings.ToLower(e) == "true")
}
func (this *ServiceMenu) getAll_Cache(ctx context.Context, request *em_protobuf.Empty, enableCache bool) (*em_protobuf.Response, error) {
	// Get the menu from cache
	// 从缓存中获取menu
	ctx_json, err := em.Cache.GetString(application.Cache_MenuGetAll)
	if err != nil {
		if err == redis.Nil {
			return this.getAll_NoCache(ctx, request, enableCache)
		}
		return em.ErrorRpc(codes.InvalidArgument, em.ERROR_Code, em.I18n.TranslateFromRequest(ctx, "ERROR_Get"), nil, em.LogError.OutputAndReturnError(em.MessageWithLineNum(err.Error())))
	}

	return em.SuccessRpc(em.SUCCESS_Code, em.I18n.TranslateFromRequest(ctx, "SUCCESS_Get"), ctx_json)
}
func (this *ServiceMenu) getAll_NoCache(ctx context.Context, request *em_protobuf.Empty, enableCache bool) (*em_protobuf.Response, error) {
	ctx_json, err := em.Kv.ReadKey(define.MakeServiceConfField(em.Micro.Config.Service.RpcName, application.KvServiceMenu))
	if err != nil || len(ctx_json) == 0 || !json.Valid([]byte(ctx_json)) {
		return em.ErrorRpc(codes.InvalidArgument, em.ERROR_Code, "The app/menu is not configured or incorrect.", nil, em.LogError.OutputAndReturnError(em.MessageWithLineNum(err.Error())))
	}
	/* ctx_json, err := ioutil.ReadFile("./storage/menu/menu.json")
	if err != nil {
		return em.ErrorRpc(codes.InvalidArgument, em.ERROR_Code, em.I18n.TranslateFromRequest(ctx, "ERROR_Get"), nil, em.LogError.OutputAndReturnError(em.MessageWithLineNum(err.Error())))
	} */
	// Save menu
	// 储存菜单
	if enableCache {
		em.Cache.SetString(application.Cache_MenuGetAll, ctx_json, 0)
	}

	return em.SuccessRpc(em.SUCCESS_Code, em.I18n.TranslateFromRequest(ctx, "SUCCESS_Get"), string(ctx_json))
}

// Create Menu
// 创建菜单
type validate_MenuCreate struct {
	Menu string `json:"menu" validate:"required"`
}
func (this *ServiceMenu) Create(ctx context.Context, request *protobuf.MenuCreate) (*em_protobuf.Response, error) {
	// Validate
	var vd validate_MenuCreate
	err := em.ChangeType(request, &vd)
	if err != nil {
		em.LogError.Output(em.MessageWithLineNum(err.Error()))
		return em.ErrorRpc(codes.InvalidArgument, em.ERROR_Code, em.I18n.TranslateFromRequest(ctx, "ERROR_Create"), nil, err)
	}
	err = em.Validator.ValidateStruct(vd)
	if err != nil {
		return em.ErrorRpc(codes.InvalidArgument, em.ERROR_Code, err.Error(), nil, err)
	}

	// Backup old menu
	oldMenu, err := em.Kv.ReadKey(define.MakeServiceConfField(em.Micro.Config.Service.RpcName, application.KvServiceMenu))
	if err != nil {
		return em.ErrorRpc(codes.InvalidArgument, em.ERROR_Code, em.I18n.TranslateFromRequest(ctx, "ERROR_Create"), nil, em.LogError.OutputAndReturnError(em.MessageWithLineNum(err.Error())))
	}

	err = em.Kv.CrateOrUpdateKey(define.MakeServiceConfField(em.Micro.Config.Service.RpcName, application.KvServiceMenuBackup), oldMenu)
	if err != nil {
		return em.ErrorRpc(codes.InvalidArgument, em.ERROR_Code, em.I18n.TranslateFromRequest(ctx, "ERROR_Create"), nil, em.LogError.OutputAndReturnError(em.MessageWithLineNum(err.Error())))
	}

	// Write new menu
	err = em.Kv.CrateOrUpdateKey(define.MakeServiceConfField(em.Micro.Config.Service.RpcName, application.KvServiceMenu), request.Menu)
	if err != nil {
		return em.ErrorRpc(codes.InvalidArgument, em.ERROR_Code, em.I18n.TranslateFromRequest(ctx, "ERROR_Create"), nil, em.LogError.OutputAndReturnError(em.MessageWithLineNum(err.Error())))
	}


	// Move files
	// 移动文件
	/* err = os.Rename("storage/menu/menu.json", "storage/menu/menu.json.bak")
	if err != nil {
		em.LogError.Output(em.MessageWithLineNum(err.Error()))
		return em.ErrorRpc(codes.InvalidArgument, em.ERROR_Code, em.I18n.TranslateFromRequest(ctx, "ERROR_Create"), nil, em.LogError.OutputAndReturnError(em.MessageWithLineNum(err.Error())))
	}*/

	// Write file
	// 写入文件
	/*var s = []byte(request.Menu)
	err = ioutil.WriteFile("storage/menu/menu.json", s, 0666)
	if err != nil {
		em.LogError.Output(em.MessageWithLineNum("Failed to write menu file!" + err.Error()))

		// Restore history menu
		// 还原历史菜单
		err2 := os.Rename("storage/menu/menu.json.bak", "storage/menu/menu.json")
		if err2 != nil {
			em.LogError.Output(em.MessageWithLineNum("Failed to restore the backup menu file!" + err2.Error()))
		}

		return em.ErrorRpc(codes.InvalidArgument, em.ERROR_Code, em.I18n.TranslateFromRequest(ctx, "ERROR_Create"), nil, err)
	}*/

	// Delete Cache
	// 删除缓存
	e, _ := em.Kv.ReadKey(define.KvCacheEnable)
	if strings.ToLower(e) == "true" {
		em.Cache.DeleteString(application.Cache_MenuGetAll)
	}

	return em.SuccessRpc(em.SUCCESS_Code, em.I18n.TranslateFromRequest(ctx, "SUCCESS_Create"), nil)
}