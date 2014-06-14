package lib_test

import (
	"bytes"
	"encoding/binary"
	"reflect"
	"strings"
	"testing"

	"github.com/calsol/teleserver/can"
	"github.com/calsol/teleserver/lib"
	"github.com/calsol/teleserver/messages"
)

func TestXSPDecode(t *testing.T) {
	tests := []struct {
		data string
		want []string
	}{
		{
			data: "abc",
			want: []string{"abc"},
		},
		{
			data: "sdfa\xe7giu\x92e43\xe7jfeiu\x00o3q\xe7aau\x14a",
			want: []string{"sdfa", "gi\xe7e43", "jfeiuo3q", "aaaa"},
		},
		{
			data: "\xe7giu\x92e43\xe7jfeiu\x00o3q\xe7aau\x14a\xe7",
			want: []string{"gi\xe7e43", "jfeiuo3q", "aaaa"},
		},
	}

	for i, c := range tests {
		scanner := lib.NewXSPScanner(strings.NewReader(c.data))
		for j, want := range c.want {
			if !scanner.Scan() {
				t.Errorf("%d: failed on scan %d", i, j)
			}
			if want != scanner.Text() {
				t.Errorf("%d: got %s, want %s", i, scanner.Text(), want)
			}
		}
		if err := scanner.Err(); err != nil {
			t.Error(err)
		}
	}
}

type msgAndErr struct {
	msg messages.CAN
	err string
}

func TestReadCAN(t *testing.T) {
	tests := []struct {
		data string
		want []msgAndErr
	}{
		{
			data: "bad",
			want: []msgAndErr{
				{err: "packet 0x616: payload size 2 != actual size 1: [100]"},
			},
		},
		{
			data: "bad\xe7beef",
			want: []msgAndErr{
				{err: "packet 0x616: payload size 2 != actual size 1: [100]"},
				{err: "packet 0x656: unknown id, size 2: [101 102]"},
			},
		},
		{
			data: "\x18\x50" + "\xcd\xcc\x44\x41\x66\x66\x36\x42",
			want: []msgAndErr{
				{msg: &messages.MotorDriveCommand{MotorCurrent: 12.3, MotorVelocity: 45.6}},
			},
		},
		{
			data: "\x28\x50" + "\x66\x66\x36\x42\xcd\xcc\x44\x41" + "\xe7bad\xe7" +
				"\x38\x40" + "\x00\x80\xad\x43\x41\xb1\x2d\x42",
			want: []msgAndErr{
				{msg: &messages.MotorPowerCommand{BusCurrent: 12.3}},
				{err: "packet 0x616: payload size 2 != actual size 1: [100]"},
				{msg: &messages.VelocityMeasurement{MotorVelocity: 347, VehicleVelocity: 43.4231}},
			},
		},
	}

	for i, c := range tests {
		r := lib.NewXSPCANReader(strings.NewReader(c.data))
		for j, want := range c.want {
			msg, err := r.Read()
			if err != nil && err.Error() != want.err {
				t.Errorf("%d: on %d got error %q, want %q", i, j, err, want.err)
			}
			if !reflect.DeepEqual(msg, want.msg) {
				t.Errorf("%d: on %d got %v, want %v", i, j, msg, want.msg)
			}
		}
	}
}

func TestReadSocketCAN(t *testing.T) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, can.NewFrame(0x502, []byte{0x66, 0x66, 0x36, 0x42, 0xcd, 0xcc, 0x44, 0x41}))
	binary.Write(buf, binary.LittleEndian, can.NewFrame(0x403, []byte{0x00, 0x80, 0xad, 0x43, 0x41, 0xb1, 0x2d, 0x42}))
	test4Data := buf.String()

	tests := []struct {
		data string
		want []msgAndErr
	}{
		{
			data: "bad",
			want: []msgAndErr{
				{err: "got 3 bytes, want 16"},
			},
		},
		{
			data: "this is 16 chars",
			want: []msgAndErr{
				{err: "packet 0x73696874: payload size 8 is greater than 32: [49 54 32 99 104 97 114 115]"},
				{err: "EOF"},
			},
		},
		{
			data: "thi\x00s is more than 16 chars",
			want: []msgAndErr{
				{err: "packet 0x696874: payload size 8 is greater than 115: [32 109 111 114 101 32 116 104]"},
				{err: "got 11 bytes, want 16"},
			},
		},
		{
			data: "\x01\x05\x00\x00" + "\x08" + "\x00\x00\x00" + "\xcd\xcc\x44\x41\x66\x66\x36\x42",
			want: []msgAndErr{
				{msg: &messages.MotorDriveCommand{MotorCurrent: 12.3, MotorVelocity: 45.6}},
			},
		},
		{
			data: test4Data,
			want: []msgAndErr{
				{msg: &messages.MotorPowerCommand{BusCurrent: 12.3}},
				{msg: &messages.VelocityMeasurement{MotorVelocity: 347, VehicleVelocity: 43.4231}},
			},
		},
	}

	for i, c := range tests {
		r := lib.NewSocketCANReader(strings.NewReader(c.data))
		for j, want := range c.want {
			msg, err := r.Read()
			if err != nil && err.Error() != want.err {
				t.Errorf("%d: on %d got error %q, want %q", i, j, err, want.err)
			}
			if !reflect.DeepEqual(msg, want.msg) {
				t.Errorf("%d: on %d got %v, want %v", i, j, msg, want.msg)
			}
		}
	}
}
