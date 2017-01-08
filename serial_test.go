package serial

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

func TestUint32Marshal(t *testing.T) {
	ints := []uint32{0, 1, 31415, 4294967295}
	for _, i := range ints {
		b := make([]byte, 4)
		MarshalUint32(i, b)
		assert.Equal(t, i, UnmarshalUint32(b))
	}
}

func TestUint16Marshal(t *testing.T) {
	ints := []uint16{0, 1, 31415, 65535}
	for _, i := range ints {
		b := make([]byte, 4)
		MarshalUint16(i, b)
		assert.Equal(t, i, UnmarshalUint16(b))
	}
}

func TestByteSlice(t *testing.T) {
	maxLen := 500
	b := make([]byte, (maxLen/8)+5)
	for i := 0; i < 100; i++ {
		l := rand.Int() % maxLen
		bls := make([]bool, l)
		for j := range bls {
			bls[j] = rand.Int()%2 == 0
		}
		MarshalBoolSlice(bls, b)
		if !assert.Equal(t, bls, UnmarshalBoolSlice(b)) {
			return
		}
	}
}
