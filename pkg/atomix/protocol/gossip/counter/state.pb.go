// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: atomix/protocol/gossip/counter/state.proto

package counter

import (
	fmt "fmt"
	meta "github.com/atomix/api/go/atomix/primitive/meta"
	_ "github.com/atomix/go-framework/pkg/atomix/protocol/gossip"
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

type CounterState struct {
	meta.ObjectMeta `protobuf:"bytes,1,opt,name=meta,proto3,embedded=meta" json:"meta"`
	Increments      map[string]int64 `protobuf:"bytes,2,rep,name=increments,proto3" json:"increments,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"`
	Decrements      map[string]int64 `protobuf:"bytes,3,rep,name=decrements,proto3" json:"decrements,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"`
}

func (m *CounterState) Reset()         { *m = CounterState{} }
func (m *CounterState) String() string { return proto.CompactTextString(m) }
func (*CounterState) ProtoMessage()    {}
func (*CounterState) Descriptor() ([]byte, []int) {
	return fileDescriptor_f8b29713b11736d7, []int{0}
}
func (m *CounterState) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *CounterState) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_CounterState.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *CounterState) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CounterState.Merge(m, src)
}
func (m *CounterState) XXX_Size() int {
	return m.Size()
}
func (m *CounterState) XXX_DiscardUnknown() {
	xxx_messageInfo_CounterState.DiscardUnknown(m)
}

var xxx_messageInfo_CounterState proto.InternalMessageInfo

func (m *CounterState) GetIncrements() map[string]int64 {
	if m != nil {
		return m.Increments
	}
	return nil
}

func (m *CounterState) GetDecrements() map[string]int64 {
	if m != nil {
		return m.Decrements
	}
	return nil
}

func init() {
	proto.RegisterType((*CounterState)(nil), "atomix.protocol.gossip.counter.CounterState")
	proto.RegisterMapType((map[string]int64)(nil), "atomix.protocol.gossip.counter.CounterState.DecrementsEntry")
	proto.RegisterMapType((map[string]int64)(nil), "atomix.protocol.gossip.counter.CounterState.IncrementsEntry")
}

func init() {
	proto.RegisterFile("atomix/protocol/gossip/counter/state.proto", fileDescriptor_f8b29713b11736d7)
}

