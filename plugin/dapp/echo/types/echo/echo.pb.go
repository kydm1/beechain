// Code generated by protoc-gen-go. DO NOT EDIT.
// source: echo.proto

/*
Package echo is a generated protocol buffer package.

It is generated from these files:
	echo.proto

It has these top-level messages:
	Ping
	Pang
	EchoAction
	PingLog
	PangLog
	Query
	QueryResult
*/
package echo

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

// ping操作action
type Ping struct {
	Msg string `protobuf:"bytes,1,opt,name=msg" json:"msg,omitempty"`
}

func (m *Ping) Reset()                    { *m = Ping{} }
func (m *Ping) String() string            { return proto.CompactTextString(m) }
func (*Ping) ProtoMessage()               {}
func (*Ping) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Ping) GetMsg() string {
	if m != nil {
		return m.Msg
	}
	return ""
}

// pang操作action
type Pang struct {
	Msg string `protobuf:"bytes,1,opt,name=msg" json:"msg,omitempty"`
}

func (m *Pang) Reset()                    { *m = Pang{} }
func (m *Pang) String() string            { return proto.CompactTextString(m) }
func (*Pang) ProtoMessage()               {}
func (*Pang) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *Pang) GetMsg() string {
	if m != nil {
		return m.Msg
	}
	return ""
}

// 本执行器的统一Action结构
type EchoAction struct {
	// Types that are valid to be assigned to Value:
	//	*EchoAction_Ping
	//	*EchoAction_Pang
	Value isEchoAction_Value `protobuf_oneof:"value"`
	Ty    int32              `protobuf:"varint,3,opt,name=ty" json:"ty,omitempty"`
}

func (m *EchoAction) Reset()                    { *m = EchoAction{} }
func (m *EchoAction) String() string            { return proto.CompactTextString(m) }
func (*EchoAction) ProtoMessage()               {}
func (*EchoAction) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

type isEchoAction_Value interface {
	isEchoAction_Value()
}

type EchoAction_Ping struct {
	Ping *Ping `protobuf:"bytes,1,opt,name=ping,oneof"`
}
type EchoAction_Pang struct {
	Pang *Pang `protobuf:"bytes,2,opt,name=pang,oneof"`
}

func (*EchoAction_Ping) isEchoAction_Value() {}
func (*EchoAction_Pang) isEchoAction_Value() {}

func (m *EchoAction) GetValue() isEchoAction_Value {
	if m != nil {
		return m.Value
	}
	return nil
}

func (m *EchoAction) GetPing() *Ping {
	if x, ok := m.GetValue().(*EchoAction_Ping); ok {
		return x.Ping
	}
	return nil
}

func (m *EchoAction) GetPang() *Pang {
	if x, ok := m.GetValue().(*EchoAction_Pang); ok {
		return x.Pang
	}
	return nil
}

func (m *EchoAction) GetTy() int32 {
	if m != nil {
		return m.Ty
	}
	return 0
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*EchoAction) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _EchoAction_OneofMarshaler, _EchoAction_OneofUnmarshaler, _EchoAction_OneofSizer, []interface{}{
		(*EchoAction_Ping)(nil),
		(*EchoAction_Pang)(nil),
	}
}

func _EchoAction_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*EchoAction)
	// value
	switch x := m.Value.(type) {
	case *EchoAction_Ping:
		b.EncodeVarint(1<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Ping); err != nil {
			return err
		}
	case *EchoAction_Pang:
		b.EncodeVarint(2<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Pang); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("EchoAction.Value has unexpected type %T", x)
	}
	return nil
}

func _EchoAction_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*EchoAction)
	switch tag {
	case 1: // value.ping
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Ping)
		err := b.DecodeMessage(msg)
		m.Value = &EchoAction_Ping{msg}
		return true, err
	case 2: // value.pang
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Pang)
		err := b.DecodeMessage(msg)
		m.Value = &EchoAction_Pang{msg}
		return true, err
	default:
		return false, nil
	}
}

func _EchoAction_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*EchoAction)
	// value
	switch x := m.Value.(type) {
	case *EchoAction_Ping:
		s := proto.Size(x.Ping)
		n += proto.SizeVarint(1<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *EchoAction_Pang:
		s := proto.Size(x.Pang)
		n += proto.SizeVarint(2<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

// ping操作生成的日志结构
type PingLog struct {
	Msg   string `protobuf:"bytes,1,opt,name=msg" json:"msg,omitempty"`
	Echo  string `protobuf:"bytes,2,opt,name=echo" json:"echo,omitempty"`
	Count int32  `protobuf:"varint,3,opt,name=count" json:"count,omitempty"`
}

func (m *PingLog) Reset()                    { *m = PingLog{} }
func (m *PingLog) String() string            { return proto.CompactTextString(m) }
func (*PingLog) ProtoMessage()               {}
func (*PingLog) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *PingLog) GetMsg() string {
	if m != nil {
		return m.Msg
	}
	return ""
}

func (m *PingLog) GetEcho() string {
	if m != nil {
		return m.Echo
	}
	return ""
}

func (m *PingLog) GetCount() int32 {
	if m != nil {
		return m.Count
	}
	return 0
}

// pang操作生成的日志结构
type PangLog struct {
	Msg   string `protobuf:"bytes,1,opt,name=msg" json:"msg,omitempty"`
	Echo  string `protobuf:"bytes,2,opt,name=echo" json:"echo,omitempty"`
	Count int32  `protobuf:"varint,3,opt,name=count" json:"count,omitempty"`
}

func (m *PangLog) Reset()                    { *m = PangLog{} }
func (m *PangLog) String() string            { return proto.CompactTextString(m) }
func (*PangLog) ProtoMessage()               {}
func (*PangLog) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *PangLog) GetMsg() string {
	if m != nil {
		return m.Msg
	}
	return ""
}

func (m *PangLog) GetEcho() string {
	if m != nil {
		return m.Echo
	}
	return ""
}

func (m *PangLog) GetCount() int32 {
	if m != nil {
		return m.Count
	}
	return 0
}

// 查询请求结构
type Query struct {
	Msg string `protobuf:"bytes,1,opt,name=msg" json:"msg,omitempty"`
}

func (m *Query) Reset()                    { *m = Query{} }
func (m *Query) String() string            { return proto.CompactTextString(m) }
func (*Query) ProtoMessage()               {}
func (*Query) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *Query) GetMsg() string {
	if m != nil {
		return m.Msg
	}
	return ""
}

