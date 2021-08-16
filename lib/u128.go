package lib

import (
	"math"
	"math/bits"
)

// U128 is a unsigned 128-bit number type in big-endian form.
//
// The endianness of `Hi` and `Lo` themselves is machine-dependent.
type U128 struct {
	Hi uint64
	Lo uint64
}

// Copyright 2021 The Inet.Af AUTHORS. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// U128BitMasks are bitmasks with the topmost n bits of a 128-bit number, where
// n is the array index.
//
// generated with https://play.golang.org/p/64XKxaUSa_9
var U128BitMasks = [...]U128{
	0:   {0x0000000000000000, 0x0000000000000000},
	1:   {0x8000000000000000, 0x0000000000000000},
	2:   {0xc000000000000000, 0x0000000000000000},
	3:   {0xe000000000000000, 0x0000000000000000},
	4:   {0xf000000000000000, 0x0000000000000000},
	5:   {0xf800000000000000, 0x0000000000000000},
	6:   {0xfc00000000000000, 0x0000000000000000},
	7:   {0xfe00000000000000, 0x0000000000000000},
	8:   {0xff00000000000000, 0x0000000000000000},
	9:   {0xff80000000000000, 0x0000000000000000},
	10:  {0xffc0000000000000, 0x0000000000000000},
	11:  {0xffe0000000000000, 0x0000000000000000},
	12:  {0xfff0000000000000, 0x0000000000000000},
	13:  {0xfff8000000000000, 0x0000000000000000},
	14:  {0xfffc000000000000, 0x0000000000000000},
	15:  {0xfffe000000000000, 0x0000000000000000},
	16:  {0xffff000000000000, 0x0000000000000000},
	17:  {0xffff800000000000, 0x0000000000000000},
	18:  {0xffffc00000000000, 0x0000000000000000},
	19:  {0xffffe00000000000, 0x0000000000000000},
	20:  {0xfffff00000000000, 0x0000000000000000},
	21:  {0xfffff80000000000, 0x0000000000000000},
	22:  {0xfffffc0000000000, 0x0000000000000000},
	23:  {0xfffffe0000000000, 0x0000000000000000},
	24:  {0xffffff0000000000, 0x0000000000000000},
	25:  {0xffffff8000000000, 0x0000000000000000},
	26:  {0xffffffc000000000, 0x0000000000000000},
	27:  {0xffffffe000000000, 0x0000000000000000},
	28:  {0xfffffff000000000, 0x0000000000000000},
	29:  {0xfffffff800000000, 0x0000000000000000},
	30:  {0xfffffffc00000000, 0x0000000000000000},
	31:  {0xfffffffe00000000, 0x0000000000000000},
	32:  {0xffffffff00000000, 0x0000000000000000},
	33:  {0xffffffff80000000, 0x0000000000000000},
	34:  {0xffffffffc0000000, 0x0000000000000000},
	35:  {0xffffffffe0000000, 0x0000000000000000},
	36:  {0xfffffffff0000000, 0x0000000000000000},
	37:  {0xfffffffff8000000, 0x0000000000000000},
	38:  {0xfffffffffc000000, 0x0000000000000000},
	39:  {0xfffffffffe000000, 0x0000000000000000},
	40:  {0xffffffffff000000, 0x0000000000000000},
	41:  {0xffffffffff800000, 0x0000000000000000},
	42:  {0xffffffffffc00000, 0x0000000000000000},
	43:  {0xffffffffffe00000, 0x0000000000000000},
	44:  {0xfffffffffff00000, 0x0000000000000000},
	45:  {0xfffffffffff80000, 0x0000000000000000},
	46:  {0xfffffffffffc0000, 0x0000000000000000},
	47:  {0xfffffffffffe0000, 0x0000000000000000},
	48:  {0xffffffffffff0000, 0x0000000000000000},
	49:  {0xffffffffffff8000, 0x0000000000000000},
	50:  {0xffffffffffffc000, 0x0000000000000000},
	51:  {0xffffffffffffe000, 0x0000000000000000},
	52:  {0xfffffffffffff000, 0x0000000000000000},
	53:  {0xfffffffffffff800, 0x0000000000000000},
	54:  {0xfffffffffffffc00, 0x0000000000000000},
	55:  {0xfffffffffffffe00, 0x0000000000000000},
	56:  {0xffffffffffffff00, 0x0000000000000000},
	57:  {0xffffffffffffff80, 0x0000000000000000},
	58:  {0xffffffffffffffc0, 0x0000000000000000},
	59:  {0xffffffffffffffe0, 0x0000000000000000},
	60:  {0xfffffffffffffff0, 0x0000000000000000},
	61:  {0xfffffffffffffff8, 0x0000000000000000},
	62:  {0xfffffffffffffffc, 0x0000000000000000},
	63:  {0xfffffffffffffffe, 0x0000000000000000},
	64:  {0xffffffffffffffff, 0x0000000000000000},
	65:  {0xffffffffffffffff, 0x8000000000000000},
	66:  {0xffffffffffffffff, 0xc000000000000000},
	67:  {0xffffffffffffffff, 0xe000000000000000},
	68:  {0xffffffffffffffff, 0xf000000000000000},
	69:  {0xffffffffffffffff, 0xf800000000000000},
	70:  {0xffffffffffffffff, 0xfc00000000000000},
	71:  {0xffffffffffffffff, 0xfe00000000000000},
	72:  {0xffffffffffffffff, 0xff00000000000000},
	73:  {0xffffffffffffffff, 0xff80000000000000},
	74:  {0xffffffffffffffff, 0xffc0000000000000},
	75:  {0xffffffffffffffff, 0xffe0000000000000},
	76:  {0xffffffffffffffff, 0xfff0000000000000},
	77:  {0xffffffffffffffff, 0xfff8000000000000},
	78:  {0xffffffffffffffff, 0xfffc000000000000},
	79:  {0xffffffffffffffff, 0xfffe000000000000},
	80:  {0xffffffffffffffff, 0xffff000000000000},
	81:  {0xffffffffffffffff, 0xffff800000000000},
	82:  {0xffffffffffffffff, 0xffffc00000000000},
	83:  {0xffffffffffffffff, 0xffffe00000000000},
	84:  {0xffffffffffffffff, 0xfffff00000000000},
	85:  {0xffffffffffffffff, 0xfffff80000000000},
	86:  {0xffffffffffffffff, 0xfffffc0000000000},
	87:  {0xffffffffffffffff, 0xfffffe0000000000},
	88:  {0xffffffffffffffff, 0xffffff0000000000},
	89:  {0xffffffffffffffff, 0xffffff8000000000},
	90:  {0xffffffffffffffff, 0xffffffc000000000},
	91:  {0xffffffffffffffff, 0xffffffe000000000},
	92:  {0xffffffffffffffff, 0xfffffff000000000},
	93:  {0xffffffffffffffff, 0xfffffff800000000},
	94:  {0xffffffffffffffff, 0xfffffffc00000000},
	95:  {0xffffffffffffffff, 0xfffffffe00000000},
	96:  {0xffffffffffffffff, 0xffffffff00000000},
	97:  {0xffffffffffffffff, 0xffffffff80000000},
	98:  {0xffffffffffffffff, 0xffffffffc0000000},
	99:  {0xffffffffffffffff, 0xffffffffe0000000},
	100: {0xffffffffffffffff, 0xfffffffff0000000},
	101: {0xffffffffffffffff, 0xfffffffff8000000},
	102: {0xffffffffffffffff, 0xfffffffffc000000},
	103: {0xffffffffffffffff, 0xfffffffffe000000},
	104: {0xffffffffffffffff, 0xffffffffff000000},
	105: {0xffffffffffffffff, 0xffffffffff800000},
	106: {0xffffffffffffffff, 0xffffffffffc00000},
	107: {0xffffffffffffffff, 0xffffffffffe00000},
	108: {0xffffffffffffffff, 0xfffffffffff00000},
	109: {0xffffffffffffffff, 0xfffffffffff80000},
	110: {0xffffffffffffffff, 0xfffffffffffc0000},
	111: {0xffffffffffffffff, 0xfffffffffffe0000},
	112: {0xffffffffffffffff, 0xffffffffffff0000},
	113: {0xffffffffffffffff, 0xffffffffffff8000},
	114: {0xffffffffffffffff, 0xffffffffffffc000},
	115: {0xffffffffffffffff, 0xffffffffffffe000},
	116: {0xffffffffffffffff, 0xfffffffffffff000},
	117: {0xffffffffffffffff, 0xfffffffffffff800},
	118: {0xffffffffffffffff, 0xfffffffffffffc00},
	119: {0xffffffffffffffff, 0xfffffffffffffe00},
	120: {0xffffffffffffffff, 0xffffffffffffff00},
	121: {0xffffffffffffffff, 0xffffffffffffff80},
	122: {0xffffffffffffffff, 0xffffffffffffffc0},
	123: {0xffffffffffffffff, 0xffffffffffffffe0},
	124: {0xffffffffffffffff, 0xfffffffffffffff0},
	125: {0xffffffffffffffff, 0xfffffffffffffff8},
	126: {0xffffffffffffffff, 0xfffffffffffffffc},
	127: {0xffffffffffffffff, 0xfffffffffffffffe},
	128: {0xffffffffffffffff, 0xffffffffffffffff},
}

