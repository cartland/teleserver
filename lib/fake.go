package lib

import (
	"io"
	"log"
	"math"
	"time"

	"github.com/calsol/teleserver/msgs"
	"github.com/stvnrhodes/broadcaster"
)

const (
	// Time between fake data readings
	metricPeriod = 500 * time.Millisecond
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
	msg.ArrayVoltage = uint16(5000 + 500*math.Sin(t.Seconds()+fid))
	msg.ArrayCurrent = uint16(5000 + 500*math.Cos(t.Seconds()+fid))
	msg.BatteryVoltage = uint16(10000 + 1000*math.Cos(t.Seconds()/2+fid))
	msg.Temperature = uint16(2500 + 5000*math.Cos(t.Seconds()/3+fid))
	return msgs.NewCANPlus(msg)
}

// readTill takes bytes from the reader until it sees b.
func readTill(r io.Reader, b byte) {
	p := make([]byte, 1)
	for _, err := r.Read(p); p[0] != b; _, err = r.Read(p) {
		if err != nil {
			log.Fatal(err)
		}
	}
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
