package msgs

type PowerHubHeartbeat struct {
	ID       uint16 `json:"-"`
	Location string
	Time     int32 `binpack:"0-3"`
}

func (p PowerHubHeartbeat) New() CAN      { return &PowerHubHeartbeat{ID: p.ID, Location: p.Location} }
func (p PowerHubHeartbeat) canID() uint16 { return p.ID }
