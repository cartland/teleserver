package msgs

type CANAccelPos struct {
	PedalPosition int32
}

func (CANAccelPos) New() CAN { return &CANAccelPos{} }

type CANBrakePos struct {
	PedalPosition int32
}

func (CANBrakePos) New() CAN { return &CANBrakePos{} }

type PanelSwitchPos struct {
	HVSwitch     bool `binpack:"0.0"`
	SolarSwitch  bool `binpack:"0.1"`
	LightsSwitch bool `binpack:"0.2"`
	ExtraSwitch  bool `binpack:"0.2"`
}

func (PanelSwitchPos) New() CAN { return &PanelSwitchPos{} }
