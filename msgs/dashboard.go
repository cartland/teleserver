package msgs

type CANAccelPos struct {
	PedalPosition float64
}

func (CANAccelPos) New() CAN { return &CANAccelPos{} }

type DashHeartbeat struct {
	Time int32
}

func (DashHeartbeat) New() CAN { return &DashHeartbeat{} }

type CANBrakePos struct {
	PedalPosition float64
}

func (CANBrakePos) New() CAN { return &CANBrakePos{} }

type PanelSwitchPos struct {
	HVSwitch      bool `binpack:"0.0"`
	SolarSwitch   bool `binpack:"0.1"`
	LightsSwitch  bool `binpack:"0.2"`
	ExtraSwitch   bool `binpack:"0.3"`
	ForwardToggle bool `binpack:"0.4"`
	NeutralToggle bool `binpack:"0.5"`
	ReverseToggle bool `binpack:"0.6"`
	RegenToggle   bool `binpack:"0.7"`
	HazardsToggle bool `binpack:"1.0"`
	PowerEcoMode  bool `binpack:"1.1"`
}

func (PanelSwitchPos) New() CAN { return &PanelSwitchPos{} }
