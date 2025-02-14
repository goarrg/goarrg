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
	Target    toolchain.Target
	SDL       SDLConfig
	VkHeaders VkHeadersConfig
}

func Install(d Dependencies) {
	if !golang.ValidTarget(d.Target) || !golang.CgoSupported(d.Target) {
		panic(debug.Errorf("Invalid Target: %+v", d.Target))
	}
	if d.VkHeaders.Install {
		err := installVkHeaders(d.VkHeaders)
		if err != nil {
			panic(debug.ErrorWrapf(err, "Failed to install vulkan-headers"))
		}
	}
	if d.SDL.Install {
		err := installSDL(d.Target, d.SDL)
		if err != nil {
			panic(debug.ErrorWrapf(err, "Failed to install SDL"))
		}
	}
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

func BuildTags(b BuildOptions) string {
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