// 查询结果结构
type QueryResult struct {
	Msg   string `protobuf:"bytes,1,opt,name=msg" json:"msg,omitempty"`
	Count int32  `protobuf:"varint,2,opt,name=count" json:"count,omitempty"`
}

func (m *QueryResult) Reset()                    { *m = QueryResult{} }
func (m *QueryResult) String() string            { return proto.CompactTextString(m) }
func (*QueryResult) ProtoMessage()               {}
func (*QueryResult) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *QueryResult) GetMsg() string {
	if m != nil {
		return m.Msg
	}
	return ""
}

func (m *QueryResult) GetCount() int32 {
	if m != nil {
		return m.Count
	}
	return 0
}

func init() {
	proto.RegisterType((*Ping)(nil), "echo.Ping")
	proto.RegisterType((*Pang)(nil), "echo.Pang")
	proto.RegisterType((*EchoAction)(nil), "echo.EchoAction")
	proto.RegisterType((*PingLog)(nil), "echo.PingLog")
	proto.RegisterType((*PangLog)(nil), "echo.PangLog")
	proto.RegisterType((*Query)(nil), "echo.Query")
	proto.RegisterType((*QueryResult)(nil), "echo.QueryResult")
}

func init() { proto.RegisterFile("echo.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 215 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x91, 0x31, 0x4b, 0xc7, 0x30,
	0x10, 0xc5, 0x6d, 0xda, 0x58, 0x7a, 0x05, 0x91, 0xe0, 0x10, 0xb7, 0xd2, 0xa9, 0x53, 0x07, 0xc5,
	0x0f, 0xa0, 0x50, 0x70, 0x70, 0xd0, 0x7c, 0x83, 0x18, 0x42, 0x1b, 0xa8, 0x49, 0x69, 0x13, 0xa1,
	0xdf, 0x5e, 0x72, 0x2d, 0xa2, 0x90, 0xed, 0xbf, 0xbd, 0xf0, 0x1e, 0xbf, 0x97, 0xbb, 0x03, 0xd0,
	0x6a, 0x72, 0xfd, 0xb2, 0x3a, 0xef, 0x58, 0x11, 0x75, 0xcb, 0xa1, 0x78, 0x37, 0x76, 0x64, 0xb7,
	0x90, 0x7f, 0x6d, 0x23, 0xcf, 0x9a, 0xac, 0xab, 0x44, 0x94, 0xe8, 0xc8, 0xa4, 0x63, 0x00, 0x06,
	0x35, 0xb9, 0x67, 0xe5, 0x8d, 0xb3, 0xac, 0x81, 0x62, 0x31, 0xf6, 0x08, 0xd4, 0x0f, 0xd0, 0x63,
	0x45, 0x64, 0xbe, 0x5e, 0x09, 0x74, 0x30, 0x21, 0xed, 0xc8, 0xc9, 0xbf, 0x84, 0x3c, 0x13, 0xb1,
	0xe3, 0x06, 0x88, 0xdf, 0x79, 0xde, 0x64, 0x1d, 0x15, 0xc4, 0xef, 0x2f, 0x25, 0xd0, 0x6f, 0x39,
	0x07, 0xdd, 0x0e, 0x50, 0x46, 0xd4, 0x9b, 0x4b, 0xfc, 0x83, 0x31, 0xc0, 0x19, 0x90, 0x5b, 0x09,
	0xd4, 0xec, 0x0e, 0xa8, 0x72, 0xc1, 0xfa, 0x13, 0x76, 0x3c, 0x10, 0x23, 0x2f, 0xc7, 0xdc, 0x03,
	0xfd, 0x08, 0x7a, 0xdd, 0x13, 0x3b, 0x79, 0x82, 0x1a, 0x2d, 0xa1, 0xb7, 0x30, 0xfb, 0x44, 0xcb,
	0x2f, 0x91, 0xfc, 0x21, 0x7e, 0x5e, 0xe3, 0x2d, 0x1e, 0x7f, 0x02, 0x00, 0x00, 0xff, 0xff, 0x91,
	0xcb, 0x59, 0x09, 0x99, 0x01, 0x00, 0x00,
}
