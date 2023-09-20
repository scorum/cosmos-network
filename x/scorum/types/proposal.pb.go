// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: network/scorum/v1/proposal.proto

package types

import (
	fmt "fmt"
	types "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// MintProposal is a gov proposal for minting assets.
type MintProposal struct {
	Title       string     `protobuf:"bytes,1,opt,name=title,proto3" json:"title,omitempty"`
	Description string     `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	Recipient   string     `protobuf:"bytes,3,opt,name=recipient,proto3" json:"recipient,omitempty"`
	Amount      types.Coin `protobuf:"bytes,4,opt,name=amount,proto3" json:"amount"`
}

func (m *MintProposal) Reset()         { *m = MintProposal{} }
func (m *MintProposal) String() string { return proto.CompactTextString(m) }
func (*MintProposal) ProtoMessage()    {}
func (*MintProposal) Descriptor() ([]byte, []int) {
	return fileDescriptor_44f7827f81ae208f, []int{0}
}
func (m *MintProposal) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MintProposal) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MintProposal.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MintProposal) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MintProposal.Merge(m, src)
}
func (m *MintProposal) XXX_Size() int {
	return m.Size()
}
func (m *MintProposal) XXX_DiscardUnknown() {
	xxx_messageInfo_MintProposal.DiscardUnknown(m)
}

var xxx_messageInfo_MintProposal proto.InternalMessageInfo

func (m *MintProposal) GetTitle() string {
	if m != nil {
		return m.Title
	}
	return ""
}

func (m *MintProposal) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *MintProposal) GetRecipient() string {
	if m != nil {
		return m.Recipient
	}
	return ""
}

func (m *MintProposal) GetAmount() types.Coin {
	if m != nil {
		return m.Amount
	}
	return types.Coin{}
}

func init() {
	proto.RegisterType((*MintProposal)(nil), "network.scorum.v1.MintProposal")
}

func init() { proto.RegisterFile("network/scorum/v1/proposal.proto", fileDescriptor_44f7827f81ae208f) }

var fileDescriptor_44f7827f81ae208f = []byte{
	// 268 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x4c, 0x90, 0xbd, 0x4e, 0xc3, 0x30,
	0x14, 0x85, 0x63, 0x28, 0x95, 0xea, 0xb2, 0x10, 0x75, 0x08, 0x15, 0x32, 0x11, 0x53, 0x26, 0x5b,
	0x81, 0x81, 0x3d, 0xcc, 0x48, 0xa8, 0x23, 0x5b, 0x62, 0xac, 0x60, 0xd1, 0xf8, 0x5a, 0xf6, 0x4d,
	0x80, 0xb7, 0xe0, 0x05, 0x78, 0x9f, 0x8e, 0x1d, 0x99, 0x10, 0x4a, 0x5e, 0x04, 0x35, 0x3f, 0x82,
	0xcd, 0x3e, 0xdf, 0x27, 0xeb, 0xf8, 0xd0, 0xd8, 0x28, 0x7c, 0x05, 0xf7, 0x22, 0xbc, 0x04, 0x57,
	0x57, 0xa2, 0x49, 0x85, 0x75, 0x60, 0xc1, 0xe7, 0x5b, 0x6e, 0x1d, 0x20, 0x84, 0x67, 0xa3, 0xc1,
	0x07, 0x83, 0x37, 0xe9, 0x7a, 0x55, 0x42, 0x09, 0x3d, 0x15, 0x87, 0xd3, 0x20, 0xae, 0x99, 0x04,
	0x5f, 0x81, 0x17, 0x45, 0xee, 0x95, 0x68, 0xd2, 0x42, 0x61, 0x9e, 0x0a, 0x09, 0xda, 0x0c, 0xfc,
	0xea, 0x93, 0xd0, 0xd3, 0x7b, 0x6d, 0xf0, 0x61, 0x7c, 0x3f, 0x5c, 0xd1, 0x13, 0xd4, 0xb8, 0x55,
	0x11, 0x89, 0x49, 0xb2, 0xd8, 0x0c, 0x97, 0x30, 0xa6, 0xcb, 0x27, 0xe5, 0xa5, 0xd3, 0x16, 0x35,
	0x98, 0xe8, 0xa8, 0x67, 0xff, 0xa3, 0xf0, 0x82, 0x2e, 0x9c, 0x92, 0xda, 0x6a, 0x65, 0x30, 0x3a,
	0xee, 0xf9, 0x5f, 0x10, 0xde, 0xd2, 0x79, 0x5e, 0x41, 0x6d, 0x30, 0x9a, 0xc5, 0x24, 0x59, 0x5e,
	0x9f, 0xf3, 0xa1, 0x17, 0x3f, 0xf4, 0xe2, 0x63, 0x2f, 0x7e, 0x07, 0xda, 0x64, 0xb3, 0xdd, 0xf7,
	0x65, 0xb0, 0x19, 0xf5, 0x2c, 0xdb, 0xb5, 0x8c, 0xec, 0x5b, 0x46, 0x7e, 0x5a, 0x46, 0x3e, 0x3a,
	0x16, 0xec, 0x3b, 0x16, 0x7c, 0x75, 0x2c, 0x78, 0x4c, 0x4a, 0x8d, 0xcf, 0x75, 0xc1, 0x25, 0x54,
	0xd3, 0x4e, 0xd3, 0x6c, 0x6f, 0x53, 0x80, 0xef, 0x56, 0xf9, 0x62, 0xde, 0x7f, 0xf5, 0xe6, 0x37,
	0x00, 0x00, 0xff, 0xff, 0x2f, 0xe9, 0x2c, 0xf8, 0x57, 0x01, 0x00, 0x00,
}

func (m *MintProposal) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MintProposal) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MintProposal) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size, err := m.Amount.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintProposal(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x22
	if len(m.Recipient) > 0 {
		i -= len(m.Recipient)
		copy(dAtA[i:], m.Recipient)
		i = encodeVarintProposal(dAtA, i, uint64(len(m.Recipient)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Description) > 0 {
		i -= len(m.Description)
		copy(dAtA[i:], m.Description)
		i = encodeVarintProposal(dAtA, i, uint64(len(m.Description)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Title) > 0 {
		i -= len(m.Title)
		copy(dAtA[i:], m.Title)
		i = encodeVarintProposal(dAtA, i, uint64(len(m.Title)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintProposal(dAtA []byte, offset int, v uint64) int {
	offset -= sovProposal(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *MintProposal) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Title)
	if l > 0 {
		n += 1 + l + sovProposal(uint64(l))
	}
	l = len(m.Description)
	if l > 0 {
		n += 1 + l + sovProposal(uint64(l))
	}
	l = len(m.Recipient)
	if l > 0 {
		n += 1 + l + sovProposal(uint64(l))
	}
	l = m.Amount.Size()
	n += 1 + l + sovProposal(uint64(l))
	return n
}

func sovProposal(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozProposal(x uint64) (n int) {
	return sovProposal(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *MintProposal) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowProposal
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: MintProposal: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MintProposal: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Title", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProposal
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthProposal
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthProposal
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Title = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Description", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProposal
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthProposal
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthProposal
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Description = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Recipient", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProposal
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthProposal
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthProposal
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Recipient = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Amount", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProposal
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthProposal
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthProposal
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Amount.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipProposal(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthProposal
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipProposal(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowProposal
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowProposal
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowProposal
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthProposal
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupProposal
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthProposal
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthProposal        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowProposal          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupProposal = fmt.Errorf("proto: unexpected end of group")
)
