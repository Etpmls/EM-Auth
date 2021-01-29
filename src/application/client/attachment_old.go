// https://github.com/grpc/grpc-go/blob/15a78f19307d5faf10cfdd9d4e664c65a387cbd1/examples/helloworld/greeter_client/main.go
package client

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Etpmls/EM-Auth/src/application"
	"github.com/Etpmls/EM-Auth/src/application/protobuf"
	em "github.com/Etpmls/Etpmls-Micro/v2"
	em_library "github.com/Etpmls/Etpmls-Micro/v2/library"
)

func (this *client) User_GetAvatar(token string, owner_id uint32, owner_type string) (string, error) {
	// 1.Connect Service
	cl := em.Micro.Client.NewClient()
	err := cl.ConnectService(application.Service_Attachment)
	if err != nil {
		return "", err
	}
	c := protobuf.NewAttachmentClient(cl.Conn)

	// 2. Set Header
	cl.Header = map[string]string{"token": token}

	// 3.Do
	var path string
	err = cl.Sync(func() error {

		// - 1. Request
		r, err := c.GetOne(*cl.Context, &protobuf.AttachmentGetOne{
			Service:       em_library.Config.Service.RpcName,
			OwnerId:       owner_id,
			OwnerType:     owner_type,
		})
		if err != nil {
			em.LogError.Output(em.MessageWithLineNum(err.Error()))
			return err
		}
		// - 2. Get Response
		type res struct {
			Path string	`json:"path"`
		}
		var rsps res
		err = json.Unmarshal([]byte(r.GetData()), &rsps)
		if err != nil {
			return err
		}
		path = rsps.Path

		return nil

	},nil)

	return path, err
}

func (this *client) User_CreateAvatar(ctx context.Context, path string, owner_id uint32, owner_type string) error {
	// 1.Connect Service
	cl := em.Micro.Client.NewClient()
	err := cl.ConnectService(application.Service_Attachment)
	if err != nil {
		return err
	}
	c := protobuf.NewAttachmentClient(cl.Conn)

	// 2. Set Header
	// Get token By Request
	cl.Context = &ctx
	token, err := em.Micro.Auth.GetTokenFromCtx(ctx)
	if err != nil {
		return err
	}
	cl.Header = map[string]string{"token": token}

	// 3.Do
	err = cl.Sync(func() error {

		r, err := c.Create(ctx, &protobuf.AttachmentCreate{
			Service:       em_library.Config.Service.RpcName,
			Path:		   path,
			OwnerId:       owner_id,
			OwnerType:     owner_type,
		})
		if err != nil {
			em.LogError.Output(em.MessageWithLineNum_OneRecord(err.Error()))
			return err
		}

		if r.GetStatus() == em.SUCCESS_Status {
			return nil
		} else {
			em.LogError.Output(em.MessageWithLineNum("Create failed!"))
			return errors.New("Create failed!")
		}

	}, nil)

	return err
}

func (this *client) Setting_DiskCleanUp(ctx context.Context) error {
	// 1.Connect Service
	cl := em.Micro.Client.NewClient()
	err := cl.ConnectService(application.Service_Attachment)
	if err != nil {
		return err
	}
	c := protobuf.NewAttachmentClient(cl.Conn)

	// 2. Set Header
	// Get token By Request
	cl.Context = &ctx
	token, err := em.Micro.Auth.GetTokenFromCtx(ctx)
	if err != nil {
		return err
	}
	cl.Header = map[string]string{"token": token}

	// 3.Do
	err = cl.Sync(func() error {

		r, err := c.DiskCleanUp(ctx, &protobuf.AttachmentDiskCleanUp{
			Service:       em_library.Config.Service.RpcName,
		})
		if err != nil {
			em.LogError.Output(em.MessageWithLineNum(err.Error()))
			return err
		}

		if r.GetStatus() == em.SUCCESS_Status {
			return nil
		} else {
			em.LogError.Output(em.MessageWithLineNum("Delete Failed!"))
			return errors.New("Delete Failed!")
		}

	}, nil)

	return err
}
