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
	"runtime"
	"strings"

	"goarrg.com/toolchain"
	"goarrg.com/toolchain/cgodep"
)

func main() {
	cFlags := false
	libs := false
	static := false
	version := false
	exists := false
	flag.BoolVar(&cFlags, "cflags", false, "")
	flag.BoolVar(&libs, "libs", false, "")
	flag.BoolVar(&static, "static", false, "")
	flag.BoolVar(&version, "version", false, "")
	flag.BoolVar(&exists, "exists", false, "")
	flag.Bool("print-errors", false, "")
	flag.Bool("short-errors", false, "")
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
		fmt.Println(err)
		os.Exit(1)
	}
	// we need to escape "\"
	for _, f := range outputFlags {
		fmt.Print(strings.ReplaceAll(f, "\\", "\\\\"), " ")
	}
}
