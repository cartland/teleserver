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