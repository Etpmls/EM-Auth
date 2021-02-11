package middleware

import (
	"context"
	"github.com/Etpmls/EM-Auth/src/application/protobuf"
	em "github.com/Etpmls/Etpmls-Micro/v2"
	"github.com/Etpmls/Etpmls-Micro/v2/define"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"strings"
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
		m, _ := em.Kv.List(define.KvCaptcha)
		if strings.ToLower(m[define.KvCaptchaEnable]) != "true" {
			return handler(ctx, req)
		}

		// Verify captcha
		// 验证码验证
		if len(m[define.KvCaptchaSecret]) == 0 {
			return em.ErrorRpc(codes.InvalidArgument, em.ERROR_Code, "Please configure captcha/secret", nil, err)
		}

		ok2, response := em.Captcha.VerifyV2(m[define.KvCaptchaSecret], v.Captcha)
		if ok2 {
			return handler(ctx, req)
		}

		// Debug Error Message
		em.LogDebug.OutputSimplePath(response)
		return em.ErrorRpc(codes.InvalidArgument, em.ERROR_Code, em.I18n.TranslateFromRequest(ctx, "ERROR_MESSAGE_CaptchaVerificationFailed"), nil, err)
	}
}