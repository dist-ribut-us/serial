// Package serial provides helper function for serializing data.
package serial

import (
	"github.com/dist-ribut-us/errors"
)

// MarshalUint32 takes a unit and a byte slice and writes the uint32 to the
// first 4 bytes of the slice. It does not check the slice length and will panic
// if the slice does not have a length of at least 4.
func MarshalUint32(ui uint32, b []byte) []byte { return marshalUint(uint64(ui), 4, b) }

// UnmarshalUint32 reads the first 4 bytes of a slice into a uint32. It does not
// check the slice length and will panic if the slice does not have a length of
// at least 4.
func UnmarshalUint32(b []byte) uint32 { return uint32(unmarshalUint(4, b)) }

// MarshalUint16 takes a unit and a byte slice and writes the uint16 to the
// first 2 bytes of the slice. It does not check the slice length and will panic
// if the slice does not have a length of at least 2.
func MarshalUint16(ui uint16, b []byte) []byte { return marshalUint(uint64(ui), 2, b) }

// UnmarshalUint16 reads the first 2 bytes of a slice into a uint16. It does not
// check the slice length and will panic if the slice does not have a length of
// at least 2.
func UnmarshalUint16(b []byte) uint16 { return uint16(unmarshalUint(2, b)) }

func marshalUint(ui uint64, l int, b []byte) []byte {
	if b == nil {
		b = make([]byte, l)
	}
	i := 0
	for ; i < l && ui > 0; i++ {
		b[i] = byte(ui)
		ui >>= 8
	}
	for ; i < l; i++ {
		b[i] = 0
	}
	return b
}

func unmarshalUint(l int, b []byte) uint64 {
	var ui uint64
	i := l - 1
	for {
		ui += uint64(b[i])
		if i == 0 {
			return ui
		}
		i--
		ui <<= 8
	}
}

// MarshalBoolSlice marshals a slice of bools into a byte slice. Each bool takes
// only one bit. It does not check the slice length and will panic if the slice
// is not long enough. The required length with be the length of the bool slice
// divided by 8 rounded up plus 4.
func MarshalBoolSlice(bls []bool, b []byte) {
	ln := len(bls)
	MarshalUint32(uint32(ln), b)
	var bt byte
	i := 4
	for j, bl := range bls {
		bt <<= 1
		if bl {
			bt |= 1
		}
		if (j+1)%8 == 0 {
			b[i] = bt
			bt = 0
			i++
		}
	}
	if os := ln % 8; os != 0 {
		// fix offset
		for ; os < 8; os++ {
			bt <<= 1
		}
		b[i] = bt
	}
}

// UnmarshalBoolSlice will unmarshal a slice of bools from a byte slice. If the
// byte slice is malformed, it may panic.
func UnmarshalBoolSlice(b []byte) []bool {
	bls := make([]bool, int(UnmarshalUint32(b)))
	i := 4
	var bt byte
	for j := range bls {
		if j%8 == 0 {
			bt = b[i]
			i++
		}
		bls[j] = bt >= 128
		bt <<= 1
	}
	return bls
}

const (
	// ErrLengthsDoNotMatch is returned when a ByteSlicesPrefixer is given data
	// with the wrong number of byte slices.
	ErrLengthsDoNotMatch = errors.String("Lengths do not match")
	// ErrLengthTooLong is returned if a ByteSlicesPrefixer length is above 8.
	ErrLengthTooLong = errors.String("Length cannot exceed 8 bytes")
	// ErrIncorrectZero is returned when a BySlicesPrefixer has a zero in any but
	// the last position.
	ErrIncorrectZero = errors.String("Zero is only valid as final length")
	// ErrBadFormat is returned during ByteSlicesPrefixer Unmarshal if the format
	// does not match the description.
	ErrBadFormat = errors.String("Not properly formatted")
)

// ByteSlicesPrefixer describes how to prefix a [][]byte for serialization.
type ByteSlicesPrefixer []int

// Marshal uses the ByteSlicesPrefixer to marshal a [][]byte and insert length
// information.
func (pre ByteSlicesPrefixer) Marshal(data [][]byte) ([]byte, error) {
	if len(pre) != len(data) {
		return nil, ErrLengthsDoNotMatch
	}
	sum := 0
	for i, ln := range pre {
		if ln > 8 {
			return nil, ErrLengthTooLong
		}
		if ln > 0 {
			sum += ln
		}
		sum += len(data[i])
	}
	idx := 0
	b := make([]byte, sum)
	for i, ln := range pre {
		if ln > 0 {
			marshalUint(uint64(len(data[i])), ln, b[idx:])
			idx += ln
		}
		copy(b[idx:], data[i])
		idx += len(data[i])
	}
	return b, nil
}

