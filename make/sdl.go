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

package goarrg

import (
	"bytes"
	_ "embed"
	"io"
	"os"
	"path/filepath"

	"goarrg.com/debug"
	"goarrg.com/toolchain"
	"goarrg.com/toolchain/cgodep"
	"goarrg.com/toolchain/cmake"
	"goarrg.com/toolchain/golang"
	"golang.org/x/crypto/openpgp" //nolint: staticcheck
)

const (
	sdlVersion = "2.28.4"
	sdlBuild   = sdlVersion + "-goarrg0"
	sdlSHA256  = "888b8c39f36ae2035d023d1b14ab0191eb1d26403c3cf4d4d5ede30e66a4942c"
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

var sdlCgoFlags = map[string]cgodep.Flags{
	"linux": {
		CFlags:        []string{"-D_REENTRANT"},
		LDFlags:       []string{"-lSDL2"},
		StaticLDFlags: []string{"-lSDL2-static", "-pthread", "-lm", "-lrt"},
	},
	"windows": {
		LDFlags: []string{"-mwindows", "-lSDL2"},
		StaticLDFlags: []string{
			"-mwindows", "-lSDL2-static", "-lm", "-luser32", "-lgdi32",
			"-lwinmm", "-limm32", "-lole32", "-loleaut32", "-lversion", "-luuid",
			"-ladvapi32", "-lsetupapi", "-lshell32",
		},
	},
}

var sdlLibRenames = map[toolchain.Build]map[string][][2]string{
	toolchain.BuildRelease: {
		"linux": {
			{"libSDL2.a", "libSDL2-static.a"},
		},
		"windows": {
			{"libSDL2.a", "libSDL2-static.a"},
		},
	},
	toolchain.BuildDevelopment: {
		"linux": {
			{"libSDL2.a", "libSDL2-static.a"},
		},
		"windows": {
			{"libSDL2.a", "libSDL2-static.a"},
		},
	},
	toolchain.BuildDebug: {
		"linux": {
			{"libSDL2d.a", "libSDL2-static.a"},
			{"libSDL2d.so", "libSDL2.so"},
		},
		"windows": {
			{"libSDL2d.a", "libSDL2-static.a"},
			{"libSDL2d.dll.a", "libSDL2.dll.a"},
		},
	},
}

type SDLConfig struct {
	Install bool
	Build   toolchain.Build
}

func installSDL(t toolchain.Target, c SDLConfig) error {
	installDir := cgodep.InstallDir("sdl2", t, c.Build)
	if cgodep.ReadVersion(installDir) == sdlBuild {
		return cgodep.SetActiveBuild("sdl2", t, c.Build)
	}
	if err := os.RemoveAll(installDir); err != nil {
		return err
	}

	flags, ok := sdlCgoFlags[t.OS]
	if !ok {
		return debug.Errorf("SDL2 has no build support for target: %q", t.OS)
	}
	data, err := cgodep.Get("https://github.com/libsdl-org/SDL/releases/download/release-"+sdlVersion+"/SDL2-"+sdlVersion+".tar.gz", "SDL2.tar.gz", func(target io.ReadSeeker) error {
		keyring, err := openpgp.ReadArmoredKeyRing(bytes.NewReader([]byte(sdlPK)))
		if err != nil {
			return err
		}
		if _, err = openpgp.CheckDetachedSignature(keyring, target, bytes.NewReader(sdlSig)); err != nil {
			return err
		}
		if _, err := target.Seek(0, io.SeekStart); err != nil {
			return err
		}
		return cgodep.VerifySHA256(target, sdlSHA256)
	})
	if err != nil {
		return debug.ErrorWrapf(err, "Failed to download SDL2")
	}
	defer data.Close()

	srcDir, err := os.MkdirTemp("", "goarrg-sdl")
	if err != nil {
		return debug.ErrorWrapf(err, "Failed to make temp dir: %q", srcDir)
	}
	defer os.RemoveAll(srcDir)

	debug.VPrintf("Extracting SDL2")

	if err := extractTARGZ(data, srcDir); err != nil {
		return debug.ErrorWrapf(err, "Failed to extract SDL2")
	}

	buildDir, err := os.MkdirTemp("", "goarrg-sdl-build")
	if err != nil {
		return debug.ErrorWrapf(err, "Failed to make temp dir: %q", buildDir)
	}

	defer os.RemoveAll(buildDir)

	if err := cmake.Configure(t, c.Build, srcDir, buildDir, installDir, map[string]string{
		"CMAKE_SKIP_INSTALL_RPATH": "1", "CMAKE_SKIP_RPATH": "1", "SDL_RPATH": "0",
		"SDL_DIRECTX": "0", "SDL_RENDER_D3D": "0",
	}); err != nil {
		return err
	}
	if err := cmake.Build(buildDir); err != nil {
		return err
	}
	if err := cmake.Install(buildDir); err != nil {
		return err
	}

	// rename libs to be work around SDL's weird inconsistencies
	for _, rename := range sdlLibRenames[c.Build][t.OS] {
		oldLib := filepath.Join(installDir, "lib", rename[0])
		if err := os.Rename(oldLib, filepath.Join(installDir, "lib", rename[1])); err != nil {
			return debug.ErrorWrapf(err, "Failed to rename %q", rename[0])
		}
	}

	golang.SetShouldCleanCache()
	flags.CFlags = append([]string{"-I" + filepath.Join(installDir, "include")}, flags.CFlags...)
	flags.LDFlags = append([]string{"-L" + filepath.Join(installDir, "lib")}, flags.LDFlags...)
	flags.StaticLDFlags = append([]string{"-L" + filepath.Join(installDir, "lib")}, flags.StaticLDFlags...)
	return cgodep.WriteMetaFile("sdl2", t, c.Build, cgodep.Meta{
		Version: sdlBuild,
		Flags:   flags,
	})
}
