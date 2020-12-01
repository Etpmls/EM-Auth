package register

import (
	"context"
	"github.com/Etpmls/EM-Auth/src/application/protobuf"
	"github.com/Etpmls/EM-Auth/src/application/service"
	em_protobuf "github.com/Etpmls/Etpmls-Micro/protobuf"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

// Register Rpc Service
func RegisterRpcService(s *grpc.Server)  {
	// protobuf.RegisterUserServer(s, &service.ServiceUser{})
	protobuf.RegisterUserServer(s, &service.ServiceUser{})
	protobuf.RegisterRoleServer(s, &service.ServiceRole{})
	protobuf.RegisterMenuServer(s, &service.ServiceMenu{})
	protobuf.RegisterPermissionServer(s, &service.ServicePermission{})
	protobuf.RegisterSettingServer(s, &service.ServiceSetting{})
	em_protobuf.RegisterAuthServer(s, &service.ServiceAuth{})
	return
}

// Register Http Service
func RegisterHttpService(ctx context.Context, mux *runtime.ServeMux, grpcServerEndpoint *string, opts []grpc.DialOption) error {
	/*err := protobuf.RegisterUserHandlerFromEndpoint(ctx, mux,  *grpcServerEndpoint, opts)
	if err != nil {
		return err
	}*/
	err := protobuf.RegisterUserHandlerFromEndpoint(ctx, mux,  *grpcServerEndpoint, opts)
	if err != nil {
		return err
	}
	err = protobuf.RegisterRoleHandlerFromEndpoint(ctx, mux,  *grpcServerEndpoint, opts)
	if err != nil {
		return err
	}
	err = protobuf.RegisterMenuHandlerFromEndpoint(ctx, mux,  *grpcServerEndpoint, opts)
	if err != nil {
		return err
	}
	err = protobuf.RegisterPermissionHandlerFromEndpoint(ctx, mux,  *grpcServerEndpoint, opts)
	if err != nil {
		return err
	}
	err = protobuf.RegisterSettingHandlerFromEndpoint(ctx, mux,  *grpcServerEndpoint, opts)
	if err != nil {
		return err
	}

	return nil
}
