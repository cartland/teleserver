package msgs

import (
	"fmt"
	"reflect"
	"time"

	"github.com/calsol/teleserver/binpack"
	"github.com/calsol/teleserver/can"
)

var idToMessage = map[uint16]CAN{
	0x030: &BatteryCarShutdown{},
	0x040: &BatteryHeartbeat{},
	0x041: &CutoffHeartbeat{},
	0x042: &MotorHeartbeat{ID: 0x042, Location: "Left"},
	0x043: &MotorHeartbeat{ID: 0x043, Location: "Right"},
	0x044: &DashHeartbeat{},
	0x045: &PowerHubHeartbeat{ID: 0x045, Location: "Bottom"},
	0x046: &PowerHubHeartbeat{ID: 0x046, Location: "Top"},
	0x123: &PackVoltage{},
	0x124: &PackCurrent{},
	0x125: &PackTemperature{},
	0x128: &BMSBalancing{},
	0x129: &BMSState{},
	0x130: &BatteryModule{ID: 0x130},
	0x131: &BatteryModule{ID: 0x131},
	0x132: &BatteryModule{ID: 0x132},
	0x140: &BatteryModule{ID: 0x140},
	0x141: &BatteryModule{ID: 0x141},
	0x142: &BatteryModule{ID: 0x142},
	0x150: &BatteryModule{ID: 0x150},
	0x151: &BatteryModule{ID: 0x151},
	0x152: &BatteryModule{ID: 0x152},
	0x260: &CutoffTrigger{},
	0x261: &CutoffAnalogIn{},
	0x262: &CutoffSPIIn{},
	0x280: &CANAccelPos{},
	0x281: &CANBrakePos{},
	0x310: &MotorRPM{ID: 0x310, Location: "Left"},
	0x311: &MotorRPM{ID: 0x311, Location: "Right"},

	// Tritium commands, not relevant for Zephyr
	0x402: &BusMeasurement{},
	0x403: &VelocityMeasurement{},
	0x501: &MotorDriveCommand{},
	0x502: &MotorPowerCommand{},

	0x600: &MPPTStatus{ID: 0x600, ArrayLocation: "FrontRight"},
	0x601: &MPPTStatus{ID: 0x601, ArrayLocation: "FrontLeft"},
	0x602: &MPPTStatus{ID: 0x602, ArrayLocation: "BackRight"},
	0x603: &MPPTStatus{ID: 0x603, ArrayLocation: "BackLeft"},
	0x610: &MPPTEnable{ID: 0x610, ArrayLocation: "FrontRight"},
	0x611: &MPPTEnable{ID: 0x611, ArrayLocation: "FrontLeft"},
	0x612: &MPPTEnable{ID: 0x612, ArrayLocation: "BackRight"},
	0x613: &MPPTEnable{ID: 0x613, ArrayLocation: "BackLeft"},
	0x700: &PanelSwitchPos{},
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

// NewCAN parses a CAN message to turn it into something semantically meaningful.
func NewCAN(m can.Message) (CAN, error) {
	id, body := m.GetID(), m.GetData()
	msg := IDToMessage(uint16(id))
	if err := binpack.Unmarshal(body, msg); err != nil {
		return nil, fmt.Errorf("packet 0x%x: payload %v: %v", id, body, err)
	}
	return msg, nil
}

// NewCANPlus is a convenience function to add extra info to a CAN message.
func NewCANPlus(msg CAN) CANPlus {
	return CANPlus{msg, GetID(msg), time.Now().UTC()}
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
	ID   uint16  `json:"-"`
	Data [8]byte `binpack:"0-8"`
}

func (u Unknown) New() CAN      { return &Unknown{ID: u.ID} }
func (u Unknown) canID() uint16 { return u.ID }

// We have a special function for Unknown so that it can deal with messages of
// multiple lengths.
func (u *Unknown) UnmarshalBinpack(b []byte) error {
	copy(u.Data[:], b)
	return nil
}
