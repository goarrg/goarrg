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

package cgodep

import (
	"bytes"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"

	"goarrg.com/cmd/goarrg/internal/exec"
	"goarrg.com/cmd/goarrg/internal/toolchain"
	"goarrg.com/debug"
)

const (
	sdlVR     = "SDL2-2.0.14"
	sdlBuild  = "goarrg0"
	sdlSHA256 = "d8215b571a581be1332d2106f8036fcb03d12a70bae01e20f424976d275432bc"
)

const sdlPK = `
-----BEGIN PGP PUBLIC KEY BLOCK-----
Version: GnuPG v1.0.4 (GNU/Linux)
Comment: For info see http://www.gnupg.org

mQGiBDpWOb0RBADQwd3d9mzt6KzqlsgXf9mikBuMbpKzYs1SBKYpdzUs9sRY0CnH
vCQTrL5sI57yKLnqEl6SbIiE75ZwrSWwvUDFFTh35Jew5nPZwv64en2kw2y4qrnJ
kBZCHDSU4KgfUZtoJ25Tmeru5MLNbXxCOoMszO5L5OchwMrGMtmFLRA/bwCgy5Th
d1/vJo+bej9tbgv++SJ05o0D/3MPK7EBoxWkQ0I+ScqOsvSMRQXWc/hXy4lyIp8e
xJByBApkv0LiiT3KlPpq/K2gTlDlCZ/JTt6Rv8Ug0g47R3a0aoz9kfc15UjHdiap
UOfF9MWmmbw59Lyx6+y2e0/C5xWzNOR1G4G5y4RZL/GXrp67xz/0fEhI85R+eASq
AEfSBAC5ZxwnBwyl+h+PXeJYKrPQjSUlgtSAkKp7PNBywwlue1LcSb7j4cc+cmgH
QMVuM883LPE59btNzFTAZjlzzIMiaXf5h9EkDARTGQ1wFiO3V5vIbVLh4kAoNfpT
egy7bYn3UrlbKg3V2DbCdEXm1zQufZzK7T0yenA5Ps8xXX7mNrQhU2FtIExhbnRp
bmdhIDxzbG91a2VuQGxpYnNkbC5vcmc+iFcEExECABcFAjpWOb0FCwcKAwQDFQMC
AxYCAQIXgAAKCRAwpZN3p3Y75t9RAJ48WI+nOPes0WK7t381Ij4JfSYxWQCgjpMa
Dg3/ah23HZhYtTKtHUzD9zi5AQ0EOlY5wxAEAPvjB0B5RNAj8hBF/Lq78w5rJ1/f
5RqWXmdfxApuEE/9OEFXUSUXms9f/IWvySdyf48Pk4t2h8b8i7F0f3R+tcCp6m0P
t1BSNHYumfmtonTy5FHqpwBVlEi7I0s5mD3kxO+k8PQbATHH5smFnoz2UTc+MzQj
UdtTzXUkUgqvf9zTAAMGA/9Y/h6rhi3YYXeI6SmbXqcmzsQKzaWVhLew67szejnY
sKIJ1ja4MefYlthCXgmIBriNftxIGtBI0Pcmzwpn0eknRNK6NgpmESbGKCWh59Je
iAK5hdBPe47LSFVct5zSO9vQhRDyLzhzPPtB3XeoKTUkLWxBSLbeZVwcHPIK/wLa
l4hGBBgRAgAGBQI6VjnDAAoJEDClk3endjvmxmUAn3Ne6Z3UULpie8RJP15RBt6K
2MWFAJ9hK/Ls/FeBJ9d50qxmYdZ2RrTXNg==
=toqC
-----END PGP PUBLIC KEY BLOCK-----
`

//go:embed sdl.sig
var sdlSig []byte

