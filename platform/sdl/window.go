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
	"goarrg.com/gmath"
)

//nolint:deadcode,varcheck,unused
const (
	windowEventCreated uint32 = (1 << iota)
	windowEventShown
	windowEventHidden
	windowEventRectChanged
	windowEventEnter
	windowEventLeave
	windowEventFocusGained
	windowEventFocusLost
	windowEventClose
)

type windowEvent struct {
	event uint32
}

type windowAPI interface {
	resize(int, int)
	destroy()
}

type window struct {
	rect          gmath.Recti
	bounds        gmath.Bounds3i
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
				err := debug.Errorf(C.GoString(C.SDL_GetError()))
				C.SDL_ClearError()
				Platform.logger.EPrintf("Failed to create window: %s", err.Error())
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
		err := debug.Errorf(C.GoString(C.SDL_GetError()))
		C.SDL_ClearError()
		Platform.logger.EPrintf("Failed to create window: %s", err.Error())
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
	if (e.event & windowEventCreated) != 0 {
		Platform.logger.VPrintf("Window event created")
		C.SDL_ShowWindow(window.cWindow)
		C.SDL_FlushEvent(C.SDL_WINDOWEVENT)

		cRect := C.SDL_Rect{}
		C.SDL_GetWindowPosition(window.cWindow, &cRect.x, &cRect.y)
		C.SDL_GetWindowSize(window.cWindow, &cRect.w, &cRect.h)

		window.rect = gmath.Recti{X: int(cRect.x), Y: int(cRect.y), W: int(cRect.w), H: int(cRect.h)}
		window.bounds.Min = gmath.Vector3i{X: window.rect.X, Y: window.rect.Y}
		window.bounds.Max = gmath.Vector3i{X: window.rect.W, Y: window.rect.H}.Add(window.bounds.Min)

		// we will always have keyboard focus at this point
		window.keyboardFocus = true
		/*
			we may have mouse focus, without this when we would not get a focus
			entered event until the mouse leave and re-enters. for when we
			actually do not have mouse focus, the pointInsideWindow should cover
			that since the window will be on top of everything.
		*/
		window.mouseFocus = true
	}

	if (e.event & windowEventRectChanged) != 0 {
		Platform.logger.VPrintf("Window event rect changed")
		cRect := C.SDL_Rect{}
		C.SDL_GetWindowPosition(window.cWindow, &cRect.x, &cRect.y)
		C.SDL_GetWindowSize(window.cWindow, &cRect.w, &cRect.h)

		oldRect := window.rect
		window.rect = gmath.Recti{X: int(cRect.x), Y: int(cRect.y), W: int(cRect.w), H: int(cRect.h)}
		window.bounds.Min = gmath.Vector3i{X: window.rect.X, Y: window.rect.Y}
		window.bounds.Max = gmath.Vector3i{X: window.rect.W, Y: window.rect.H}.Add(window.bounds.Min)

		if oldRect.W != window.rect.W || oldRect.H != window.rect.H {
			window.api.resize(window.rect.W, window.rect.H)
		}
	} else if (e.event & windowEventShown) != 0 {
		Platform.logger.VPrintf("Window event shown")
		window.api.resize(window.rect.W, window.rect.H)
	}

	if (e.event & windowEventFocusGained) != 0 {
		Platform.logger.VPrintf("Window event focus gained")
		window.keyboardFocus = true
	}

	if (e.event & windowEventFocusLost) != 0 {
		Platform.logger.VPrintf("Window event focus lost")
		window.keyboardFocus = false
	}

	if (e.event & windowEventEnter) != 0 {
		Platform.logger.VPrintf("Window event enter")
		window.mouseFocus = true
	}

	if (e.event & windowEventLeave) != 0 {
		Platform.logger.VPrintf("Window event leave")
		window.mouseFocus = false
	}

	if (e.event & windowEventHidden) != 0 {
		Platform.logger.VPrintf("Window event hidden")
		window.api.resize(0, 0)
	}
}

func (window *window) destroy() {
	if window.api != nil {
		window.api.destroy()
	}

	C.SDL_DestroyWindow(window.cWindow)
	if err := C.GoString(C.SDL_GetError()); err != "" {
		C.SDL_ClearError()
		Platform.logger.EPrintf("SDL_DestroyWindow failed: %s", err)
	}
}

func SetWindowTitle(title string) {
	C.setWindowTitle(Platform.display.mainWindow.cWindow, title+"\x00")
}
