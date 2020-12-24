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
	"goarrg.com/gmath"
	"goarrg.com/input"
)

type mouse struct {
	motion      input.Coords
	motionDelta input.Axis
	wheelDelta  input.Axis
	state       uint32
	events      uint32
}

func (m *mouse) Type() input.DeviceType {
	return input.DeviceTypeMouse
}

func (m *mouse) ID() input.DeviceID {
	return 0
}

func (m *mouse) StateFor(a input.DeviceAction) interface{} {
	switch a {
	case input.MouseMotion:
		return m.motion
	case input.MouseWheel:
		return input.Axis{}
	case input.MouseLeft, input.MouseMiddle, input.MouseRight, input.MouseBack, input.MouseForward:
		if m.state&(1<<a) != 0 {
			return input.Value(1)
		}

		return input.Value(0)
	}

	return nil
}

func (m *mouse) StateDeltaFor(a input.DeviceAction) interface{} {
	switch a {
	case input.MouseMotion:
		return m.motionDelta
	case input.MouseWheel:
		return m.wheelDelta
	case input.MouseLeft, input.MouseMiddle, input.MouseRight, input.MouseBack, input.MouseForward:
		if m.state&(1<<a) != 0 {
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
	return m.events&(1<<(2*a)) != 0
}

func (m *mouse) ActionEndedFor(a input.DeviceAction) bool {
	return m.events&(1<<((2*a)+1)) != 0
}

func mouseState(e C.goEvent) mouse {
	var cX, cY C.int

	state := mouse{}
	state.state = uint32(C.SDL_GetMouseState(&cX, &cY)) << 1
	state.motion = input.Coords{
		Point3f64: gmath.Point3f64{
			X: float64(cX),
			Y: float64(cY),
		},
	}

	C.SDL_GetRelativeMouseState(&cX, &cY)
	state.motionDelta = input.Axis{
		Vector3f64: gmath.Vector3f64{
			X: float64(cX),
			Y: float64(cY),
		},
	}

	if state.motionDelta != (input.Axis{}) {
		state.events |= 1 << (2 * input.MouseMotion)
	} else if Platform.input.lastState.mouse.motionDelta != (input.Axis{}) {
		state.events |= 1 << ((2 * input.MouseMotion) + 1)
	}

	if e.mouseWheelX != 0 || e.mouseWheelY != 0 {
		state.wheelDelta = input.Axis{
			Vector3f64: gmath.Vector3f64{
				X: float64(e.mouseWheelX),
				Y: float64(e.mouseWheelY),
			},
		}
		state.events |= 1 << (2 * input.MouseWheel)
	} else if Platform.input.lastState.mouse.wheelDelta != (input.Axis{}) {
		state.events |= 1 << ((2 * input.MouseWheel) + 1)
	}

	for i := input.MouseLeft; i < input.MouseWheel; i++ {
		if (state.state&(1<<i) != 0) && (Platform.input.lastState.mouse.state&(1<<i) == 0) {
			state.events |= 1 << (2 * i)
		} else if (state.state&(1<<i) == 0) && (Platform.input.lastState.mouse.state&(1<<i) != 0) {
			state.events |= 1 << ((2 * i) + 1)
		}
	}

	return state
}