var sdlCgoFlags = map[string]cgoFlags{
	"linux": {
		CFlags:        "-D_REENTRANT",
		LDFlags:       "-lSDL2 -pthread",
		StaticLDFlags: "-lSDL2-static -Wl,--no-undefined -lm -ldl -lpthread -lrt",
	},
	"windows": {
		LDFlags:       "-lmingw32 -lSDL2 -mwindows",
		StaticLDFlags: "-lmingw32 -lSDL2-static -mwindows -Wl,--no-undefined -lm -luser32 -lgdi32 -lwinmm -limm32 -lole32 -loleaut32 -lshell32 -lsetupapi -lversion -luuid",
	},
}

const (
	sdlShort = "Installs SDL2 as required by platform/sdl"
	sdlLong  = sdlShort + `

Available targets: %q
`
)

func init() {
	targetList := make([]string, 0, len(sdlCgoFlags))
	for t := range sdlCgoFlags {
		targetList = append(targetList, t)
	}

	targetList = sort.StringSlice(targetList)

	cgoDeps["sdl2"] = cgoDep{
		name:           "sdl2",
		short:          sdlShort,
		long:           fmt.Sprintf(sdlLong, targetList),
		version:        sdlVR + "-" + sdlBuild,
		targetSpecific: true,
		install:        sdlInstall,
	}
}

func sdlInstall(installDir string) cgoFlags {
	cgoFlags, ok := sdlCgoFlags[toolchain.TargetOS()]
	if !ok {
		debug.LogE("sdl2 is not available on target: %q\nRun %q to see available targets", toolchain.Target(), "go run goarrg.com/cmd/goarrg install sdl2 -h")
		os.Exit(2)
	}

	data, err := get("https://www.libsdl.org/release/"+sdlVR+".tar.gz", func(data []byte) error {
		if err := verifyPGP(data, []byte(sdlPK), sdlSig); err != nil {
			return err
		}
		return verifySHA256(data, sdlSHA256)
	})
	if err != nil {
		panic(debug.ErrorWrapf(err, "Failed to download SDL2"))
	}

	srcDir, err := os.MkdirTemp("", "goarrg-sdl")
	if err != nil {
		panic(debug.ErrorWrapf(err, "Failed to make temp dir: %q", srcDir))
	}

	defer os.RemoveAll(srcDir)

	debug.LogV("Extracting SDL2")

	if err := extractTARGZ(bytes.NewReader(data), srcDir); err != nil {
		panic(debug.ErrorWrapf(err, "Failed to extract SDL2"))
	}

	buildDir, err := os.MkdirTemp("", "goarrg-sdl-build")
	if err != nil {
		panic(debug.ErrorWrapf(err, "Failed to make temp dir: %q", buildDir))
	}

	defer os.RemoveAll(buildDir)

	cmakeArgs := []string{
		// cmake expects Linux/Windows not linux/windows
		"-DCMAKE_SYSTEM_NAME=" + strings.Title(toolchain.TargetOS()),
		"-DCMAKE_BUILD_TYPE=Release", "-DRPATH=0", "-DRENDER_D3D=0", "-DDIRECTX=0",
	}

	cmakeArgs = append(cmakeArgs, "-DCMAKE_INSTALL_PREFIX="+installDir, "-S", srcDir, "-B", buildDir)

	// windows will default to msvc even with CC/CXX set and we don't want that
	if runtime.GOOS == "windows" {
		cmakeArgs = append(cmakeArgs, "-G", "MinGW Makefiles")
	}

	if err := exec.Run("cmake", cmakeArgs...); err != nil {
		panic(err)
	}

	if err := exec.Run("cmake", "--build", buildDir, "-j", strconv.Itoa(runtime.NumCPU())); err != nil {
		panic(err)
	}

	if err := exec.Run("cmake", "--install", buildDir); err != nil {
		panic(err)
	}

	// linux builds name the static lib differently from windows, so rename it.
	// The -static suffix also makes it easier to static link without the
	// "-Wl,-Bstatic -SDL2 -Wl,-Bdynamic" nonsense.
	if toolchain.TargetOS() == "linux" {
		err := os.Rename(filepath.Join(installDir, "lib", "libSDL2.a"), filepath.Join(installDir, "lib", "libSDL2-static.a"))
		if err != nil {
			panic(err)
		}
	}

	return cgoFlags
}
