//go:build !goarrg_disable_gl
// +build !goarrg_disable_gl

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
*/
import "C"

import (
	"goarrg.com"
	"goarrg.com/debug"
)

type glInstance struct{}

func (gl *glInstance) ProcAddr() uintptr {
	return uintptr(C.SDL_GL_GetProcAddress)
}

func (gl *glInstance) SwapBuffers() {
	C.SDL_GL_SwapWindow(Platform.display.mainWindow.cWindow)
}

type glWindow struct {
	renderer goarrg.GLRenderer
	cContext C.SDL_GLContext
}

func glInit(r goarrg.GLRenderer) error {
	Platform.logger.IPrintf("Creating gl Window")

	if r == nil {
		err := debug.Errorf("Invalid renderer")
		Platform.logger.EPrintf("Failed to create window: Invalid renderer")
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
	case bool(C.SDL_GL_SetSwapInterval(-1)):
		Platform.logger.IPrintf("vsync set to late swap tearing")
	case bool(C.SDL_GL_SetSwapInterval(1)):
		C.SDL_ClearError()
		Platform.logger.IPrintf("vsync enabled")
	default:
		err := debug.ErrorWrapf(debug.Errorf("%s", C.GoString(C.SDL_GetError())), "Failed to enable vsync")
		C.SDL_ClearError()
		return err
	}

	if err := r.GLInit(platformInterface{}, &glInstance{}); err != nil {
		return err
	}

	Platform.display.mainWindow.api.resize(Platform.config.Window.Rect.W, Platform.config.Window.Rect.H)
	Platform.logger.IPrintf("Created gl window")

	return nil
}

func (glw *glWindow) resize(w int, h int) {
	glw.renderer.Resize(w, h)
}

func (glw *glWindow) destroy() {
	C.SDL_GL_DestroyContext(glw.cContext)
}