// ZeroU128 is a zero-valued U128.
var ZeroU128 U128

// MaxU128 is the largest possible U128 value.
var MaxU128 = NewU128(math.MaxUint64, math.MaxUint64)

// NewU128 returns the U128 value .
func NewU128(hi uint64, lo uint64) U128 {
	return U128{hi, lo}
}

// U128From64 converts a uint64 to a U128 value.
func U128From64(v uint64) U128 {
	return NewU128(0, v)
}

// IsZero returns whether `v` is 0.
func (v U128) IsZero() bool {
	return v.Hi|v.Lo == 0
}

// IsMax returns whether `v` is the maximum U128 number.
func (v U128) IsMax() bool {
	return v.Hi&v.Lo == math.MaxUint64
}

// Cmp compares `v1` and `v2` and returns:
//
// - -1 if `v1<v2`
// -  0 if `v1==v2`
// -  1 if `v1>v2`
func (v1 U128) Cmp(v2 U128) int {
	if v1 == v2 {
		return 0
	} else if v1.Hi < v2.Hi || (v1.Hi == v2.Hi && v1.Lo < v2.Lo) {
		return -1
	} else {
		return 1
	}
}

// Eq returns whether `v1 == v2`.
func (v1 U128) Eq(v2 U128) bool {
	return v1.Cmp(v2) == 0
}

