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
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type cgoFlags struct {
	CFlags        string
	LDFlags       string
	StaticLDFlags string
}

type cgoDepMeta struct {
	Version  string
	CgoFlags cgoFlags
}

var dataDir string

func main() {
	goos := os.Getenv("GOOS")
	goarch := os.Getenv("GOARCH")

	if goos == "" {
		goos = runtime.GOOS
	}

	if goarch == "" {
		goarch = runtime.GOARCH
	}

	target := goos + "_" + goarch

	if dataDir == "" {
		panic("Empty dataDir, please do not build this program manually")
	}

	cFlags := false
	libs := false
	static := false
	version := false
	flag.BoolVar(&cFlags, "cflags", false, "")
	flag.BoolVar(&libs, "libs", false, "")
	flag.BoolVar(&static, "static", false, "")
	flag.BoolVar(&version, "version", false, "")
	flag.Parse()

	if version {
		fmt.Println("0.0.0-goarrg0")
		os.Exit(0)
	}

	pkgList := map[string]struct{}{}

	for _, p := range flag.Args() {
		pkgList[p] = struct{}{}
	}

	for p := range pkgList {
		dir := filepath.Join(dataDir, "cgodep", p, target)
		data, err := os.ReadFile(filepath.Join(dir, "goarrg_cgodep.json"))
		if err != nil {
			if !errors.Is(err, os.ErrNotExist) {
				panic(err)
			}

			dir = filepath.Join(dataDir, "cgodep", p)
			data, err = os.ReadFile(filepath.Join(dir, "goarrg_cgodep.json"))
			if err != nil {
				if !errors.Is(err, os.ErrNotExist) {
					panic(err)
				}

				args := make([]string, 0, 4)

				switch {
				case cFlags:
					args = append(args, "--cflags", "--", p)
				case libs && !static:
					args = append(args, "--libs", "--", p)
				case libs && static:
					args = append(args, "--libs", "--static", "--", p)
				}

				cmd := exec.Command("pkg-config", args...)
				cmd.Stderr = os.Stderr
				cmd.Stdout = os.Stdout
				err := cmd.Run()
				if err != nil {
					if !errors.Is(err, exec.ErrNotFound) {
						panic(err)
					}
					panic(fmt.Errorf("%q not found", p))
				}

				continue
			}
		}

		meta := cgoDepMeta{}
		if err := json.Unmarshal(data, &meta); err != nil {
			panic(err)
		}

		// add "-I{{.InstallDir}}/include" and "-L{{.InstallDir}}/lib"
		switch {
		case cFlags:
			fmt.Print(strings.ReplaceAll(fmt.Sprint(meta.CgoFlags.CFlags, " -I"+filepath.Join(dir, "include")), "\\", "\\\\"), " ")
		case libs && !static:
			fmt.Print(strings.ReplaceAll(fmt.Sprint(meta.CgoFlags.LDFlags, " -L"+filepath.Join(dir, "lib")), "\\", "\\\\"), " ")
		case libs && static:
			fmt.Print(strings.ReplaceAll(fmt.Sprint(meta.CgoFlags.StaticLDFlags, " -L"+filepath.Join(dir, "lib")), "\\", "\\\\"), " ")
		}
	}
}
