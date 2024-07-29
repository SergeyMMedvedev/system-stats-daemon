// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.4.0
// - protoc             v5.27.1
// source: system_stats_service.proto

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.62.0 or later.
const _ = grpc.SupportPackageIsVersion8

const (
	SystemStatsService_StreamSystemStats_FullMethodName = "/system_stats.SystemStatsService/StreamSystemStats"
)

// SystemStatsServiceClient is the client API for SystemStatsService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SystemStatsServiceClient interface {
	StreamSystemStats(ctx context.Context, in *SystemStatsRequest, opts ...grpc.CallOption) (SystemStatsService_StreamSystemStatsClient, error)
}

type systemStatsServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewSystemStatsServiceClient(cc grpc.ClientConnInterface) SystemStatsServiceClient {
	return &systemStatsServiceClient{cc}
}

func (c *systemStatsServiceClient) StreamSystemStats(ctx context.Context, in *SystemStatsRequest, opts ...grpc.CallOption) (SystemStatsService_StreamSystemStatsClient, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &SystemStatsService_ServiceDesc.Streams[0], SystemStatsService_StreamSystemStats_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &systemStatsServiceStreamSystemStatsClient{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type SystemStatsService_StreamSystemStatsClient interface {
	Recv() (*SystemStatsResponse, error)
	grpc.ClientStream
}

type systemStatsServiceStreamSystemStatsClient struct {
	grpc.ClientStream
}

func (x *systemStatsServiceStreamSystemStatsClient) Recv() (*SystemStatsResponse, error) {
	m := new(SystemStatsResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// SystemStatsServiceServer is the server API for SystemStatsService service.
// All implementations must embed UnimplementedSystemStatsServiceServer
// for forward compatibility
type SystemStatsServiceServer interface {
	StreamSystemStats(*SystemStatsRequest, SystemStatsService_StreamSystemStatsServer) error
	mustEmbedUnimplementedSystemStatsServiceServer()
}

// UnimplementedSystemStatsServiceServer must be embedded to have forward compatible implementations.
type UnimplementedSystemStatsServiceServer struct {
}

func (UnimplementedSystemStatsServiceServer) StreamSystemStats(*SystemStatsRequest, SystemStatsService_StreamSystemStatsServer) error {
	return status.Errorf(codes.Unimplemented, "method StreamSystemStats not implemented")
}
func (UnimplementedSystemStatsServiceServer) mustEmbedUnimplementedSystemStatsServiceServer() {}

// UnsafeSystemStatsServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SystemStatsServiceServer will
// result in compilation errors.
type UnsafeSystemStatsServiceServer interface {
	mustEmbedUnimplementedSystemStatsServiceServer()
}

func RegisterSystemStatsServiceServer(s grpc.ServiceRegistrar, srv SystemStatsServiceServer) {
	s.RegisterService(&SystemStatsService_ServiceDesc, srv)
}

func _SystemStatsService_StreamSystemStats_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(SystemStatsRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(SystemStatsServiceServer).StreamSystemStats(m, &systemStatsServiceStreamSystemStatsServer{ServerStream: stream})
}

type SystemStatsService_StreamSystemStatsServer interface {
	Send(*SystemStatsResponse) error
	grpc.ServerStream
}

type systemStatsServiceStreamSystemStatsServer struct {
	grpc.ServerStream
}

func (x *systemStatsServiceStreamSystemStatsServer) Send(m *SystemStatsResponse) error {
	return x.ServerStream.SendMsg(m)
}

// SystemStatsService_ServiceDesc is the grpc.ServiceDesc for SystemStatsService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SystemStatsService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "system_stats.SystemStatsService",
	HandlerType: (*SystemStatsServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "StreamSystemStats",
			Handler:       _SystemStatsService_StreamSystemStats_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "system_stats_service.proto",
}