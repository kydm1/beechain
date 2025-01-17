// Code generated by protoc-gen-go. DO NOT EDIT.
// source: manage.proto

package types

import (
	fmt "fmt"
	math "math"

	types "github.com/33cn/chain33/types"
	proto "github.com/golang/protobuf/proto"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type ManageAction struct {
	// Types that are valid to be assigned to Value:
	//	*ManageAction_Modify
	Value                isManageAction_Value `protobuf_oneof:"value"`
	Ty                   int32                `protobuf:"varint,2,opt,name=Ty,proto3" json:"Ty,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *ManageAction) Reset()         { *m = ManageAction{} }
func (m *ManageAction) String() string { return proto.CompactTextString(m) }
func (*ManageAction) ProtoMessage()    {}
func (*ManageAction) Descriptor() ([]byte, []int) {
	return fileDescriptor_519fa8ed5ffbbc8f, []int{0}
}

func (m *ManageAction) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ManageAction.Unmarshal(m, b)
}
func (m *ManageAction) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ManageAction.Marshal(b, m, deterministic)
}
func (m *ManageAction) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ManageAction.Merge(m, src)
}
func (m *ManageAction) XXX_Size() int {
	return xxx_messageInfo_ManageAction.Size(m)
}
func (m *ManageAction) XXX_DiscardUnknown() {
	xxx_messageInfo_ManageAction.DiscardUnknown(m)
}

var xxx_messageInfo_ManageAction proto.InternalMessageInfo

type isManageAction_Value interface {
	isManageAction_Value()
}

type ManageAction_Modify struct {
	Modify *types.ModifyConfig `protobuf:"bytes,1,opt,name=modify,proto3,oneof"`
}

func (*ManageAction_Modify) isManageAction_Value() {}

func (m *ManageAction) GetValue() isManageAction_Value {
	if m != nil {
		return m.Value
	}
	return nil
}

func (m *ManageAction) GetModify() *types.ModifyConfig {
	if x, ok := m.GetValue().(*ManageAction_Modify); ok {
		return x.Modify
	}
	return nil
}

func (m *ManageAction) GetTy() int32 {
	if m != nil {
		return m.Ty
	}
	return 0
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*ManageAction) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _ManageAction_OneofMarshaler, _ManageAction_OneofUnmarshaler, _ManageAction_OneofSizer, []interface{}{
		(*ManageAction_Modify)(nil),
	}
}

func _ManageAction_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*ManageAction)
	// value
	switch x := m.Value.(type) {
	case *ManageAction_Modify:
		b.EncodeVarint(1<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Modify); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("ManageAction.Value has unexpected type %T", x)
	}
	return nil
}

func _ManageAction_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*ManageAction)
	switch tag {
	case 1: // value.modify
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(types.ModifyConfig)
		err := b.DecodeMessage(msg)
		m.Value = &ManageAction_Modify{msg}
		return true, err
	default:
		return false, nil
	}
}

func _ManageAction_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*ManageAction)
	// value
	switch x := m.Value.(type) {
	case *ManageAction_Modify:
		s := proto.Size(x.Modify)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

func init() {
	proto.RegisterType((*ManageAction)(nil), "types.ManageAction")
}

func init() { proto.RegisterFile("manage.proto", fileDescriptor_519fa8ed5ffbbc8f) }

var fileDescriptor_519fa8ed5ffbbc8f = []byte{
	// 136 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0xc9, 0x4d, 0xcc, 0x4b,
	0x4c, 0x4f, 0xd5, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x2d, 0xa9, 0x2c, 0x48, 0x2d, 0x96,
	0xe2, 0x4b, 0xad, 0x48, 0x4d, 0x2e, 0x2d, 0xc9, 0x2f, 0x82, 0x08, 0x2b, 0x85, 0x71, 0xf1, 0xf8,
	0x82, 0x95, 0x39, 0x26, 0x97, 0x64, 0xe6, 0xe7, 0x09, 0xe9, 0x72, 0xb1, 0xe5, 0xe6, 0xa7, 0x64,
	0xa6, 0x55, 0x4a, 0x30, 0x2a, 0x30, 0x6a, 0x70, 0x1b, 0x09, 0xeb, 0x81, 0xf5, 0xe9, 0xf9, 0x82,
	0x05, 0x9d, 0xf3, 0xf3, 0xd2, 0x32, 0xd3, 0x3d, 0x18, 0x82, 0xa0, 0x8a, 0x84, 0xf8, 0xb8, 0x98,
	0x42, 0x2a, 0x25, 0x98, 0x14, 0x18, 0x35, 0x58, 0x83, 0x98, 0x42, 0x2a, 0x9d, 0xd8, 0xb9, 0x58,
	0xcb, 0x12, 0x73, 0x4a, 0x53, 0x93, 0xd8, 0xc0, 0xc6, 0x1b, 0x03, 0x02, 0x00, 0x00, 0xff, 0xff,
	0xf1, 0x53, 0xa6, 0x25, 0x85, 0x00, 0x00, 0x00,
}
