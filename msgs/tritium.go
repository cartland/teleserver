package msgs

// The WaveSculptor motor controller must receive a Motor Drive Command frame at
// least once every 250ms. If this does not occur, the controller will assume
// that communications have failed and will halt all motor control functions,
// placing the system into neutral.
type MotorDriveCommand struct {
	// Desired motor current set point as a percentage of maximum current setting.
	MotorCurrent float32 `binpack:"0-4"`
	// Desired motor velocity set point in metres/second
	MotorVelocity float32 `binpack:"4-8"`
}

func (MotorDriveCommand) New() CAN { return &MotorDriveCommand{} }

type MotorPowerCommand struct {
	// Desired set point of current drawn from the bus by the controller as a
	// percentage of absolute bus current limit.
	BusCurrent float32 `binpack:"4-8"`
}

func (MotorPowerCommand) New() CAN { return &MotorPowerCommand{} }

type BusMeasurement struct {
	// DC Bus voltage at the controller.
	BusVoltage float32 `binpack:"0-4"`
	// Current drawn from the DC bus by the controller.
	BusCurrent float32 `binpack:"4-8"`
}

func (BusMeasurement) New() CAN { return &BusMeasurement{} }

type VelocityMeasurement struct {
	// Motor angular frequency in revolutions per minute.
	MotorVelocity float32 `binpack:"0-4"`
	// Vehicle velocity in metres / second.
	VehicleVelocity float32 `binpack:"4-8"`
}

func (VelocityMeasurement) New() CAN { return &VelocityMeasurement{} }
