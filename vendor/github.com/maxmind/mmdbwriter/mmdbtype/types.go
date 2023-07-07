// Package mmdbtype provides types used within the MaxMind DB format.
package mmdbtype

import (
	"encoding/binary"
	"fmt"
	"io"
	"math/big"
	"math/bits"
	"sort"
)

type typeNum byte

const (
	typeNumExtended typeNum = iota
	typeNumPointer
	typeNumString
	typeNumFloat64
	typeNumBytes
	typeNumUint16
	typeNumUint32
	typeNumMap
	typeNumInt32
	typeNumUint64
	typeNumUint128
	typeNumSlice
	// We don't use the next two. They are placeholders. See the spec
	// for more details.
	typeNumContainer //nolint: deadcode, varcheck // placeholder
	typeNumMarker    //nolint: deadcode, varcheck // placeholder
	typeNumBool
	typeNumFloat32
)

type writer interface {
	io.Writer
	WriteByte(byte) error
	WriteString(string) (int, error)
	WriteOrWritePointer(DataType) (int64, error)
}

// DataType represents a MaxMind DB data type.
type DataType interface {
	Copy() DataType
	size() int
	typeNum() typeNum
	WriteTo(writer) (int64, error)
}

// Bool is the MaxMind DB boolean type.
type Bool bool

// Copy the value.
func (t Bool) Copy() DataType { return t }

func (t Bool) size() int {
	if t {
		return 1
	}
	return 0
}

func (t Bool) typeNum() typeNum {
	return typeNumBool
}

// WriteTo writes the value to w.
func (t Bool) WriteTo(w writer) (int64, error) {
	return writeCtrlByte(w, t)
}

// Bytes is the MaxMind DB bytes type.
type Bytes []byte

// Copy the value.
func (t Bytes) Copy() DataType {
	nv := make(Bytes, len(t))
	copy(nv, t)
	return nv
}

func (t Bytes) size() int {
	return len(t)
}

func (t Bytes) typeNum() typeNum {
	return typeNumBytes
}

// WriteTo writes the value to w.
func (t Bytes) WriteTo(w writer) (int64, error) {
	numBytes, err := writeCtrlByte(w, t)
	if err != nil {
		return numBytes, err
	}

	written, err := w.Write(t)
	numBytes += int64(written)
	if err != nil {
		return numBytes, fmt.Errorf(`writing "%s" as bytes: %w`, t, err)
	}
	return numBytes, nil
}

// Float32 is the MaxMind DB float type.
type Float32 float32

// Copy the value.
func (t Float32) Copy() DataType { return t }

func (t Float32) size() int {
	return 4
}

func (t Float32) typeNum() typeNum {
	return typeNumFloat32
}

// WriteTo writes the value to w.
func (t Float32) WriteTo(w writer) (int64, error) {
	numBytes, err := writeCtrlByte(w, t)
	if err != nil {
		return numBytes, err
	}

	err = binary.Write(w, binary.BigEndian, t)
	if err != nil {
		return numBytes, fmt.Errorf("writing %f as float32: %w", t, err)
	}
	return numBytes + int64(t.size()), nil
}

// Float64 is the MaxMind DB double type.
type Float64 float64

// Copy the value.
func (t Float64) Copy() DataType { return t }

func (t Float64) size() int {
	return 8
}

func (t Float64) typeNum() typeNum {
	return typeNumFloat64
}

// WriteTo writes the value to w.
func (t Float64) WriteTo(w writer) (int64, error) {
	numBytes, err := writeCtrlByte(w, t)
	if err != nil {
		return numBytes, err
	}

	err = binary.Write(w, binary.BigEndian, t)
	if err != nil {
		return numBytes, fmt.Errorf("writing %f as float64: %w", t, err)
	}
	return numBytes + int64(t.size()), nil
}

// Int32 is the MaxMind DB signed 32-bit integer type.
type Int32 int32

// Copy the value.
func (t Int32) Copy() DataType { return t }

func (t Int32) size() int {
	return 4 - bits.LeadingZeros32(uint32(t))/8
}

func (t Int32) typeNum() typeNum {
	return typeNumInt32
}

