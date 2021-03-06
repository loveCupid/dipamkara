// Code generated by protoc-gen-go. DO NOT EDIT.
// source: src/kernal/http.proto

package kernal

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type HttpReq struct {
	ServiceName          string   `protobuf:"bytes,1,opt,name=service_name,json=serviceName,proto3" json:"service_name,omitempty"`
	Method               string   `protobuf:"bytes,2,opt,name=method,proto3" json:"method,omitempty"`
	Body                 string   `protobuf:"bytes,3,opt,name=body,proto3" json:"body,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *HttpReq) Reset()         { *m = HttpReq{} }
func (m *HttpReq) String() string { return proto.CompactTextString(m) }
func (*HttpReq) ProtoMessage()    {}
func (*HttpReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_6fcbda49aa2b1611, []int{0}
}

func (m *HttpReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_HttpReq.Unmarshal(m, b)
}
func (m *HttpReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_HttpReq.Marshal(b, m, deterministic)
}
func (m *HttpReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HttpReq.Merge(m, src)
}
func (m *HttpReq) XXX_Size() int {
	return xxx_messageInfo_HttpReq.Size(m)
}
func (m *HttpReq) XXX_DiscardUnknown() {
	xxx_messageInfo_HttpReq.DiscardUnknown(m)
}

var xxx_messageInfo_HttpReq proto.InternalMessageInfo

func (m *HttpReq) GetServiceName() string {
	if m != nil {
		return m.ServiceName
	}
	return ""
}

func (m *HttpReq) GetMethod() string {
	if m != nil {
		return m.Method
	}
	return ""
}

func (m *HttpReq) GetBody() string {
	if m != nil {
		return m.Body
	}
	return ""
}

type HttpRsp struct {
	Code                 int32    `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Msg                  string   `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
	Reply                string   `protobuf:"bytes,3,opt,name=reply,proto3" json:"reply,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *HttpRsp) Reset()         { *m = HttpRsp{} }
func (m *HttpRsp) String() string { return proto.CompactTextString(m) }
func (*HttpRsp) ProtoMessage()    {}
func (*HttpRsp) Descriptor() ([]byte, []int) {
	return fileDescriptor_6fcbda49aa2b1611, []int{1}
}

func (m *HttpRsp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_HttpRsp.Unmarshal(m, b)
}
func (m *HttpRsp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_HttpRsp.Marshal(b, m, deterministic)
}
func (m *HttpRsp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HttpRsp.Merge(m, src)
}
func (m *HttpRsp) XXX_Size() int {
	return xxx_messageInfo_HttpRsp.Size(m)
}
func (m *HttpRsp) XXX_DiscardUnknown() {
	xxx_messageInfo_HttpRsp.DiscardUnknown(m)
}

var xxx_messageInfo_HttpRsp proto.InternalMessageInfo

func (m *HttpRsp) GetCode() int32 {
	if m != nil {
		return m.Code
	}
	return 0
}

func (m *HttpRsp) GetMsg() string {
	if m != nil {
		return m.Msg
	}
	return ""
}

func (m *HttpRsp) GetReply() string {
	if m != nil {
		return m.Reply
	}
	return ""
}

func init() {
	proto.RegisterType((*HttpReq)(nil), "kernal.HttpReq")
	proto.RegisterType((*HttpRsp)(nil), "kernal.HttpRsp")
}

func init() { proto.RegisterFile("src/kernal/http.proto", fileDescriptor_6fcbda49aa2b1611) }

var fileDescriptor_6fcbda49aa2b1611 = []byte{
	// 197 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x2d, 0x2e, 0x4a, 0xd6,
	0xcf, 0x4e, 0x2d, 0xca, 0x4b, 0xcc, 0xd1, 0xcf, 0x28, 0x29, 0x29, 0xd0, 0x2b, 0x28, 0xca, 0x2f,
	0xc9, 0x17, 0x62, 0x83, 0x08, 0x29, 0x45, 0x70, 0xb1, 0x7b, 0x94, 0x94, 0x14, 0x04, 0xa5, 0x16,
	0x0a, 0x29, 0x72, 0xf1, 0x14, 0xa7, 0x16, 0x95, 0x65, 0x26, 0xa7, 0xc6, 0xe7, 0x25, 0xe6, 0xa6,
	0x4a, 0x30, 0x2a, 0x30, 0x6a, 0x70, 0x06, 0x71, 0x43, 0xc5, 0xfc, 0x12, 0x73, 0x53, 0x85, 0xc4,
	0xb8, 0xd8, 0x72, 0x53, 0x4b, 0x32, 0xf2, 0x53, 0x24, 0x98, 0xc0, 0x92, 0x50, 0x9e, 0x90, 0x10,
	0x17, 0x4b, 0x52, 0x7e, 0x4a, 0xa5, 0x04, 0x33, 0x58, 0x14, 0xcc, 0x56, 0x72, 0x85, 0x9a, 0x5c,
	0x5c, 0x00, 0x92, 0x4e, 0xce, 0x4f, 0x81, 0x98, 0xc8, 0x1a, 0x04, 0x66, 0x0b, 0x09, 0x70, 0x31,
	0xe7, 0x16, 0xa7, 0x43, 0xcd, 0x01, 0x31, 0x85, 0x44, 0xb8, 0x58, 0x8b, 0x52, 0x0b, 0x72, 0x60,
	0xa6, 0x40, 0x38, 0x46, 0xe6, 0x5c, 0xdc, 0x20, 0x63, 0x82, 0x21, 0xae, 0x10, 0xd2, 0xe0, 0x62,
	0x71, 0x4e, 0xcc, 0xc9, 0x11, 0xe2, 0xd7, 0x83, 0x78, 0x40, 0x0f, 0xea, 0x7a, 0x29, 0x54, 0x81,
	0xe2, 0x82, 0x24, 0x36, 0xb0, 0x47, 0x8d, 0x01, 0x01, 0x00, 0x00, 0xff, 0xff, 0x1a, 0x44, 0x37,
	0xf2, 0x01, 0x01, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// HttpServiceClient is the client API for HttpService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type HttpServiceClient interface {
	Call(ctx context.Context, in *HttpReq, opts ...grpc.CallOption) (*HttpRsp, error)
}

type httpServiceClient struct {
	cc *grpc.ClientConn
}

func NewHttpServiceClient(cc *grpc.ClientConn) HttpServiceClient {
	return &httpServiceClient{cc}
}

func (c *httpServiceClient) Call(ctx context.Context, in *HttpReq, opts ...grpc.CallOption) (*HttpRsp, error) {
	out := new(HttpRsp)
	err := c.cc.Invoke(ctx, "/kernal.HttpService/Call", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// HttpServiceServer is the server API for HttpService service.
type HttpServiceServer interface {
	Call(context.Context, *HttpReq) (*HttpRsp, error)
}

func RegisterHttpServiceServer(s *grpc.Server, srv HttpServiceServer) {
	s.RegisterService(&_HttpService_serviceDesc, srv)
}

func _HttpService_Call_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HttpReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HttpServiceServer).Call(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/kernal.HttpService/Call",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HttpServiceServer).Call(ctx, req.(*HttpReq))
	}
	return interceptor(ctx, in, info, handler)
}

var _HttpService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "kernal.HttpService",
	HandlerType: (*HttpServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Call",
			Handler:    _HttpService_Call_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "src/kernal/http.proto",
}
