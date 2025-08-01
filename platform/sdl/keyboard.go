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
	#cgo pkg-config: sdl3
	#include <SDL3/SDL.h>
	#include "event.h"
*/
import "C"

import (
	"bytes"
	"unsafe"

	"goarrg.com/input"
)

type keyboard struct {
	currentState [^input.DeviceAction(0)]byte
	lastState    [^input.DeviceAction(0)]byte
}

func (k *keyboard) Type() string {
	return input.DeviceTypeKeyboard
}

func (k *keyboard) Scan(mask input.ScanMask) input.DeviceAction {
	if !mask.HasBits(input.ScanValue) {
		return 0
	}
	if i := bytes.IndexByte(k.currentState[:], 1); i > 0 {
		return input.DeviceAction(i)
	}
	return 0
}

func (k *keyboard) StateFor(a input.DeviceAction) input.State {
	if k.currentState[a] == 1 {
		return input.Value(1)
	}
	return input.Value(0)
}

func (k *keyboard) StateDeltaFor(a input.DeviceAction) input.StateDelta {
	if k.ActionStartedFor(a) {
		return input.Value(1)
	}

	if k.ActionEndedFor(a) {
		return input.Value(-1)
	}

	return input.Value(0)
}

func (k *keyboard) ActionStartedFor(a input.DeviceAction) bool {
	return (k.currentState[a] == 1) && (k.lastState[a] == 0)
}

func (k *keyboard) ActionEndedFor(a input.DeviceAction) bool {
	return (k.lastState[a] == 1) && (k.currentState[a] == 0)
}

func (k *keyboard) update(C.goEvent) {
	k.lastState = k.currentState

	if !Platform.display.hasKeyboardFocus() {
		k.currentState = [^input.DeviceAction(0)]byte{}
		return
	}

	cNumKeys := C.int(0)
	cKB := C.SDL_GetKeyboardState(&cNumKeys)
	kb := unsafe.Slice((*byte)(unsafe.Pointer(cKB)), int(cNumKeys))
	copy(k.currentState[:], kb)
}
