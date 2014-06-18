package messages

import (
	"encoding/binary"
	"fmt"
)

// The MPPTStatus message is sent from the MPPTs at a regular interval
type MPPTStatus struct {
	// id holds the CAN ID of the message
	id uint16 `json:"-"`
	// Array location is a human-readable description of the array location.
	ArrayLocation string
	// Array voltage ís scaled by 100, or 1 count = 10mV
	ArrayVoltage uint16
	// Array current ís scaled by 1000, or 1 count = 1mA
	ArrayCurrent uint16
	// Battery voltage is scaled by 100, or 1 count = 10mV
	BatteryVoltage uint16
	// Temperature is scaled by 100, or 1 count = 10mC
	Temperature uint16
}

func (m MPPTStatus) New() CAN      { return &MPPTStatus{id: m.id, ArrayLocation: m.ArrayLocation} }
func (m MPPTStatus) canID() uint16 { return m.id }
func (m *MPPTStatus) UnmarshalBinary(b []byte) error {
	if len(b) != 8 {
		return fmt.Errorf("data is %d bytes, need %d", len(b), 8)
	}
	m.ArrayVoltage = binary.LittleEndian.Uint16(b[0:2])
	m.ArrayCurrent = binary.LittleEndian.Uint16(b[2:4])
	m.BatteryVoltage = binary.LittleEndian.Uint16(b[4:6])
	m.Temperature = binary.LittleEndian.Uint16(b[6:8])
	return nil
}

// The MPPTEnable message is sent to the MPPTs to turn them on or off
type MPPTEnable struct {
	// id holds the CAN ID of the message
	id uint16 `json:"-"`
	// Array location is a human-readable description of the array location
	ArrayLocation string
	// The message will either enable or disable the power point trackets
	Enable bool
}

func (m MPPTEnable) New() CAN      { return &MPPTEnable{ArrayLocation: m.ArrayLocation} }
func (m MPPTEnable) canID() uint16 { return m.id }
func (m *MPPTEnable) UnmarshalBinary(b []byte) error {
	if len(b) != 1 {
		return fmt.Errorf("data is %d bytes, need %d", len(b), 1)
	}
	m.Enable = b[0]&0x1 == 0x1
	return nil
}
