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
	"goarrg.com/debug"
	"goarrg.com/toolchain"
	"goarrg.com/toolchain/golang"
)

type Config struct {
	Target    toolchain.Target
	SDL       SDLConfig
	VkHeaders VkHeadersConfig
}

func Install(c Config) {
	if !golang.ValidTarget(c.Target) || !golang.CgoSupported(c.Target) {
		panic(debug.Errorf("Invalid Target: %+v", c.Target))
	}
	if c.VkHeaders.Install {
		err := installVkHeaders(c.VkHeaders)
		if err != nil {
			panic(debug.ErrorWrapf(err, "Failed to install vulkan-headers"))
		}
	}
	if c.SDL.Install {
		err := installSDL(c.Target, c.SDL)
		if err != nil {
			panic(debug.ErrorWrapf(err, "Failed to install SDL"))
		}
	}
}
