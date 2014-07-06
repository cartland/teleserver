package lib_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/calsol/teleserver/can"
	"github.com/calsol/teleserver/lib"
	"github.com/calsol/teleserver/msgs"
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
	msg msgs.CAN
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
				{msg: &msgs.Unknown{ID: 0x656, Data: [8]byte{101, 102, 0, 0, 0, 0, 0, 0}}},
			},
		},
		{
			data: "\x18\x50" + "\xcd\xcc\x44\x41\x66\x66\x36\x42",
			want: []msgAndErr{
				{msg: &msgs.MotorDriveCommand{MotorCurrent: 12.3, MotorVelocity: 45.6}},
			},
		},
		{
			data: "\x28\x50" + "\x66\x66\x36\x42\xcd\xcc\x44\x41" + "\xe7bad\xe7" +
				"\x38\x40" + "\x00\x80\xad\x43\x41\xb1\x2d\x42",
			want: []msgAndErr{
				{msg: &msgs.MotorPowerCommand{BusCurrent: 12.3}},
				{err: "packet 0x616: payload size 2 != actual size 1: [100]"},
				{msg: &msgs.VelocityMeasurement{MotorVelocity: 347, VehicleVelocity: 43.4231}},
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
				t.Errorf("%d: on %d got %#v, want %#v", i, j, msg, want.msg)
			}
		}
	}
}

// Conn holds a connection to the CAN socket.
type fakeConn chan *can.Frame

func makeFakeConn(data []*can.Frame) can.Conn {
	f := make(fakeConn, len(data))
	for _, d := range data {
		f <- d
	}
	close(f)
	return f
}
func (f fakeConn) ReadFrame() (*can.Frame, error)  { return <-f, nil }
func (f fakeConn) WriteFrame(*can.Frame) error     { return nil }
func (f fakeConn) Read([]byte) (n int, err error)  { return }
func (f fakeConn) Write([]byte) (n int, err error) { return }
func (f fakeConn) Close() error                    { return nil }

func TestReadSocketCAN(t *testing.T) {
	tests := []struct {
		data []*can.Frame
		want []msgAndErr
	}{
		{
			data: []*can.Frame{
				can.NewFrame(can.Simple{0x501, []byte{0x66, 0x66, 0x36, 0x42, 0xcd, 0xcc, 0x44, 0x41}}),
			},
			want: []msgAndErr{
				{msg: &msgs.MotorDriveCommand{MotorCurrent: 45.6, MotorVelocity: 12.3}},
			},
		},
		{
			data: []*can.Frame{
				can.NewFrame(can.Simple{0x502, []byte{0x66, 0x66, 0x36, 0x42, 0xcd, 0xcc, 0x44, 0x41}}),
				can.NewFrame(can.Simple{0x403, []byte{0x00, 0x80, 0xad, 0x43, 0x41, 0xb1, 0x2d, 0x42}}),
			},
			want: []msgAndErr{
				{msg: &msgs.MotorPowerCommand{BusCurrent: 12.3}},
				{msg: &msgs.VelocityMeasurement{MotorVelocity: 347, VehicleVelocity: 43.4231}},
			},
		},
	}

	for i, c := range tests {
		r := lib.NewSocketCANReader(makeFakeConn(c.data))
		for j, want := range c.want {
			msg, err := r.Read()
			if err != nil && err.Error() != want.err {
				t.Errorf("%d: on %d got error %q, want %q", i, j, err, want.err)
			}
			if !reflect.DeepEqual(msg, want.msg) {
				t.Errorf("%d: on %d got %#v, want %#v", i, j, msg, want.msg)
			}
		}
	}
}