// WriteTo writes the value to w.
func (t Int32) WriteTo(w writer) (int64, error) {
	numBytes, err := writeCtrlByte(w, t)
	if err != nil {
		return numBytes, err
	}

	size := t.size()
	// We ignore leading zeros
	for i := size; i > 0; i-- {
		err = w.WriteByte(byte((int32(t) >> (8 * (i - 1))) & 0xFF))
		if err != nil {
			return numBytes + int64(size-i), fmt.Errorf("writing int32: %w", err)
		}
	}
	return numBytes + int64(size), nil
}

// Map is the MaxMind DB map type.
type Map map[String]DataType

// Copy makes a deep copy of the Map.
func (t Map) Copy() DataType {
	newMap := make(Map, len(t))
	for k, v := range t {
		newMap[k] = v.Copy()
	}
	return newMap
}

func (t Map) size() int {
	return len(t)
}

func (t Map) typeNum() typeNum {
	return typeNumMap
}

// WriteTo writes the value to w.
func (t Map) WriteTo(w writer) (int64, error) {
	numBytes, err := writeCtrlByte(w, t)
	if err != nil {
		return numBytes, err
	}

	// We want database builds to be reproducible. As such, we insert
	// the map items in order by key value. In the future, we will
	// likely use a more relevant characteristic here (e.g., putting
	// fields more likely to be accessed first).
	keys := make([]string, 0, len(t))
	for k := range t {
		keys = append(keys, string(k))
	}
	sort.Strings(keys)

	for _, ks := range keys {
		k := String(ks)
		written, err := w.WriteOrWritePointer(k)
		numBytes += written
		if err != nil {
			return numBytes, err
		}
		written, err = w.WriteOrWritePointer(t[k])
		numBytes += written
		if err != nil {
			return numBytes, err
		}
	}
	return numBytes, nil
}

// Pointer is the MaxMind DB pointer type for internal use in the writer. You
// should not use this type in data structures that you pass to methods on
// mmdbwriter.Tree. Doing so may result in a corrupt database.
type Pointer uint32

// Copy the value.
func (t Pointer) Copy() DataType { return t }

const (
	pointerMaxSize0 = 1 << 11
	pointerMaxSize1 = pointerMaxSize0 + (1 << 19)
	pointerMaxSize2 = pointerMaxSize1 + (1 << 27)
)

func (t Pointer) size() int {
	switch {
	case t < pointerMaxSize0:
		return 0
	case t < pointerMaxSize1:
		return 1
	case t < pointerMaxSize2:
		return 2
	default:
		return 3
	}
}

// WrittenSize is the actual total size of the pointer in the
// database data section.
func (t Pointer) WrittenSize() int64 {
	return int64(t.size() + 2)
}

func (t Pointer) typeNum() typeNum {
	return typeNumPointer
}

// WriteTo writes the value to w.
func (t Pointer) WriteTo(w writer) (int64, error) {
	size := t.size()
	switch size {
	case 0:
		err := w.WriteByte(0b00100000 | byte(0b111&(t>>8)))
		if err != nil {
			return 0, fmt.Errorf("writing pointer: %w", err)
		}
		err = w.WriteByte(byte(0xFF & t))
		if err != nil {
			return 1, fmt.Errorf("writing pointer: %w", err)
		}
	case 1:
		v := t - pointerMaxSize0
		err := w.WriteByte(0b00101000 | byte(0b111&(v>>16)))
		if err != nil {
			return 0, fmt.Errorf("writing pointer: %w", err)
		}
		err = w.WriteByte(byte(0xFF & (v >> 8)))
		if err != nil {
			return 1, fmt.Errorf("writing pointer: %w", err)
		}
		err = w.WriteByte(byte(0xFF & v))
		if err != nil {
			return 2, fmt.Errorf("writing pointer: %w", err)
		}
	case 2:
		v := t - pointerMaxSize1
		err := w.WriteByte(0b00110000 | byte(0b111&(v>>24)))
		if err != nil {
			return 0, fmt.Errorf("writing pointer: %w", err)
		}
		err = w.WriteByte(byte(0xFF & (v >> 16)))
		if err != nil {
			return 1, fmt.Errorf("writing pointer: %w", err)
		}
		err = w.WriteByte(byte(0xFF & (v >> 8)))
		if err != nil {
			return 2, fmt.Errorf("writing pointer: %w", err)
		}
		err = w.WriteByte(byte(0xFF & v))
		if err != nil {
			return 3, fmt.Errorf("writing pointer: %w", err)
		}
	case 3:
		err := w.WriteByte(0b00111000)
		if err != nil {
			return 0, fmt.Errorf("writing pointer: %w", err)
		}
		err = w.WriteByte(byte(0xFF & (t >> 24)))
		if err != nil {
			return 1, fmt.Errorf("writing pointer: %w", err)
		}
		err = w.WriteByte(byte(0xFF & (t >> 16)))
		if err != nil {
			return 2, fmt.Errorf("writing pointer: %w", err)
		}
		err = w.WriteByte(byte(0xFF & (t >> 8)))
		if err != nil {
			return 3, fmt.Errorf("writing pointer: %w", err)
		}
		err = w.WriteByte(byte(0xFF & t))
		if err != nil {
			return 4, fmt.Errorf("writing pointer: %w", err)
		}
	}
	return t.WrittenSize(), nil
}

