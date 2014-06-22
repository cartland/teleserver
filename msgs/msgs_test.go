package msgs_test

import (
	"reflect"
	"testing"

	"github.com/calsol/teleserver/binpack"
	"github.com/calsol/teleserver/msgs"
)

func TestMessageIDs(t *testing.T) {
	for i, c := range []struct {
		msg, want msgs.CAN
		id        uint16
	}{

		{msgs.IDToMessage(0x10), &msgs.Unknown{ID: 0x10}, 0x10},

		{msgs.IDToMessage(0x501), &msgs.MotorDriveCommand{0, 0}, 0x501},

		{func() msgs.CAN {
			msg := msgs.IDToMessage(0x600)
			binpack.Unmarshal([]byte{0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc, 0xde, 0xf0}, msg)
			return msg
		}(), &msgs.MPPTStatus{0x600, "FrontRight", 0x3412, 0x7856, 0xbc9a, 0xf0de}, 0x600},

		{func() msgs.CAN {
			msg := msgs.IDToMessage(0x611)
			binpack.Unmarshal([]byte{0x01}, msg)
			return msg
		}(), &msgs.MPPTEnable{0x611, "FrontLeft", true}, 0x611},
	} {
		if !reflect.DeepEqual(c.msg, c.want) {
			t.Errorf("%d: got %#v, want %#v", i, c.msg, c.want)
		}
		if id := msgs.GetID(c.msg); id != c.id {
			t.Errorf("%d: got id %x, want %x", i, id, c.id)
		}
	}
}
