package msgs

type MotorRPM struct {
	ID         uint16 `json:"-"`
	Location   string
	CurrentRPM float32 `binpack:"0-4"`
}

func (m MotorRPM) New() CAN      { return &MotorRPM{ID: m.ID, Location: m.Location} }
func (m MotorRPM) canID() uint16 { return m.ID }

type MotorHeartbeat struct {
	ID       uint16 `json:"-"`
	Location string
	Time     int32 `binpack:"0-4"`
}

func (m MotorHeartbeat) New() CAN      { return &MotorHeartbeat{ID: m.ID, Location: m.Location} }
func (m MotorHeartbeat) canID() uint16 { return m.ID }
