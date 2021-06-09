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
	#define SDL_MAIN_HANDLED
	#include <SDL2/SDL.h>
	#include "event.h"

	extern int processEvents(goEvent*);

	static void setHints() {
		SDL_SetHint(SDL_HINT_NO_SIGNAL_HANDLERS, "1");
	}
*/
import "C"
import (
	"runtime"
	"sync"

	"goarrg.com"
	"goarrg.com/debug"
)

type platform struct {
	config  Config
	audio   audioSystem
	display displaySystem
	input   inputSystem
}

var (
	Platform = &platform{}
	initOnce = sync.Once{}
)

func init() {
	runtime.LockOSThread()
}

func (*platform) Init() error {
	err := debug.ErrorNew("Init must be called only once")

	initOnce.Do(func() {
		debug.LogV("Platform initializing")

		C.SDL_EventState(C.SDL_KEYDOWN, C.SDL_DISABLE)
		C.SDL_EventState(C.SDL_KEYUP, C.SDL_DISABLE)

		C.SDL_EventState(C.SDL_MOUSEMOTION, C.SDL_DISABLE)
		C.SDL_EventState(C.SDL_MOUSEBUTTONDOWN, C.SDL_DISABLE)
		C.SDL_EventState(C.SDL_MOUSEBUTTONUP, C.SDL_DISABLE)

		C.setHints()

		C.SDL_SetMainReady()
		if C.SDL_Init(C.SDL_INIT_VIDEO) != 0 {
			err = debug.ErrorNew(C.GoString(C.SDL_GetError()))
			C.SDL_ClearError()
			Popup("Platform init failed: %s", err.Error())
			return
		}

		Platform.input.init()

		err = nil
		debug.LogV("Platform initialized")
	})

	return err
}

func (*platform) Update() {
	Platform.audio.update()

	cEvent := C.goEvent{
		window: Platform.display.mainWindow.cID,
	}

	if C.processEvents(&cEvent) == 0 {
		goarrg.Shutdown()
	}

	if cEvent.windowState != 0 {
		Platform.display.mainWindow.processEvent(windowEvent{
			uint32(cEvent.windowState),
		})
	}

	Platform.input.update(cEvent)
}

func (*platform) Destroy() {
	Platform.audio.destroy()
	Platform.display.destroy()
	C.SDL_Quit()
}
