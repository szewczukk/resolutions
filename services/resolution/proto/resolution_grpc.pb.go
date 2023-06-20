// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.23.3
// source: proto/resolution.proto

package proto

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

// ResolutionServiceClient is the client API for ResolutionService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ResolutionServiceClient interface {
	GetAllResolutions(ctx context.Context, in *GetAllResolutionsRequest, opts ...grpc.CallOption) (*RepeatedResolutions, error)
	GetResolutionsByUserId(ctx context.Context, in *UserId, opts ...grpc.CallOption) (*RepeatedResolutions, error)
	CreateResolution(ctx context.Context, in *CreateResolutionRequest, opts ...grpc.CallOption) (*Resolution, error)
}

type resolutionServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewResolutionServiceClient(cc grpc.ClientConnInterface) ResolutionServiceClient {
	return &resolutionServiceClient{cc}
}

func (c *resolutionServiceClient) GetAllResolutions(ctx context.Context, in *GetAllResolutionsRequest, opts ...grpc.CallOption) (*RepeatedResolutions, error) {
	out := new(RepeatedResolutions)
	err := c.cc.Invoke(ctx, "/ResolutionService/GetAllResolutions", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *resolutionServiceClient) GetResolutionsByUserId(ctx context.Context, in *UserId, opts ...grpc.CallOption) (*RepeatedResolutions, error) {
	out := new(RepeatedResolutions)
	err := c.cc.Invoke(ctx, "/ResolutionService/GetResolutionsByUserId", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *resolutionServiceClient) CreateResolution(ctx context.Context, in *CreateResolutionRequest, opts ...grpc.CallOption) (*Resolution, error) {
	out := new(Resolution)
	err := c.cc.Invoke(ctx, "/ResolutionService/CreateResolution", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ResolutionServiceServer is the server API for ResolutionService service.
// All implementations must embed UnimplementedResolutionServiceServer
// for forward compatibility
type ResolutionServiceServer interface {
	GetAllResolutions(context.Context, *GetAllResolutionsRequest) (*RepeatedResolutions, error)
	GetResolutionsByUserId(context.Context, *UserId) (*RepeatedResolutions, error)
	CreateResolution(context.Context, *CreateResolutionRequest) (*Resolution, error)
	mustEmbedUnimplementedResolutionServiceServer()
}

// UnimplementedResolutionServiceServer must be embedded to have forward compatible implementations.
type UnimplementedResolutionServiceServer struct {
}

func (UnimplementedResolutionServiceServer) GetAllResolutions(context.Context, *GetAllResolutionsRequest) (*RepeatedResolutions, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllResolutions not implemented")
}
func (UnimplementedResolutionServiceServer) GetResolutionsByUserId(context.Context, *UserId) (*RepeatedResolutions, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetResolutionsByUserId not implemented")
}
func (UnimplementedResolutionServiceServer) CreateResolution(context.Context, *CreateResolutionRequest) (*Resolution, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateResolution not implemented")
}
func (UnimplementedResolutionServiceServer) mustEmbedUnimplementedResolutionServiceServer() {}

// UnsafeResolutionServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ResolutionServiceServer will
// result in compilation errors.
type UnsafeResolutionServiceServer interface {
	mustEmbedUnimplementedResolutionServiceServer()
}

func RegisterResolutionServiceServer(s grpc.ServiceRegistrar, srv ResolutionServiceServer) {
	s.RegisterService(&ResolutionService_ServiceDesc, srv)
}

func _ResolutionService_GetAllResolutions_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAllResolutionsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ResolutionServiceServer).GetAllResolutions(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ResolutionService/GetAllResolutions",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ResolutionServiceServer).GetAllResolutions(ctx, req.(*GetAllResolutionsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ResolutionService_GetResolutionsByUserId_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserId)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ResolutionServiceServer).GetResolutionsByUserId(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ResolutionService/GetResolutionsByUserId",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ResolutionServiceServer).GetResolutionsByUserId(ctx, req.(*UserId))
	}
	return interceptor(ctx, in, info, handler)
}

func _ResolutionService_CreateResolution_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateResolutionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ResolutionServiceServer).CreateResolution(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ResolutionService/CreateResolution",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ResolutionServiceServer).CreateResolution(ctx, req.(*CreateResolutionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ResolutionService_ServiceDesc is the grpc.ServiceDesc for ResolutionService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ResolutionService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "ResolutionService",
	HandlerType: (*ResolutionServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetAllResolutions",
			Handler:    _ResolutionService_GetAllResolutions_Handler,
		},
		{
			MethodName: "GetResolutionsByUserId",
			Handler:    _ResolutionService_GetResolutionsByUserId_Handler,
		},
		{
			MethodName: "CreateResolution",
			Handler:    _ResolutionService_CreateResolution_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/resolution.proto",
}
