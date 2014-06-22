package msgs

// The MPPTStatus message is sent from the MPPTs at a regular interval
type MPPTStatus struct {
	// ID holds the CAN ID of the message
	ID uint16 `json:"-"`
	// Array location is a human-readable description of the array location.
	ArrayLocation string
	// Array voltage ís scaled by 100, or 1 count = 10mV
	ArrayVoltage uint16 `binpack:"0-2"`
	// Array current ís scaled by 1000, or 1 count = 1mA
	ArrayCurrent uint16 `binpack:"2-4"`
	// Battery voltage is scaled by 100, or 1 count = 10mV
	BatteryVoltage uint16 `binpack:"4-6"`
	// Temperature is scaled by 100, or 1 count = 10mC
	Temperature uint16 `binpack:"6-8"`
}

func (m MPPTStatus) New() CAN      { return &MPPTStatus{ID: m.ID, ArrayLocation: m.ArrayLocation} }
func (m MPPTStatus) canID() uint16 { return m.ID }

// The MPPTEnable message is sent to the MPPTs to turn them on or off
type MPPTEnable struct {
	// id holds the CAN ID of the message
	ID uint16 `json:"-"`
	// Array location is a human-readable description of the array location
	ArrayLocation string
	// The message will either enable or disable the power point trackets
	Enable bool `binpack:"0.0"`
}

func (m MPPTEnable) New() CAN      { return &MPPTEnable{ID: m.ID, ArrayLocation: m.ArrayLocation} }
func (m MPPTEnable) canID() uint16 { return m.ID }
