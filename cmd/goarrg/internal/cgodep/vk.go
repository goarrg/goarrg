/*
Copyright 2021 The goARRG Authors.

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
	"errors"
	"os"
	"os/exec"
	"path/filepath"

	"goarrg.com/cmd/goarrg/internal/toolchain"
	"goarrg.com/debug"
)

const (
	vkBuild = "goarrg0"
	vkShort = `Installs the vulkan goarrg-config file to be used with "#cgo pkg-config: vulkan".`
	vkLong  = vkShort + `

Searches ${VULKAN_SDK} before asking the system default C compiler for the location.
When cross compiling, if ${VULKAN_SDK} is unset, only the headers will be available.
Config files will be automatically updated when the SDK/header location changes based on the rules above.`
)

var (
	vkSDK     string
	vkHeaders string
	vkErr     error
)

func init() {
	vkSDK = os.Getenv("VULKAN_SDK")

	if vkSDK == "" {
		// search using system default as we do not care when the cross compiler has the headers
		gcc, _ := exec.LookPath("gcc")
		clang, _ := exec.LookPath("clang")

		var errGCC, errClang error

		if gcc != "" {
			vkHeaders, errGCC = toolchain.FindHeader(gcc, "vulkan/vulkan.h")
		}

		if vkHeaders == "" && clang != "" {
			vkHeaders, errClang = toolchain.FindHeader(clang, "vulkan/vulkan.h")
		}

		if vkHeaders == "" {
			switch {
			case errGCC != nil && errClang != nil:
				vkErr = debug.Errorf("Failed to find \"vulkan/vulkan.h\" with gcc and clang:\n%s\n%s", errGCC.Error(), errClang.Error())
			case errGCC != nil:
				vkErr = errGCC
			case errClang != nil:
				vkErr = errClang
			default:
				vkErr = debug.Errorf("Unable to find gcc or clang")
			}
		}
	}

	// give each one a different version specific to sdk/header location so we
	// can detect changes and update the file.
	// err also gets a version so we can panic during cgodep.Check() if somehow
	// the sdk/headers were found at one time but not anymore.
	// since we don't actually write the meta file, it should always panic.

	if vkSDK != "" {
		cgoDeps["vulkan"] = cgoDep{
			name:    "vulkan",
			short:   vkShort,
			long:    vkLong,
			version: vkSDK + "-" + vkBuild,
			install: vulkanInstall,
		}
		return
	}
	if vkHeaders != "" {
		cgoDeps["vulkan"] = cgoDep{
			name:    "vulkan",
			short:   vkShort,
			long:    vkLong,
			version: vkHeaders + "-" + vkBuild,
			install: vulkanInstall,
		}
		return
	}
	if vkErr != nil {
		cgoDeps["vulkan"] = cgoDep{
			name:    "vulkan",
			short:   vkShort,
			long:    vkLong,
			version: vkBuild,
			install: vulkanInstall,
		}
		return
	}

	panic("Should not be here")
}

func vulkanInstall(installDir string) cgoFlags {
	if vkSDK != "" {
		debug.LogI("${VULKAN_SDK} set to: %q", vkSDK)
		// vulkan caps the first letter on the windows SDK ... why!
		dir := filepath.Join(vkSDK, "Include")
		stat, err := os.Stat(dir)
		if err != nil {
			if !errors.Is(err, os.ErrNotExist) {
				panic(err)
			}
			dir = filepath.Join(vkSDK, "include")
			if stat, err := os.Stat(dir); err != nil {
				panic(debug.ErrorWrapf(err, "Failed to find include directory: %q", dir))
			} else if !stat.IsDir() {
				panic(debug.Errorf("Expected %q to be a directory", dir))
			}

			return cgoFlags{
				CFlags:        "-I" + filepath.Join(vkSDK, "include"),
				LDFlags:       "-L" + filepath.Join(vkSDK, "lib"),
				StaticLDFlags: "-L" + filepath.Join(vkSDK, "lib"),
			}
		}

		if !stat.IsDir() {
			panic(debug.Errorf("Expected %q to be a directory", dir))
		}

		return cgoFlags{
			CFlags:        "-I" + filepath.Join(vkSDK, "Include"),
			LDFlags:       "-L" + filepath.Join(vkSDK, "Lib"),
			StaticLDFlags: "-L" + filepath.Join(vkSDK, "Lib"),
		}
	}

	debug.LogI("${VULKAN_SDK} unset, searching system")

	if vkHeaders != "" {
		debug.LogI("Found vulkan headers at: %q", vkHeaders)
		if err := os.MkdirAll(filepath.Join(installDir, "include"), 0o755); err != nil {
			panic(err)
		}
		if err := os.Symlink(vkHeaders, filepath.Join(installDir, "include", "vulkan")); err != nil {
			panic(err)
		}
		return cgoFlags{}
	}

	if vkErr != nil {
		panic(vkErr)
	}

	panic("Should not be here")
}

func HaveVK() bool {
	// by the time we are here, all checks should've completed so no need to check
	// if they are actually valid locations
	return vkSDK != "" || vkHeaders != ""
}
