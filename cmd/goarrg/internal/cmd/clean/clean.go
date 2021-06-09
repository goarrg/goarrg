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

package clean

import (
	"goarrg.com/cmd/goarrg/internal/cgodep"
	"goarrg.com/cmd/goarrg/internal/cmd"
	"goarrg.com/cmd/goarrg/internal/exec"
	"goarrg.com/cmd/goarrg/internal/toolchain"
)

const (
	short = `Wrapper for "go clean [go args]".`
	long  = short + ``
)

var CMD = &cmd.CMD{
	Run:   run,
	Name:  "clean",
	Usage: "-- [go args]",
	Short: short,
	Long:  long,
}

var (
	flagCgoDep      bool
	flagCgoDepCache bool
)

func init() {
	CMD.Flags.BoolVar(&flagCgoDep, "cgodep", false, "Also remove all built C dependencies, they will be rebuilt as needed.")
	CMD.Flags.BoolVar(&flagCgoDepCache, "cgodepcache", false, "Also remove the C dependencies downloaded files, they will be redownloaded as needed.")
}

func run(args []string) bool {
	toolchain.Setup()

	args = append([]string{"clean"}, args...)
	if err := exec.Run("go", args...); err != nil {
		panic(err)
	}

	if flagCgoDep {
		cgodep.Clean()
	}
	if flagCgoDepCache {
		cgodep.CleanCache()
	}

	return true
}
