// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proxy/trojan/config.proto

package trojan

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
	protocol "v2ray.com/core/common/protocol"
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

type Account struct {
	Password             string   `protobuf:"bytes,1,opt,name=password,proto3" json:"password,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Account) Reset()         { *m = Account{} }
func (m *Account) String() string { return proto.CompactTextString(m) }
func (*Account) ProtoMessage()    {}
func (*Account) Descriptor() ([]byte, []int) {
	return fileDescriptor_c2cd7d26d2c3a1c9, []int{0}
}

func (m *Account) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Account.Unmarshal(m, b)
}
func (m *Account) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Account.Marshal(b, m, deterministic)
}
func (m *Account) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Account.Merge(m, src)
}
func (m *Account) XXX_Size() int {
	return xxx_messageInfo_Account.Size(m)
}
func (m *Account) XXX_DiscardUnknown() {
	xxx_messageInfo_Account.DiscardUnknown(m)
}

var xxx_messageInfo_Account proto.InternalMessageInfo

func (m *Account) GetPassword() string {
	if m != nil {
		return m.Password
	}
	return ""
}

type Fallback struct {
	Type                 string   `protobuf:"bytes,1,opt,name=type,proto3" json:"type,omitempty"`
	Dest                 string   `protobuf:"bytes,2,opt,name=dest,proto3" json:"dest,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Fallback) Reset()         { *m = Fallback{} }
func (m *Fallback) String() string { return proto.CompactTextString(m) }
func (*Fallback) ProtoMessage()    {}
func (*Fallback) Descriptor() ([]byte, []int) {
	return fileDescriptor_c2cd7d26d2c3a1c9, []int{1}
}

func (m *Fallback) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Fallback.Unmarshal(m, b)
}
func (m *Fallback) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Fallback.Marshal(b, m, deterministic)
}
func (m *Fallback) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Fallback.Merge(m, src)
}
func (m *Fallback) XXX_Size() int {
	return xxx_messageInfo_Fallback.Size(m)
}
func (m *Fallback) XXX_DiscardUnknown() {
	xxx_messageInfo_Fallback.DiscardUnknown(m)
}

var xxx_messageInfo_Fallback proto.InternalMessageInfo

func (m *Fallback) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *Fallback) GetDest() string {
	if m != nil {
		return m.Dest
	}
	return ""
}

