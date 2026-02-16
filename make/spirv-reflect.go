/*
Copyright 2026 The goARRG Authors.

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
	spirvReflectURL   = "https://github.com/KhronosGroup/SPIRV-Reflect.git"
	spirvReflectBuild = "-goarrg0"
)

type SPIRVReflectConfig struct {
	Install bool
	Build   toolchain.Build
}

func installSPIRVReflect(t toolchain.Target, c SPIRVReflectConfig, tag string) error {
	srcDir := filepath.Join(cgodep.CacheDir(), "spirv-reflect")
	installDir := cgodep.InstallDir("spirv-reflect", t, c.Build)
	spirvReflectVersion := tag + spirvReflectBuild
	var rebuild bool
	{
		installedVersion := cgodep.ReadVersion(installDir)
		if installedVersion != spirvReflectVersion {
			golang.SetShouldCleanCache()
		}
		rebuild = !strings.HasPrefix(installedVersion, tag) || !strings.HasSuffix(installedVersion, spirvReflectBuild)
	}
	if rebuild {
		if err := os.RemoveAll(installDir); err != nil {
			return err
		}

		{
			refs, err := git.GetRemoteTags(spirvReflectURL, tag)
			if err != nil {
				return err
			}
			if err := git.CloneOrFetch(spirvReflectURL, srcDir, refs[0]); err != nil {
				return err
			}
		}

		buildDir, err := os.MkdirTemp("", "goarrg-spirv-reflect-build")
		if err != nil {
			return debug.ErrorWrapf(err, "Failed to make temp dir: %q", buildDir)
		}
		defer os.RemoveAll(buildDir)

		defs := map[string]string{
			"CMAKE_SKIP_INSTALL_RPATH": "1", "CMAKE_SKIP_RPATH": "1",
			"SPIRV_REFLECT_EXECUTABLE": "0", "SPIRV_REFLECT_STATIC_LIB": "1",
		}
		if err := cmake.Configure(t, c.Build, srcDir, buildDir, installDir, defs); err != nil {
			return err
		}
		if err := cmake.Build(buildDir); err != nil {
			return err
		}
		if err := cmake.Install(buildDir); err != nil {
			return err
		}
	}
	return cgodep.WriteMetaFile("spirv-reflect", t, c.Build, cgodep.Meta{
		Version: spirvReflectVersion,
		Flags: cgodep.Flags{
			CFlags:        []string{"-I" + filepath.Join(installDir, "include")},
			LDFlags:       []string{"-L" + filepath.Join(installDir, "lib"), "-lspirv-reflect-static"},
			StaticLDFlags: []string{"-L" + filepath.Join(installDir, "lib"), "-lspirv-reflect-static"},
		},
	})
}
