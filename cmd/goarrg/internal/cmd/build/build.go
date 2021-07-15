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
	"os"

	"goarrg.com/cmd/goarrg/internal/cgodep"
	"goarrg.com/cmd/goarrg/internal/cmd"
	"goarrg.com/cmd/goarrg/internal/exec"
	"goarrg.com/cmd/goarrg/internal/toolchain"
	"goarrg.com/debug"
)

const (
	buildShort = `Compile game and goarrg C dependencies if needed.\nWrapper for "go build [go args]".`
	buildLong  = buildShort + ``
)

var CMD = &cmd.CMD{
	Run:   Run,
	Name:  "build",
	Usage: "-- [go args]",
	Short: buildShort,
	Long:  buildLong,
	CMDs:  map[string]*cmd.CMD{},
}

func init() {
	AddFlags(&CMD.Flags)
	toolchain.AddFlags(&CMD.Flags)
}

func Run(args []string) bool {
	if !cmd.PackageMain() {
		debug.EPrintf("Current directory is not a package main")
		os.Exit(2)
	}

	toolchain.Setup()
	cgodep.Check()

	args = append([]string{"build"}, args...)

	if DisableVK() {
		debug.IPrintf("Vulkan disabled")
		args = toolchain.AppendTag(args, "disable_vk")
	}

	if DisableGL() {
		debug.IPrintf("GL disabled")
		args = toolchain.AppendTag(args, "disable_gl")
	}

	debug.IPrintf("Building project")

	if err := exec.Run("go", args...); err != nil {
		panic(err)
	}

	debug.IPrintf("Done building project")

	return true
}