// Slice is the MaxMind DB array type.
type Slice []DataType

// Copy makes a deep copy of the Slice.
func (t Slice) Copy() DataType {
	newSlice := make(Slice, len(t))
	for k, v := range t {
		newSlice[k] = v.Copy()
	}
	return newSlice
}

func (t Slice) size() int {
	return len(t)
}

func (t Slice) typeNum() typeNum {
	return typeNumSlice
}

// WriteTo writes the value to w.
func (t Slice) WriteTo(w writer) (int64, error) {
	numBytes, err := writeCtrlByte(w, t)
	if err != nil {
		return numBytes, err
	}

	for _, e := range t {
		written, err := w.WriteOrWritePointer(e)
		numBytes += written
		if err != nil {
			return numBytes, err
		}
	}
	return numBytes, nil
}

// String is the MaxMind DB string type.
type String string

// Copy the value.
func (t String) Copy() DataType { return t }

func (t String) size() int {
	return len(t)
}

func (t String) typeNum() typeNum {
	return typeNumString
}

// WriteTo writes the value to w.
func (t String) WriteTo(w writer) (int64, error) {
	numBytes, err := writeCtrlByte(w, t)
	if err != nil {
		return numBytes, err
	}

	written, err := w.WriteString(string(t))
	numBytes += int64(written)
	if err != nil {
		return numBytes, fmt.Errorf(`writing "%s" as a string: %w`, t, err)
	}
	return numBytes, nil
}

// Uint16 is the MaxMind DB unsigned 16-bit integer type.
type Uint16 uint16

// Copy the value.
func (t Uint16) Copy() DataType { return t }

func (t Uint16) size() int {
	return 2 - bits.LeadingZeros16(uint16(t))/8
}

func (t Uint16) typeNum() typeNum {
	return typeNumUint16
}

// WriteTo writes the value to w.
func (t Uint16) WriteTo(w writer) (int64, error) {
	numBytes, err := writeCtrlByte(w, t)
	if err != nil {
		return numBytes, err
	}

	size := t.size()
	// We ignore leading zeros
	for i := size; i > 0; i-- {
		err = w.WriteByte(byte(t >> (8 * (i - 1)) & 0xFF))
		if err != nil {
			return numBytes + int64(size-i), fmt.Errorf("writing uint16: %w", err)
		}
	}
	return numBytes + int64(size), nil
}

// Uint32 is the MaxMind DB unsigned 32-bit integer type.
type Uint32 uint32

// Copy the value.
func (t Uint32) Copy() DataType { return t }

func (t Uint32) size() int {
	return 4 - bits.LeadingZeros32(uint32(t))/8
}

func (t Uint32) typeNum() typeNum {
	return typeNumUint32
}

// WriteTo writes the value to w.
func (t Uint32) WriteTo(w writer) (int64, error) {
	numBytes, err := writeCtrlByte(w, t)
	if err != nil {
		return numBytes, err
	}

	size := t.size()
	// We ignore leading zeros
	for i := size; i > 0; i-- {
		err = w.WriteByte(byte(t >> (8 * (i - 1)) & 0xFF))
		if err != nil {
			return numBytes + int64(size-i), fmt.Errorf("writing uint32: %w", err)
		}
	}
	return numBytes + int64(size), nil
}

