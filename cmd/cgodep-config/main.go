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

package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"goarrg.com/toolchain"
	"goarrg.com/toolchain/cgodep"
)

func main() {
	cFlags := false
	cFlagsOnlyPath := false
	cFlagsOnlyOther := false
	libs := false
	libsOnlyPath := false
	libsOnly := false
	libsOnlyOther := false
	static := false
	version := false
	exists := false
	modVersion := false
	variable := ""
	flag.BoolVar(&cFlags, "cflags", false, "")
	flag.BoolVar(&cFlagsOnlyPath, "cflags-only-I", false, "")
	flag.BoolVar(&cFlagsOnlyOther, "cflags-only-other", false, "")
	flag.BoolVar(&libs, "libs", false, "")
	flag.BoolVar(&libsOnlyPath, "libs-only-L", false, "")
	flag.BoolVar(&libsOnly, "libs-only-l", false, "")
	flag.BoolVar(&libsOnlyOther, "libs-only-other", false, "")
	flag.BoolVar(&static, "static", false, "")
	flag.BoolVar(&version, "version", false, "")
	flag.BoolVar(&exists, "exists", false, "")
	flag.BoolVar(&modVersion, "modversion", false, "")
	flag.Bool("print-errors", false, "")
	flag.Bool("short-errors", false, "")
	flag.StringVar(&variable, "variable", "", "")
	flag.Parse()

	if version {
		fmt.Println("0.0.0-goarrg0")
		os.Exit(0)
	}

	mode := cgodep.ResolveMode(0)

	if cFlags {
		mode |= cgodep.ResolveCFlags
	}
	if libs {
		mode |= cgodep.ResolveLDFlags
	}
	if static {
		mode |= cgodep.ResolveStaticFlags
	}
	if exists {
		mode |= cgodep.ResolveExists
	}

	t := toolchain.Target{
		OS:   os.Getenv("GOOS"),
		Arch: os.Getenv("GOARCH"),
	}

	if t.OS == "" {
		t.OS = runtime.GOOS
	}

	if t.Arch == "" {
		t.Arch = runtime.GOARCH
	}

	outputFlags, err := cgodep.Resolve(t, mode, flag.Args()...)
	if err != nil {
		fallback := toolchain.EnvGet("CGODEP_PKG_CONFIG")
		if fallback == "" {
			fallback = "pkg-config"
		}
		args := []string{}

		if cFlags {
			args = append(args, "--cflags")
		}
		if cFlagsOnlyPath {
			args = append(args, "--cflags-only-I")
		}
		if cFlagsOnlyOther {
			args = append(args, "--cflags-only-other")
		}

		if libs {
			args = append(args, "--libs")
		}
		if libsOnlyPath {
			args = append(args, "--libs-only-L")
		}
		if libsOnly {
			args = append(args, "--libs-only-l")
		}
		if libsOnlyOther {
			args = append(args, "--libs-only-other")
		}

		if static {
			args = append(args, "--static")
		}
		if modVersion {
			args = append(args, "--modversion")
		}
		if variable != "" {
			args = append(args, "--variable", variable)
		}
		args = append(args, flag.Args()...)

		ex := exec.Command(fallback, args...)
		ex.Env = os.Environ()
		ex.Stdout = os.Stdout
		ex.Stderr = os.Stderr
		if ex.Run() != nil {
			os.Exit(1)
		}
		os.Exit(0)
	}
	// we need to escape "\"
	for _, f := range outputFlags {
		fmt.Print(strings.ReplaceAll(f, "\\", "\\\\"), " ")
	}
}
