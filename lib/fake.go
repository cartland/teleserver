package lib

import (
	"math"
	"time"

	"github.com/calsol/teleserver/msgs"
	"github.com/stvnrhodes/broadcaster"
)

const (
	// Time between fake data readings
	metricPeriod = 200 * time.Millisecond
)

var startTime = time.Now()

func getSpeed() msgs.CANPlus {
	t := time.Since(startTime)
	v := float32(50 + 20*math.Cos(t.Seconds()))
	return msgs.NewCANPlus(&msgs.VelocityMeasurement{VehicleVelocity: v})
}

func getPower() msgs.CANPlus {
	t := time.Since(startTime)
	v := float32(100 + 10*math.Sin(t.Seconds()))
	a := v/4 + float32(10*math.Sin(t.Seconds()/1.8))
	return msgs.NewCANPlus(&msgs.BusMeasurement{BusVoltage: v, BusCurrent: a})
}

func getMPPT(id uint16) msgs.CANPlus {
	t := time.Since(startTime)
	msg := msgs.IDToMessage(id).(*msgs.MPPTStatus)
	fid := float64(id)
	msg.ArrayVoltage = uint16(5000 + 500*math.Sin(t.Seconds()/20+fid))
	msg.ArrayCurrent = uint16(5000 + 500*math.Cos(t.Seconds()/25+fid))
	msg.BatteryVoltage = uint16(10000 + 1000*math.Cos(t.Seconds()/2+fid))
	msg.Temperature = uint16(2500 + 5000*math.Cos(t.Seconds()/3+fid))
	return msgs.NewCANPlus(msg)
}

// GenFake broadcasts fake data for speed, voltage, and power.
func GenFake(b broadcaster.Caster) {
	for {
		b.Cast(getSpeed())
		b.Cast(getPower())
		b.Cast(getMPPT(0x600))
		b.Cast(getMPPT(0x601))
		b.Cast(getMPPT(0x602))
		b.Cast(getMPPT(0x603))
		time.Sleep(metricPeriod)
	}
}
