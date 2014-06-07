package messages

import "encoding"

// IDToMessage provides a mapping from message ids to message types
var IDToMessage = map[uint16]CAN{
	0x501: &MotorDriveCommand{},
	0x502: &MotorPowerCommand{},
	// 0x400: IdentificationInformation{},
	// 0x401: StatusInformation{},
	0x402: &BusMeasurement{},
	0x403: &VelocityMeasurement{},
}

// CAN describes the data stored inside messages from the CAN bus.
type CAN interface {
	encoding.BinaryUnmarshaler
	// Because types aren't first class, we use New to create new messages.
	New() CAN
}
