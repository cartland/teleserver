package lib_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/calsol/teleserver/lib"
	"github.com/calsol/teleserver/messages"
)

func TestDemuxDecode(t *testing.T) {
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
