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
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strings"

	"goarrg.com/cmd/goarrg/internal/base"
	"goarrg.com/cmd/goarrg/internal/dep"
	"goarrg.com/cmd/goarrg/internal/exec"
	"goarrg.com/debug"
)

func yourself(args []string) bool {
	if len(args) > 0 {
		return false
	}

	gocache, err := ioutil.TempDir("", "goarrg")

	if err != nil {
		panic(err)
	}

	if err := os.MkdirAll(gocache, 0o755); err != nil {
		panic(err)
	}

	defer os.RemoveAll(gocache)
	if err := os.Setenv("GOCACHE", gocache); err != nil {
		panic(err)
	}

	dep.Build()
	debug.LogI("Building goarrg")

	tags := ""

	if dep.FlagDisableVK() {
		tags += ",disable_vk"
	}

	if dep.FlagDisableGL() {
		tags += ",disable_gl"
	}

	var pkgs []string

	if out, err := exec.RunOutput("go", "list", "-f", "{{.ImportPath}}", "goarrg.com/..."); err != nil {
		panic(err)
	} else {
		tmp := strings.Fields(string(out))
		pkgs = tmp[:0]

		for _, pkg := range tmp {
			if !(strings.Contains(pkg, "/example") || strings.Contains(pkg, "/test")) {
				pkgs = append(pkgs, pkg)
			}
		}
	}

	if err := exec.Run("go", append([]string{"build", "-tags=" + tags}, pkgs...)...); err != nil {
		panic(err)
	}

	if err := exec.Run("go", append([]string{"build", "-tags=" + tags + ",debug"}, pkgs...)...); err != nil {
		panic(err)
	}

	if runtime.GOOS == base.GOOS() {
		os.Setenv("GODEBUG", "cgocheck=2")

		if out, err := exec.RunOutput("go", "test", "-tags="+tags, "-race", "-v", "-count=1", "goarrg.com/..."); err != nil {
			panic(fmt.Sprintf("%v\n\n%s", err, out))
		} else if base.IsVeryVerbose() {
			fmt.Println(string(out))
		}

		if out, err := exec.RunOutput("go", "test", "-tags="+tags+",debug", "-race", "-v", "-count=1", "goarrg.com/..."); err != nil {
			panic(fmt.Sprintf("%v\n\n%s", err, out))
		} else if base.IsVeryVerbose() {
			fmt.Println(string(out))
		}

		os.Setenv("GODEBUG", "")
	}

	debug.LogI("Done building goarrg")

	return true
}
