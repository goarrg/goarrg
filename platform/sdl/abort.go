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

	int Popup(_GoString_ err) {
		return SDL_ShowSimpleMessageBox(SDL_MESSAGEBOX_ERROR, "Error", _GoStringPtr(err), NULL);
	}
*/
import "C"

import (
	"fmt"
	"runtime"
)

func Abort() {
	panic("Fatal Error")
}

func AbortPopup(format string, args ...interface{}) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	err := fmt.Sprintf(format, args...)
	Platform.logger.IPrintf("Displaying AbortPopup with message:\n%s", err)
	if C.Popup(err+"\x00") != 0 {
		Platform.logger.EPrintf("Failed to create popup: %s", C.GoString(C.SDL_GetError()))
		C.SDL_ClearError()
	}
	panic("Fatal Error")
}
