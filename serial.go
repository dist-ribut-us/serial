// Package serial provides helper function for serializing data.
package serial

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
