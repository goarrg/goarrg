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

const (
	_ DeviceAction = iota
	MouseLeft
	MouseMiddle
	MouseRight
	MouseBack
	MouseForward
	MouseWheel
	MouseMotion
)

type SystemCursor uint32

const (
	SystemCursorDefault SystemCursor = iota
	SystemCursorText
	SystemCursorWait
	SystemCursorCrosshair
	SystemCursorProgress

	SystemCursorResizeLeft
	SystemCursorResizeRight

	SystemCursorResizeTop
	SystemCursorResizeBottom

	SystemCursorResizeTopLeft
	SystemCursorResizeTopRight

	SystemCursorResizeBottomLeft
	SystemCursorResizeBottomRight

	SystemCursorMove
	SystemCursorNotAllowed
	SystemCursorPointer
)

type MouseMode uint32

const (
	MouseModeDefault MouseMode = iota
	/*
		MouseModeHidden hides the cursor such that it is not visible but otherwise functional.
		Implementations should to limit this to within the window and take precautions to ensure
		the cursor is useable on debug pause. (Windows handles it for you but linux doesn't)
	*/
	MouseModeHidden
	/*
		MouseModeGrabbed prevents the cursor from leaving the window but otherwise functional.
		Implementations should take precautions to ensure the cursor is useable on debug pause.
		(Windows handles it for you but linux doesn't)
	*/
	MouseModeGrabbed
	/*
		MouseModeRelative is a mode combining hidden and grabbed where the application only needs relative mouse movement.
		This is the only mode that continues to generate MousePosition deltas when hitting the edge of the window.
		Implementations should take precautions to ensure the cursor is useable on debug pause.
		(Windows handles it for you but linux doesn't)
	*/
	MouseModeRelative
)

type Mouse interface {
	Device
	SetSystemCursor(SystemCursor)
	SetMode(MouseMode)
}
