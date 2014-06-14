package messages

import (
	"encoding"
	"reflect"
	"time"
)

// IDToMessage provides a mapping from message ids to message types.
var IDToMessage = map[uint16]CAN{
	0x501: &MotorDriveCommand{},
	0x502: &MotorPowerCommand{},
	// 0x400: IdentificationInformation{},
	// 0x401: StatusInformation{},
	0x402: &BusMeasurement{},
	0x403: &VelocityMeasurement{},
}

// GetID returns the ID of a CAN message based on the mapping in IDToMessage.
func GetID(msg CAN) uint16 {
	if u, ok := msg.(*Unknown); ok {
		return u.CANID
	}
	msgType := reflect.TypeOf(msg)
	for id, typ := range IDToMessage {
		if reflect.TypeOf(typ) == msgType {
			return id
		}
	}
	return 0
}

// CAN describes the data stored inside messages from the CAN bus.
type CAN interface {
	encoding.BinaryUnmarshaler
	// Because types aren't first class, we use New to create new messages.
	New() CAN
}

// NewCANPlus is a convenience function to add extra info to a CAN message.
func NewCANPlus(msg CAN) CANPlus {
	return CANPlus{msg, GetID(msg), time.Now()}
}

// CANPlus is CAN with some extra stuff
type CANPlus struct {
	CAN   CAN
	CANID uint16    `json:"canID"`
	Time  time.Time `json:"time"`
}

// Unknown is used if no id is recognized.
type Unknown struct {
	// When creating an Unknown message, be sure to include the CANID so that
	// GetID will return the correct id.
	CANID uint16 `json:"-"`
	Data  [8]byte
}

func (u Unknown) New() CAN { return &Unknown{CANID: u.CANID} }
func (u *Unknown) UnmarshalBinary(b []byte) error {
	copy(u.Data[:], b)
	return nil
}
