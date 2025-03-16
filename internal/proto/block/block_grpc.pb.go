// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.3
// source: block.proto

package block

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	BlockService_FetchChunk_FullMethodName  = "/block.BlockService/FetchChunk"
	BlockService_UpdateBlock_FullMethodName = "/block.BlockService/UpdateBlock"
	BlockService_StreamChunk_FullMethodName = "/block.BlockService/StreamChunk"
)

// BlockServiceClient is the client API for BlockService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BlockServiceClient interface {
	FetchChunk(ctx context.Context, in *FetchChunkRequest, opts ...grpc.CallOption) (*FetchChunkResponse, error)
	UpdateBlock(ctx context.Context, in *UpdateBlockRequest, opts ...grpc.CallOption) (*UpdateBlockResponse, error)
	StreamChunk(ctx context.Context, in *ChunkRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[ChunkUpdate], error)
}

type blockServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewBlockServiceClient(cc grpc.ClientConnInterface) BlockServiceClient {
	return &blockServiceClient{cc}
}

func (c *blockServiceClient) FetchChunk(ctx context.Context, in *FetchChunkRequest, opts ...grpc.CallOption) (*FetchChunkResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(FetchChunkResponse)
	err := c.cc.Invoke(ctx, BlockService_FetchChunk_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *blockServiceClient) UpdateBlock(ctx context.Context, in *UpdateBlockRequest, opts ...grpc.CallOption) (*UpdateBlockResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdateBlockResponse)
	err := c.cc.Invoke(ctx, BlockService_UpdateBlock_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *blockServiceClient) StreamChunk(ctx context.Context, in *ChunkRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[ChunkUpdate], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &BlockService_ServiceDesc.Streams[0], BlockService_StreamChunk_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[ChunkRequest, ChunkUpdate]{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type BlockService_StreamChunkClient = grpc.ServerStreamingClient[ChunkUpdate]

// BlockServiceServer is the server API for BlockService service.
// All implementations must embed UnimplementedBlockServiceServer
// for forward compatibility.
type BlockServiceServer interface {
	FetchChunk(context.Context, *FetchChunkRequest) (*FetchChunkResponse, error)
	UpdateBlock(context.Context, *UpdateBlockRequest) (*UpdateBlockResponse, error)
	StreamChunk(*ChunkRequest, grpc.ServerStreamingServer[ChunkUpdate]) error
	mustEmbedUnimplementedBlockServiceServer()
}

// UnimplementedBlockServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedBlockServiceServer struct{}

func (UnimplementedBlockServiceServer) FetchChunk(context.Context, *FetchChunkRequest) (*FetchChunkResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FetchChunk not implemented")
}
func (UnimplementedBlockServiceServer) UpdateBlock(context.Context, *UpdateBlockRequest) (*UpdateBlockResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateBlock not implemented")
}
func (UnimplementedBlockServiceServer) StreamChunk(*ChunkRequest, grpc.ServerStreamingServer[ChunkUpdate]) error {
	return status.Errorf(codes.Unimplemented, "method StreamChunk not implemented")
}
func (UnimplementedBlockServiceServer) mustEmbedUnimplementedBlockServiceServer() {}
func (UnimplementedBlockServiceServer) testEmbeddedByValue()                      {}

// UnsafeBlockServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BlockServiceServer will
// result in compilation errors.
type UnsafeBlockServiceServer interface {
	mustEmbedUnimplementedBlockServiceServer()
}

func RegisterBlockServiceServer(s grpc.ServiceRegistrar, srv BlockServiceServer) {
	// If the following call pancis, it indicates UnimplementedBlockServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&BlockService_ServiceDesc, srv)
}

func _BlockService_FetchChunk_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FetchChunkRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BlockServiceServer).FetchChunk(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BlockService_FetchChunk_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BlockServiceServer).FetchChunk(ctx, req.(*FetchChunkRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BlockService_UpdateBlock_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateBlockRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BlockServiceServer).UpdateBlock(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BlockService_UpdateBlock_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BlockServiceServer).UpdateBlock(ctx, req.(*UpdateBlockRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BlockService_StreamChunk_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ChunkRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(BlockServiceServer).StreamChunk(m, &grpc.GenericServerStream[ChunkRequest, ChunkUpdate]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type BlockService_StreamChunkServer = grpc.ServerStreamingServer[ChunkUpdate]

// BlockService_ServiceDesc is the grpc.ServiceDesc for BlockService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var BlockService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "block.BlockService",
	HandlerType: (*BlockServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "FetchChunk",
			Handler:    _BlockService_FetchChunk_Handler,
		},
		{
			MethodName: "UpdateBlock",
			Handler:    _BlockService_UpdateBlock_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "StreamChunk",
			Handler:       _BlockService_StreamChunk_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "block.proto",
}