// Unmarshal uses the ByteSlicesPrefixer to unmarshal a []byte int a [][]byte by
// extracting length information.
func (pre ByteSlicesPrefixer) Unmarshal(b []byte) ([][]byte, error) {
	data := make([][]byte, len(pre))
	idx := 0
	ln := 0
	for i, pln := range pre {
		if pln > 8 {
			return nil, ErrLengthTooLong
		}
		if pln > 0 {
			if idx+pln > len(b) {
				return nil, ErrBadFormat
			}
			ln = int(unmarshalUint(pln, b[idx:]))
			idx += pln
		} else if pln < 0 {
			ln = -pln
		} else {
			if i != len(pre)-1 {
				return nil, ErrIncorrectZero
			}
			data[i] = b[idx:]
			break
		}

		if idx+ln > len(b) {
			return nil, ErrBadFormat
		}
		data[i] = b[idx : idx+ln]
		idx += ln
	}
	return data, nil
}

// MarshalByteSlices takes a slice of byte slices and marshals them into a
// a single byte slice. The prefixLens determine how many bytes to use for
// length headers. Positive values set the length in bytes, values less than or
// equal to 0 will result in no prefix being added.
func MarshalByteSlices(prefixLens []int, data [][]byte) ([]byte, error) {
	if len(prefixLens) != len(data) {
		return nil, ErrLengthsDoNotMatch
	}
	sum := 0
	for i, ln := range prefixLens {
		if ln > 8 {
			return nil, ErrLengthTooLong
		}
		if ln > 0 {
			sum += ln
		}
		sum += len(data[i])
	}
	idx := 0
	b := make([]byte, sum)
	for i, ln := range prefixLens {
		if ln > 0 {
			marshalUint(uint64(len(data[i])), ln, b[idx:])
			idx += ln
		}
		copy(b[idx:], data[i])
		idx += len(data[i])
	}
	return b, nil
}

// UnmarshalByteSlices takes a byte slice and breaks it into a slice of byte
// slices. The prefixes determine what length headers to expect. Postive values
// are interpreted as the number of bytes to read as the length header. Negative
// values indicate the absolute length (so -5 means the slice is 5 bytes long)
// and 0 is only valid for the final value and means consume the rest.
func UnmarshalByteSlices(prefixLens []int, b []byte) ([][]byte, error) {
	data := make([][]byte, len(prefixLens))
	idx := 0
	ln := 0
	for i, pln := range prefixLens {
		if pln > 8 {
			return nil, ErrLengthTooLong
		}
		if pln > 0 {
			if idx+pln > len(b) {
				return nil, ErrBadFormat
			}
			ln = int(unmarshalUint(pln, b[idx:]))
			idx += pln
		} else if pln < 0 {
			ln = -pln
		} else {
			if i != len(prefixLens)-1 {
				return nil, ErrIncorrectZero
			}
			data[i] = b[idx:]
			break
		}

		if idx+ln > len(b) {
			return nil, ErrBadFormat
		}
		data[i] = b[idx : idx+ln]
		idx += ln
	}
	return data, nil
}

// SlicesPacker describes the number of bytes to use to describe the count of
// slices and Size describes the number of bytes to use to describe the size of
// each slice
type SlicesPacker struct {
	Count int
	Size  int
}

// Marshal a [][]byte into []byte inserting a header with total number of slices
// and prefixing the length of each slice.
func (s SlicesPacker) Marshal(data [][]byte) ([]byte, error) {
	if s.Count > 8 || s.Size > 8 {
		return nil, ErrLengthTooLong
	}

	l := s.Count
	for _, d := range data {
		l += s.Size + len(d)
	}

	b := make([]byte, s.Count, l)
	marshalUint(uint64(len(data)), s.Count, b)

	buf := make([]byte, s.Size)
	for _, d := range data {
		marshalUint(uint64(len(d)), s.Size, buf)
		b = append(b, buf...)
		b = append(b, d...)
	}

	return b, nil
}

// Unmarshal a []byte to a [][]byte using the SlicesPacker to extract the total
// count and the size of each.
func (s SlicesPacker) Unmarshal(data []byte) ([][]byte, error) {
	if s.Count > 8 || s.Size > 8 {
		return nil, ErrLengthTooLong
	}

	c := unmarshalUint(s.Count, data)
	data = data[s.Count:]

	b := make([][]byte, c)
	for i := range b {
		sz := unmarshalUint(s.Size, data)
		data = data[s.Size:]
		b[i] = data[:sz]
		data = data[sz:]
	}

	return b, nil
}
