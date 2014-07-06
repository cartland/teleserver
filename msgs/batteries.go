package msgs

type BatteryHeartbeat struct {
	Time uint16
}

func (BatteryHeartbeat) New() CAN { return &BatteryHeartbeat{} }

type BatteryModule struct {
	ID       uint16 `json:"-"`
	Voltage0 uint16 `binpack:"0-2"`
	Voltage1 uint16 `binpack:"2-4"`
	Voltage2 uint16 `binpack:"4-6"`
	Voltage3 uint16 `binpack:"6-8"`
}

func (b BatteryModule) New() CAN      { return &BatteryModule{ID: b.ID} }
func (b BatteryModule) canID() uint16 { return b.ID }

type PackVoltage struct {
	Voltage float32
}

func (PackVoltage) New() CAN { return &PackVoltage{} }

type PackCurrent struct {
	Current float32
}

func (PackCurrent) New() CAN { return &PackCurrent{} }

type PackTemperature struct {
	Highest int8
}

func (PackTemperature) New() CAN { return &PackTemperature{} }

type BPSBalancing struct {
	Cell00 bool `binpack:"0.0"`
	Cell01 bool `binpack:"0.1"`
	Cell02 bool `binpack:"0.2"`
	Cell03 bool `binpack:"0.3"`
	Cell04 bool `binpack:"0.4"`
	Cell05 bool `binpack:"0.5"`
	Cell06 bool `binpack:"0.6"`
	Cell07 bool `binpack:"0.7"`
	Cell08 bool `binpack:"1.0"`
	Cell09 bool `binpack:"1.1"`
	Cell10 bool `binpack:"1.2"`
	Cell11 bool `binpack:"1.3"`
	Cell12 bool `binpack:"1.4"`
	Cell13 bool `binpack:"1.5"`
	Cell14 bool `binpack:"1.6"`
	Cell15 bool `binpack:"1.7"`
	Cell16 bool `binpack:"2.0"`
	Cell17 bool `binpack:"2.1"`
	Cell18 bool `binpack:"2.2"`
	Cell19 bool `binpack:"2.3"`
	Cell20 bool `binpack:"2.4"`
	Cell21 bool `binpack:"2.5"`
	Cell22 bool `binpack:"2.6"`
	Cell23 bool `binpack:"2.7"`
	Cell24 bool `binpack:"3.0"`
	Cell25 bool `binpack:"3.1"`
	Cell26 bool `binpack:"3.2"`
	Cell27 bool `binpack:"3.3"`
}

func (BPSBalancing) New() CAN { return &BPSBalancing{} }

type BPSState struct {
	State int8
}

func (BPSState) New() CAN { return &BPSState{} }
