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

type DeviceAction uint8

const (
	/*
		DeviceTypeKeyboard keys will be where the hardware says they are
		e.g. "W S A D" will always have the same physical location
		but may not be labeled "W S A D".
	*/
	DeviceTypeKeyboard = "keyboard"

	/*
		DeviceTypeKeyboardMapped keys will be where the OS says they are.
		e.g. "I" will always refer to the "I" key
	*/
	// DeviceTypeKeyboardMapped

	/*
		DeviceTypeMouse represents a typical 5 button plus scrollwheel mouse.
		MouseMotion coordinates would be relative to the drawable surface where
		top left is {0, 0, 0} and bottom right is {w, h, 0}.
	*/
	DeviceTypeMouse = "mouse"
)

type Device interface {
	/*
		Returns a string used to identify the device type. If type is one of the
		predefined DeviceType* strings then the implementation must follow the
		documentation for said type and its' defined DeviceAction codes.
	*/
	Type() string

	/*
		Returns the current frame's state of the action as either a input.Value,
		input.Axis or a input.Coords.
		Returns nil if DeviceAction is invalid.

		On device disconnect, the device should act as if the player stopped
		interacting with said device. For Value and Axis actions this means they
		are in the default state. For Coords this means they are frozen in place.
		StateDeltaFor and ActionEndedFor should be triggered as expected
		for the state change.

		If each space separated number is the value returned a frame,
		then pressing a button would go:
		default - pressed - released - default
		0 0 0 0 0 1 1 1 1 1 0 0 0 0 0 0 0 0 0
	*/
	StateFor(DeviceAction) State

	/*
		Returns the delta state of the current and previous frame for the action
		as either a input.Value or a input.Axis. Does not return a input.Coords
		as the delta of 2 input.Coords would be a input.Axis (vector).
		Returns nil if DeviceAction is invalid.

		On device disconnect, the device should act as if the player stopped
		interacting with said device. For buttons that were pressed the frame
		previous to the disconnect, the delta would reflect a button release.
		The delta of Coords actions would be zero as they are frozen in place.
		While the delta of Axis actions would reflect the return to the default state.

		If each space separated number is the value returned a frame,
		then pressing a button would go:
		default - pressed - released - default
		0 0 0 0 0 1 0 0 0 0 -1 0 0 0 0 0 0 0 0
	*/
	StateDeltaFor(DeviceAction) StateDelta

	/*
		Returns true if this is the first frame a action started,
		Returns false if DeviceAction is invalid.

		On device disconnect, the device should act as if the player stopped
		interacting with said device, nothing should return true.

		If each space separated number is the value returned a frame,
		then pressing a button would go:
		default - pressed - released - default
		0 0 0 0 0 1 0 0 0 0 0 0 0 0 0 0 0 0 0
	*/
	ActionStartedFor(DeviceAction) bool

	/*
		Returns true if this is the first frame a action stopped,
		Returns false if DeviceAction is invalid.

		On device disconnect, the device should act as if the player stopped
		interacting with said device. All actions that were started the frame
		previous to the disconnect shall be considered ended.

		If each space separated number is the value returned a frame,
		then pressing a button would go:
		default - pressed - released - default
		0 0 0 0 0 0 0 0 0 0 1 0 0 0 0 0 0 0 0
	*/
	ActionEndedFor(DeviceAction) bool
}
