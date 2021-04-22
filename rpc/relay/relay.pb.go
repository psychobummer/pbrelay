// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.14.0
// source: rpc/relay/relay.proto

package relay

import (
	context "context"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type StreamMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id        string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Data      []byte `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
	Signature []byte `protobuf:"bytes,3,opt,name=signature,proto3" json:"signature,omitempty"`
}

func (x *StreamMessage) Reset() {
	*x = StreamMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_relay_relay_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StreamMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StreamMessage) ProtoMessage() {}

func (x *StreamMessage) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_relay_relay_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StreamMessage.ProtoReflect.Descriptor instead.
func (*StreamMessage) Descriptor() ([]byte, []int) {
	return file_rpc_relay_relay_proto_rawDescGZIP(), []int{0}
}

func (x *StreamMessage) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *StreamMessage) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *StreamMessage) GetSignature() []byte {
	if x != nil {
		return x.Signature
	}
	return nil
}

type GetStreamRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *GetStreamRequest) Reset() {
	*x = GetStreamRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_relay_relay_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetStreamRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetStreamRequest) ProtoMessage() {}

func (x *GetStreamRequest) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_relay_relay_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetStreamRequest.ProtoReflect.Descriptor instead.
func (*GetStreamRequest) Descriptor() ([]byte, []int) {
	return file_rpc_relay_relay_proto_rawDescGZIP(), []int{1}
}

func (x *GetStreamRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type Empty struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Empty) Reset() {
	*x = Empty{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_relay_relay_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Empty) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Empty) ProtoMessage() {}

func (x *Empty) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_relay_relay_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Empty.ProtoReflect.Descriptor instead.
func (*Empty) Descriptor() ([]byte, []int) {
	return file_rpc_relay_relay_proto_rawDescGZIP(), []int{2}
}

var File_rpc_relay_relay_proto protoreflect.FileDescriptor

var file_rpc_relay_relay_proto_rawDesc = []byte{
	0x0a, 0x15, 0x72, 0x70, 0x63, 0x2f, 0x72, 0x65, 0x6c, 0x61, 0x79, 0x2f, 0x72, 0x65, 0x6c, 0x61,
	0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x72, 0x65, 0x6c, 0x61, 0x79, 0x22, 0x51,
	0x0a, 0x0d, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12,
	0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12,
	0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x64,
	0x61, 0x74, 0x61, 0x12, 0x1c, 0x0a, 0x09, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72,
	0x65, 0x22, 0x22, 0x0a, 0x10, 0x47, 0x65, 0x74, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x02, 0x69, 0x64, 0x22, 0x07, 0x0a, 0x05, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x32, 0x86,
	0x01, 0x0a, 0x0c, 0x52, 0x65, 0x6c, 0x61, 0x79, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12,
	0x36, 0x0a, 0x0c, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x12,
	0x14, 0x2e, 0x72, 0x65, 0x6c, 0x61, 0x79, 0x2e, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x1a, 0x0c, 0x2e, 0x72, 0x65, 0x6c, 0x61, 0x79, 0x2e, 0x45, 0x6d,
	0x70, 0x74, 0x79, 0x22, 0x00, 0x28, 0x01, 0x12, 0x3e, 0x0a, 0x09, 0x47, 0x65, 0x74, 0x53, 0x74,
	0x72, 0x65, 0x61, 0x6d, 0x12, 0x17, 0x2e, 0x72, 0x65, 0x6c, 0x61, 0x79, 0x2e, 0x47, 0x65, 0x74,
	0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x14, 0x2e,
	0x72, 0x65, 0x6c, 0x61, 0x79, 0x2e, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x22, 0x00, 0x30, 0x01, 0x42, 0x0b, 0x5a, 0x09, 0x72, 0x70, 0x63, 0x2f, 0x72,
	0x65, 0x6c, 0x61, 0x79, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_rpc_relay_relay_proto_rawDescOnce sync.Once
	file_rpc_relay_relay_proto_rawDescData = file_rpc_relay_relay_proto_rawDesc
)

func file_rpc_relay_relay_proto_rawDescGZIP() []byte {
	file_rpc_relay_relay_proto_rawDescOnce.Do(func() {
		file_rpc_relay_relay_proto_rawDescData = protoimpl.X.CompressGZIP(file_rpc_relay_relay_proto_rawDescData)
	})
	return file_rpc_relay_relay_proto_rawDescData
}

var file_rpc_relay_relay_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_rpc_relay_relay_proto_goTypes = []interface{}{
	(*StreamMessage)(nil),    // 0: relay.StreamMessage
	(*GetStreamRequest)(nil), // 1: relay.GetStreamRequest
	(*Empty)(nil),            // 2: relay.Empty
}
var file_rpc_relay_relay_proto_depIdxs = []int32{
	0, // 0: relay.RelayService.CreateStream:input_type -> relay.StreamMessage
	1, // 1: relay.RelayService.GetStream:input_type -> relay.GetStreamRequest
	2, // 2: relay.RelayService.CreateStream:output_type -> relay.Empty
	0, // 3: relay.RelayService.GetStream:output_type -> relay.StreamMessage
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_rpc_relay_relay_proto_init() }
func file_rpc_relay_relay_proto_init() {
	if File_rpc_relay_relay_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_rpc_relay_relay_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StreamMessage); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_rpc_relay_relay_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetStreamRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_rpc_relay_relay_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Empty); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_rpc_relay_relay_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_rpc_relay_relay_proto_goTypes,
		DependencyIndexes: file_rpc_relay_relay_proto_depIdxs,
		MessageInfos:      file_rpc_relay_relay_proto_msgTypes,
	}.Build()
	File_rpc_relay_relay_proto = out.File
	file_rpc_relay_relay_proto_rawDesc = nil
	file_rpc_relay_relay_proto_goTypes = nil
	file_rpc_relay_relay_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// RelayServiceClient is the client API for RelayService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type RelayServiceClient interface {
	CreateStream(ctx context.Context, opts ...grpc.CallOption) (RelayService_CreateStreamClient, error)
	GetStream(ctx context.Context, in *GetStreamRequest, opts ...grpc.CallOption) (RelayService_GetStreamClient, error)
}

type relayServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewRelayServiceClient(cc grpc.ClientConnInterface) RelayServiceClient {
	return &relayServiceClient{cc}
}

func (c *relayServiceClient) CreateStream(ctx context.Context, opts ...grpc.CallOption) (RelayService_CreateStreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &_RelayService_serviceDesc.Streams[0], "/relay.RelayService/CreateStream", opts...)
	if err != nil {
		return nil, err
	}
	x := &relayServiceCreateStreamClient{stream}
	return x, nil
}

type RelayService_CreateStreamClient interface {
	Send(*StreamMessage) error
	CloseAndRecv() (*Empty, error)
	grpc.ClientStream
}

type relayServiceCreateStreamClient struct {
	grpc.ClientStream
}

func (x *relayServiceCreateStreamClient) Send(m *StreamMessage) error {
	return x.ClientStream.SendMsg(m)
}

func (x *relayServiceCreateStreamClient) CloseAndRecv() (*Empty, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(Empty)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *relayServiceClient) GetStream(ctx context.Context, in *GetStreamRequest, opts ...grpc.CallOption) (RelayService_GetStreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &_RelayService_serviceDesc.Streams[1], "/relay.RelayService/GetStream", opts...)
	if err != nil {
		return nil, err
	}
	x := &relayServiceGetStreamClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type RelayService_GetStreamClient interface {
	Recv() (*StreamMessage, error)
	grpc.ClientStream
}

type relayServiceGetStreamClient struct {
	grpc.ClientStream
}

func (x *relayServiceGetStreamClient) Recv() (*StreamMessage, error) {
	m := new(StreamMessage)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// RelayServiceServer is the server API for RelayService service.
type RelayServiceServer interface {
	CreateStream(RelayService_CreateStreamServer) error
	GetStream(*GetStreamRequest, RelayService_GetStreamServer) error
}

// UnimplementedRelayServiceServer can be embedded to have forward compatible implementations.
type UnimplementedRelayServiceServer struct {
}

func (*UnimplementedRelayServiceServer) CreateStream(RelayService_CreateStreamServer) error {
	return status.Errorf(codes.Unimplemented, "method CreateStream not implemented")
}
func (*UnimplementedRelayServiceServer) GetStream(*GetStreamRequest, RelayService_GetStreamServer) error {
	return status.Errorf(codes.Unimplemented, "method GetStream not implemented")
}

func RegisterRelayServiceServer(s *grpc.Server, srv RelayServiceServer) {
	s.RegisterService(&_RelayService_serviceDesc, srv)
}

func _RelayService_CreateStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(RelayServiceServer).CreateStream(&relayServiceCreateStreamServer{stream})
}

type RelayService_CreateStreamServer interface {
	SendAndClose(*Empty) error
	Recv() (*StreamMessage, error)
	grpc.ServerStream
}

type relayServiceCreateStreamServer struct {
	grpc.ServerStream
}

func (x *relayServiceCreateStreamServer) SendAndClose(m *Empty) error {
	return x.ServerStream.SendMsg(m)
}

func (x *relayServiceCreateStreamServer) Recv() (*StreamMessage, error) {
	m := new(StreamMessage)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _RelayService_GetStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(GetStreamRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(RelayServiceServer).GetStream(m, &relayServiceGetStreamServer{stream})
}

type RelayService_GetStreamServer interface {
	Send(*StreamMessage) error
	grpc.ServerStream
}

type relayServiceGetStreamServer struct {
	grpc.ServerStream
}

func (x *relayServiceGetStreamServer) Send(m *StreamMessage) error {
	return x.ServerStream.SendMsg(m)
}

var _RelayService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "relay.RelayService",
	HandlerType: (*RelayServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "CreateStream",
			Handler:       _RelayService_CreateStream_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "GetStream",
			Handler:       _RelayService_GetStream_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "rpc/relay/relay.proto",
}
