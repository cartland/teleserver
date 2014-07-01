// Package binpack packs slices of bytes into structs according to the tags on
// the structs. If there are no tags, it tries reading directly into the struct
// with a little endian binary.Read
package binpack

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

const tagname = "binpack"

var byteorder = binary.LittleEndian

type BinpackUnmarshaler interface {
	UnmarshalBinpack(data []byte) error
}

func binpackSize(v reflect.Value) int {
	size := 0
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if tag := f.Tag.Get(tagname); tag != "" {
			if bytes := strings.Split(tag, "."); len(bytes) > 1 {
				if s, err := strconv.Atoi(bytes[0]); err == nil && s+1 > size {
					size = s + 1
				}
			}
			bytes := strings.Split(tag, "-")
			if s, err := strconv.Atoi(bytes[len(bytes)-1]); err == nil && s > size {
				size = s
			}
		}
	}
	return size
}

// Unmarshal unmarshals the binaries according to the struct tags.
// Example tags:
//    "0" - Take the first byte
//    "0-4" - Take the first 4 bytes
//    "0.2 - Take the third bit of the first byte
func Unmarshal(data []byte, v interface{}) error {
	if m, ok := v.(BinpackUnmarshaler); ok {
		return m.UnmarshalBinpack(data)
	}

	if d := reflect.ValueOf(v); d.Kind() == reflect.Ptr {
		if rv := d.Elem(); rv.Kind() == reflect.Struct {
			if s := binpackSize(rv); s > 0 {
				dec := &decoder{data}
				return dec.value(rv)
			}
		}
	}

	return binary.Read(bytes.NewReader(data), byteorder, v)
}

type decoder struct{ data []byte }

func (d *decoder) value(v reflect.Value) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if tag := f.Tag.Get(tagname); tag != "" {
			if err := d.decodeTag(v.Field(i), tag); err != nil {
				return err
			}
		}
	}
	return nil
}
func (d *decoder) decodeTag(v reflect.Value, tag string) error {
	field := v.Addr().Interface()
	var start, end int
	var bit uint
	if s, err := strconv.Atoi(tag); err == nil && inRange(d.data, s, s+1) {
		return binary.Read(bytes.NewReader(d.data[s:s+1]), binary.LittleEndian, field)
	} else if n, err := fmt.Sscanf(tag, "%d-%d", &start, &end); n > 0 && err == nil && inRange(d.data, start, end) {
		return binary.Read(bytes.NewReader(d.data[start:end]), binary.LittleEndian, field)
	} else if n, err := fmt.Sscanf(tag, "%d.%d", &start, &bit); n > 0 && err == nil && start < len(d.data) {
		v.SetBool(d.data[start]&(1<<bit) != 0)
	}
	return nil
}

func inRange(s []byte, a, b int) bool {
	return (0 <= a) && (a <= b) && (b <= len(s))
}
