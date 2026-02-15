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

package golang

import (
	"encoding/json"
	"runtime"
	"sync"

	"goarrg.com/debug"
	"goarrg.com/toolchain"
)

type platform struct {
	cgoSupported bool
}

var (
	platforms          = map[toolchain.Target]platform{}
	setupPlatformsOnce = sync.Once{}
)

func setupPlatforms() {
	setupPlatformsOnce.Do(func() {
		j, err := toolchain.RunOutput("go", "tool", "dist", "list", "-json")
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
			t := toolchain.Target{
				OS: entry.GOOS, Arch: entry.GOARCH,
			}
			platforms[t] = platform{
				cgoSupported: entry.CgoSupported,
			}
		}
	})
}

func ValidTarget(t toolchain.Target) bool {
	if t == (toolchain.Target{}) {
		t = toolchain.Target{
			OS:   runtime.GOOS,
			Arch: runtime.GOARCH,
		}
	}
	setupPlatforms()
	_, ok := platforms[t]
	return ok
}

func CgoSupported(t toolchain.Target) bool {
	if t == (toolchain.Target{}) {
		t = toolchain.Target{
			OS:   runtime.GOOS,
			Arch: runtime.GOARCH,
		}
	}
	setupPlatforms()
	return platforms[t].cgoSupported
}
