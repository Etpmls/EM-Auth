package service

import (
	"context"
	"github.com/Etpmls/EM-Auth/src/application"
	"github.com/Etpmls/EM-Auth/src/application/middleware"
	"github.com/Etpmls/EM-Auth/src/application/model"
	em "github.com/Etpmls/Etpmls-Micro"
	em_protobuf "github.com/Etpmls/Etpmls-Micro/protobuf"
	"github.com/Etpmls/Etpmls-Micro/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ServiceAuth struct {
	em_protobuf.UnimplementedAuthServer
}

func (this *ServiceAuth) Check(ctx context.Context, request *em_protobuf.AuthCheck) (response *em_protobuf.AuthCheckResponse, err error) {

	var permission model.Permission
	ps, err := permission.GetAll()
	if err != nil {
		em.LogError.Output(em_utils.MessageWithLineNum(err.Error()))
		return nil, status.Error(codes.Internal, err.Error())
	}

	var p model.Permission
	for _, v := range ps {
		if v.Path == request.GetService() {
			p = v
			break
		}
	}
	if p.ID == 0 {
		return &em_protobuf.AuthCheckResponse{
			Success: false,
		}, nil
	}


	switch p.Auth {
	case application.Auth_NoVerify:
		return &em_protobuf.AuthCheckResponse{
			Success: true,
		}, nil
	case application.Auth_BasicVerify:
		return &em_protobuf.AuthCheckResponse{
			Success: true,
		}, nil
	case application.Auth_AdvancedVerify:
		err := middleware.NewAuth().PermissionVerify(request.GetService(), string(request.GetUserId()))
		if err == nil {
			return &em_protobuf.AuthCheckResponse{
				Success: true,
			}, nil
		}
	}

	return &em_protobuf.AuthCheckResponse{
		Success: false,
	}, nil
}