type ClientConfig struct {
	Server               []*protocol.ServerEndpoint `protobuf:"bytes,1,rep,name=server,proto3" json:"server,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                   `json:"-"`
	XXX_unrecognized     []byte                     `json:"-"`
	XXX_sizecache        int32                      `json:"-"`
}

func (m *ClientConfig) Reset()         { *m = ClientConfig{} }
func (m *ClientConfig) String() string { return proto.CompactTextString(m) }
func (*ClientConfig) ProtoMessage()    {}
func (*ClientConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_c2cd7d26d2c3a1c9, []int{2}
}

func (m *ClientConfig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ClientConfig.Unmarshal(m, b)
}
func (m *ClientConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ClientConfig.Marshal(b, m, deterministic)
}
func (m *ClientConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ClientConfig.Merge(m, src)
}
func (m *ClientConfig) XXX_Size() int {
	return xxx_messageInfo_ClientConfig.Size(m)
}
func (m *ClientConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_ClientConfig.DiscardUnknown(m)
}

var xxx_messageInfo_ClientConfig proto.InternalMessageInfo

func (m *ClientConfig) GetServer() []*protocol.ServerEndpoint {
	if m != nil {
		return m.Server
	}
	return nil
}

type ServerConfig struct {
	Users                []*protocol.User `protobuf:"bytes,1,rep,name=users,proto3" json:"users,omitempty"`
	Fallback             *Fallback        `protobuf:"bytes,2,opt,name=fallback,proto3" json:"fallback,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *ServerConfig) Reset()         { *m = ServerConfig{} }
func (m *ServerConfig) String() string { return proto.CompactTextString(m) }
func (*ServerConfig) ProtoMessage()    {}
func (*ServerConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_c2cd7d26d2c3a1c9, []int{3}
}

func (m *ServerConfig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ServerConfig.Unmarshal(m, b)
}
func (m *ServerConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ServerConfig.Marshal(b, m, deterministic)
}
func (m *ServerConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ServerConfig.Merge(m, src)
}
func (m *ServerConfig) XXX_Size() int {
	return xxx_messageInfo_ServerConfig.Size(m)
}
func (m *ServerConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_ServerConfig.DiscardUnknown(m)
}

var xxx_messageInfo_ServerConfig proto.InternalMessageInfo

func (m *ServerConfig) GetUsers() []*protocol.User {
	if m != nil {
		return m.Users
	}
	return nil
}

func (m *ServerConfig) GetFallback() *Fallback {
	if m != nil {
		return m.Fallback
	}
	return nil
}

func init() {
	proto.RegisterType((*Account)(nil), "v2ray.core.proxy.trojan.Account")
	proto.RegisterType((*Fallback)(nil), "v2ray.core.proxy.trojan.Fallback")
	proto.RegisterType((*ClientConfig)(nil), "v2ray.core.proxy.trojan.ClientConfig")
	proto.RegisterType((*ServerConfig)(nil), "v2ray.core.proxy.trojan.ServerConfig")
}

func init() {
	proto.RegisterFile("proxy/trojan/config.proto", fileDescriptor_c2cd7d26d2c3a1c9)
}

var fileDescriptor_c2cd7d26d2c3a1c9 = []byte{
	// 310 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x91, 0xb1, 0x4e, 0xf3, 0x30,
	0x14, 0x85, 0x95, 0xfe, 0x3f, 0xa5, 0xb8, 0x9d, 0xbc, 0xb4, 0xa4, 0x4b, 0x1b, 0x09, 0xa9, 0x62,
	0xb0, 0xa5, 0x20, 0xb1, 0x31, 0xd0, 0x0a, 0xe6, 0xca, 0x40, 0x07, 0x16, 0xe4, 0x3a, 0x2e, 0x0a,
	0x24, 0xbe, 0x96, 0xed, 0x16, 0xf2, 0x00, 0xbc, 0x0c, 0x4f, 0x89, 0x62, 0x27, 0x55, 0x85, 0x54,
	0xb6, 0x9b, 0x7b, 0xbe, 0x73, 0x73, 0x4e, 0x82, 0xce, 0xb5, 0x81, 0xcf, 0x8a, 0x3a, 0x03, 0x6f,
	0x5c, 0x51, 0x01, 0x6a, 0x93, 0xbf, 0x12, 0x6d, 0xc0, 0x01, 0x1e, 0xee, 0x52, 0xc3, 0x2b, 0x22,
	0xc0, 0x48, 0xe2, 0x29, 0x12, 0xa8, 0x38, 0x16, 0x50, 0x96, 0xa0, 0xa8, 0xc7, 0x04, 0x14, 0x74,
	0x6b, 0xa5, 0x09, 0xa6, 0x78, 0xfa, 0x5b, 0xb3, 0xd2, 0xec, 0xa4, 0x79, 0xb1, 0x5a, 0x8a, 0x80,
	0x24, 0x17, 0xe8, 0xf4, 0x56, 0x08, 0xd8, 0x2a, 0x87, 0x63, 0xd4, 0xd3, 0xdc, 0xda, 0x0f, 0x30,
	0xd9, 0x28, 0x9a, 0x44, 0xb3, 0x33, 0xb6, 0x7f, 0x4e, 0x52, 0xd4, 0xbb, 0xe7, 0x45, 0xb1, 0xe6,
	0xe2, 0x1d, 0x63, 0xf4, 0xdf, 0x55, 0x5a, 0x36, 0x8c, 0x9f, 0xeb, 0x5d, 0x26, 0xad, 0x1b, 0x75,
	0xc2, 0xae, 0x9e, 0x13, 0x86, 0x06, 0x8b, 0x22, 0x97, 0xca, 0x2d, 0x7c, 0x11, 0x3c, 0x47, 0xdd,
	0xf0, 0xfe, 0x51, 0x34, 0xf9, 0x37, 0xeb, 0xa7, 0x97, 0xe4, 0xa0, 0x53, 0x48, 0x4a, 0xda, 0xa4,
	0xe4, 0xc1, 0x93, 0x77, 0x2a, 0xd3, 0x90, 0x2b, 0xc7, 0x1a, 0x67, 0xf2, 0x15, 0xa1, 0x41, 0x90,
	0x9a, 0xa3, 0xd7, 0xe8, 0xa4, 0x2e, 0x6c, 0x9b, 0x9b, 0x93, 0xbf, 0x6e, 0x3e, 0x59, 0x69, 0x58,
	0xc0, 0xf1, 0x0d, 0xea, 0x6d, 0x9a, 0x42, 0x3e, 0x74, 0x3f, 0x9d, 0x92, 0x23, 0x9f, 0x98, 0xb4,
	0xcd, 0xd9, 0xde, 0x32, 0x5f, 0xa1, 0xb1, 0x80, 0xf2, 0x98, 0x63, 0x19, 0x3d, 0x8f, 0x5b, 0xa9,
	0xa4, 0xb5, 0x4c, 0x0f, 0xff, 0xec, 0x77, 0x67, 0xb8, 0x4a, 0x19, 0xaf, 0xc8, 0xa2, 0x36, 0x2e,
	0xbd, 0xf1, 0xd1, 0x2b, 0xeb, 0xae, 0x0f, 0x7b, 0xf5, 0x13, 0x00, 0x00, 0xff, 0xff, 0x7e, 0x15,
	0xd5, 0xac, 0x0a, 0x02, 0x00, 0x00,
}
