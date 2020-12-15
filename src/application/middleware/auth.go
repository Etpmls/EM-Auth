package middleware

import (
	"context"
	"github.com/Etpmls/EM-Auth/src/application/protobuf"
	em "github.com/Etpmls/Etpmls-Micro"
	em_library "github.com/Etpmls/Etpmls-Micro/library"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Only Verify Token
// 仅验证token
// Ensure the security of the intranet (without going through the API gateway)
// 保证内网安全（不经过API网关）
func (this *middleware) Auth() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		// If it is the route of Baisc, skip
		// 如果是基础的路由，则跳过
		switch req.(type) {
		case *protobuf.UserLogin:
			return handler(ctx, req)
		case *protobuf.UserGetCurrent:
			return handler(ctx, req)
		}

		// Get token from header
		token, err:= em.Micro.Auth.Rpc_GetTokenFromHeader(ctx)
		if err != nil || token == "" {
			return nil, status.Error(codes.PermissionDenied, em_library.I18n.TranslateFromRequest(ctx, "ERROR_MESSAGE_PermissionDenied"))
		}

		b, _ := em.Micro.Auth.VerifyToken(token)
		if !b {
			return nil, status.Error(codes.PermissionDenied, em_library.I18n.TranslateFromRequest(ctx, "ERROR_MESSAGE_PermissionDenied"))
		}

		// Pass the token to the method
		// 把token传递到方法中
		ctx = context.WithValue(ctx,"token", token)

		em.LogDebug.Output("Auth middleware runs successfully!")	// Debug
		return handler(ctx, req)
	}
}

