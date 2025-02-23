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
	"strings"

	"goarrg.com/debug"
	"goarrg.com/toolchain"
	"goarrg.com/toolchain/golang"
)

type Dependencies struct {
	SDL       SDLConfig
	VkHeaders VkHeadersConfig
}

/*
DebugFeatures contains build options that affects debugging, they may or may not require toolchain.BuildDebug.
*/
type DebugFeatures struct {
	/*
		If true, activates the debug.Trace* functions, without it they do nothing.
		Does not require toolchain.DebugBuild
	*/
	Trace bool
}

/*
DisableFeatures contains build options that are enabled by default but are otherwise optional.
*/
type DisableFeatures struct {
	OpenGL bool // If true, platform packages will not allow the initialization of OpenGL apps.
	Vulkan bool // If true, platform packages will not allow the initialization of Vulkan apps.
}
type BuildOptions struct {
	Build   toolchain.Build
	Debug   DebugFeatures
	Disable DisableFeatures
}

type Config struct {
	Target       toolchain.Target
	Dependencies Dependencies
	BuildOptions BuildOptions
}

func buildTags(b BuildOptions) string {
	var str string

	switch b.Build {
	case toolchain.BuildRelease:
		str = "goarrg_build_release,"
	case toolchain.BuildDevelopment:
		str = "goarrg_build_development,"
	case toolchain.BuildDebug:
		str = "goarrg_build_debug,"
	default:
		panic(debug.Errorf("Invalid build: %+v", b))
	}

	{
		if b.Debug.Trace {
			str += "goarrg_debug_enable_trace,"
		}
	}
	{
		if b.Disable.OpenGL {
			str += "goarrg_disable_gl,"
		}
		if b.Disable.Vulkan {
			str += "goarrg_disable_vk,"
		}
	}

	return strings.TrimSuffix(str, ",")
}

/*
Install will install all optional dependencies indicated and does any required
setup and then returns the required build tags to pass to `go build -tags=...`.
*/
func Install(c Config) string {
	if !golang.ValidTarget(c.Target) || !golang.CgoSupported(c.Target) {
		panic(debug.Errorf("Invalid Target: %+v", c.Target))
	}
	if c.Dependencies.VkHeaders.Install {
		err := installVkHeaders(c.Dependencies.VkHeaders)
		if err != nil {
			panic(debug.ErrorWrapf(err, "Failed to install vulkan-headers"))
		}
	}
	if c.Dependencies.SDL.Install {
		err := installSDL(c.Target, c.Dependencies.SDL)
		if err != nil {
			panic(debug.ErrorWrapf(err, "Failed to install SDL"))
		}
	}
	return buildTags(c.BuildOptions)
}
