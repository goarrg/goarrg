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
	"flag"
	"os"
	"path/filepath"
	"strings"

	"goarrg.com/cmd/goarrg/internal/base"
	"goarrg.com/cmd/goarrg/internal/dep"
	"goarrg.com/cmd/goarrg/internal/exec"
	"goarrg.com/debug"
)

var CMD = &base.CMD{
	Run:   run,
	Name:  "test",
	Short: "Tests project in current directory",
	Long:  "",
	CMDs:  map[string]*base.CMD{},
}

func init() {
	setFlags := func(f *flag.FlagSet) {
		dep.SetFlags(f)
	}

	setFlags(&CMD.Flag)

	for _, cmd := range CMD.CMDs {
		setFlags(&cmd.Flag)
	}
}

func appendTag(args []string, tag string) []string {
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

func run(args []string) bool {
	{
		gocache, err := filepath.Abs(filepath.Join(".", ".goarrg", "gocache"))

		if err != nil {
			panic(err)
		}

		if err := os.MkdirAll(gocache, 0o755); err != nil {
			panic(err)
		}

		if err := os.Setenv("GOCACHE", gocache); err != nil {
			panic(err)
		}
	}

	dep.Build()
	debug.LogI("Testing project")

	for i, arg := range args {
		if arg == "--" {
			args = args[:i]
			break
		}
	}

	args = append([]string{"test", ".", "-race", "-v"}, args...)

	if dep.FlagDisableVK() {
		args = appendTag(args, "disable_vk")
	}

	if dep.FlagDisableGL() {
		args = appendTag(args, "disable_gl")
	}

	if err := exec.Run("go", args...); err != nil {
		panic(err)
	}

	debug.LogI("Done testing project")

	return true
}
