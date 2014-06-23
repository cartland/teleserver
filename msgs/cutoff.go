package msgs

type CutoffHeartbeat struct {
	Time int32
}

func (CutoffHeartbeat) New() CAN { return &CutoffHeartbeat{} }

type CutoffTrigger struct {
	SolarRelay  bool `binpack:"0.0"`
	SolarToggle bool `binpack:"0.1"`
	MotorRelay  bool `binpack:"0.2"`
	MotorToggle bool `binpack:"0.3"`
}

func (CutoffTrigger) New() CAN { return &CutoffTrigger{} }

type CutoffAnalogIn struct {
	Ain0 int16
	Ain1 int16
	Ain2 int16
	Ain3 int16
}

func (CutoffAnalogIn) New() CAN { return &CutoffAnalogIn{} }

type CutoffSPIIn struct {
	Ain0 int16
	Ain1 int16
	Ain2 int16
}

func (CutoffSPIIn) New() CAN { return &CutoffSPIIn{} }
