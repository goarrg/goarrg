/*
Copyright 2020 The goARRG Authors.

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

package build

import (
	"flag"
	"os"

	"goarrg.com/cmd/goarrg/internal/cgodep"
	"goarrg.com/cmd/goarrg/internal/toolchain"
	"goarrg.com/debug"
)

var (
	flagDisableVK bool
	flagDisableGL bool
)

func AddFlags(f *flag.FlagSet) {
	f.BoolVar(&flagDisableVK, "disable_vk", false, "Disables all Vulkan related functionality")
	f.BoolVar(&flagDisableGL, "disable_gl", false, "Disables all OpenGL related functionality")
}

func DisableVK() bool {
	if !cgodep.HaveVK() {
		debug.IPrintf("Vulkan SDK not found")
		return true
	}
	return flagDisableVK
}

func DisableGL() bool {
	// only check for glu.h since there should be no way to install it without gl.h
	// check here and not init as we need to know what the actual compiler is
	// and we won't know that at init
	_, err := toolchain.FindHeader(os.Getenv("CC"), "GL/glu.h")
	if err != nil {
		debug.IPrintf("Unable to find: %q", "GL/glu.h")
		return true
	}
	return flagDisableGL
}
