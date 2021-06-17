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

package toolchain

import (
	"encoding/json"
	"os/exec"

	"goarrg.com/debug"
)

type toolchain struct {
	cc  string
	cxx string
	ar  string
}

type platform struct {
	cgoSupported bool
	toolchain    toolchain
}

var platforms = map[string]platform{}

func gccArch(goarch string) string {
	switch goarch {
	case "amd64":
		return "x86_64"
	case "386":
		return "i686"
	}
	return ""
}

func gccABI(goos string) string {
	switch goos {
	case "linux":
		// go does not support musl, it does support clang but we have to pick one
		// and clang does not have a arch/abi specific name nor is it easy to find
		// the supported targets so leave it to the user to setup.
		return "linux-gnu"
	case "windows":
		// go does not support msvc, clang (on windows) or mingw32.
		return "w64-mingw32"
	}
	return ""
}

func init() {
	cmd := exec.Command("go", "tool", "dist", "list", "-json")
	j, err := cmd.Output()
	if err != nil {
		panic(debug.ErrorWrapf(err, "Failed to get list of platforms"))
	}

	list := []struct {
		GOOS         string
		GOARCH       string
		CgoSupported bool
	}{}
	if err := json.Unmarshal(j, &list); err != nil {
		panic(debug.ErrorWrapf(err, "Failed to unmarshal json"))
	}

	for _, entry := range list {
		name := entry.GOOS + "_" + entry.GOARCH
		platforms[name] = platform{
			cgoSupported: entry.CgoSupported,
		}

		arch := gccArch(entry.GOARCH)
		if arch == "" {
			continue
		}

		abi := gccABI(entry.GOOS)
		if abi == "" {
			continue
		}

		p := platforms[name]
		p.toolchain = toolchain{
			cc:  arch + "-" + abi + "-gcc",
			cxx: arch + "-" + abi + "-g++",
			ar:  arch + "-" + abi + "-gcc-ar",
		}

		platforms[name] = p
	}
}

func ValidPlatform(p string) bool {
	_, ok := platforms[p]
	return ok
}

func CgoSupported() bool {
	return platforms[flagTarget].cgoSupported
}
