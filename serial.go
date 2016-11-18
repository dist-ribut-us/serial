// Package serial provides helper function for serializing data.
package serial

// MarshalUint32 takes a unit and a byte slice and writes the uint32 to the
// first 4 bytes of the slice. It does not check the slice length and will panic
// if the slice does not have a length of at least 4.
func MarshalUint32(ui uint32, b []byte) { marshalUint(uint64(ui), 4, b) }

// UnmarshalUint32 reads the first 4 bytes of a slice into a uint32. It does not
// check the slice length and will panic if the slice does not have a length of
// at least 4.
func UnmarshalUint32(b []byte) uint32 { return uint32(unmarshalUint(4, b)) }

// MarshalUint16 takes a unit and a byte slice and writes the uint16 to the
// first 2 bytes of the slice. It does not check the slice length and will panic
// if the slice does not have a length of at least 2.
func MarshalUint16(ui uint16, b []byte) { marshalUint(uint64(ui), 2, b) }

// UnmarshalUint16 reads the first 2 bytes of a slice into a uint16. It does not
// check the slice length and will panic if the slice does not have a length of
// at least 2.
func UnmarshalUint16(b []byte) uint16 { return uint16(unmarshalUint(2, b)) }

func marshalUint(ui uint64, l int, b []byte) {
	for i := 0; i < l; i++ {
		b[i] = byte(ui)
		ui >>= 8
	}
}

func unmarshalUint(l int, b []byte) uint64 {
	var ui uint64
	for i := l - 1; i >= 0; i-- {
		ui <<= 8
		ui += uint64(b[i])
	}
	return ui
}
