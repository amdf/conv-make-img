// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.6.1
// source: svc.proto

package svc

import (
	context "context"
	httpbody "google.golang.org/genproto/googleapis/api/httpbody"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// TengwarConverterClient is the client API for TengwarConverter service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TengwarConverterClient interface {
	ConvertText(ctx context.Context, in *SimpleConvertRequest, opts ...grpc.CallOption) (*ConvertResponse, error)
	MakeImage(ctx context.Context, in *ConvertRequest, opts ...grpc.CallOption) (*httpbody.HttpBody, error)
}

type tengwarConverterClient struct {
	cc grpc.ClientConnInterface
}

func NewTengwarConverterClient(cc grpc.ClientConnInterface) TengwarConverterClient {
	return &tengwarConverterClient{cc}
}

func (c *tengwarConverterClient) ConvertText(ctx context.Context, in *SimpleConvertRequest, opts ...grpc.CallOption) (*ConvertResponse, error) {
	out := new(ConvertResponse)
	err := c.cc.Invoke(ctx, "/imgtengwar.TengwarConverter/ConvertText", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tengwarConverterClient) MakeImage(ctx context.Context, in *ConvertRequest, opts ...grpc.CallOption) (*httpbody.HttpBody, error) {
	out := new(httpbody.HttpBody)
	err := c.cc.Invoke(ctx, "/imgtengwar.TengwarConverter/MakeImage", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TengwarConverterServer is the server API for TengwarConverter service.
// All implementations must embed UnimplementedTengwarConverterServer
// for forward compatibility
type TengwarConverterServer interface {
	ConvertText(context.Context, *SimpleConvertRequest) (*ConvertResponse, error)
	MakeImage(context.Context, *ConvertRequest) (*httpbody.HttpBody, error)
	mustEmbedUnimplementedTengwarConverterServer()
}

// UnimplementedTengwarConverterServer must be embedded to have forward compatible implementations.
type UnimplementedTengwarConverterServer struct {
}

func (UnimplementedTengwarConverterServer) ConvertText(context.Context, *SimpleConvertRequest) (*ConvertResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ConvertText not implemented")
}
func (UnimplementedTengwarConverterServer) MakeImage(context.Context, *ConvertRequest) (*httpbody.HttpBody, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MakeImage not implemented")
}
func (UnimplementedTengwarConverterServer) mustEmbedUnimplementedTengwarConverterServer() {}

// UnsafeTengwarConverterServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TengwarConverterServer will
// result in compilation errors.
type UnsafeTengwarConverterServer interface {
	mustEmbedUnimplementedTengwarConverterServer()
}

func RegisterTengwarConverterServer(s grpc.ServiceRegistrar, srv TengwarConverterServer) {
	s.RegisterService(&TengwarConverter_ServiceDesc, srv)
}

func _TengwarConverter_ConvertText_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SimpleConvertRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TengwarConverterServer).ConvertText(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/imgtengwar.TengwarConverter/ConvertText",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TengwarConverterServer).ConvertText(ctx, req.(*SimpleConvertRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TengwarConverter_MakeImage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ConvertRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TengwarConverterServer).MakeImage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/imgtengwar.TengwarConverter/MakeImage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TengwarConverterServer).MakeImage(ctx, req.(*ConvertRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// TengwarConverter_ServiceDesc is the grpc.ServiceDesc for TengwarConverter service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var TengwarConverter_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "imgtengwar.TengwarConverter",
	HandlerType: (*TengwarConverterServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ConvertText",
			Handler:    _TengwarConverter_ConvertText_Handler,
		},
		{
			MethodName: "MakeImage",
			Handler:    _TengwarConverter_MakeImage_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "svc.proto",
}
