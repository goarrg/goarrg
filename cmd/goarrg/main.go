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
	"sort"
	"strings"

	"goarrg.com/cmd/goarrg/internal/cmd"
	"goarrg.com/cmd/goarrg/internal/cmd/build"
	"goarrg.com/cmd/goarrg/internal/cmd/clean"
	"goarrg.com/cmd/goarrg/internal/cmd/install"
	"goarrg.com/cmd/goarrg/internal/cmd/run"
	"goarrg.com/cmd/goarrg/internal/cmd/test"
	"goarrg.com/debug"
)

var cmds = map[string]*cmd.CMD{
	build.CMD.Name:   build.CMD,
	clean.CMD.Name:   clean.CMD,
	install.CMD.Name: install.CMD,
	run.CMD.Name:     run.CMD,
	test.CMD.Name:    test.CMD,
}

func main() {
	if len(os.Args) < 2 {
		help()
	}

	if cmd, ok := cmds[os.Args[1]]; ok {
		cmd.Exec(os.Args[2:])
	} else if os.Args[1] == "-h" {
		help()
	} else {
		debug.LogE("Invalid command %q", os.Args[1])
		help()
		os.Exit(2)
	}
}

func help() {
	fmt.Fprintf(os.Stderr, "Usage:\n\t"+filepath.Base(os.Args[0])+" [command] [arguments]\n\nCommands:\n")

	sortedCMDs := make([]string, 0, len(cmds))

	for _, cmd := range cmds {
		sortedCMDs = append(sortedCMDs, cmd.Name)
	}

	sort.Strings(sortedCMDs)

	for _, cmd := range sortedCMDs {
		fmt.Fprintf(os.Stderr, "\t%s\n\t\t%s\n", cmds[cmd].Name, strings.ReplaceAll(strings.TrimSpace(cmds[cmd].Short), "\n", "\n\t\t"))
	}

	fmt.Fprintf(os.Stderr, "\nUse \""+filepath.Base(os.Args[0])+" [command] -h\" for more information about that command.\n")
}
