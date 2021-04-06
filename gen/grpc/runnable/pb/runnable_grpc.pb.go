// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package runnablepb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// RunnableClient is the client API for Runnable service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RunnableClient interface {
	// Retrieve information about runnables registered in FuseML.
	List(ctx context.Context, in *ListRequest, opts ...grpc.CallOption) (*ListResponse, error)
	// Register a runnable with the FuseML runnable store.
	Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterResponse, error)
	// Retrieve an Runnable from FuseML.
	Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetResponse, error)
}

type runnableClient struct {
	cc grpc.ClientConnInterface
}

func NewRunnableClient(cc grpc.ClientConnInterface) RunnableClient {
	return &runnableClient{cc}
}

func (c *runnableClient) List(ctx context.Context, in *ListRequest, opts ...grpc.CallOption) (*ListResponse, error) {
	out := new(ListResponse)
	err := c.cc.Invoke(ctx, "/runnable.Runnable/List", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *runnableClient) Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterResponse, error) {
	out := new(RegisterResponse)
	err := c.cc.Invoke(ctx, "/runnable.Runnable/Register", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *runnableClient) Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetResponse, error) {
	out := new(GetResponse)
	err := c.cc.Invoke(ctx, "/runnable.Runnable/Get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RunnableServer is the server API for Runnable service.
// All implementations must embed UnimplementedRunnableServer
// for forward compatibility
type RunnableServer interface {
	// Retrieve information about runnables registered in FuseML.
	List(context.Context, *ListRequest) (*ListResponse, error)
	// Register a runnable with the FuseML runnable store.
	Register(context.Context, *RegisterRequest) (*RegisterResponse, error)
	// Retrieve an Runnable from FuseML.
	Get(context.Context, *GetRequest) (*GetResponse, error)
	mustEmbedUnimplementedRunnableServer()
}

// UnimplementedRunnableServer must be embedded to have forward compatible implementations.
type UnimplementedRunnableServer struct {
}

func (UnimplementedRunnableServer) List(context.Context, *ListRequest) (*ListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method List not implemented")
}
func (UnimplementedRunnableServer) Register(context.Context, *RegisterRequest) (*RegisterResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Register not implemented")
}
func (UnimplementedRunnableServer) Get(context.Context, *GetRequest) (*GetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedRunnableServer) mustEmbedUnimplementedRunnableServer() {}

// UnsafeRunnableServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RunnableServer will
// result in compilation errors.
type UnsafeRunnableServer interface {
	mustEmbedUnimplementedRunnableServer()
}

func RegisterRunnableServer(s grpc.ServiceRegistrar, srv RunnableServer) {
	s.RegisterService(&Runnable_ServiceDesc, srv)
}

func _Runnable_List_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RunnableServer).List(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/runnable.Runnable/List",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RunnableServer).List(ctx, req.(*ListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Runnable_Register_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RunnableServer).Register(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/runnable.Runnable/Register",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RunnableServer).Register(ctx, req.(*RegisterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Runnable_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RunnableServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/runnable.Runnable/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RunnableServer).Get(ctx, req.(*GetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Runnable_ServiceDesc is the grpc.ServiceDesc for Runnable service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Runnable_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "runnable.Runnable",
	HandlerType: (*RunnableServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "List",
			Handler:    _Runnable_List_Handler,
		},
		{
			MethodName: "Register",
			Handler:    _Runnable_Register_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _Runnable_Get_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "runnable.proto",
}
