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

package input

type DeviceType uint8
type DeviceID uint8
type DeviceAction uint8

const (
	_ DeviceType = iota

	/*
		DeviceTypeKeyboard keys will be where the hardware says they are
		e.g. "W S A D" will always have the same physical location
		but may not be labeled "W S A D".
	*/
	DeviceTypeKeyboard

	/*
		DeviceTypeKeyboardMapped keys will be where the OS says they are.
		e.g. "I" will always refer to the "I" key
	*/
	// DeviceTypeKeyboardMapped

	DeviceTypeMouse
	DeviceTypeCount
)

type Device interface {
	Type() DeviceType
	ID() DeviceID

	/*
		Returns the current state of the action as either a input.Value, input.Axis or a input.Coords.
		e.g. buttons return 1 for pressed and 0 when at the default position.
	*/
	StateFor(DeviceAction) interface{}

	/*
		Returns the delta of current and previous state of the action
		as either a input.Value or a input.Axis.
		e.g. pressing a button would go
		default - pressed - released - default
		0 ... 0   1 ... 1     -1       0 ... 0
	*/
	StateDeltaFor(DeviceAction) interface{}

	/*
		Returns true if this is the first snapshot a action started,
		e.g. if a button was pressed
	*/
	ActionStartedFor(DeviceAction) bool

	/*
		Returns true if this is the first snapshot a action stopped,
		e.g. if a button was released
	*/
	ActionEndedFor(DeviceAction) bool
}
