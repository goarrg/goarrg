/*
Copyright 2021 The goARRG Authors.

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

package toolchain

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"goarrg.com/cmd/goarrg/internal/cmd"
	"goarrg.com/cmd/goarrg/internal/exec"
	"goarrg.com/debug"
)

func CrossCompiling() bool {
	return os.Getenv("GOOS") != runtime.GOOS || os.Getenv("GOARCH") != runtime.GOARCH
}

func printEnv() {
	env, err := exec.RunOutput("go", "env")
	if err != nil {
		panic(err)
	}

	if targetOS == "windows" {
		if runtime.GOOS == "windows" {
			debug.LogI("go env\n%sset RC=%s", env, os.Getenv("RC"))
		} else {
			debug.LogI("go env\n%sRC=%s", env, os.Getenv("RC"))
		}
	} else {
		debug.LogI("go env\n%s", env)
	}
}

func lookPathSetEnv(env, value string) {
	if _, err := exec.LookPath(value); err != nil {
		panic(debug.ErrorWrap(err, "Unable to find: %q", value))
	}
	if err := os.Setenv(env, value); err != nil {
		panic(err)
	}
}

func Setup() {
	debug.LogI("Setting up go")

	{
		gocache := filepath.Join(cmd.ModuleDataPath(), "gocache")
		if err := os.MkdirAll(gocache, 0o755); err != nil {
			panic(err)
		}
		if err := os.Setenv("GOCACHE", gocache); err != nil {
			panic(err)
		}
	}

	p := platforms[flagTarget]

	if !p.cgoSupported {
		if err := os.Setenv("GOOS", targetOS); err != nil {
			panic(err)
		}
		if err := os.Setenv("GOARCH", targetArch); err != nil {
			panic(err)
		}
		debug.LogW("cgo unsupported on target: %q", flagTarget)
		printEnv()
		return
	}

	err := os.Setenv("CGO_ENABLED", "1")
	if err != nil {
		panic(err)
	}

	// must be called before setting GOOS/GOARCH
	setupCgoTools()

	{
		if err := os.Setenv("GOOS", targetOS); err != nil {
			panic(err)
		}
		if err := os.Setenv("GOARCH", targetArch); err != nil {
			panic(err)
		}
	}

	if !CrossCompiling() {
		// set CC/CXX/AR for easy access
		if cc, err := exec.RunOutput("go", "env", "CC"); err != nil {
			panic(err)
		} else if err := os.Setenv("CC", strings.TrimSpace(string(cc))); err != nil {
			panic(err)
		}

		if cxx, err := exec.RunOutput("go", "env", "CXX"); err != nil {
			panic(err)
		} else if err := os.Setenv("CXX", strings.TrimSpace(string(cxx))); err != nil {
			panic(err)
		}

		if ar, err := exec.RunOutput("go", "env", "AR"); err != nil {
			panic(err)
		} else if err := os.Setenv("AR", strings.TrimSpace(string(ar))); err != nil {
			panic(err)
		}

		printEnv()
		return
	}

	debug.LogI("Detected cross compiling target: %s", flagTarget)

	{
		ccSet := os.Getenv("CC") != ""
		cxxSet := os.Getenv("CXX") != ""
		arSet := os.Getenv("AR") != ""
		rcSet := os.Getenv("RC") != ""

		// we only need RC on windows
		if targetOS != "windows" {
			rcSet = true
		}

		if ccSet && cxxSet && arSet && rcSet {
			printEnv()
			return
		}
	}

	if p.toolchain == (toolchain{}) {
		panic(debug.ErrorNew("CC/CXX/AR unset, no known defaults for target: %q", flagTarget))
	}

	debug.LogI("Setting cgo toolchain for target: %s", flagTarget)

	lookPathSetEnv("CC", p.toolchain.cc)
	lookPathSetEnv("CXX", p.toolchain.cxx)
	lookPathSetEnv("AR", p.toolchain.ar)

	// C/C++ deps may require RC when targeting windows.
	if targetOS == "windows" {
		arch := gccArch(targetArch)
		abi := gccABI(targetOS)

		rc := arch + "-" + abi + "-windres"
		lookPathSetEnv("RC", rc)
	}

	printEnv()
}

func AppendTag(args []string, tag string) []string {
	haveTagsArg := false
	for i, arg := range args {
		if strings.HasPrefix(arg, "-tags") {
			arg = strings.TrimPrefix(arg, "-tags")
			if strings.HasPrefix(arg, "=") {
				arg = strings.TrimPrefix(arg, "=")
				args[i] = "-tags=" + strings.ReplaceAll(arg, " ", ",") + "," + tag
			} else {
				arg = args[i+1]
				args[i+1] = strings.ReplaceAll(arg, " ", ",") + "," + tag
			}
			haveTagsArg = true
			break
		}
	}

	if !haveTagsArg {
		args = append(args, "-tags="+tag)
	}

	return args
}
