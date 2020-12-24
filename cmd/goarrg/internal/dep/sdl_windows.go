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

package dep

import (
	"runtime"
	"strconv"

	"goarrg.com/cmd/goarrg/internal/exec"
)

const sdlConfigLinux = ``

const sdlStaticConfigLinux = ``

const sdlConfigWindows = `Libs: -L${libdir} -lmingw32 -lSDL2main -lSDL2 -mwindows
Cflags: -I${includedir} -Dmain=SDL_main`

const sdlStaticConfigWindows = `Libs: -L${libdir} -lmingw32 -lSDL2main ${libdir}/libSDL2-static.a -mwindows -Wl,--no-undefined -lm -luser32 -lgdi32 -lwinmm -limm32 -lole32 -loleaut32 -lshell32 -lsetupapi -lversion -luuid
Cflags: -I${includedir} -Dmain=SDL_main`

func sdlWindows() {
	cmakeArgs := []string{
		"-G", "MinGW Makefiles", "-DCMAKE_BUILD_TYPE=Release", "-DRPATH:BOOL=0", "-DRENDER_D3D:BOOL=0", "-DDIRECTX:BOOL=0",
		"-DCMAKE_INSTALL_PREFIX:PATH=" + usrData, "..",
	}

	if err := exec.Run("cmake", cmakeArgs...); err != nil {
		panic(err)
	}

	if err := exec.Run("mingw32-make", "-j", strconv.Itoa(runtime.NumCPU())); err != nil {
		panic(err)
	}

	if err := exec.Run("mingw32-make", "install"); err != nil {
		panic(err)
	}
}

func sdlLinux() {
	panic("No support for target os: linux")
}