var fileDescriptor_f8b29713b11736d7 = []byte{
	// 319 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xd2, 0x4a, 0x2c, 0xc9, 0xcf,
	0xcd, 0xac, 0xd0, 0x2f, 0x28, 0xca, 0x2f, 0xc9, 0x4f, 0xce, 0xcf, 0xd1, 0x4f, 0xcf, 0x2f, 0x2e,
	0xce, 0x2c, 0xd0, 0x4f, 0xce, 0x2f, 0xcd, 0x2b, 0x49, 0x2d, 0xd2, 0x2f, 0x2e, 0x49, 0x2c, 0x49,
	0xd5, 0x03, 0xcb, 0x0a, 0xc9, 0x41, 0xd4, 0xea, 0xc1, 0xd4, 0xea, 0x41, 0xd4, 0xea, 0x41, 0xd5,
	0x4a, 0x29, 0xc1, 0xcd, 0xca, 0xcc, 0xcd, 0x2c, 0xc9, 0x2c, 0x4b, 0xd5, 0xcf, 0x4d, 0x2d, 0x49,
	0xd4, 0xcf, 0x4f, 0xca, 0x4a, 0x4d, 0x2e, 0x81, 0xe8, 0x42, 0x52, 0x83, 0x6a, 0x1f, 0x92, 0x3d,
	0x52, 0x22, 0xe9, 0xf9, 0xe9, 0xf9, 0x60, 0xa6, 0x3e, 0x88, 0x05, 0x11, 0x55, 0xea, 0x65, 0xe6,
	0xe2, 0x71, 0x86, 0xd8, 0x14, 0x0c, 0x52, 0x2c, 0xe4, 0xcc, 0xc5, 0x02, 0x32, 0x5f, 0x82, 0x51,
	0x81, 0x51, 0x83, 0xdb, 0x48, 0x51, 0x0f, 0xee, 0x3a, 0xa8, 0xed, 0x7a, 0x20, 0x59, 0x3d, 0x7f,
	0xb0, 0xed, 0xbe, 0xa9, 0x25, 0x89, 0x4e, 0x3c, 0x27, 0xee, 0xc9, 0x33, 0x5c, 0xb8, 0x27, 0xcf,
	0xd8, 0xd1, 0xaa, 0xc1, 0x18, 0x04, 0xd6, 0x2c, 0x14, 0xc3, 0xc5, 0x95, 0x99, 0x97, 0x5c, 0x94,
	0x9a, 0x9b, 0x9a, 0x57, 0x52, 0x2c, 0xc1, 0xa4, 0xc0, 0xac, 0xc1, 0x6d, 0x64, 0xa3, 0x87, 0xdf,
	0xa3, 0x7a, 0xc8, 0xce, 0xd0, 0xf3, 0x84, 0x6b, 0x77, 0xcd, 0x2b, 0x29, 0xaa, 0x0c, 0x42, 0x32,
	0x0f, 0x64, 0x7a, 0x4a, 0x2a, 0xdc, 0x74, 0x66, 0x32, 0x4c, 0x77, 0x49, 0x45, 0x33, 0x1d, 0x61,
	0x9e, 0x94, 0x2d, 0x17, 0x3f, 0x9a, 0xe5, 0x42, 0x02, 0x5c, 0xcc, 0xd9, 0xa9, 0x95, 0xe0, 0x20,
	0xe1, 0x0c, 0x02, 0x31, 0x85, 0x44, 0xb8, 0x58, 0xcb, 0x12, 0x73, 0x4a, 0x53, 0x25, 0x98, 0x14,
	0x18, 0x35, 0x98, 0x83, 0x20, 0x1c, 0x2b, 0x26, 0x0b, 0x46, 0x90, 0x76, 0x34, 0xd3, 0x49, 0xd1,
	0xee, 0x24, 0x71, 0xe2, 0x91, 0x1c, 0xe3, 0x85, 0x47, 0x72, 0x8c, 0x0f, 0x1e, 0xc9, 0x31, 0x4e,
	0x78, 0x2c, 0xc7, 0x70, 0xe1, 0xb1, 0x1c, 0xc3, 0x8d, 0xc7, 0x72, 0x0c, 0x49, 0x6c, 0x60, 0x9f,
	0x19, 0x03, 0x02, 0x00, 0x00, 0xff, 0xff, 0x77, 0x84, 0x78, 0x3d, 0x5c, 0x02, 0x00, 0x00,
}

