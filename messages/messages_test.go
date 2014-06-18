package messages

import (
	"reflect"
	"testing"
)

func TestMessageIDs(t *testing.T) {
	for i, c := range []struct {
		msg, want CAN
		id        uint16
	}{

		{IDToMessage(0x10), &Unknown{id: 0x10}, 0x10},

		{IDToMessage(0x501), &MotorDriveCommand{0, 0}, 0x501},

		{func() CAN {
			msg := IDToMessage(0x600)
			msg.UnmarshalBinary([]byte{0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc, 0xde, 0xf0})
			return msg
		}(), &MPPTStatus{0x600, "Front Right", 0x3412, 0x7856, 0xbc9a, 0xf0de}, 0x600},

		{func() CAN {
			msg := IDToMessage(0x611)
			msg.UnmarshalBinary([]byte{0x01})
			return msg
		}(), &MPPTEnable{0x611, "Front Left", true}, 0x611},
	} {
		if !reflect.DeepEqual(c.msg, c.want) {
			t.Errorf("%d: got %#v, want %#v", i, c.msg, c.want)
		}
		if id := GetID(c.msg); id != c.id {
			t.Errorf("%d: got id %x, want %x", i, id, c.id)
		}
	}
}
