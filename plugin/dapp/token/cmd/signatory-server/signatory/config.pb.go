// Code generated by protoc-gen-go. DO NOT EDIT.
// source: config.proto

/*
Package signatory is a generated protocol buffer package.

It is generated from these files:
	config.proto

It has these top-level messages:
	Config
*/
package signatory

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// Config token相关cmd配置
type Config struct {
	Whitelist    []string `protobuf:"bytes,1,rep,name=whitelist" json:"whitelist,omitempty"`
	JrpcBindAddr string   `protobuf:"bytes,2,opt,name=jrpcBindAddr" json:"jrpcBindAddr,omitempty"`
	Privkey      string   `protobuf:"bytes,3,opt,name=privkey" json:"privkey,omitempty"`
}

func (m *Config) Reset()                    { *m = Config{} }
func (m *Config) String() string            { return proto.CompactTextString(m) }
func (*Config) ProtoMessage()               {}
func (*Config) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Config) GetWhitelist() []string {
	if m != nil {
		return m.Whitelist
	}
	return nil
}

func (m *Config) GetJrpcBindAddr() string {
	if m != nil {
		return m.JrpcBindAddr
	}
	return ""
}

func (m *Config) GetPrivkey() string {
	if m != nil {
		return m.Privkey
	}
	return ""
}

func init() {
	proto.RegisterType((*Config)(nil), "signatory.Config")
}

func init() { proto.RegisterFile("config.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 125 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x49, 0xce, 0xcf, 0x4b,
	0xcb, 0x4c, 0xd7, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x2c, 0xce, 0x4c, 0xcf, 0x4b, 0x2c,
	0xc9, 0x2f, 0xaa, 0x54, 0x4a, 0xe1, 0x62, 0x73, 0x06, 0x4b, 0x09, 0xc9, 0x70, 0x71, 0x96, 0x67,
	0x64, 0x96, 0xa4, 0xe6, 0x64, 0x16, 0x97, 0x48, 0x30, 0x2a, 0x30, 0x6b, 0x70, 0x06, 0x21, 0x04,
	0x84, 0x94, 0xb8, 0x78, 0xb2, 0x8a, 0x0a, 0x92, 0x9d, 0x32, 0xf3, 0x52, 0x1c, 0x53, 0x52, 0x8a,
	0x24, 0x98, 0x14, 0x18, 0x35, 0x38, 0x83, 0x50, 0xc4, 0x84, 0x24, 0xb8, 0xd8, 0x0b, 0x8a, 0x32,
	0xcb, 0xb2, 0x53, 0x2b, 0x25, 0x98, 0xc1, 0xd2, 0x30, 0x6e, 0x12, 0x1b, 0xd8, 0x5e, 0x63, 0x40,
	0x00, 0x00, 0x00, 0xff, 0xff, 0xd0, 0x7d, 0xfe, 0x99, 0x87, 0x00, 0x00, 0x00,
}
