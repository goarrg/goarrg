//+build !disable_gl

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
*/
import "C"
import (
	"goarrg.com"
	"goarrg.com/debug"
)

type glInstance struct {
}

func (gl *glInstance) ProcAddr() uintptr {
	return uintptr(C.SDL_GL_GetProcAddress)
}

func (gl *glInstance) SwapBuffers() {
	C.SDL_GL_SwapWindow(Platform.display.mainWindow.cWindow)
}

type glWindow struct {
	renderer goarrg.GLRenderer
	cContext C.SDL_GLContext

	windowW int
	windowH int
}

func glInit(r goarrg.GLRenderer) error {
	debug.LogV("SDL creating gl Window")

	if r == nil {
		err := debug.ErrorNew("Invalid renderer")
		debug.LogE("SDL failed to create window: Invalid renderer")
		return err
	}

	glCfg := r.GLConfig()

	switch glCfg.Profile {
	case goarrg.GLProfileCore:
		C.SDL_GL_SetAttribute(C.SDL_GL_CONTEXT_PROFILE_MASK, C.SDL_GL_CONTEXT_PROFILE_CORE)
	case goarrg.GLProfileCompat:
		C.SDL_GL_SetAttribute(C.SDL_GL_CONTEXT_PROFILE_MASK, C.SDL_GL_CONTEXT_PROFILE_COMPATIBILITY)
	case goarrg.GLProfileES:
		C.SDL_GL_SetAttribute(C.SDL_GL_CONTEXT_PROFILE_MASK, C.SDL_GL_CONTEXT_PROFILE_ES)
	}

	if glCfg.Major > 0 {
		C.SDL_GL_SetAttribute(C.SDL_GL_CONTEXT_MAJOR_VERSION, C.int(glCfg.Major))
	}

	C.SDL_GL_SetAttribute(C.SDL_GL_CONTEXT_MINOR_VERSION, C.int(glCfg.Minor))
	C.SDL_GL_SetAttribute(C.SDL_GL_SHARE_WITH_CURRENT_CONTEXT, 0)
	C.SDL_GL_SetAttribute(C.SDL_GL_DOUBLEBUFFER, 1)

	err := createWindow(C.SDL_WINDOW_OPENGL)

	if err != nil {
		return err
	}

	Platform.display.mainWindow.api = &glWindow{
		renderer: r,
		cContext: C.SDL_GL_CreateContext(Platform.display.mainWindow.cWindow),
	}

	switch {
	case C.SDL_GL_SetSwapInterval(-1) == 0:
		debug.LogI("vsync set to late swap tearing")
	case C.SDL_GL_SetSwapInterval(1) == 0:
		C.SDL_ClearError()
		debug.LogI("vsync enabled")
	default:
		err := debug.ErrorWrap(debug.ErrorNew(C.GoString(C.SDL_GetError())), "Failed to enable vsync")
		C.SDL_ClearError()
		return err
	}

	if err := r.GLInit(&glInstance{}); err != nil {
		return err
	}

	Platform.display.mainWindow.api.resize(Platform.config.Window.Rect.W, Platform.config.Window.Rect.H)
	debug.LogV("SDL created gl window")

	return nil
}

func (glw *glWindow) resize(w int, h int) {
	glw.windowW = w
	glw.windowH = h

	var cW, cH C.int
	C.SDL_GL_GetDrawableSize(Platform.display.mainWindow.cWindow, &cW, &cH)
	glw.renderer.Resize(int(cW), int(cH))
}

func (glw *glWindow) destroy() {
	C.SDL_GL_DeleteContext(glw.cContext)
}
