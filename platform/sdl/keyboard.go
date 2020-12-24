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
	"reflect"
	"unsafe"

	"goarrg.com/input"
)

type keyboard struct {
	state  [^input.DeviceAction(0)]bool
	events [int(^input.DeviceAction(0)) * 2]bool
}

func (k *keyboard) Type() input.DeviceType {
	return input.DeviceTypeKeyboard
}

func (k *keyboard) ID() input.DeviceID {
	return 0
}

func (k *keyboard) StateFor(a input.DeviceAction) interface{} {
	if k.state[a] {
		return input.Value(1)
	}

	return input.Value(0)
}

func (k *keyboard) StateDeltaFor(a input.DeviceAction) interface{} {
	if k.state[a] {
		return input.Value(1)
	}

	if k.ActionEndedFor(a) {
		return input.Value(-1)
	}

	return input.Value(0)
}

func (k *keyboard) ActionStartedFor(a input.DeviceAction) bool {
	return k.events[2*a]
}

func (k *keyboard) ActionEndedFor(a input.DeviceAction) bool {
	return k.events[(2*a)+1]
}

func keyboardState(C.goEvent) keyboard {
	state := keyboard{}
	cNumKeys := C.int(0)
	cKB := C.SDL_GetKeyboardState(&cNumKeys)
	kb := *(*[]uint8)(unsafe.Pointer(&reflect.SliceHeader{
		uintptr(unsafe.Pointer(cKB)), int(cNumKeys), int(cNumKeys),
	}))

	for i := input.KeyA; i < input.KeyRightGUI; i++ {
		state.state[i] = kb[i] == 1

		if state.state[i] {
			if !Platform.input.lastState.keyboard.state[i] {
				state.events[2*i] = true
			}
		} else if Platform.input.lastState.keyboard.state[i] {
			state.events[(2*i)+1] = true
		}
	}

	return state
}
