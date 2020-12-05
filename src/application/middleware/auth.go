package middleware

import (
	"context"
	"errors"
	"github.com/Etpmls/EM-Auth/src/application"
	"github.com/Etpmls/EM-Auth/src/application/model"
	em "github.com/Etpmls/Etpmls-Micro"
	em_library "github.com/Etpmls/Etpmls-Micro/library"
	em_protobuf "github.com/Etpmls/Etpmls-Micro/protobuf"
	em_utils "github.com/Etpmls/Etpmls-Micro/utils"
	"github.com/dgrijalva/jwt-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"path/filepath"
	"strconv"
)

func (this *middleware) Auth() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		// fullMethodName: /protobuf.User/GetCurrent
		service := em_library.NewGrpc().GetServiceName(info.FullMethod)

		// // If it is the route of UserLogin, skip
		// 如果是AuthCheck的路由，则跳过
		_, ok := req.(*em_protobuf.AuthCheck)
		if ok {
			return handler(ctx, req)
		}

		var permission model.Permission
		ps, err := permission.GetAll()
		if err != nil {
			em.LogError.Output(em_utils.MessageWithLineNum(err.Error()))
			return em.ErrorRpc(codes.Internal, em.ERROR_Code, em_library.I18n.TranslateFromRequest(ctx, "ERROR_MESSAGE_GetPermission"), nil, err)
		}
		var p model.Permission
		for _, v := range ps {
			if v.Path == service {
				p = v
				break
			}
		}
		if p.ID == 0 {
			return em.ErrorRpc(codes.InvalidArgument, em.ERROR_Code, em_library.I18n.TranslateFromRequest(ctx, "ERROR_MESSAGE_PermissionNotRegistered"), nil, errors.New("Permission not registered!"))
		}

		// No auth
		switch p.Auth {
		case application.Auth_NoVerify:
			return handler(ctx, req)
		case application.Auth_BasicVerify:
			return NewAuth().BasicVerify(ctx, req, handler)
		case application.Auth_AdvancedVerify:
			return NewAuth().AdvancedVerify(ctx, req, handler, service)
		}

		return em.ErrorRpc(codes.InvalidArgument, em.ERROR_Code, em_library.I18n.TranslateFromRequest(ctx, "ERROR_MESSAGE_PermissionDenied"), nil, errors.New("ERROR_MESSAGE_PermissionDenied"))
	}
}

type auth struct {}

func NewAuth() *auth {
	return &auth{}
}

// Only verify token
// 仅检查token
func (this *auth) BasicVerify(ctx context.Context, req interface{}, handler grpc.UnaryHandler)  (resp interface{}, err error) {
	// Get token
	// 获取令牌
	g := em_library.NewGrpc()
	token, err := g.ExtractHeader(ctx, "token")
	if err != nil {
		return em.ErrorRpc(codes.InvalidArgument, em.ERROR_Code, em_library.I18n.TranslateFromRequest(ctx, "ERROR_MESSAGE_GetToken"), nil, err)
	}

	_, err = em_library.JwtToken.ParseToken(token)
	if err != nil {
		return em.ErrorRpc(codes.Unauthenticated, em.ERROR_Code, em_library.I18n.TranslateFromRequest(ctx, "ERROR_MESSAGE_TokenVerificationFailed"), nil, err)
	}

	// Pass the token to the method
	// 把token传递到方法中
	ctx = context.WithValue(ctx,"token", token)

	return handler(ctx, req)
}

func (this *auth) AdvancedVerify(ctx context.Context, req interface{}, handler grpc.UnaryHandler, fullMethodName string) (resp interface{}, err error) {
	// Get token
	// 获取令牌
	g := em_library.NewGrpc()
	token, err := g.ExtractHeader(ctx, "token")
	if err != nil {
		return em.ErrorRpc(codes.InvalidArgument, em.ERROR_Code,  fullMethodName + ":" + em_library.I18n.TranslateFromRequest(ctx, "ERROR_MESSAGE_GetToken"), nil, err)
	}

	// Get Claims
	// 获取Claims
	tmp, err := em_library.JwtToken.ParseToken(token)
	tk, ok := tmp.(*jwt.Token)
	if !ok || err != nil {
		return em.ErrorRpc(codes.Unauthenticated, em.ERROR_Code,  fullMethodName + ":" + em_library.I18n.TranslateFromRequest(ctx, "ERROR_MESSAGE_TokenVerificationFailed"), nil, err)
	}

	// Determine whether the role has the corresponding permissions
	// 判断所属角色是否有相应的权限
	if claims,ok := tk.Claims.(jwt.MapClaims); ok && tk.Valid {
		if userId, ok := claims["jti"].(string); ok {
			err := NewAuth().PermissionVerify(fullMethodName, userId)
			if err == nil {
				// Pass the token to the method
				// 把token传递到方法中
				ctx = context.WithValue(ctx,"token", token)
				return handler(ctx, req)
			}
		}
	}

	return em.ErrorRpc(codes.InvalidArgument, em.ERROR_Code, fullMethodName + ":" + em_library.I18n.TranslateFromRequest(ctx, "ERROR_MESSAGE_PermissionDenied"), nil, errors.New("Permission denied"))
}

func (this *auth) PermissionVerify(service string, idStr string) (err error) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return err
	}

	// 1. Get user ID
	// 1.获取用户ID
	var u model.User
	em.DB.Preload("Roles").First(&u, id)
	var ids []uint
	for _, v := range u.Roles {
		// If it is an administrator group
		// 如果为管理员组
		if v.ID == 1 {
			return nil
		}
		ids = append(ids, v.ID)
	}
	// Get role related permissions
	// 获取角色相关权限
	var r []model.Role
	em.DB.Preload("Permissions").Where(ids).Find(&r)

	// Determine whether there is a request permission
	// 判断是否有请求权限
	for _, v := range r {
		for _, subv := range v.Permissions {
			// Path comparison
			// 路径对比
			b, _ := filepath.Match(subv.Path, service)
			if b {
				return nil
			}
		}
	}

	return errors.New("Permission verification failed")
}

