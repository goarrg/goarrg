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

package cmake

import (
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"goarrg.com/toolchain"
	"goarrg.com/toolchain/cc"
)

/*
Configure configures the cmake project with the given args.
Requires cmake 3.13 for "-S" and "-B"
*/
func Configure(t toolchain.Target, b toolchain.Build, srcDir, buildDir, installDir string, defines map[string]string) error {
	cmakeArgs := []string{
		// cmake expects Linux/Windows not linux/windows
		"-DCMAKE_SYSTEM_NAME=" + strings.Title(t.OS), // nolint: staticcheck
		// "-DCMAKE_SYSTEM_PROCESSOR=",
		"-DCMAKE_INSTALL_PREFIX=" + installDir,
		"-DCMAKE_COLOR_MAKEFILE=0", // these persists even after cmake exits and are annoying
	}

	// cmake does not always use CC/CXX so we have to set these flags too
	if cc, err := exec.LookPath(toolchain.EnvGet("CC")); err == nil {
		cmakeArgs = append(cmakeArgs, "-DCMAKE_C_COMPILER="+filepath.ToSlash(cc))
	}
	if cxx, err := exec.LookPath(toolchain.EnvGet("CXX")); err == nil {
		cmakeArgs = append(cmakeArgs, "-DCMAKE_CXX_COMPILER="+filepath.ToSlash(cxx))
	}

	/*
		if toolchain.EnvGet("CFLAGS") != "" {
			cmakeArgs = append(cmakeArgs, "-DCMAKE_C_FLAGS_DEBUG=", "-DCMAKE_C_FLAGS_RELWITHDEBINFO=", "-DCMAKE_C_FLAGS_MINSIZEREL=", "-DCMAKE_C_FLAGS_RELEASE=")
		}
		if toolchain.EnvGet("CXXFLAGS") != "" {
			cmakeArgs = append(cmakeArgs, "-DCMAKE_CXX_FLAGS_DEBUG=", "-DCMAKE_CXX_FLAGS_RELWITHDEBINFO=", "-DCMAKE_CXX_FLAGS_MINSIZEREL=", "-DCMAKE_CXX_FLAGS_RELEASE=")
		}
	*/

	switch b {
	case toolchain.BuildRelease:
		cmakeArgs = append(cmakeArgs, "-DCMAKE_BUILD_TYPE=Release")
	case toolchain.BuildDevelopment:
		cmakeArgs = append(cmakeArgs, "-DCMAKE_BUILD_TYPE=RelWithDebInfo")
	case toolchain.BuildDebug:
		cmakeArgs = append(cmakeArgs, "-DCMAKE_BUILD_TYPE=Debug")
	}

	for k, v := range defines {
		cmakeArgs = append(cmakeArgs, "-D"+k+"="+v)
	}

	// these must come after -D args
	cmakeArgs = append(cmakeArgs, "-S", srcDir, "-B", buildDir)

	// windows will default to msvc even with CC/CXX set and we don't want that
	if runtime.GOOS == "windows" {
		isMingw32, _ := cc.FindMacro(cc.Config{}, "__MINGW32__")
		isMingw64, _ := cc.FindMacro(cc.Config{}, "__MINGW64__")
		if isMingw32 || isMingw64 {
			cmakeArgs = append(cmakeArgs, "-G", "MinGW Makefiles")
		}
	}

	return toolchain.Run("cmake", cmakeArgs...)
}

/*
Build will build the pre-configured cmake project located at buildDir.
Requires cmake 3.0 for "--build"
*/
func Build(buildDir string) error {
	return toolchain.Run("cmake", "--build", buildDir, "-j", strconv.Itoa(runtime.NumCPU()))
}

/*
Install will install the prebuilt cmake project located at buildDir.
Requires cmake 3.15 for "--install"
*/
func Install(buildDir string) error {
	return toolchain.Run("cmake", "--install", buildDir)
}
