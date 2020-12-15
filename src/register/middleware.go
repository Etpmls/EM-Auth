package register

import (
	"github.com/Etpmls/EM-Auth/src/application/middleware"
	em "github.com/Etpmls/Etpmls-Micro"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"google.golang.org/grpc"
)

// Register GRPC middleware
// 注册GRPC中间件
func RegisterGrpcMiddleware() *grpc.Server {
	s := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				// Panic recover
				grpc_recovery.UnaryServerInterceptor(),
				// I18n
				em.DefaultMiddleware().I18n(),
				// Captcha auth
				middleware.NewMiddleware().Captcha(),
				middleware.NewMiddleware().Auth(),
			),
		),
	)
	return s
}