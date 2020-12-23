package middleware

import (
	"context"
	"github.com/Etpmls/EM-Auth/src/application"
	"github.com/Etpmls/EM-Auth/src/application/protobuf"
	em "github.com/Etpmls/Etpmls-Micro"
	em_library "github.com/Etpmls/Etpmls-Micro/library"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)


func (this *middleware) Captcha() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {

		// // If it is not the route of UserLogin, skip
		// 如果不是UserLogin的路由，则跳过
		v, ok := req.(*protobuf.UserLogin)
		if !ok {
			return handler(ctx, req)
		}

		// If the verification code function is not turned on
		// 如果未开启验证码功能
		if application.ServiceConfig.Service.Captcha == false {
			return handler(ctx, req)
		}

		// Verify captcha
		// 验证码验证
		ok2 := em.Captcha.Verify(em_library.Config.Captcha.Secret, v.Captcha)
		if ok2 {
			return handler(ctx, req)
		}

		return em.ErrorRpc(codes.InvalidArgument, em.ERROR_Code, em.I18n.TranslateFromRequest(ctx, "ERROR_MESSAGE_CaptchaVerificationFailed"), nil, err)
	}
}