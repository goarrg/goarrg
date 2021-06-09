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

package run

import (
	"os"
	"path/filepath"

	"goarrg.com/cmd/goarrg/internal/cmd"
	"goarrg.com/cmd/goarrg/internal/cmd/build"
	"goarrg.com/cmd/goarrg/internal/exec"
	"goarrg.com/debug"
)

const (
	short = "Compile game and goarrg C dependencies if needed, then execute"
	long  = short + `

Wrapper for "go run goarrg.com/cmd/goarrg build [go args] -o {{.TmpDir}}/{{.TmpFile}}" followed by "{{.TmpDir}}/{{.TmpFile}} [game args]"
`
)

var CMD = &cmd.CMD{
	Run:   run,
	Name:  "run",
	Usage: "-- [go args] -- [game args]",
	Short: short,
	Long:  long,
	CMDs:  map[string]*cmd.CMD{},
}

func init() {
	build.AddFlags(&CMD.Flags)
}

func run(args []string) bool {
	if !cmd.PackageMain() {
		debug.LogE("Current directory is not a package main")
		os.Exit(2)
	}

	var execArgs []string

	for i, arg := range args {
		if arg == "--" {
			execArgs = args[i+1:]
			args = args[:i]
			break
		}
	}

	buildDir, err := os.MkdirTemp("", "goarrg")
	if err != nil {
		panic(err)
	}

	defer os.RemoveAll(buildDir)
	args = append(args, "-o", filepath.Join(buildDir, "goarrg.test"))

	if !build.Run(args) {
		return false
	}

	if ret := exec.RunExit(filepath.Join(buildDir, "goarrg.test"), execArgs...); ret != 0 {
		os.Exit(ret)
	}

	debug.LogI("Done running project")

	return true
}