// Uint64 is the MaxMind DB unsigned 64-bit integer type.
type Uint64 uint64

// Copy the value.
func (t Uint64) Copy() DataType { return t }

func (t Uint64) size() int {
	return 8 - bits.LeadingZeros64(uint64(t))/8
}

func (t Uint64) typeNum() typeNum {
	return typeNumUint64
}

// WriteTo writes the value to w.
func (t Uint64) WriteTo(w writer) (int64, error) {
	numBytes, err := writeCtrlByte(w, t)
	if err != nil {
		return numBytes, err
	}

	size := t.size()

	// We ignore leading zeros
	for i := size; i > 0; i-- {
		err = w.WriteByte(byte(t >> (8 * (i - 1)) & 0xFF))
		if err != nil {
			return numBytes + int64(size-i), fmt.Errorf("writing uint64: %w", err)
		}
	}
	return numBytes + int64(size), nil
}

// Uint128 is the MaxMind DB unsigned 128-bit integer type.
type Uint128 big.Int

// Copy make a deep copy of the Uint128.
func (t *Uint128) Copy() DataType {
	nv := big.Int{}
	nv.Set((*big.Int)(t))
	uv := Uint128(nv)
	return &uv
}

func (t *Uint128) size() int {
	// We add 7 here as we want the ceiling of the division operation rather
	// than the floor.
	return ((*big.Int)(t).BitLen() + 7) / 8
}

func (t *Uint128) typeNum() typeNum {
	return typeNumUint128
}

// WriteTo writes the value to w.
func (t *Uint128) WriteTo(w writer) (int64, error) {
	numBytes, err := writeCtrlByte(w, t)
	if err != nil {
		return numBytes, err
	}

	written, err := w.Write((*big.Int)(t).Bytes())
	numBytes += int64(written)
	if err != nil {
		return numBytes, fmt.Errorf("writing uint128: %w", err)
	}
	return numBytes, nil
}

const (
	firstSize  = 29
	secondSize = firstSize + 256
	thirdSize  = secondSize + (1 << 16)
	maxSize    = thirdSize + (1 << 24)
)

func writeCtrlByte(w writer, t DataType) (int64, error) {
	size := t.size()

	typeNum := t.typeNum()

	var firstByte byte
	var secondByte byte

	if typeNum < 8 {
		firstByte = byte(typeNum << 5)
	} else {
		firstByte = byte(typeNumExtended << 5)
		secondByte = byte(typeNum - 7)
	}

	leftOver := 0
	leftOverSize := 0
	switch {
	case size < firstSize:
		firstByte |= byte(size)
	case size <= secondSize:
		firstByte |= 29
		leftOver = size - firstSize
		leftOverSize = 1
	case size <= thirdSize:
		firstByte |= 30
		leftOver = size - secondSize
		leftOverSize = 2
	case size <= maxSize:
		firstByte |= 31
		leftOver = size - thirdSize
		leftOverSize = 3
	default:
		return 0, fmt.Errorf(
			"cannot store %d bytes; max size is %d",
			size,
			maxSize,
		)
	}

	err := w.WriteByte(firstByte)
	if err != nil {
		return 0, fmt.Errorf(
			"writing first ctrl byte (type: %d, size: %d): %w",
			typeNum,
			size,
			err,
		)
	}
	numBytes := int64(1)

	if secondByte != 0 {
		err = w.WriteByte(secondByte)
		if err != nil {
			return numBytes, fmt.Errorf(
				"writing second ctrl byte (type: %d, size: %d): %w",
				typeNum,
				size,
				err,
			)
		}
		numBytes++
	}

	for i := leftOverSize - 1; i >= 0; i-- {
		v := byte((leftOver >> (8 * i)) & 0xFF)
		err = w.WriteByte(v)
		if err != nil {
			return numBytes, fmt.Errorf(
				"writing remaining ctrl bytes (type: %d, size: %d, value: %d): %w",
				typeNum,
				size,
				v,
				err,
			)
		}
		numBytes++
	}
	return numBytes, nil
}
