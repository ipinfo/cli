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
