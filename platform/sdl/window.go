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
	#include <stdint.h>

	void setWindowTitle(SDL_Window *window, _GoString_ title) {
		SDL_SetWindowTitle(window, _GoStringPtr(title));
	}

	void pushWindowCreatedEvent(uint32_t window) {
		SDL_Event e = {};
		e.type = SDL_WINDOWEVENT;
		e.window.windowID = window;
		SDL_PushEvent(&e);
	}
*/
import "C"
import (
	"unsafe"

	"goarrg.com/debug"
)

//nolint:deadcode,varcheck,unused
const (
	windowEventCreated uint32 = (1 << iota)
	windowEventShown
	windowEventHidden
	windowEventMoved
	windowEventResized
	windowEventEnter
	windowEventLeave
	windowEventFocusGained
	windowEventFocusLost
	windowEventClose
)

type windowEvent struct {
	event uint32
	x     int32
	y     int32
	w     int32
	h     int32
}

type windowAPI interface {
	resize(int, int)
	destroy()
}

type window struct {
	cWindow       *C.SDL_Window
	api           windowAPI
	cID           C.uint32_t
	keyboardFocus bool
	mouseFocus    bool
}

func createWindow(flags C.uint32_t) error {
	if Platform.config.Window.Rect.X < 0 {
		Platform.config.Window.Rect.X = C.SDL_WINDOWPOS_UNDEFINED
	}

	if Platform.config.Window.Rect.Y < 0 {
		Platform.config.Window.Rect.Y = C.SDL_WINDOWPOS_UNDEFINED
	}

	rect := Platform.config.Window.Rect

	switch Platform.config.Window.Mode {
	case WindowModeBorderless:
		{
			var cRect C.SDL_Rect

			if C.SDL_GetDisplayBounds(0, &cRect) != 0 {
				err := debug.ErrorNew(C.GoString(C.SDL_GetError()))
				C.SDL_ClearError()
				debug.LogE("SDL failed to create window: %s", err.Error())
			}

			flags |= C.SDL_WINDOW_BORDERLESS

			rect.X = int(cRect.x)
			rect.Y = int(cRect.y)
			rect.W = int(cRect.w)
			rect.H = int(cRect.h)
		}

	case WindowModeFullscreen:
		flags |= C.SDL_WINDOW_FULLSCREEN
	}

	cTitle := C.CString(Platform.config.Window.Title)
	defer C.free(unsafe.Pointer(cTitle))

	cWindow := C.SDL_CreateWindow(
		cTitle,
		C.int(rect.X),
		C.int(rect.Y),
		C.int(rect.W),
		C.int(rect.H),
		C.SDL_WINDOW_HIDDEN|C.SDL_WINDOW_RESIZABLE|C.SDL_WINDOW_ALLOW_HIGHDPI|flags,
	)

	if cWindow == nil {
		err := debug.ErrorNew(C.GoString(C.SDL_GetError()))
		C.SDL_ClearError()
		debug.LogE("SDL failed to create window: %s", err.Error())
		return err
	}

	Platform.display.mainWindow = &window{
		cWindow: cWindow,
		cID:     C.SDL_GetWindowID(cWindow),
	}

	C.pushWindowCreatedEvent(Platform.display.mainWindow.cID)

	return nil
}

func (window *window) processEvent(e windowEvent) {
	switch {
	case (e.event & windowEventResized) != 0:
		window.api.resize(int(e.w), int(e.h))
	case (e.event & windowEventCreated) != 0:
		C.SDL_ShowWindow(window.cWindow)

	case (e.event & windowEventFocusGained) != 0:
		window.keyboardFocus = true
	case (e.event & windowEventFocusLost) != 0:
		window.keyboardFocus = false

	case (e.event & windowEventEnter) != 0:
		window.mouseFocus = true
	case (e.event & windowEventLeave) != 0:
		window.mouseFocus = false
	}
}

func (window *window) destroy() {
	if window.api != nil {
		window.api.destroy()
	}

	C.SDL_DestroyWindow(window.cWindow)
	if err := C.GoString(C.SDL_GetError()); err != "" {
		C.SDL_ClearError()
		debug.LogE("SDL_DestroyWindow failed: %s", err)
	}
}
