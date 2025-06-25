/*
Copyright 2020 The goARRG Authors.

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

package sdl

/*
	#cgo pkg-config: sdl2
	#include <SDL2/SDL.h>
	#include "event.h"
*/
import "C"

import (
	"math/bits"

	"goarrg.com/gmath"
	"goarrg.com/input"
)

type mouse struct {
	motion       input.Coords
	motionDelta  input.Axis
	wheel        input.Axis
	wheelDelta   input.Axis
	currentState uint8
	lastState    uint8
}

func (m *mouse) Type() string {
	return input.DeviceTypeMouse
}

func (m *mouse) Scan(mask input.ScanMask) input.DeviceAction {
	if mask.HasBits(input.ScanValue) {
		i := input.DeviceAction(bits.TrailingZeros8(m.currentState))
		if i > 0 && i <= input.MouseForward {
			return i
		}
	}
	if mask.HasBits(input.ScanAxis) {
		if m.ActionStartedFor(input.MouseWheel) {
			return input.MouseWheel
		}
	}
	if mask.HasBits(input.ScanCoords) {
		if m.ActionStartedFor(input.MouseMotion) {
			return input.MouseMotion
		}
	}
	return 0
}

func (m *mouse) StateFor(a input.DeviceAction) input.State {
	switch a {
	case input.MouseMotion:
		return m.motion
	case input.MouseWheel:
		return m.wheel
	case input.MouseLeft, input.MouseMiddle, input.MouseRight, input.MouseBack, input.MouseForward:
		if m.currentState&(1<<a) != 0 {
			return input.Value(1)
		}

		return input.Value(0)
	}

	return nil
}

func (m *mouse) StateDeltaFor(a input.DeviceAction) input.StateDelta {
	switch a {
	case input.MouseMotion:
		return m.motionDelta
	case input.MouseWheel:
		return m.wheelDelta
	case input.MouseLeft, input.MouseMiddle, input.MouseRight, input.MouseBack, input.MouseForward:
		if m.ActionStartedFor(a) {
			return input.Value(1)
		}

		if m.ActionEndedFor(a) {
			return input.Value(-1)
		}

		return input.Value(0)
	}

	return nil
}

func (m *mouse) ActionStartedFor(a input.DeviceAction) bool {
	if a > input.MouseMotion {
		return false
	}

	mask := uint8(1 << a)
	return ((m.currentState & mask) - (m.lastState & mask)) == mask
}

func (m *mouse) ActionEndedFor(a input.DeviceAction) bool {
	if a > input.MouseMotion {
		return false
	}

	mask := uint8(1 << a)
	return ((m.lastState & mask) - (m.currentState & mask)) == mask
}

func (m *mouse) update(e C.goEvent) {
	var cX, cY C.int

	m.lastState = m.currentState

	if !Platform.display.hasMouseFocus() {
		m.currentState = 0
		m.motionDelta = input.Axis{}
		m.wheelDelta = input.Axis{}
		return
	}

	m.currentState = uint8(C.SDL_GetGlobalMouseState(&cX, &cY)) << 1
	pos := gmath.Point3int{X: int(cX), Y: int(cY)}

	if !Platform.display.pointInsideWindow(pos) {
		m.currentState = 0
		m.motionDelta = input.Axis{}
		m.wheelDelta = input.Axis{}
		return
	}

	pos = Platform.display.globalPointToRelativePoint(pos)
	m.motion = input.Coords{
		Point3f64: gmath.Point3f64{
			X: float64(pos.X),
			Y: float64(pos.Y),
		},
	}

	C.SDL_GetRelativeMouseState(&cX, &cY)
	m.motionDelta = input.Axis{
		Vector3f64: gmath.Vector3f64{
			X: float64(cX),
			Y: float64(cY),
		},
	}

	if m.motionDelta != (input.Axis{}) {
		m.currentState |= uint8(1 << input.MouseMotion)
	}

	oldWheel := m.wheel
	m.wheel = input.Axis{
		Vector3f64: gmath.Vector3f64{
			X: float64(e.mouseWheelX),
			Y: float64(e.mouseWheelY),
		},
	}

	m.wheelDelta = input.Axis{
		Vector3f64: m.wheel.Vector3f64.Subtract(oldWheel.Vector3f64),
	}

	if m.wheelDelta != (input.Axis{}) {
		m.currentState |= uint8(1 << input.MouseWheel)
	}
}
