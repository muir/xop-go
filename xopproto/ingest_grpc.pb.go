// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.7
// source: ingest.proto

package xopproto

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

// IngestClient is the client API for Ingest service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type IngestClient interface {
	Ping(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error)
	UploadFragment(ctx context.Context, in *IngestFragment, opts ...grpc.CallOption) (*Error, error)
	PrepareToStream(ctx context.Context, in *SourceIdentity, opts ...grpc.CallOption) (*ReadyToStream, error)
	Stream(ctx context.Context, opts ...grpc.CallOption) (Ingest_StreamClient, error)
}

type ingestClient struct {
	cc grpc.ClientConnInterface
}

func NewIngestClient(cc grpc.ClientConnInterface) IngestClient {
	return &ingestClient{cc}
}

func (c *ingestClient) Ping(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/Ingest/Ping", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ingestClient) UploadFragment(ctx context.Context, in *IngestFragment, opts ...grpc.CallOption) (*Error, error) {
	out := new(Error)
	err := c.cc.Invoke(ctx, "/Ingest/UploadFragment", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ingestClient) PrepareToStream(ctx context.Context, in *SourceIdentity, opts ...grpc.CallOption) (*ReadyToStream, error) {
	out := new(ReadyToStream)
	err := c.cc.Invoke(ctx, "/Ingest/PrepareToStream", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ingestClient) Stream(ctx context.Context, opts ...grpc.CallOption) (Ingest_StreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &Ingest_ServiceDesc.Streams[0], "/Ingest/Stream", opts...)
	if err != nil {
		return nil, err
	}
	x := &ingestStreamClient{stream}
	return x, nil
}

type Ingest_StreamClient interface {
	Send(*FragmentInStream) error
	Recv() (*FragmentAck, error)
	grpc.ClientStream
}

type ingestStreamClient struct {
	grpc.ClientStream
}

func (x *ingestStreamClient) Send(m *FragmentInStream) error {
	return x.ClientStream.SendMsg(m)
}

func (x *ingestStreamClient) Recv() (*FragmentAck, error) {
	m := new(FragmentAck)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// IngestServer is the server API for Ingest service.
// All implementations must embed UnimplementedIngestServer
// for forward compatibility
type IngestServer interface {
	Ping(context.Context, *Empty) (*Empty, error)
	UploadFragment(context.Context, *IngestFragment) (*Error, error)
	PrepareToStream(context.Context, *SourceIdentity) (*ReadyToStream, error)
	Stream(Ingest_StreamServer) error
	mustEmbedUnimplementedIngestServer()
}

// UnimplementedIngestServer must be embedded to have forward compatible implementations.
type UnimplementedIngestServer struct {
}

func (UnimplementedIngestServer) Ping(context.Context, *Empty) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (UnimplementedIngestServer) UploadFragment(context.Context, *IngestFragment) (*Error, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UploadFragment not implemented")
}
func (UnimplementedIngestServer) PrepareToStream(context.Context, *SourceIdentity) (*ReadyToStream, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PrepareToStream not implemented")
}
func (UnimplementedIngestServer) Stream(Ingest_StreamServer) error {
	return status.Errorf(codes.Unimplemented, "method Stream not implemented")
}
func (UnimplementedIngestServer) mustEmbedUnimplementedIngestServer() {}

// UnsafeIngestServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to IngestServer will
// result in compilation errors.
type UnsafeIngestServer interface {
	mustEmbedUnimplementedIngestServer()
}

func RegisterIngestServer(s grpc.ServiceRegistrar, srv IngestServer) {
	s.RegisterService(&Ingest_ServiceDesc, srv)
}

func _Ingest_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IngestServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Ingest/Ping",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IngestServer).Ping(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Ingest_UploadFragment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IngestFragment)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IngestServer).UploadFragment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Ingest/UploadFragment",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IngestServer).UploadFragment(ctx, req.(*IngestFragment))
	}
	return interceptor(ctx, in, info, handler)
}

func _Ingest_PrepareToStream_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SourceIdentity)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IngestServer).PrepareToStream(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Ingest/PrepareToStream",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IngestServer).PrepareToStream(ctx, req.(*SourceIdentity))
	}
	return interceptor(ctx, in, info, handler)
}

func _Ingest_Stream_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(IngestServer).Stream(&ingestStreamServer{stream})
}

type Ingest_StreamServer interface {
	Send(*FragmentAck) error
	Recv() (*FragmentInStream, error)
	grpc.ServerStream
}

type ingestStreamServer struct {
	grpc.ServerStream
}

func (x *ingestStreamServer) Send(m *FragmentAck) error {
	return x.ServerStream.SendMsg(m)
}

func (x *ingestStreamServer) Recv() (*FragmentInStream, error) {
	m := new(FragmentInStream)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Ingest_ServiceDesc is the grpc.ServiceDesc for Ingest service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Ingest_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "Ingest",
	HandlerType: (*IngestServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Ping",
			Handler:    _Ingest_Ping_Handler,
		},
		{
			MethodName: "UploadFragment",
			Handler:    _Ingest_UploadFragment_Handler,
		},
		{
			MethodName: "PrepareToStream",
			Handler:    _Ingest_PrepareToStream_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Stream",
			Handler:       _Ingest_Stream_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "ingest.proto",
}