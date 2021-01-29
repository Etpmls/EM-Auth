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

// PermissionClient is the client API for Permission service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PermissionClient interface {
	GetAll(ctx context.Context, in *protobuf.Pagination, opts ...grpc.CallOption) (*protobuf.Response, error)
	Create(ctx context.Context, in *PermissionCreate, opts ...grpc.CallOption) (*protobuf.Response, error)
	Edit(ctx context.Context, in *PermissionEdit, opts ...grpc.CallOption) (*protobuf.Response, error)
	Delete(ctx context.Context, in *PermissionDelete, opts ...grpc.CallOption) (*protobuf.Response, error)
	GetAdvancedVerify(ctx context.Context, in *protobuf.Empty, opts ...grpc.CallOption) (*protobuf.Response, error)
}

type permissionClient struct {
	cc grpc.ClientConnInterface
}

func NewPermissionClient(cc grpc.ClientConnInterface) PermissionClient {
	return &permissionClient{cc}
}

func (c *permissionClient) GetAll(ctx context.Context, in *protobuf.Pagination, opts ...grpc.CallOption) (*protobuf.Response, error) {
	out := new(protobuf.Response)
	err := c.cc.Invoke(ctx, "/protobuf.Permission/GetAll", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *permissionClient) Create(ctx context.Context, in *PermissionCreate, opts ...grpc.CallOption) (*protobuf.Response, error) {
	out := new(protobuf.Response)
	err := c.cc.Invoke(ctx, "/protobuf.Permission/Create", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *permissionClient) Edit(ctx context.Context, in *PermissionEdit, opts ...grpc.CallOption) (*protobuf.Response, error) {
	out := new(protobuf.Response)
	err := c.cc.Invoke(ctx, "/protobuf.Permission/Edit", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *permissionClient) Delete(ctx context.Context, in *PermissionDelete, opts ...grpc.CallOption) (*protobuf.Response, error) {
	out := new(protobuf.Response)
	err := c.cc.Invoke(ctx, "/protobuf.Permission/Delete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *permissionClient) GetAdvancedVerify(ctx context.Context, in *protobuf.Empty, opts ...grpc.CallOption) (*protobuf.Response, error) {
	out := new(protobuf.Response)
	err := c.cc.Invoke(ctx, "/protobuf.Permission/GetAdvancedVerify", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PermissionServer is the server API for Permission service.
// All implementations must embed UnimplementedPermissionServer
// for forward compatibility
type PermissionServer interface {
	GetAll(context.Context, *protobuf.Pagination) (*protobuf.Response, error)
	Create(context.Context, *PermissionCreate) (*protobuf.Response, error)
	Edit(context.Context, *PermissionEdit) (*protobuf.Response, error)
	Delete(context.Context, *PermissionDelete) (*protobuf.Response, error)
	GetAdvancedVerify(context.Context, *protobuf.Empty) (*protobuf.Response, error)
	mustEmbedUnimplementedPermissionServer()
}

// UnimplementedPermissionServer must be embedded to have forward compatible implementations.
type UnimplementedPermissionServer struct {
}

func (UnimplementedPermissionServer) GetAll(context.Context, *protobuf.Pagination) (*protobuf.Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAll not implemented")
}
func (UnimplementedPermissionServer) Create(context.Context, *PermissionCreate) (*protobuf.Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedPermissionServer) Edit(context.Context, *PermissionEdit) (*protobuf.Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Edit not implemented")
}
func (UnimplementedPermissionServer) Delete(context.Context, *PermissionDelete) (*protobuf.Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedPermissionServer) GetAdvancedVerify(context.Context, *protobuf.Empty) (*protobuf.Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAdvancedVerify not implemented")
}
func (UnimplementedPermissionServer) mustEmbedUnimplementedPermissionServer() {}

// UnsafePermissionServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PermissionServer will
// result in compilation errors.
type UnsafePermissionServer interface {
	mustEmbedUnimplementedPermissionServer()
}

func RegisterPermissionServer(s grpc.ServiceRegistrar, srv PermissionServer) {
	s.RegisterService(&_Permission_serviceDesc, srv)
}

func _Permission_GetAll_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(protobuf.Pagination)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PermissionServer).GetAll(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protobuf.Permission/GetAll",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PermissionServer).GetAll(ctx, req.(*protobuf.Pagination))
	}
	return interceptor(ctx, in, info, handler)
}

func _Permission_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PermissionCreate)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PermissionServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protobuf.Permission/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PermissionServer).Create(ctx, req.(*PermissionCreate))
	}
	return interceptor(ctx, in, info, handler)
}

func _Permission_Edit_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PermissionEdit)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PermissionServer).Edit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protobuf.Permission/Edit",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PermissionServer).Edit(ctx, req.(*PermissionEdit))
	}
	return interceptor(ctx, in, info, handler)
}

func _Permission_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PermissionDelete)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PermissionServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protobuf.Permission/Delete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PermissionServer).Delete(ctx, req.(*PermissionDelete))
	}
	return interceptor(ctx, in, info, handler)
}

func _Permission_GetAdvancedVerify_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(protobuf.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PermissionServer).GetAdvancedVerify(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protobuf.Permission/GetAdvancedVerify",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PermissionServer).GetAdvancedVerify(ctx, req.(*protobuf.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

var _Permission_serviceDesc = grpc.ServiceDesc{
	ServiceName: "protobuf.Permission",
	HandlerType: (*PermissionServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetAll",
			Handler:    _Permission_GetAll_Handler,
		},
		{
			MethodName: "Create",
			Handler:    _Permission_Create_Handler,
		},
		{
			MethodName: "Edit",
			Handler:    _Permission_Edit_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _Permission_Delete_Handler,
		},
		{
			MethodName: "GetAdvancedVerify",
			Handler:    _Permission_GetAdvancedVerify_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "permission.proto",
}
