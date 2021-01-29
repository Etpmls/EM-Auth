package client

import (
	"context"
	"github.com/Etpmls/EM-Auth/src/application"
	"github.com/Etpmls/EM-Auth/src/application/protobuf"
	em "github.com/Etpmls/Etpmls-Micro/v2"
	em_protobuf "github.com/Etpmls/Etpmls-Micro/v2/protobuf"
)

func (this *client) Attachment_Delete(ctx *context.Context, owner_ids []uint32, owner_type string) error {
	cl := em.Micro.Client.NewClient()
	err := cl.ConnectServiceWithToken(application.Service_Attachment, ctx)
	if err != nil {
		em.LogWarn.OutputSimplePath(err)
		return err
	}

	c := protobuf.NewAttachmentClient(cl.Conn)
	return cl.Sync_Simple(func() (response *em_protobuf.Response, e error) {
		return c.Delete(*ctx, &protobuf.AttachmentDelete{
			Service:       em.Micro.Config.Service.RpcName,
			OwnerIds:       owner_ids,
			OwnerType:     owner_type,
		})
	},nil)
}

func (this *client) Attachment_CreateMany(ctx *context.Context, paths []string, owner_id uint32, owner_type string) error {
	cl := em.Micro.Client.NewClient()
	err := cl.ConnectServiceWithToken(application.Service_Attachment, ctx)
	if err != nil {
		em.LogWarn.OutputSimplePath(err)
		return err
	}

	c := protobuf.NewAttachmentClient(cl.Conn)
	return cl.Sync_Simple(func() (response *em_protobuf.Response, e error) {
		return c.CreateMany(*ctx, &protobuf.AttachmentCreateMany{
			Service:       em.Micro.Config.Service.RpcName,
			Paths:         paths,
			OwnerId:       owner_id,
			OwnerType:     owner_type,
		})
	},nil)
}

func (this *client) Attachment_GetMany(ctx *context.Context, owner_ids []uint32, owner_type string) ([]byte, error) {
	cl := em.Micro.Client.NewClient()
	err := cl.ConnectServiceWithToken(application.Service_Attachment, ctx)
	if err != nil {
		em.LogWarn.OutputSimplePath(err)
		return nil, err
	}

	c := protobuf.NewAttachmentClient(cl.Conn)
	return cl.Sync_SimpleV2(func() (response *em_protobuf.Response, e error) {
		return c.GetMany(*ctx, &protobuf.AttachmentGetMany{
			Service:       em.Micro.Config.Service.RpcName,
			OwnerIds:       owner_ids,
			OwnerType:     owner_type,
		})
	},nil)
}

func (this *client) Attachment_Append(ctx *context.Context, paths []string, owner_id uint32, owner_type string, cb func(error) error) ([]byte, error) {
	cl := em.Micro.Client.NewClient()
	err := cl.ConnectServiceWithToken(application.Service_Attachment, ctx)
	if err != nil {
		em.LogWarn.OutputSimplePath(err)
		return nil, err
	}

	c := protobuf.NewAttachmentClient(cl.Conn)
	return cl.Sync_SimpleV2(func() (response *em_protobuf.Response, e error) {
		return c.Append(*ctx, &protobuf.AttachmentAppend{
			Service:       em.Micro.Config.Service.RpcName,
			Paths:         paths,
			OwnerId:       owner_id,
			OwnerType:     owner_type,
		})
	},cb)
}