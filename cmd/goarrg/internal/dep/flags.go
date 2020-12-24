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

package dep

import "flag"

var flagDep bool
var flagNoDep bool
var flagDisableVK bool
var flagDisableGL bool

func SetFlags(f *flag.FlagSet) {
	f.BoolVar(&flagDep, "dep", false, "Builds C dependencies")
	f.BoolVar(&flagNoDep, "nodep", false, "Do not build C dependencies")
	f.BoolVar(&flagDisableVK, "disable_vk", false, "Disables all Vulkan related functionality")
	f.BoolVar(&flagDisableGL, "disable_gl", false, "Disables all OpenGL related functionality")
}

func FlagDisableVK() bool {
	return flagDisableVK
}

func FlagDisableGL() bool {
	return flagDisableGL
}
