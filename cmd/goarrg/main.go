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

package main

import (
	"fmt"
	"os"
	"path/filepath"

	"goarrg.com/cmd/goarrg/internal/base"
	"goarrg.com/cmd/goarrg/internal/cmd/build"
	"goarrg.com/cmd/goarrg/internal/cmd/clean"
	"goarrg.com/cmd/goarrg/internal/cmd/run"
	"goarrg.com/cmd/goarrg/internal/cmd/test"
)

var cmds = map[string]*base.CMD{
	clean.CMD.Name: clean.CMD,
	build.CMD.Name: build.CMD,
	run.CMD.Name:   run.CMD,
	test.CMD.Name:  test.CMD,
}

func main() {
	if len(os.Args) < 2 {
		help()
	}

	if cmd, ok := cmds[os.Args[1]]; ok {
		cmd.Exec(os.Args[2:])
	} else {
		if os.Args[1] != "-h" {
			fmt.Fprintf(os.Stderr, "Invalid command: %s\n\n", os.Args[1])
		}
		help()
	}
}

func help() {
	fmt.Fprintf(os.Stderr, "Usage:\n\t"+filepath.Base(os.Args[0])+" [command] [arguments]\n\nCommands:\n")

	for _, cmd := range cmds {
		fmt.Fprintf(os.Stderr, "\t%-8s\n\t\t%s\n", cmd.Name, cmd.Short)
	}

	fmt.Fprintf(os.Stderr, "\nUse \""+filepath.Base(os.Args[0])+" [command] -h\" for more information about that command.\n")
	os.Exit(2)
}
