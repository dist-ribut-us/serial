package serial

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUint32Marshal(t *testing.T) {
	ints := []uint32{0, 1, 31415}
	for _, i := range ints {
		b := make([]byte, 4)
		MarshalUint32(i, b)
		assert.Equal(t, i, UnmarshalUint32(b))
	}
}

func TestUint16Marshal(t *testing.T) {
	ints := []uint16{0, 1, 31415}
	for _, i := range ints {
		b := make([]byte, 4)
		MarshalUint16(i, b)
		assert.Equal(t, i, UnmarshalUint16(b))
	}
}
