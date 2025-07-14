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
	#include <SDL3/SDL_mouse.h>
	#include "event.h"
*/
import "C"

import (
	"math/bits"

	"goarrg.com/gmath"
	"goarrg.com/input"
)

type mouse struct {
	cCursor      *C.SDL_Cursor
	mode         input.MouseMode
	motion       input.Coords
	motionDelta  input.Axis
	wheel        input.Axis
	wheelDelta   input.Axis
	currentState uint8
	lastState    uint8
}

var _ input.Mouse = (*mouse)(nil)

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

func (m *mouse) SetSystemCursor(c input.SystemCursor) {
	f := func(sc C.SDL_SystemCursor) {
		go func() {
			Platform.taskChan <- func() {
				C.SDL_DestroyCursor(m.cCursor)
				m.cCursor = C.SDL_CreateSystemCursor(sc)
				C.SDL_SetCursor(m.cCursor)
			}
		}()
	}
	switch c {
	case input.SystemCursorDefault:
		f(C.SDL_SYSTEM_CURSOR_DEFAULT)
	case input.SystemCursorText:
		f(C.SDL_SYSTEM_CURSOR_TEXT)
	case input.SystemCursorWait:
		f(C.SDL_SYSTEM_CURSOR_WAIT)
	case input.SystemCursorCrosshair:
		f(C.SDL_SYSTEM_CURSOR_CROSSHAIR)
	case input.SystemCursorProgress:
		f(C.SDL_SYSTEM_CURSOR_PROGRESS)

	case input.SystemCursorResizeHorizontal:
		f(C.SDL_SYSTEM_CURSOR_EW_RESIZE)
	case input.SystemCursorResizeVertical:
		f(C.SDL_SYSTEM_CURSOR_NS_RESIZE)
	case input.SystemCursorResizeDiagonalBackward:
		f(C.SDL_SYSTEM_CURSOR_NWSE_RESIZE)
	case input.SystemCursorResizeDiagonalForward:
		f(C.SDL_SYSTEM_CURSOR_NESW_RESIZE)

	case input.SystemCursorMove:
		f(C.SDL_SYSTEM_CURSOR_MOVE)
	case input.SystemCursorNotAllowed:
		f(C.SDL_SYSTEM_CURSOR_NOT_ALLOWED)
	case input.SystemCursorPointer:
		f(C.SDL_SYSTEM_CURSOR_POINTER)
	}
}

func (m *mouse) SetMode(mode input.MouseMode) {
	m.mode = input.MouseModeDefault

	switch mode {
	case input.MouseModeDefault:
		go func() {
			Platform.taskChan <- func() {
				C.SDL_SetWindowMouseGrab(Platform.display.mainWindow.cWindow, false)
				C.SDL_SetWindowRelativeMouseMode(Platform.display.mainWindow.cWindow, false)
				C.SDL_ShowCursor()
			}
		}()
	case input.MouseModeHidden:
		go func() {
			Platform.taskChan <- func() {
				C.SDL_SetWindowMouseGrab(Platform.display.mainWindow.cWindow, false)
				C.SDL_SetWindowRelativeMouseMode(Platform.display.mainWindow.cWindow, false)
				if !Platform.config.Debug.SafeMouse {
					C.SDL_HideCursor()
				}
			}
		}()
	case input.MouseModeGrabbed:
		go func() {
			Platform.taskChan <- func() {
				if Platform.config.Debug.SafeMouse || !bool(C.SDL_SetWindowMouseGrab(Platform.display.mainWindow.cWindow, true)) {
					m.mode = input.MouseModeGrabbed
				}
				C.SDL_SetWindowRelativeMouseMode(Platform.display.mainWindow.cWindow, false)
				C.SDL_ShowCursor()
			}
		}()
	case input.MouseModeRelative:
		go func() {
			Platform.taskChan <- func() {
				if Platform.config.Debug.SafeMouse {
					m.mode = input.MouseModeRelative
					C.SDL_WarpMouseInWindow(Platform.display.mainWindow.cWindow,
						C.float(Platform.display.mainWindow.rect.W)/2, C.float(Platform.display.mainWindow.rect.H)/2)
					C.SDL_ShowCursor()
				} else if !C.SDL_SetWindowRelativeMouseMode(Platform.display.mainWindow.cWindow, true) {
					m.mode = input.MouseModeRelative
					C.SDL_HideCursor()
					C.SDL_WarpMouseInWindow(Platform.display.mainWindow.cWindow,
						C.float(Platform.display.mainWindow.rect.W)/2, C.float(Platform.display.mainWindow.rect.H)/2)
				}
			}
		}()
	}
}

func (m *mouse) update(e C.goEvent) {
	var cX, cY C.float

	m.lastState = m.currentState
	if (m.mode == input.MouseModeDefault && !Platform.display.hasMouseFocus()) || !Platform.display.hasKeyboardFocus() {
		m.currentState = 0
		m.motionDelta = input.Axis{}
		m.wheelDelta = input.Axis{}
		return
	}

	m.currentState = uint8(C.SDL_GetGlobalMouseState(&cX, &cY)) << 1
	pos := gmath.Point3f64{X: float64(cX), Y: float64(cY)}
	flushMouse := false

	switch {
	case m.mode == input.MouseModeGrabbed:
		flushMouse = true
		pos = Platform.display.mainWindow.bounds.ClampPoint(pos)
		pos = Platform.display.globalPointToRelativePoint(pos)
		C.SDL_WarpMouseInWindow(Platform.display.mainWindow.cWindow,
			C.float(pos.X), C.float(pos.Y))
	case m.mode == input.MouseModeRelative:
		flushMouse = true
		pos = Platform.display.globalPointToRelativePoint(pos)
		C.SDL_WarpMouseInWindow(Platform.display.mainWindow.cWindow,
			C.float(Platform.display.mainWindow.rect.W)/2, C.float(Platform.display.mainWindow.rect.H)/2)
	case !Platform.display.mainWindow.bounds.CheckPoint(pos):
		m.currentState = 0
		m.motionDelta = input.Axis{}
		m.wheelDelta = input.Axis{}
		return
	default:
		pos = Platform.display.globalPointToRelativePoint(pos)
	}

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

	if flushMouse {
		C.SDL_PumpEvents()
		C.SDL_GetRelativeMouseState(&cX, &cY)
	}
}
