package msgs

import (
	"encoding"
	"reflect"
	"time"
)

var idToMessage = map[uint16]CAN{
	0x402: &BusMeasurement{},
	0x403: &VelocityMeasurement{},
	0x501: &MotorDriveCommand{},
	0x502: &MotorPowerCommand{},
	0x600: &MPPTStatus{ID: 0x600, ArrayLocation: "Front Right"},
	0x601: &MPPTStatus{ID: 0x601, ArrayLocation: "Front Left"},
	0x602: &MPPTStatus{ID: 0x602, ArrayLocation: "Back Right"},
	0x603: &MPPTStatus{ID: 0x603, ArrayLocation: "Back Left"},
	0x610: &MPPTEnable{ID: 0x610, ArrayLocation: "Front Right"},
	0x611: &MPPTEnable{ID: 0x611, ArrayLocation: "Front Left"},
	0x612: &MPPTEnable{ID: 0x612, ArrayLocation: "Back Right"},
	0x613: &MPPTEnable{ID: 0x613, ArrayLocation: "Back Left"},
}

// IDToMessage provides a mapping from message ids to message types.
func IDToMessage(id uint16) CAN {
	if msg, ok := idToMessage[id]; ok {
		// Copy the message so that we don't modify the map
		return msg.New()
	}
	return &Unknown{ID: id}
}

// CAN describes the data stored inside messages from the CAN bus.
type CAN interface {
	encoding.BinaryUnmarshaler
	// New will create a new message, preserving id and other meta information.
	New() CAN
}

func GetID(msg CAN) uint16 {
	if msg, ok := msg.(ider); ok {
		return msg.canID()
	}
	msgType := reflect.TypeOf(msg)
	for id, typ := range idToMessage {
		if reflect.TypeOf(typ) == msgType {
			return id
		}
	}
	return 0
}

// NewCANPlus is a convenience function to add extra info to a CAN message.
func NewCANPlus(msg CAN) CANPlus {
	return CANPlus{msg, GetID(msg), time.Now()}
}

type ider interface {
	canID() uint16
}

// CANPlus is CAN with some extra stuff
type CANPlus struct {
	CAN   CAN
	CANID uint16    `json:"canID"`
	Time  time.Time `json:"time"`
}

// Unknown is used if no id is recognized.
type Unknown struct {
	ID   uint16 `json:"-"`
	Data [8]byte
}

func (u Unknown) New() CAN      { return &Unknown{ID: u.ID} }
func (u Unknown) canID() uint16 { return u.ID }
func (u *Unknown) UnmarshalBinary(b []byte) error {
	copy(u.Data[:], b)
	return nil
}
