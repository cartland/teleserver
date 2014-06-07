package messages

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

func twoFloats(b []byte) (float32, float32, error) {
	if len(b) != 8 {
		return 0, 0, fmt.Errorf("data is %d bytes, need %d", len(b), 8)
	}
	var f1, f2 float32
	if err := binary.Read(bytes.NewReader(b[:4]), binary.LittleEndian, &f1); err != nil {
		return 0, 0, err
	}
	if err := binary.Read(bytes.NewReader(b[4:]), binary.LittleEndian, &f2); err != nil {
		return 0, 0, err
	}
	return f1, f2, nil
}

// The WaveSculptor motor controller must receive a Motor Drive Command frame at
// least once every 250ms. If this does not occur, the controller will assume
// that communications have failed and will halt all motor control functions,
// placing the system into neutral.
type MotorDriveCommand struct {
	// Desired motor current set point as a percentage of maximum current setting.
	MotorCurrent float32 `json:"Motor Current"`
	// Desired motor velocity set point in metres/second
	MotorVelocity float32 `json:"Motor Velocity"`
}

func (MotorDriveCommand) New() CAN { return &MotorDriveCommand{} }
func (m *MotorDriveCommand) UnmarshalBinary(b []byte) error {
	var err error
	m.MotorCurrent, m.MotorVelocity, err = twoFloats(b)
	return err
}

type MotorPowerCommand struct {
	// Desired set point of current drawn from the bus by the controller as a
	// percentage of absolute bus current limit.
	BusCurrent float32 `json:"Bus Current"`
}

func (MotorPowerCommand) New() CAN { return &MotorPowerCommand{} }
func (m *MotorPowerCommand) UnmarshalBinary(b []byte) error {
	var err error
	_, m.BusCurrent, err = twoFloats(b)
	return err
}

type BusMeasurement struct {
	// DC Bus voltage at the controller.
	BusVoltage float32 `json:"Bus Voltage`
	// Current drawn from the DC bus by the controller.
	BusCurrent float32 `json:"Bus Current"`
}

func (BusMeasurement) New() CAN { return &BusMeasurement{} }
func (m *BusMeasurement) UnmarshalBinary(b []byte) error {
	var err error
	m.BusVoltage, m.BusCurrent, err = twoFloats(b)
	return err
}

type VelocityMeasurement struct {
	// Motor angular frequency in revolutions per minute.
	MotorVelocity float32 `json:"Motor Velocity"`
	// Vehicle velocity in metres / second.
	VehicleVelocity float32 `json:"Vehicle Velocity"`
}

func (VelocityMeasurement) New() CAN { return &VelocityMeasurement{} }
func (v *VelocityMeasurement) UnmarshalBinary(b []byte) error {
	var err error
	v.MotorVelocity, v.VehicleVelocity, err = twoFloats(b)
	return err
}
