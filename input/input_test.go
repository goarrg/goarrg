//go:build !goarrg_build_debug
// +build !goarrg_build_debug

/*
Copyright 2021 The goARRG Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package input

import (
	"bytes"
	"math/rand"
	"testing"
)

type dummyDevice struct {
	currentState [^DeviceAction(0)]byte
	lastState    [^DeviceAction(0)]byte
}

func (d *dummyDevice) Type() string {
	return "dummy"
}

func (d *dummyDevice) Scan(mask ScanMask) DeviceAction {
	if !mask.HasBits(ScanValue) {
		return 0
	}
	if i := bytes.IndexByte(d.currentState[:], 1); i > 0 {
		return DeviceAction(i)
	}
	return 0
}

func (d *dummyDevice) StateFor(a DeviceAction) State {
	if d.currentState[a] == 1 {
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
	return (d.currentState[a] == 1) && (d.lastState[a] == 0)
}

func (d *dummyDevice) ActionEndedFor(a DeviceAction) bool {
	return (d.lastState[a] == 1) && (d.currentState[a] == 0)
}

var (
	device Device
	action DeviceAction
)

func TestScan(t *testing.T) {
	devices := [16]dummyDevice{}

	for i := range devices {
		RegisterDevice(&devices[i])
	}

	i := rand.Intn(16)
	a := DeviceAction(rand.Intn(int(^DeviceAction(0))))
	devices[i].currentState[a] = 1

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