func (m *CounterState) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *CounterState) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *CounterState) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Decrements) > 0 {
		for k := range m.Decrements {
			v := m.Decrements[k]
			baseI := i
			i = encodeVarintState(dAtA, i, uint64(v))
			i--
			dAtA[i] = 0x10
			i -= len(k)
			copy(dAtA[i:], k)
			i = encodeVarintState(dAtA, i, uint64(len(k)))
			i--
			dAtA[i] = 0xa
			i = encodeVarintState(dAtA, i, uint64(baseI-i))
			i--
			dAtA[i] = 0x1a
		}
	}
	if len(m.Increments) > 0 {
		for k := range m.Increments {
			v := m.Increments[k]
			baseI := i
			i = encodeVarintState(dAtA, i, uint64(v))
			i--
			dAtA[i] = 0x10
			i -= len(k)
			copy(dAtA[i:], k)
			i = encodeVarintState(dAtA, i, uint64(len(k)))
			i--
			dAtA[i] = 0xa
			i = encodeVarintState(dAtA, i, uint64(baseI-i))
			i--
			dAtA[i] = 0x12
		}
	}
	{
		size, err := m.ObjectMeta.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintState(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func encodeVarintState(dAtA []byte, offset int, v uint64) int {
	offset -= sovState(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *CounterState) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.ObjectMeta.Size()
	n += 1 + l + sovState(uint64(l))
	if len(m.Increments) > 0 {
		for k, v := range m.Increments {
			_ = k
			_ = v
			mapEntrySize := 1 + len(k) + sovState(uint64(len(k))) + 1 + sovState(uint64(v))
			n += mapEntrySize + 1 + sovState(uint64(mapEntrySize))
		}
	}
	if len(m.Decrements) > 0 {
		for k, v := range m.Decrements {
			_ = k
			_ = v
			mapEntrySize := 1 + len(k) + sovState(uint64(len(k))) + 1 + sovState(uint64(v))
			n += mapEntrySize + 1 + sovState(uint64(mapEntrySize))
		}
	}
	return n
}

func sovState(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozState(x uint64) (n int) {
	return sovState(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *CounterState) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowState
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
			return fmt.Errorf("proto: CounterState: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: CounterState: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ObjectMeta", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowState
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
				return ErrInvalidLengthState
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthState
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.ObjectMeta.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Increments", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowState
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
				return ErrInvalidLengthState
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthState
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Increments == nil {
				m.Increments = make(map[string]int64)
			}
			var mapkey string
			var mapvalue int64
			for iNdEx < postIndex {
				entryPreIndex := iNdEx
				var wire uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowState
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
				if fieldNum == 1 {
					var stringLenmapkey uint64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowState
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						stringLenmapkey |= uint64(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					intStringLenmapkey := int(stringLenmapkey)
					if intStringLenmapkey < 0 {
						return ErrInvalidLengthState
					}
					postStringIndexmapkey := iNdEx + intStringLenmapkey
					if postStringIndexmapkey < 0 {
						return ErrInvalidLengthState
					}
					if postStringIndexmapkey > l {
						return io.ErrUnexpectedEOF
					}
					mapkey = string(dAtA[iNdEx:postStringIndexmapkey])
					iNdEx = postStringIndexmapkey
				} else if fieldNum == 2 {
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowState
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						mapvalue |= int64(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
				} else {
					iNdEx = entryPreIndex
					skippy, err := skipState(dAtA[iNdEx:])
					if err != nil {
						return err
					}
					if skippy < 0 {
						return ErrInvalidLengthState
					}
					if (iNdEx + skippy) > postIndex {
						return io.ErrUnexpectedEOF
					}
					iNdEx += skippy
				}
			}
			m.Increments[mapkey] = mapvalue
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Decrements", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowState
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
				return ErrInvalidLengthState
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthState
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Decrements == nil {
				m.Decrements = make(map[string]int64)
			}
			var mapkey string
			var mapvalue int64
			for iNdEx < postIndex {
				entryPreIndex := iNdEx
				var wire uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowState
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
				if fieldNum == 1 {
					var stringLenmapkey uint64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowState
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						stringLenmapkey |= uint64(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					intStringLenmapkey := int(stringLenmapkey)
					if intStringLenmapkey < 0 {
						return ErrInvalidLengthState
					}
					postStringIndexmapkey := iNdEx + intStringLenmapkey
					if postStringIndexmapkey < 0 {
						return ErrInvalidLengthState
					}
					if postStringIndexmapkey > l {
						return io.ErrUnexpectedEOF
					}
					mapkey = string(dAtA[iNdEx:postStringIndexmapkey])
					iNdEx = postStringIndexmapkey
				} else if fieldNum == 2 {
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowState
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						mapvalue |= int64(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
				} else {
					iNdEx = entryPreIndex
					skippy, err := skipState(dAtA[iNdEx:])
					if err != nil {
						return err
					}
					if skippy < 0 {
						return ErrInvalidLengthState
					}
					if (iNdEx + skippy) > postIndex {
						return io.ErrUnexpectedEOF
					}
					iNdEx += skippy
				}
			}
			m.Decrements[mapkey] = mapvalue
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipState(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthState
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthState
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
func skipState(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowState
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
					return 0, ErrIntOverflowState
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
					return 0, ErrIntOverflowState
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
				return 0, ErrInvalidLengthState
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupState
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthState
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthState        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowState          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupState = fmt.Errorf("proto: unexpected end of group")
)