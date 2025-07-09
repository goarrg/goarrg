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
	#define SDL_MAIN_HANDLED 1
	#include <SDL3/SDL.h>
	#include <SDL3/SDL_main.h>
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
	logger  *debug.Logger
	config  Config
	audio   audioSystem
	display displaySystem
	input   inputSystem
}

var (
	Platform                 = &platform{logger: debug.NewLogger("sdl")}
	_        goarrg.Platform = Platform
	initOnce                 = sync.Once{}
)

func init() {
	runtime.LockOSThread()
}

func (*platform) Init() (goarrg.PlatformInterface, error) {
	err := debug.Errorf("Init must be called only once")

	initOnce.Do(func() {
		Platform.logger.IPrintf("Platform initializing")

		C.SDL_SetEventEnabled(C.SDL_EVENT_KEY_DOWN, false)
		C.SDL_SetEventEnabled(C.SDL_EVENT_KEY_UP, false)

		C.SDL_SetEventEnabled(C.SDL_EVENT_MOUSE_MOTION, false)
		C.SDL_SetEventEnabled(C.SDL_EVENT_MOUSE_BUTTON_DOWN, false)
		C.SDL_SetEventEnabled(C.SDL_EVENT_MOUSE_BUTTON_UP, false)

		C.setHints()

		C.SDL_SetMainReady()
		if !C.SDL_Init(C.SDL_INIT_VIDEO) {
			err = debug.Errorf("%s", C.GoString(C.SDL_GetError()))
			C.SDL_ClearError()
			return
		}

		Platform.input.init()

		err = nil
		Platform.logger.IPrintf("Platform initialized")
	})

	if err != nil {
		return nil, err
	}

	return platformInterface{}, nil
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

type platformInterface struct{}

func (platformInterface) Abort() {
	Abort()
}

func (platformInterface) AbortPopup(format string, args ...interface{}) {
	AbortPopup(format, args...)
}
