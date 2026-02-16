/*
Copyright 2025 The goARRG Authors.

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
	spirvCrossURL   = "https://github.com/KhronosGroup/SPIRV-Cross.git"
	spirvCrossBuild = "-goarrg0"
)

type SPIRVCrossConfig struct {
	Install     bool
	ForceStatic bool
	Build       toolchain.Build
}

func installSPIRVCross(t toolchain.Target, c SPIRVCrossConfig, tag string) error {
	srcDir := filepath.Join(cgodep.CacheDir(), "spirv-cross")
	installDir := cgodep.InstallDir("spirv-cross", t, c.Build)
	spirvCrossVersion := tag
	if c.ForceStatic {
		spirvCrossVersion += "-static"
	}
	spirvCrossVersion += spirvCrossBuild
	var rebuild bool
	{
		installedVersion := cgodep.ReadVersion(installDir)
		if installedVersion != spirvCrossVersion {
			golang.SetShouldCleanCache()
		}
		rebuild = !strings.HasPrefix(installedVersion, tag) || !strings.HasSuffix(installedVersion, spirvCrossBuild)
	}
	if rebuild {
		if err := os.RemoveAll(installDir); err != nil {
			return err
		}

		{
			refs, err := git.GetRemoteTags(spirvCrossURL, tag)
			if err != nil {
				return err
			}
			if err := git.CloneOrFetch(spirvCrossURL, srcDir, refs[0]); err != nil {
				return err
			}
		}

		buildDir, err := os.MkdirTemp("", "goarrg-spirv-cross-build")
		if err != nil {
			return debug.ErrorWrapf(err, "Failed to make temp dir: %q", buildDir)
		}
		defer os.RemoveAll(buildDir)

		args := map[string]string{
			"CMAKE_SKIP_INSTALL_RPATH": "1", "CMAKE_SKIP_RPATH": "1",
			"SPIRV_CROSS_SHARED": "1",
			"SPIRV_CROSS_CLI":    "0", "SPIRV_CROSS_ENABLE_TESTS": "0",
			"SPIRV_CROSS_ENABLE_CPP": "1", "SPIRV_CROSS_ENABLE_C_API": "1",
			"SPIRV_CROSS_ENABLE_HLSL": "0", "SPIRV_CROSS_ENABLE_MSL": "0",
		}
		if err := cmake.Configure(t, toolchain.BuildRelease, srcDir, buildDir, installDir, args); err != nil {
			return err
		}
		if err := cmake.Build(buildDir); err != nil {
			return err
		}
		if err := cmake.Install(buildDir); err != nil {
			return err
		}
	}
	if c.ForceStatic {
		ldflags := []string{
			"-L" + filepath.Join(installDir, "lib"), "-lspirv-cross-c", "-lspirv-cross-cpp",
			"-lspirv-cross-glsl",
			"-lspirv-cross-util", "-lspirv-cross-core", "-lspirv-cross-reflect",
		}
		return cgodep.WriteMetaFile("spirv-cross", t, c.Build, cgodep.Meta{
			Version: spirvCrossVersion, Flags: cgodep.Flags{
				CFlags:        []string{"-I" + filepath.Join(installDir, "include")},
				LDFlags:       ldflags,
				StaticLDFlags: ldflags,
			},
		})
	} else {
		ldflags := []string{
			"-L" + filepath.Join(installDir, "lib"),
		}
		return cgodep.WriteMetaFile("spirv-cross", t, c.Build, cgodep.Meta{
			Version: spirvCrossVersion, Flags: cgodep.Flags{
				CFlags:  []string{"-I" + filepath.Join(installDir, "include")},
				LDFlags: append(ldflags, "-lspirv-cross-c-shared"),
				StaticLDFlags: append(ldflags, "-lspirv-cross-c", "-lspirv-cross-cpp",
					"-lspirv-cross-glsl",
					"-lspirv-cross-util", "-lspirv-cross-core", "-lspirv-cross-reflect"),
			},
		})
	}
}
