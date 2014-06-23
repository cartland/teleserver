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
	Ain0 int16 `binpack:"0-2"`
	Ain1 int16 `binpack:"2-4"`
	Ain2 int16 `binpack:"4-6"`
	Ain3 int16 `binpack:"6-8"`
}

func (CutoffAnalogIn) New() CAN { return &CutoffAnalogIn{} }

type CutoffSPIIn struct {
	Ain0 int16 `binpack:"0-2"`
	Ain1 int16 `binpack:"2-4"`
	Ain2 int16 `binpack:"4-6"`
}

func (CutoffSPIIn) New() CAN { return &CutoffSPIIn{} }
