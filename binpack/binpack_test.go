package binpack_test

import (
	"reflect"
	"testing"

	"github.com/calsol/teleserver/binpack"
)

type NormalUnpack struct {
	A, B int32
}
type UnpackFloats struct {
	A float32 `binpack:"0-4"`
	B float32 `binpack:"4-8"`
}
type UnpackFloatsRev struct {
	A float32 `binpack:"4-8"`
	B float32 `binpack:"0-4"`
}
type UnpackBools struct {
	A bool `binpack:"0.1"`
	B bool `binpack:"0.0"`
}
type UnpackLargerThanNeeded struct {
	A bool `binpack:"0.1"`
	B bool `binpack:"0.0"`
}

type UnpackSparseBytes struct {
	A byte `binpack:"0"`
	B byte `binpack:"6"`
}

func TestUnmarshal(t *testing.T) {
	tests := []struct {
		name        string
		start, want interface{}
		bytes       []byte
	}{
		{
			name:  "Normal binary unmarshal",
			start: &NormalUnpack{},
			bytes: []byte{1, 2, 3, 4, 5, 6, 7, 8},
			want:  &NormalUnpack{0x4030201, 0x8070605},
		},
		{
			name:  "Unpack floats",
			start: &UnpackFloats{},
			bytes: []byte{0xcd, 0xcc, 0x44, 0x41, 0x66, 0x66, 0x36, 0x42},
			want:  &UnpackFloats{12.3, 45.6},
		},
		{
			name:  "Unpack floats backwards",
			start: &UnpackFloatsRev{},
			bytes: []byte{0xcd, 0xcc, 0x44, 0x41, 0x66, 0x66, 0x36, 0x42},
			want:  &UnpackFloatsRev{45.6, 12.3},
		},
		{
			name:  "Unpack bools",
			start: &UnpackBools{},
			bytes: []byte{0x2},
			want:  &UnpackBools{true, false},
		},
		{
			name:  "Unpack larger than needed",
			start: &UnpackLargerThanNeeded{},
			bytes: []byte{0x2, 0x3},
			want:  &UnpackLargerThanNeeded{true, false},
		},
		{
			name:  "Unpack sparse bytes",
			start: &UnpackSparseBytes{},
			bytes: []byte{0xcd, 0xcc, 0x44, 0x41, 0x66, 0x66, 0x36},
			want:  &UnpackSparseBytes{0xcd, 0x36},
		},
		{
			name:  "Unpack too many sparse bytes",
			start: &UnpackSparseBytes{},
			bytes: []byte{0xcd, 0xcc, 0x44, 0x41, 0x66, 0x66, 0x36, 0x42},
			want:  &UnpackSparseBytes{0xcd, 0x36},
		},
		{
			name:  "Unpack too few sparse bytes",
			start: &UnpackSparseBytes{},
			bytes: []byte{0xcd, 0xcc, 0x44, 0x41, 0x66, 0x66},
			want:  &UnpackSparseBytes{0xcd, 0x00},
		},
	}

	for _, c := range tests {
		got := c.start
		if err := binpack.Unmarshal(c.bytes, got); err != nil {
			t.Errorf("%s: %v", c.name, err)
		}
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("%s: got %v, want %v", c.name, got, c.want)
		}
	}
}
