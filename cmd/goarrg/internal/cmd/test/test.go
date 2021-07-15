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

package test

import (
	"goarrg.com/cmd/goarrg/internal/cgodep"
	"goarrg.com/cmd/goarrg/internal/cmd"
	"goarrg.com/cmd/goarrg/internal/cmd/build"
	"goarrg.com/cmd/goarrg/internal/exec"
	"goarrg.com/cmd/goarrg/internal/toolchain"
	"goarrg.com/debug"
)

const (
	short = "Wrapper for \"go test [go args]\"."
	long  = short + ``
)

var CMD = &cmd.CMD{
	Run:   run,
	Name:  "test",
	Usage: "-- [go args]",
	Short: short,
	Long:  long,
	CMDs:  map[string]*cmd.CMD{},
}

func init() {
	build.AddFlags(&CMD.Flags)
}

func run(args []string) bool {
	toolchain.Setup()
	cgodep.Check()

	args = append([]string{"test", "-v", "-count=1"}, args...)

	if build.DisableVK() {
		debug.IPrintf("Vulkan disabled")
		args = toolchain.AppendTag(args, "disable_vk")
	}

	if build.DisableGL() {
		debug.IPrintf("GL disabled")
		args = toolchain.AppendTag(args, "disable_gl")
	}

	debug.IPrintf("Testing project")

	if err := exec.Run("go", args...); err != nil {
		panic(err)
	}

	debug.IPrintf("Done testing project")

	return true
}
