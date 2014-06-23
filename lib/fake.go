package lib

import (
	"math"
	"time"

	"github.com/calsol/teleserver/msgs"
	"github.com/stvnrhodes/broadcaster"
)

const (
	// Time between fake data readings
	metricPeriod = 20 * time.Millisecond
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

func getBattery(id uint16) msgs.CANPlus {
	t := time.Since(startTime)
	fid := float64(id)
	a := uint16(33000 + 3000*math.Sin(t.Seconds()/20+fid))
	b := uint16(33000 + 3000*math.Sin(t.Seconds()/10+fid))
	c := uint16(33000 + 3000*math.Sin(t.Seconds()/15+fid))
	d := uint16(33000 + 3000*math.Sin(t.Seconds()/25+fid))
	return msgs.NewCANPlus(&msgs.BatteryModule{
		ID:       id,
		Voltage0: a,
		Voltage1: b,
		Voltage2: c,
		Voltage3: d,
	})
}

func getMPPT(id uint16) msgs.CANPlus {
	t := time.Since(startTime)
	msg := msgs.IDToMessage(id).(*msgs.MPPTStatus)
	fid := float64(id)
	msg.ArrayVoltage = uint16(5000 + 500*math.Sin(t.Seconds()/20+fid))
	msg.ArrayCurrent = uint16(5000 + 500*math.Cos(t.Seconds()/25+fid))
	msg.BatteryVoltage = uint16(10000 + 1000*math.Cos(t.Seconds()/2+fid))
	msg.Temperature = uint16(2500 + 500*math.Cos(t.Seconds()/3+fid))
	return msgs.NewCANPlus(msg)
}

// GenFake broadcasts fake data for speed, voltage, and power.
func GenFake(b broadcaster.Caster) {
	for {
		b.Cast(getSpeed())
		time.Sleep(metricPeriod)
		b.Cast(getPower())
		time.Sleep(metricPeriod)
		b.Cast(getMPPT(0x600))
		time.Sleep(metricPeriod)
		b.Cast(getMPPT(0x601))
		time.Sleep(metricPeriod)
		b.Cast(getMPPT(0x602))
		time.Sleep(metricPeriod)
		b.Cast(getMPPT(0x603))
		time.Sleep(metricPeriod)
		b.Cast(getBattery(0x130))
		time.Sleep(metricPeriod)
		b.Cast(getBattery(0x131))
		time.Sleep(metricPeriod)
		b.Cast(getBattery(0x132))
		time.Sleep(metricPeriod)
		b.Cast(getBattery(0x140))
		time.Sleep(metricPeriod)
		b.Cast(getBattery(0x141))
		time.Sleep(metricPeriod)
		b.Cast(getBattery(0x142))
		time.Sleep(metricPeriod)
		b.Cast(getBattery(0x150))
		time.Sleep(metricPeriod)
		b.Cast(getBattery(0x151))
		time.Sleep(metricPeriod)
		b.Cast(getBattery(0x152))
		time.Sleep(metricPeriod)
	}
}
