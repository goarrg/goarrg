/*
Copyright 2023 The goARRG Authors.

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

package debug

/*
	#include <stdint.h>
	#include <stdio.h>
	#include <io.h>

	__attribute__((visibility("hidden"))) static inline void init_std(uintptr_t* hIn, uintptr_t* hOut, uintptr_t* hErr) {
		freopen("CONIN$", "r", stdin);
		freopen("CONOUT$", "w", stdout);
		freopen("CONOUT$", "w", stderr);

		*hIn = _get_osfhandle(_fileno(stdin));
		*hOut = _get_osfhandle(_fileno(stdout));
		*hErr = _get_osfhandle(_fileno(stderr));
	}
*/
import "C"

import (
	"os"

	"golang.org/x/sys/windows"
)

func init() {
	// if we have a working stdout, do nothing, Attach/Alloc will mess up pipes
	if _, err := os.Stdout.Stat(); err == nil {
		return
	}

	if initTerminal() {
		var hIn, hOut, hErr C.uintptr_t
		C.init_std(&hIn, &hOut, &hErr)

		// we must use the handles we get from _get_osfhandle cause GetStdHandle does not update/work sometimes
		_ = windows.SetStdHandle(windows.STD_INPUT_HANDLE, windows.Handle(hIn))
		_ = windows.SetStdHandle(windows.STD_OUTPUT_HANDLE, windows.Handle(hOut))
		_ = windows.SetStdHandle(windows.STD_ERROR_HANDLE, windows.Handle(hErr))

		os.Stdin = os.NewFile(uintptr(hIn), "/dev/stdin")
		os.Stdout = os.NewFile(uintptr(hOut), "/dev/stdout")
		os.Stderr = os.NewFile(uintptr(hErr), "/dev/stderr")

		loggerOut = os.Stderr
	}
}
