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
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"goarrg.com/cmd/goarrg/internal/cgodep"
	"goarrg.com/cmd/goarrg/internal/cmd"
	"goarrg.com/cmd/goarrg/internal/exec"
	"goarrg.com/cmd/goarrg/internal/toolchain"
	"goarrg.com/debug"
	"golang.org/x/tools/go/packages"
)

const (
	yourselfShort = `Installs all available C dependencies, does a clean build, and test goarrg.
Tests are skipped while cross compiling.`
	yourselfLong = yourselfShort + ``
)

type testEvent struct {
	Time    time.Time // encodes as an RFC3339-format string
	Action  string
	Package string
	Test    string
	Elapsed float64 // seconds
	Output  string
}

func parseTestFailures(data []byte) string {
	decoder := json.NewDecoder(bytes.NewReader(data))
	testOutputs := map[string]string{}
	output := ""

	for decoder.More() {
		e := testEvent{}
		err := decoder.Decode(&e)
		if err != nil {
			panic(err)
		}

		// ignore extra "FAIL\n" output
		if e.Output == "FAIL\n" {
			continue
		}
		if e.Action == "run" {
			continue
		}
		if e.Action == "output" {
			testOutputs[(e.Package + e.Test)] += e.Output
		}
		if e.Action == "fail" {
			output += testOutputs[(e.Package + e.Test)]
		}
	}

	if output == "" {
		panic(fmt.Sprintf("No test failures found: %s", data))
	}

	return output
}

func init() {
	CMD.CMDs["yourself"] = &cmd.CMD{
		Run:   runYourself,
		Name:  "yourself",
		Short: yourselfShort,
		Long:  yourselfLong,
	}
	AddFlags(&CMD.CMDs["yourself"].Flags)
	toolchain.AddFlags(&CMD.CMDs["yourself"].Flags)
}

func runYourself(args []string) bool {
	toolchain.Setup()

	// force goarrg to install SDL2 if not already
	err := os.MkdirAll(filepath.Join(cmd.ModuleDataPath(), "cgodep", "sdl2", toolchain.TargetOS()+"_"+toolchain.TargetArch()), 0o755)
	if err != nil {
		panic(err)
	}

	args = append([]string{"build"}, args...)

	if DisableVK() {
		debug.IPrintf("Vulkan disabled")
		args = toolchain.AppendTag(args, "disable_vk")
	} else {
		// force goarrg to install vulkan if not already
		err := os.MkdirAll(filepath.Join(cmd.ModuleDataPath(), "cgodep", "vulkan"), 0o755)
		if err != nil {
			panic(err)
		}
	}

	if DisableGL() {
		debug.IPrintf("GL disabled")
		args = toolchain.AppendTag(args, "disable_gl")
	}

	cgodep.Check()

	// if we are running "build yourself" then we are likely developing goarrg
	// setting up a new environment or having issues.
	// so do a clean build everytime to avoid cache issues
	if err := exec.Run("go", "clean", "-cache"); err != nil {
		panic(err)
	}

	debug.IPrintf("Building goarrg")

	var pkgs []string

	list, err := packages.Load(&packages.Config{Mode: packages.NeedName}, "goarrg.com/...")
	if err != nil {
		panic(err)
	}

	for _, pkg := range list {
		// ignore examples and tests
		// examples and tests "should" run later with "go test ..."
		if !strings.Contains(pkg.PkgPath, "/example") && !strings.Contains(pkg.PkgPath, "/test") {
			pkgs = append(pkgs, pkg.PkgPath)
		}
	}

	if err := exec.Run("go", append(args, pkgs...)...); err != nil {
		panic(err)
	}

	if err := exec.Run("go", append(args, pkgs...)...); err != nil {
		panic(err)
	}

	if !toolchain.CrossCompiling() {
		os.Setenv("GODEBUG", "cgocheck=2")

		if cmd.VeryVerbose() {
			// we need to run with "-p=1" to get streaming test output, tho this disables parallel testing of packages
			args[0] = "test"
			args = append(args, "-v", "-count=1", "-p=1", "goarrg.com/...")
			if err := exec.Run("go", args...); err != nil {
				os.Exit(2)
			}

			args = toolchain.AppendTag(args, "debug")
			if err := exec.Run("go", args...); err != nil {
				os.Exit(2)
			}
		} else {
			args[0] = "test"
			args = append(args, "-count=1", "-json", "goarrg.com/...")
			if out, err := exec.RunOutput("go", args...); err != nil {
				debug.EPrintf("Tests failed:\n%s", parseTestFailures(out))
				os.Exit(2)
			}

			args = toolchain.AppendTag(args, "debug")
			if out, err := exec.RunOutput("go", args...); err != nil {
				debug.EPrintf("Tests failed:\n%s", parseTestFailures(out))
				os.Exit(2)
			}
		}

		os.Setenv("GODEBUG", "")
	}

	debug.IPrintf("Done building goarrg")

	return true
}
