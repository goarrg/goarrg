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

	"goarrg.com/cmd/goarrg/internal/cgo"
	"goarrg.com/cmd/goarrg/internal/exec"
)

const sdlConfigLinux = `Libs: -L${libdir} -lSDL2
Cflags: -I${includedir} -D_REENTRANT`

const sdlStaticConfigLinux = `Libs: -L${libdir} -Wl,-Bstatic -lSDL2 -Wl,-Bdynamic -Wl,--no-undefined -lm -ldl -lpthread -lrt
Cflags: -I${includedir} -D_REENTRANT`

const sdlConfigWindows = `Libs: -L${libdir} -lmingw32 -lSDL2main -lSDL2 -mwindows
Cflags: -I${includedir} -Dmain=SDL_main`

const sdlStaticConfigWindows = `Libs: -L${libdir} -lmingw32 -lSDL2main -Wl,-Bstatic -lSDL2 -Wl,-Bdynamic -mwindows -Wl,--no-undefined -lm -luser32 -lgdi32 -lwinmm -limm32 -lole32 -loleaut32 -lshell32 -lsetupapi -lversion -luuid
Cflags: -I${includedir} -Dmain=SDL_main`

func sdlLinux() {
	if err := exec.Run("../configure", "--disable-rpath", "--prefix="+usrData); err != nil {
		panic(err)
	}

	if err := exec.Run("make", "-j", strconv.Itoa(runtime.NumCPU())); err != nil {
		panic(err)
	}

	if err := exec.Run("make", "install"); err != nil {
		panic(err)
	}
}

func sdlWindows() {
	if err := exec.Run("../configure", "--disable-rpath", "--disable-render-d3d", "--disable-directx",
		"--host="+cgo.GCCArch()+"-w64-mingw32", "--prefix="+usrData); err != nil {
		panic(err)
	}

	if err := exec.Run("make", "-j", strconv.Itoa(runtime.NumCPU())); err != nil {
		panic(err)
	}

	if err := exec.Run("make", "install"); err != nil {
		panic(err)
	}
}
