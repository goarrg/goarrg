/*
Copyright 2022 The goARRG Authors.

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
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"goarrg.com/debug"
	"goarrg.com/toolchain"
	"goarrg.com/toolchain/cc"
	"goarrg.com/toolchain/cgodep"
	"goarrg.com/toolchain/cmake"
	"goarrg.com/toolchain/git"
	"goarrg.com/toolchain/golang"
	"goarrg.com/toolchain/web"
)

const (
	shadercURL   = "https://github.com/google/shaderc.git"
	shadercBuild = "-goarrg0"
)

type ShadercConfig struct {
	Install     bool
	ForceStatic bool
	Build       toolchain.Build
}

func installShaderc(t toolchain.Target, c ShadercConfig, commit string) error {
	type repo struct {
		Name    string
		Site    string
		SubRepo string
		Commit  string
	}
	shadercCommit := ""
	thirdParty := map[string]repo{}
	{
		type knownGood struct {
			Commits []repo
		}
		var out knownGood
		err := web.GetJSON(fmt.Sprintf("https://raw.githubusercontent.com/google/shaderc/%s/known_good.json", commit), &out)
		if err != nil {
			return debug.ErrorWrapf(err, "Failed to get list of shaderc dependencies")
		}
		for _, r := range out.Commits {
			switch r.Name {
			case "shaderc":
				shadercCommit = r.Commit
			case "glslang", "spirv-tools", "spirv-headers":
				if r.Site != "github" {
					return debug.Errorf("shaderc dependency %q is hosted on %q while we only coded for github", r.Name, r.Site)
				}
				thirdParty[r.Name] = r
			}
		}
	}

	buildID := "shaderc-" + shadercCommit +
		"-glslang-" + thirdParty["glslang"].Commit +
		"-spirv-tools-" + thirdParty["spirv-tools"].Commit +
		"-spirv-headers-" + thirdParty["spirv-headers"].Commit
	shadercVersion := buildID
	if c.ForceStatic {
		shadercVersion += "-static"
	}
	shadercVersion += shadercBuild

	installDir := cgodep.InstallDir("shaderc", t, c.Build)
	var rebuild bool
	{
		installedVersion := cgodep.ReadVersion(installDir)
		if installedVersion != shadercVersion {
			golang.SetShouldCleanCache()
		}
		rebuild = !strings.HasPrefix(installedVersion, buildID) || !strings.HasSuffix(installedVersion, shadercBuild)
	}
	if rebuild {
		if err := os.RemoveAll(installDir); err != nil {
			return err
		}

		srcDir := filepath.Join(cgodep.CacheDir(), "shaderc")
		if err := git.CloneOrFetch(shadercURL, srcDir, git.Ref{Hash: shadercCommit}); err != nil {
			return err
		}
		for _, repo := range thirdParty {
			url := "https://github.com/" + repo.SubRepo + ".git"
			dir := filepath.Join(srcDir, "third_party", repo.Name)
			if err := git.CloneOrFetch(url, dir, git.Ref{Hash: repo.Commit}); err != nil {
				return err
			}
		}

		buildDir, err := os.MkdirTemp("", "goarrg-shaderc-build")
		if err != nil {
			return debug.ErrorWrapf(err, "Failed to make temp dir: %q", buildDir)
		}
		defer os.RemoveAll(buildDir)

		args := map[string]string{
			"CMAKE_SKIP_INSTALL_RPATH": "1", "CMAKE_SKIP_RPATH": "1",
			"ENABLE_GLSLANG_BINARIES":      "0",
			"SHADERC_SKIP_COPYRIGHT_CHECK": "1",
			"SHADERC_SKIP_EXAMPLES":        "1", "SHADERC_SKIP_EXECUTABLES": "1", "SHADERC_SKIP_TESTS": "1",
			"SPIRV_SKIP_EXECUTABLES": "1", "SPIRV_SKIP_TESTS": "1",
		}
		if t.OS == "windows" {
			isMingw32, _ := cc.FindMacro(cc.Config{}, "__MINGW32__")
			isMingw64, _ := cc.FindMacro(cc.Config{}, "__MINGW64__")
			if isMingw32 || isMingw64 {
				args["CMAKE_TOOLCHAIN_FILE"] = filepath.Join(srcDir, "cmake", "linux-mingw-toolchain.cmake")
			}
		}

		if err := cmake.Configure(t, c.Build, srcDir, buildDir, installDir, args); err != nil {
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
		return cgodep.WriteMetaFile("shaderc", t, c.Build, cgodep.Meta{
			Version: shadercVersion,
			Flags: cgodep.Flags{
				CFlags:        []string{"-I" + filepath.Join(installDir, "include")},
				LDFlags:       []string{"-L" + filepath.Join(installDir, "lib"), "-lshaderc_combined"},
				StaticLDFlags: []string{"-L" + filepath.Join(installDir, "lib"), "-lshaderc_combined"},
			},
		})
	} else {
		return cgodep.WriteMetaFile("shaderc", t, c.Build, cgodep.Meta{
			Version: shadercVersion,
			Flags: cgodep.Flags{
				CFlags:        []string{"-I" + filepath.Join(installDir, "include")},
				LDFlags:       []string{"-L" + filepath.Join(installDir, "lib"), "-lshaderc_shared"},
				StaticLDFlags: []string{"-L" + filepath.Join(installDir, "lib"), "-lshaderc_combined"},
			},
		})
	}
}
