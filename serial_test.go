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

func TestBoolSlice(t *testing.T) {
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

func TestReturn(t *testing.T) {
	assert.Equal(t, uint32(12345), UnmarshalUint32(MarshalUint32(12345, nil)))
}

func TestByteSlices(t *testing.T) {
	tt := []struct {
		name       string
		data       [][]byte
		prefixLens []int
	}{
		{
			name: "Basic",
			data: [][]byte{
				{1, 2, 3, 4},
				{5, 6, 7, 8},
				{9, 10, 11, 12},
			},
			prefixLens: []int{2, 4, 8},
		},
		{
			name: "Zero",
			data: [][]byte{
				{1, 2, 3, 4},
				{5, 6, 7, 8},
				{9, 10, 11, 12},
			},
			prefixLens: []int{2, 4, 0},
		},
		{
			name: "Negative",
			data: [][]byte{
				{1, 2, 3, 4},
				{5, 6, 7, 8},
				{9, 10, 11, 12},
			},
			prefixLens: []int{2, -4, 2},
		},
		{
			name: "Long",
			data: [][]byte{
				{1, 2, 3, 4},
				{1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4},
				{9, 10, 11, 12},
			},
			prefixLens: []int{2, 2, 0},
		},
		{
			name: "One byte zero",
			data: [][]byte{
				{1, 2, 3, 4},
				{1, 2, 3},
				{0},
			},
			prefixLens: []int{2, 2, 0},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			b, err := MarshalByteSlices(tc.prefixLens, tc.data)
			assert.NoError(t, err)
			data, err := UnmarshalByteSlices(tc.prefixLens, b)
			assert.NoError(t, err)
			assert.Equal(t, tc.data, data)

			p := ByteSlicesPrefixer(tc.prefixLens)
			b, err = p.Marshal(tc.data)
			assert.NoError(t, err)
			data, err = p.Unmarshal(b)
			assert.NoError(t, err)
			assert.Equal(t, tc.data, data)
		})
	}
}

func TestSlicesPacker(t *testing.T) {
	tt := []struct {
		name         string
		data         [][]byte
		slicesPacker SlicesPacker
	}{
		{
			name: "Basic",
			data: [][]byte{
				{1, 2, 3, 4},
				{5, 6, 7, 8},
				{9, 10, 11, 12},
			},
			slicesPacker: SlicesPacker{2, 2},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			b, err := tc.slicesPacker.Marshal(tc.data)
			assert.NoError(t, err)
			data, err := tc.slicesPacker.Unmarshal(b)
			assert.NoError(t, err)
			assert.Equal(t, tc.data, data)
		})
	}
}