// Gt returns whether `v1 > v2`.
func (v1 U128) Gt(v2 U128) bool {
	return v1.Cmp(v2) > 0
}

// Gte returns whether `v1 >= v2`.
func (v1 U128) Gte(v2 U128) bool {
	return v1.Cmp(v2) >= 0
}

// Lt returns whether `v1 < v2`.
func (v1 U128) Lt(v2 U128) bool {
	return v1.Cmp(v2) < 0
}

// Lte returns whether `v1 <= v2`.
func (v1 U128) Lte(v2 U128) bool {
	return v1.Cmp(v2) <= 0
}

// And returns `v1 & v2`.
func (v1 U128) And(v2 U128) U128 {
	return U128{v1.Hi & v2.Hi, v1.Lo & v2.Lo}
}

// And64 returns `v1 & v2`.
func (v1 U128) And64(v2 uint64) U128 {
	return U128{v1.Hi & 0, v1.Lo & v2}
}

// Or returns `v1 | v2`.
func (v1 U128) Or(v2 U128) U128 {
	return U128{v1.Hi | v2.Hi, v1.Lo | v2.Lo}
}

// Or64 returns `v1 | v2`.
func (v1 U128) Or64(v2 uint64) U128 {
	return U128{v1.Hi | 0, v1.Lo | v2}
}

// Xor returns `v1 ^ v2`.
func (v1 U128) Xor(v2 U128) U128 {
	return U128{v1.Hi ^ v2.Hi, v1.Lo ^ v2.Lo}
}

// Xor64 returns `v1 ^ v2`.
func (v1 U128) Xor64(v2 uint64) U128 {
	return U128{v1.Hi ^ 0, v1.Lo ^ v2}
}

// Not returns `^v`.
func (v U128) Not() U128 {
	return U128{^v.Hi, ^v.Lo}
}

// Add returns `v1 + v2` and any carry on overflow.
// The carry is guaranteed to be 0 or 1.
func (v1 U128) Add(v2 U128) (U128, uint64) {
	lo, carry := bits.Add64(v1.Lo, v2.Lo, 0)
	hi, carry := bits.Add64(v1.Hi, v2.Hi, carry)
	return U128{hi, lo}, carry
}

