/*
Copyright The KubeVault Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: kubevault.dev/operator/apis/config/v1alpha1/generated.proto

package v1alpha1

import (
	fmt "fmt"
	io "io"
	math "math"
	math_bits "math/bits"
	reflect "reflect"
	strings "strings"

	proto "github.com/gogo/protobuf/proto"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

func (m *VaultServerConfiguration) Reset()      { *m = VaultServerConfiguration{} }
func (*VaultServerConfiguration) ProtoMessage() {}
func (*VaultServerConfiguration) Descriptor() ([]byte, []int) {
	return fileDescriptor_c0900205865bff39, []int{0}
}
func (m *VaultServerConfiguration) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *VaultServerConfiguration) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (m *VaultServerConfiguration) XXX_Merge(src proto.Message) {
	xxx_messageInfo_VaultServerConfiguration.Merge(m, src)
}
func (m *VaultServerConfiguration) XXX_Size() int {
	return m.Size()
}
func (m *VaultServerConfiguration) XXX_DiscardUnknown() {
	xxx_messageInfo_VaultServerConfiguration.DiscardUnknown(m)
}

var xxx_messageInfo_VaultServerConfiguration proto.InternalMessageInfo

func init() {
	proto.RegisterType((*VaultServerConfiguration)(nil), "kubevault.dev.operator.apis.config.v1alpha1.VaultServerConfiguration")
}

func init() {
	proto.RegisterFile("kubevault.dev/operator/apis/config/v1alpha1/generated.proto", fileDescriptor_c0900205865bff39)
}

var fileDescriptor_c0900205865bff39 = []byte{
	// 442 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x93, 0x41, 0x8b, 0xd3, 0x40,
	0x14, 0xc7, 0x13, 0xad, 0xb2, 0x06, 0x4f, 0xc3, 0x1e, 0x42, 0x91, 0x69, 0xf0, 0x62, 0x41, 0x9c,
	0x61, 0x41, 0x44, 0xf0, 0xb4, 0xad, 0x08, 0x0a, 0x4a, 0x49, 0xd5, 0x83, 0x78, 0x99, 0xa6, 0x6f,
	0x93, 0xa1, 0x49, 0x5e, 0x9c, 0xcc, 0x8c, 0xec, 0xcd, 0x8f, 0xe0, 0x57, 0xf0, 0xdb, 0xf4, 0xb8,
	0xc7, 0x3d, 0x2d, 0x36, 0x7e, 0x11, 0xe9, 0x34, 0x3d, 0x6c, 0xdb, 0x6c, 0x6f, 0xf3, 0x78, 0xff,
	0xdf, 0xef, 0xbd, 0x77, 0x98, 0xe0, 0xcd, 0xc2, 0xcc, 0xc0, 0x0a, 0x93, 0x6b, 0x36, 0x07, 0xcb,
	0xb1, 0x02, 0x25, 0x34, 0x2a, 0x2e, 0x2a, 0x59, 0xf3, 0x04, 0xcb, 0x0b, 0x99, 0x72, 0x7b, 0x26,
	0xf2, 0x2a, 0x13, 0x67, 0x3c, 0x85, 0x72, 0xdd, 0x86, 0x39, 0xab, 0x14, 0x6a, 0x24, 0xcf, 0x6f,
	0xc1, 0x6c, 0x0b, 0xb3, 0x35, 0xcc, 0x36, 0x30, 0xdb, 0xc2, 0xfd, 0x17, 0xa9, 0xd4, 0x99, 0x99,
	0xb1, 0x04, 0x0b, 0x9e, 0x62, 0x8a, 0xdc, 0x39, 0x66, 0xe6, 0xc2, 0x55, 0xae, 0x70, 0xaf, 0x8d,
	0xbb, 0xff, 0x72, 0xf1, 0xba, 0x66, 0x12, 0xd7, 0x8b, 0x14, 0x22, 0xc9, 0x64, 0x09, 0xea, 0x92,
	0x57, 0x8b, 0x74, 0xb3, 0x59, 0x01, 0x5a, 0x70, 0xbb, 0xb7, 0x51, 0x9f, 0x77, 0x51, 0xca, 0x94,
	0x5a, 0x16, 0xb0, 0x07, 0xbc, 0x3a, 0x06, 0xd4, 0x49, 0x06, 0x85, 0xd8, 0xe5, 0x9e, 0xfe, 0xe9,
	0x05, 0xe1, 0xd7, 0xf5, 0xe5, 0x53, 0x50, 0x16, 0xd4, 0xd8, 0x1d, 0x6b, 0x94, 0xd0, 0x12, 0x4b,
	0x12, 0x05, 0xbd, 0x4a, 0xe8, 0x2c, 0xf4, 0x23, 0x7f, 0xf8, 0x68, 0xf4, 0x78, 0x79, 0x33, 0xf0,
	0x9a, 0x9b, 0x41, 0x6f, 0x22, 0x74, 0x16, 0xbb, 0x0e, 0xf9, 0x10, 0x90, 0x1a, 0x94, 0x95, 0x09,
	0x9c, 0x27, 0x09, 0x9a, 0x52, 0x7f, 0x12, 0x05, 0x84, 0xf7, 0x5c, 0xbe, 0xdf, 0xe6, 0xc9, 0x74,
	0x2f, 0x11, 0x1f, 0xa0, 0xc8, 0x8f, 0x60, 0xa0, 0x71, 0x01, 0x65, 0x0c, 0x56, 0xc2, 0x4f, 0x50,
	0xfb, 0x58, 0x78, 0xdf, 0x89, 0x9f, 0xb5, 0xe2, 0xc1, 0xe7, 0xbb, 0xe3, 0xf1, 0x31, 0x1f, 0x99,
	0x04, 0xa7, 0x15, 0xe6, 0x32, 0xb9, 0x1c, 0x63, 0xa9, 0x15, 0xe6, 0x39, 0xa8, 0x18, 0x73, 0x08,
	0x7b, 0x6e, 0xce, 0x93, 0x76, 0xce, 0xe9, 0xe4, 0x40, 0x26, 0x3e, 0x48, 0x92, 0xef, 0x41, 0x28,
	0x8c, 0xce, 0x3e, 0x82, 0xce, 0x70, 0xbe, 0x63, 0x7d, 0xe0, 0xac, 0x51, 0x6b, 0x0d, 0xcf, 0x3b,
	0x72, 0x71, 0xa7, 0x81, 0xe8, 0x20, 0x32, 0x35, 0x4c, 0x70, 0x7e, 0xfb, 0x96, 0x77, 0xa8, 0xc6,
	0xb5, 0x7c, 0xab, 0xa4, 0x05, 0x15, 0x3e, 0x8c, 0xfc, 0xe1, 0xc9, 0x68, 0xd8, 0x4e, 0x89, 0xbe,
	0x74, 0xe4, 0xa7, 0xef, 0x37, 0xf9, 0xf8, 0xa8, 0x71, 0xc4, 0x96, 0x2b, 0xea, 0x5d, 0xad, 0xa8,
	0x77, 0xbd, 0xa2, 0xde, 0xaf, 0x86, 0xfa, 0xcb, 0x86, 0xfa, 0x57, 0x0d, 0xf5, 0xaf, 0x1b, 0xea,
	0xff, 0x6d, 0xa8, 0xff, 0xfb, 0x1f, 0xf5, 0xbe, 0x9d, 0x6c, 0x7f, 0xc8, 0xff, 0x00, 0x00, 0x00,
	0xff, 0xff, 0x58, 0x72, 0x75, 0x2f, 0x8c, 0x03, 0x00, 0x00,
}

func (m *VaultServerConfiguration) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *VaultServerConfiguration) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *VaultServerConfiguration) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	i--
	if m.UsePodServiceAccountForCSIDriver {
		dAtA[i] = 1
	} else {
		dAtA[i] = 0
	}
	i--
	dAtA[i] = 0x30
	i -= len(m.AuthMethodControllerRole)
	copy(dAtA[i:], m.AuthMethodControllerRole)
	i = encodeVarintGenerated(dAtA, i, uint64(len(m.AuthMethodControllerRole)))
	i--
	dAtA[i] = 0x2a
	i -= len(m.PolicyControllerRole)
	copy(dAtA[i:], m.PolicyControllerRole)
	i = encodeVarintGenerated(dAtA, i, uint64(len(m.PolicyControllerRole)))
	i--
	dAtA[i] = 0x22
	i -= len(m.TokenReviewerServiceAccountName)
	copy(dAtA[i:], m.TokenReviewerServiceAccountName)
	i = encodeVarintGenerated(dAtA, i, uint64(len(m.TokenReviewerServiceAccountName)))
	i--
	dAtA[i] = 0x1a
	i -= len(m.ServiceAccountName)
	copy(dAtA[i:], m.ServiceAccountName)
	i = encodeVarintGenerated(dAtA, i, uint64(len(m.ServiceAccountName)))
	i--
	dAtA[i] = 0x12
	i -= len(m.Path)
	copy(dAtA[i:], m.Path)
	i = encodeVarintGenerated(dAtA, i, uint64(len(m.Path)))
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func encodeVarintGenerated(dAtA []byte, offset int, v uint64) int {
	offset -= sovGenerated(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *VaultServerConfiguration) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Path)
	n += 1 + l + sovGenerated(uint64(l))
	l = len(m.ServiceAccountName)
	n += 1 + l + sovGenerated(uint64(l))
	l = len(m.TokenReviewerServiceAccountName)
	n += 1 + l + sovGenerated(uint64(l))
	l = len(m.PolicyControllerRole)
	n += 1 + l + sovGenerated(uint64(l))
	l = len(m.AuthMethodControllerRole)
	n += 1 + l + sovGenerated(uint64(l))
	n += 2
	return n
}

func sovGenerated(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozGenerated(x uint64) (n int) {
	return sovGenerated(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (this *VaultServerConfiguration) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&VaultServerConfiguration{`,
		`Path:` + fmt.Sprintf("%v", this.Path) + `,`,
		`ServiceAccountName:` + fmt.Sprintf("%v", this.ServiceAccountName) + `,`,
		`TokenReviewerServiceAccountName:` + fmt.Sprintf("%v", this.TokenReviewerServiceAccountName) + `,`,
		`PolicyControllerRole:` + fmt.Sprintf("%v", this.PolicyControllerRole) + `,`,
		`AuthMethodControllerRole:` + fmt.Sprintf("%v", this.AuthMethodControllerRole) + `,`,
		`UsePodServiceAccountForCSIDriver:` + fmt.Sprintf("%v", this.UsePodServiceAccountForCSIDriver) + `,`,
		`}`,
	}, "")
	return s
}
func valueToStringGenerated(v interface{}) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("*%v", pv)
}
func (m *VaultServerConfiguration) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGenerated
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
			return fmt.Errorf("proto: VaultServerConfiguration: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: VaultServerConfiguration: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Path", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenerated
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
				return ErrInvalidLengthGenerated
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGenerated
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Path = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ServiceAccountName", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenerated
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
				return ErrInvalidLengthGenerated
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGenerated
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ServiceAccountName = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TokenReviewerServiceAccountName", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenerated
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
				return ErrInvalidLengthGenerated
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGenerated
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.TokenReviewerServiceAccountName = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PolicyControllerRole", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenerated
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
				return ErrInvalidLengthGenerated
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGenerated
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.PolicyControllerRole = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field AuthMethodControllerRole", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenerated
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
				return ErrInvalidLengthGenerated
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGenerated
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.AuthMethodControllerRole = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 6:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field UsePodServiceAccountForCSIDriver", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenerated
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.UsePodServiceAccountForCSIDriver = bool(v != 0)
		default:
			iNdEx = preIndex
			skippy, err := skipGenerated(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthGenerated
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthGenerated
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
func skipGenerated(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowGenerated
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
					return 0, ErrIntOverflowGenerated
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
			return iNdEx, nil
		case 1:
			iNdEx += 8
			return iNdEx, nil
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowGenerated
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
				return 0, ErrInvalidLengthGenerated
			}
			iNdEx += length
			if iNdEx < 0 {
				return 0, ErrInvalidLengthGenerated
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowGenerated
					}
					if iNdEx >= l {
						return 0, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					innerWire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				innerWireType := int(innerWire & 0x7)
				if innerWireType == 4 {
					break
				}
				next, err := skipGenerated(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
				if iNdEx < 0 {
					return 0, ErrInvalidLengthGenerated
				}
			}
			return iNdEx, nil
		case 4:
			return iNdEx, nil
		case 5:
			iNdEx += 4
			return iNdEx, nil
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
	}
	panic("unreachable")
}

var (
	ErrInvalidLengthGenerated = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowGenerated   = fmt.Errorf("proto: integer overflow")
)
