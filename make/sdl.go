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
	"bufio"
	"os"
	"path/filepath"
	"strings"

	"goarrg.com/debug"
	"goarrg.com/toolchain"
	"goarrg.com/toolchain/cgodep"
	"goarrg.com/toolchain/cmake"
	"goarrg.com/toolchain/git"
	"goarrg.com/toolchain/golang"
)

const (
	sdlURL   = "https://github.com/libsdl-org/SDL.git"
	sdlBuild = "-goarrg0"
)

type SDLConfig struct {
	Install bool
	Build   toolchain.Build
	Tag     string
}

func installSDL(t toolchain.Target, c SDLConfig) error {
	var sdlVersion string
	srcDir := filepath.Join(cgodep.CacheDir(), "sdl3")
	installDir := cgodep.InstallDir("sdl3", t, c.Build)

	{
		var ref git.Ref
		if c.Tag == "" {
			refs, err := git.GetRemoteTags(sdlURL, "release-*")
			if err != nil {
				return err
			}
			ref = refs[0]
		} else {
			refs, err := git.GetRemoteTags(sdlURL, c.Tag)
			if err != nil {
				return err
			}
			ref = refs[0]
		}
		sdlVersion = strings.TrimPrefix(ref.Name, "refs/tags/") + sdlBuild
		if cgodep.ReadVersion(installDir) == sdlVersion {
			return cgodep.SetActiveBuild("sdl3", t, c.Build)
		}
		err := git.CloneOrFetch(sdlURL, srcDir, ref)
		if err != nil {
			return err
		}
	}

	{
		if err := os.RemoveAll(installDir); err != nil {
			return err
		}

		buildDir, err := os.MkdirTemp("", "goarrg-sdl-build")
		if err != nil {
			return debug.ErrorWrapf(err, "Failed to make temp dir: %q", buildDir)
		}
		defer os.RemoveAll(buildDir)

		if err := cmake.Configure(t, c.Build, srcDir, buildDir, installDir, map[string]string{
			"CMAKE_SKIP_INSTALL_RPATH": "1", "CMAKE_SKIP_RPATH": "1", "SDL_RPATH": "0",
			"SDL_STATIC": "1",
			"SDL_CAMERA": "0", "SDL_RENDER": "0", "SDL_GPU": "0",
			"SDL_DUMMYAUDIO": "0", "SDL_DUMMYVIDEO": "0",
			"SDL_TEST_LIBRARY": "0",
			"CPACK_SOURCE_ZIP": "0", "CPACK_SOURCE_7Z": "0", "SDL_INSTALL_CPACK": "0",
		}); err != nil {
			return err
		}
		if err := cmake.Build(buildDir); err != nil {
			return err
		}
		if err := cmake.Install(buildDir); err != nil {
			return err
		}
	}

	// rename libs to be work around static linking weirdness
	{
		sdlLibRenames := [][2]string{
			{"libSDL3.a", "libSDL3-static.a"},
		}
		for _, rename := range sdlLibRenames {
			oldLib := filepath.Join(installDir, "lib", rename[0])
			if err := os.Rename(oldLib, filepath.Join(installDir, "lib", rename[1])); err != nil {
				return debug.ErrorWrapf(err, "Failed to rename %q", rename[0])
			}
		}
	}

	golang.SetShouldCleanCache()

	m := cgodep.Meta{
		Version: sdlVersion,
		Flags: cgodep.Flags{
			CFlags:        []string{"-I" + filepath.Join(installDir, "include")},
			LDFlags:       []string{"-L" + filepath.Join(installDir, "lib")},
			StaticLDFlags: []string{"-lSDL3-static"},
		},
	}
	{
		fIn, err := os.Open(filepath.Join(installDir, "lib", "pkgconfig", "sdl3.pc"))
		if err != nil {
			panic(err)
		}

		scanner := bufio.NewScanner(fIn)
		for scanner.Scan() {
			s := scanner.Text()
			switch {
			case strings.HasPrefix(s, "Libs:"):
				s = strings.TrimSpace(strings.TrimPrefix(s, "Libs:"))
				for _, arg := range strings.Split(s, " ") {
					if arg != "" && !strings.HasPrefix(arg, "-L") {
						m.Flags.LDFlags = append(m.Flags.LDFlags, arg)
					}
				}
			case strings.HasPrefix(s, "Libs.private:"):
				s = strings.TrimSpace(strings.TrimPrefix(s, "Libs.private:"))
				for _, arg := range strings.Split(s, " ") {
					if arg != "" && !strings.HasPrefix(arg, "-L") {
						m.Flags.StaticLDFlags = append(m.Flags.StaticLDFlags, arg)
					}
				}
			case strings.HasPrefix(s, "Cflags:"):
				s = strings.TrimSpace(strings.TrimPrefix(s, "Cflags:"))
				for _, arg := range strings.Split(s, " ") {
					if arg != "" && !strings.HasPrefix(arg, "-I") {
						m.Flags.CFlags = append(m.Flags.CFlags, arg)
					}
				}
			}
		}
	}
	m.Flags.StaticLDFlags = append(m.Flags.LDFlags, m.Flags.StaticLDFlags...)
	return cgodep.WriteMetaFile("sdl3", t, c.Build, m)
}