// Add64 returns `v1 + v2` and any carry on overflow.
// The carry is guaranteed to be 0 or 1.
func (v1 U128) Add64(v2 uint64) (U128, uint64) {
	lo, carry := bits.Add64(v1.Lo, v2, 0)
	hi, carry := bits.Add64(v1.Hi, 0, carry)
	return U128{hi, lo}, carry
}

// AddOne returns `v1 + 1`.
// If overflow occurred, then IsZero will be true on the result.
func (v U128) AddOne() U128 {
	lo, carry := bits.Add64(v.Lo, 1, 0)
	return U128{v.Hi + carry, lo}
}

// Sub returns `v1 - v2` and any borrow on underflow.
// The borrow is guaranteed to be 0 or 1.
func (v1 U128) Sub(v2 U128) (U128, uint64) {
	lo, borrow := bits.Sub64(v1.Lo, v2.Lo, 0)
	hi, borrow := bits.Sub64(v1.Hi, v2.Hi, borrow)
	return U128{hi, lo}, borrow
}

// Sub64 returns `v1 - v2` and any borrow on underflow.
// The borrow is guaranteed to be 0 or 1.
func (v1 U128) Sub64(v2 uint64) (U128, uint64) {
	lo, borrow := bits.Sub64(v1.Lo, v2, 0)
	hi, borrow := bits.Sub64(v1.Hi, 0, borrow)
	return U128{hi, lo}, borrow
}

// SubOne returns `v1 - 1`.
// If underflow occurred, then IsMax will be true on the result.
func (v U128) SubOne() U128 {
	lo, borrow := bits.Sub64(v.Lo, 1, 0)
	return U128{v.Hi - borrow, lo}
}

// LeadingZeros returns the number of leading zero bits in `v`.
// Returns 128 for `v == 0`.
func (v U128) LeadingZeros() int {
	if v.Hi > 0 {
		return bits.LeadingZeros64(v.Hi)
	}
	return 64 + bits.LeadingZeros64(v.Lo)
}

// TrailingZeros returns the number of trailing zero bits in v.
// Returns 128 for `v == 0`.
func (v U128) TrailingZeros() int {
	if v.Lo > 0 {
		return bits.TrailingZeros64(v.Lo)
	}
	return 64 + bits.TrailingZeros64(v.Hi)
}

// OnesCount returns the number of "1" bits in `v`.
// This is also sometimes referred to as the "population count" of `v`.
func (v U128) OnesCount() int {
	return bits.OnesCount64(v.Hi) + bits.OnesCount64(v.Lo)
}

// Reverse returns a version of `v` with all bits reversed.
func (v U128) Reverse() U128 {
	return U128{bits.Reverse64(v.Lo), bits.Reverse64(v.Hi)}
}

// ReverseBytes returns a version of `v` with all bytes reversed.
func (v U128) ReverseBytes() U128 {
	return U128{bits.ReverseBytes64(v.Lo), bits.ReverseBytes64(v.Hi)}
}

// Len returns the minimum number of bits required to represent `v`.
// Returns 0 for `v == 0`.
func (v U128) Len() int {
	return 128 - v.LeadingZeros()
}

// SetBitsFrom sets all bits from and including the `i`th bit, where `i` starts
// from 0 which refers to the most-significant bit.
func (v U128) SetBitsFrom(i uint8) U128 {
	return v.Or(U128BitMasks[i].Not())
}

// SetBitsUpto sets all bits from the most-significant bit up to the `ith` bit,
// where `i` starts from 0 which refers to the most-significant bit.
func (v U128) SetBitsUpto(i uint8) U128 {
	return v.Or(U128BitMasks[i])
}

// ClearBitsFrom clears all bits from and including the `i`th bit, where `i`
// starts from 0 which refers to the most-significant bit.
func (v U128) ClearBitsFrom(i uint8) U128 {
	return v.And(U128BitMasks[i])
}

// ClearBitsUpto clears all bits from the most-significant bit up to the `ith`
// bit, where `i` starts from 0 which refers to the most-significant bit.
func (v U128) ClearBitsUpto(i uint8) U128 {
	return v.And(U128BitMasks[i].Not())
}
