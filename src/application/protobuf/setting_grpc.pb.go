// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package protobuf

import (
	context "context"
	protobuf "github.com/Etpmls/Etpmls-Micro/v2/protobuf"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion7

// SettingClient is the client API for Setting service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SettingClient interface {
	CacheClear(ctx context.Context, in *protobuf.Empty, opts ...grpc.CallOption) (*protobuf.Response, error)
	DiskCleanUp(ctx context.Context, in *protobuf.Empty, opts ...grpc.CallOption) (*protobuf.Response, error)
}

type settingClient struct {
	cc grpc.ClientConnInterface
}

func NewSettingClient(cc grpc.ClientConnInterface) SettingClient {
	return &settingClient{cc}
}

func (c *settingClient) CacheClear(ctx context.Context, in *protobuf.Empty, opts ...grpc.CallOption) (*protobuf.Response, error) {
	out := new(protobuf.Response)
	err := c.cc.Invoke(ctx, "/protobuf.Setting/CacheClear", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *settingClient) DiskCleanUp(ctx context.Context, in *protobuf.Empty, opts ...grpc.CallOption) (*protobuf.Response, error) {
	out := new(protobuf.Response)
	err := c.cc.Invoke(ctx, "/protobuf.Setting/DiskCleanUp", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SettingServer is the server API for Setting service.
// All implementations must embed UnimplementedSettingServer
// for forward compatibility
type SettingServer interface {
	CacheClear(context.Context, *protobuf.Empty) (*protobuf.Response, error)
	DiskCleanUp(context.Context, *protobuf.Empty) (*protobuf.Response, error)
	mustEmbedUnimplementedSettingServer()
}

// UnimplementedSettingServer must be embedded to have forward compatible implementations.
type UnimplementedSettingServer struct {
}

func (UnimplementedSettingServer) CacheClear(context.Context, *protobuf.Empty) (*protobuf.Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CacheClear not implemented")
}
func (UnimplementedSettingServer) DiskCleanUp(context.Context, *protobuf.Empty) (*protobuf.Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DiskCleanUp not implemented")
}
func (UnimplementedSettingServer) mustEmbedUnimplementedSettingServer() {}

// UnsafeSettingServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SettingServer will
// result in compilation errors.
type UnsafeSettingServer interface {
	mustEmbedUnimplementedSettingServer()
}

func RegisterSettingServer(s grpc.ServiceRegistrar, srv SettingServer) {
	s.RegisterService(&_Setting_serviceDesc, srv)
}

func _Setting_CacheClear_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(protobuf.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SettingServer).CacheClear(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protobuf.Setting/CacheClear",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SettingServer).CacheClear(ctx, req.(*protobuf.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Setting_DiskCleanUp_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(protobuf.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SettingServer).DiskCleanUp(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protobuf.Setting/DiskCleanUp",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SettingServer).DiskCleanUp(ctx, req.(*protobuf.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

var _Setting_serviceDesc = grpc.ServiceDesc{
	ServiceName: "protobuf.Setting",
	HandlerType: (*SettingServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CacheClear",
			Handler:    _Setting_CacheClear_Handler,
		},
		{
			MethodName: "DiskCleanUp",
			Handler:    _Setting_DiskCleanUp_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "setting.proto",
}
