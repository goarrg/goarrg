package input

import (
	"math/rand"
	"testing"
)

type dummyDevice struct {
	currentState [^DeviceAction(0)]bool
	lastState    [^DeviceAction(0)]bool
}

func (d *dummyDevice) Type() string {
	return "dummy"
}

func (d *dummyDevice) StateFor(a DeviceAction) State {
	if d.currentState[a] {
		return Value(1)
	}

	return Value(0)
}

func (d *dummyDevice) StateDeltaFor(a DeviceAction) StateDelta {
	if d.ActionStartedFor(a) {
		return Value(1)
	}

	if d.ActionEndedFor(a) {
		return Value(-1)
	}

	return Value(0)
}

func (d *dummyDevice) ActionStartedFor(a DeviceAction) bool {
	return d.currentState[a] && !d.lastState[a]
}

func (d *dummyDevice) ActionEndedFor(a DeviceAction) bool {
	return d.lastState[a] && !d.currentState[a]
}

var device Device
var action DeviceAction

func TestScan(t *testing.T) {
	devices := [16]dummyDevice{}

	for i := range devices {
		RegisterDevice(&devices[i])
	}

	i := rand.Intn(16)
	a := DeviceAction(rand.Intn(int(^DeviceAction(0))))
	devices[i].currentState[a] = true

	d, a2 := Scan(ScanValue)

	if d != &devices[i] {
		t.Fatal("Device mismatch")
	}

	if a != a2 {
		t.Fatal("Action mismatch")
	}
}

func BenchmarkScan(b *testing.B) {
	devices := [16]dummyDevice{}

	for i := range devices {
		RegisterDevice(&(devices[i]))
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		device, action = Scan(8)
	}

	if device != nil {
		b.Fail()
	}
}
