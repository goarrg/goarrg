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
	#include "event.h"
*/
import "C"
import (
	"goarrg.com/debug"
	"goarrg.com/input"
)

type inputState struct {
	mouse    mouse
	keyboard keyboard
}

type inputSystem struct {
	lastState inputState
}

func (i *inputState) Device(d input.DeviceType) (input.Device, error) {
	switch d {
	case input.DeviceTypeMouse:
		return &i.mouse, nil
	case input.DeviceTypeKeyboard:
		return &i.keyboard, nil
	default:
		return nil, debug.ErrorNew("Unknown device ID")
	}
}

func (i *inputSystem) update(e C.goEvent) *inputState {
	state := inputState{
		mouse:    mouseState(e),
		keyboard: keyboardState(e),
	}

	i.lastState = state
	return &state
}
